package file

import (
	"errors"
	"sync"
)

type MemoryMetaStorage struct {
	entries      map[string]*File
	entriesMutex sync.Mutex
}

func (m *MemoryMetaStorage) ListPage(page int) ([]*File, error) {
	var metas []*File

	metas = make([]*File, 0, len(m.entries))

	for _, m := range m.entries {
		if !m.Finished() {
			continue
		}
		metas = append(metas, m)
	}

	return metas, nil
}

func NewMemoryMetaStorage() *MemoryMetaStorage {
	return &MemoryMetaStorage{
		entries: make(map[string]*File, 0),
	}
}

func (m *MemoryMetaStorage) DeleteById(id int) error {
	return nil
}

func (m *MemoryMetaStorage) FetchByHash(hash string) (*File, error) {
	meta, ok := m.entries[hash]

	if !ok {
		return nil, errors.New("file not found")
	}

	return meta, nil
}

func (m *MemoryMetaStorage) FetchBySlug(slug string) (*File, error) {
	for _, e := range m.entries {
		if e.Slug == slug {
			return e, nil
		}
	}

	return nil, errors.New("file not found")
}

func (s *MemoryMetaStorage) Create(m *File) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	if m.ID == 0 {
		m.ID = nextId()
	}

	return nil
}

var idMutex sync.Mutex
var currentId int

func nextId() int {
	idMutex.Lock()
	currentId++
	idMutex.Unlock()
	return currentId
}

func (s *MemoryMetaStorage) UpdateMeta(m *File) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	return nil
}

func (m *MemoryMetaStorage) List(size int) ([]*File, error) {
	return nil, nil
}

func (m *MemoryMetaStorage) ListFilesByUserId(size, offset int) ([]*File, error) {
	files := make([]*File, 0)
	for _, f := range m.entries {
		files = append(files, f)
	}

	return files, nil
}

func (m *MemoryMetaStorage) RemoveThumbnails(x *File) error {
	return nil
}

func (m *MemoryMetaStorage) StoreThumbnail(s string, sz int, x *File) error {
	return nil
}
