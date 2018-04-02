package filedb

import (
	"errors"
	"sync"
)

type MemoryThumbnailStorage struct {
	entries      map[int]Thumbnail
	entriesMutex sync.Mutex
}

func NewMemoryThumbnailStorage() MemoryThumbnailStorage {
	return MemoryThumbnailStorage{
		entries: make(map[int]Thumbnail, 0),
	}
}

func (s MemoryThumbnailStorage) StoreThumbnail(t Thumbnail) error {
	if t.MetaID == 0 {
		return errors.New("thumbnail missing required meta id")
	}

	s.entriesMutex.Lock()
	s.entries[t.MetaID] = t
	s.entriesMutex.Unlock()

	return nil
}

func contains(n int, hs []int) bool {
	for i := 0; i < len(hs); i++ {
		if hs[i] == n {
			return true
		}
	}

	return false
}

func (s MemoryThumbnailStorage) FetchThumbnails(ids []int) (map[int]Thumbnail, error) {
	ts := make(map[int]Thumbnail, 0)

	for _, t := range s.entries {
		if contains(t.MetaID, ids) {
			ts[t.MetaID] = t
		}
	}

	return ts, nil
}
