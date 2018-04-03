package filedb

import (
	"errors"
	"sync"
)

type MemoryMetaStorage struct {
	entries      map[string]*Meta
	entriesMutex sync.Mutex
}

func (m *MemoryMetaStorage) ListPage(page int) ([]*Meta, error) {
	var metas []*Meta

	metas = make([]*Meta, 0, len(m.entries))

	for _, m := range m.entries {
		if !m.finished() {
			continue
		}
		metas = append(metas, m)
	}

	return metas, nil
}

func NewMemoryMetaStorage() *MemoryMetaStorage {
	return &MemoryMetaStorage{
		entries: make(map[string]*Meta, 0),
	}
}

func (m *MemoryMetaStorage) FetchMeta(hash string) (*Meta, error) {
	meta, ok := m.entries[hash]

	if !ok {
		return nil, errors.New("file not found")
	}

	return meta, nil
}

func (m *MemoryMetaStorage) FetchMetaWithSlug(slug string) (*Meta, error) {
	for _, e := range m.entries {
		if e.Slug == slug {
			return e, nil
		}
	}

	return nil, errors.New("file not found")
}

func (s *MemoryMetaStorage) StoreMeta(m *Meta) error {
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
