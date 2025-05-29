package file

import (
	"errors"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type MemoryMetaStorage struct {
	entries      map[string]*Meta
	entriesMutex sync.Mutex
}

func (m *MemoryMetaStorage) ListPage(page int) ([]*Meta, error) {
	var metas []*Meta

	metas = make([]*Meta, 0, len(m.entries))

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
		entries: make(map[string]*Meta, 0),
	}
}

func (m *MemoryMetaStorage) DeleteMetaById(id int) error {
	return nil
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

	spew.Dump(s.entries)

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

func (s *MemoryMetaStorage) UpdateMeta(m *Meta) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	return nil
}

func (m *MemoryMetaStorage) List(size int) ([]*Meta, error) {
	return nil, nil
}

func (m *MemoryMetaStorage) ListFilesByUserId(size, offset int) ([]*Meta, error) {
	files := make([]*Meta, 0)
	for _, f := range m.entries {
		files = append(files, f)
	}

	return files, nil
}

func (m *MemoryMetaStorage) RemoveThumbnails(x *Meta) error {
	return nil
}

func (m *MemoryMetaStorage) StoreThumbnail(s string, sz int, x *Meta) error {
	return nil
}
