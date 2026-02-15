package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/zerolog"

	"github.com/zqz/web/backend/internal/config"
	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/handler"
	"github.com/zqz/web/backend/internal/service"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	userKey   contextKey = "user"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userSvc *service.UserService
	logger  *zerolog.Logger
	store   *sessions.CookieStore
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userSvc *service.UserService, logger *zerolog.Logger, cfg *config.Config) *AuthHandler {
	// Setup session store
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = cfg.IsProduction()
	store.Options.SameSite = http.SameSiteDefaultMode

	gothic.Store = store

	// Setup OAuth providers
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		goth.UseProviders(
			google.New(
				cfg.GoogleClientID,
				cfg.GoogleClientSecret,
				cfg.GoogleCallbackURL,
			),
		)
		logger.Info().Msg("Google OAuth provider configured")
	}

	return &AuthHandler{
		userSvc: userSvc,
		logger:  logger,
		store:   store,
	}
}

// BeginAuth starts the OAuth flow
func (h *AuthHandler) BeginAuth(w http.ResponseWriter, r *http.Request) {
	// Set provider in query if not present
	q := r.URL.Query()
	if q.Get("provider") == "" {
		q.Set("provider", "google")
		r.URL.RawQuery = q.Encode()
	}

	gothic.BeginAuthHandler(w, r)
}

// CallbackAuth handles the OAuth callback
func (h *AuthHandler) CallbackAuth(w http.ResponseWriter, r *http.Request) {
	authUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to complete auth")
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Get or create user
	// Use email as name fallback if name is empty
	name := authUser.Name
	if name == "" {
		name = authUser.Email
	}

	user, err := h.userSvc.GetOrCreateUser(r.Context(), domain.CreateUserRequest{
		Email:      authUser.Email,
		Name:       name,
		Provider:   "google",
		ProviderID: authUser.UserID,
		Role:       "member",
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get or create user")
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Store user ID in session
	session, err := h.store.Get(r, "auth-session")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get session")
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = user.ID
	if err := session.Save(r, w); err != nil {
		h.logger.Error().Err(err).Msg("failed to save session")
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int32("user_id", user.ID).Str("email", user.Email).Msg("user logged in")

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Logout logs out the user
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session
	session, err := h.store.Get(r, "auth-session")
	if err == nil {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}

	// Clear gothic session
	gothic.Logout(w, r)

	h.logger.Info().Msg("user logged out")

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// meResponse is the JSON shape for GET /auth/me
type meResponse struct {
	Authenticated bool   `json:"authenticated"`
	ID            int32  `json:"id,omitempty"`
	Email         string `json:"email,omitempty"`
	Name          string `json:"name,omitempty"`
	Admin         bool   `json:"admin,omitempty"`
	DisplayTag    string `json:"display_tag,omitempty"`
	Colour        string `json:"colour,omitempty"`
}

// CurrentUser returns the current user info as JSON
func (h *AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == nil {
		handler.SetContentType(w, handler.ContentTypeJSON)
		json.NewEncoder(w).Encode(meResponse{Authenticated: false})
		return
	}

	user, err := h.userSvc.GetUserByID(r.Context(), *userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get user")
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	handler.SetContentType(w, handler.ContentTypeJSON)
	json.NewEncoder(w).Encode(meResponse{
		Authenticated: true,
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Admin:         user.IsAdmin(),
		DisplayTag:    user.DisplayTag,
		Colour:        user.Colour,
	})
}

// UpdateProfileRequest is the JSON body for PATCH /auth/profile
type UpdateProfileRequest struct {
	DisplayTag *string `json:"display_tag"`
	Colour     *string `json:"colour"`
}

// UpdateProfile updates the current user's display tag and colour
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tag := ""
	if req.DisplayTag != nil {
		tag = *req.DisplayTag
	}
	colour := ""
	if req.Colour != nil {
		colour = *req.Colour
	}

	user, err := h.userSvc.UpdateProfile(r.Context(), *userID, tag, colour)
	if err != nil {
		h.logger.Warn().Err(err).Msg("update profile failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handler.SetContentType(w, handler.ContentTypeJSON)
	json.NewEncoder(w).Encode(meResponse{
		Authenticated: true,
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Admin:         user.IsAdmin(),
		DisplayTag:    user.DisplayTag,
		Colour:        user.Colour,
	})
}

// AuthMiddleware checks if user is authenticated and loads user into context
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.store.Get(r, "auth-session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userIDVal, ok := session.Values["user_id"]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		var userID int32
		switch v := userIDVal.(type) {
		case int:
			userID = int32(v)
		case int32:
			userID = v
		case int64:
			userID = int32(v)
		default:
			next.ServeHTTP(w, r)
			return
		}

		// Load user from database
		user, err := h.userSvc.GetUserByID(r.Context(), userID)
		if err != nil {
			h.logger.Warn().Err(err).Int32("user_id", userID).Msg("failed to load user from session")
			next.ServeHTTP(w, r)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth middleware requires authentication. If unauthorizedHandler is non-nil and request accepts HTML, it is used instead of plain 401.
func (h *AuthHandler) RequireAuth(next http.Handler, unauthorizedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserIDFromContext(r.Context()) == nil {
			if unauthorizedHandler != nil && acceptsHTML(r) {
				unauthorizedHandler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin middleware requires admin role. If forbiddenHandler is non-nil and request accepts HTML, it is used instead of plain 403.
func (h *AuthHandler) RequireAdmin(next http.Handler, forbiddenHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil || !user.IsAdmin() {
			if forbiddenHandler != nil && acceptsHTML(r) {
				forbiddenHandler.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func acceptsHTML(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

// NewRouter creates auth routes
func NewRouter(h *AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/login", h.BeginAuth)
	r.Get("/google/callback", h.CallbackAuth)
	r.Get("/logout", h.Logout)
	r.Get("/me", h.CurrentUser)
	r.Patch("/profile", h.UpdateProfile)

	return r
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx context.Context) *int32 {
	if userID, ok := ctx.Value(userIDKey).(int32); ok {
		return &userID
	}
	return nil
}

// GetUserFromContext retrieves user from context
func GetUserFromContext(ctx context.Context) *domain.User {
	if user, ok := ctx.Value(userKey).(*domain.User); ok {
		return user
	}
	return nil
}
