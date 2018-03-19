package filedb

import (
	"errors"
	"sync"
)

type MemoryMetaStorage struct {
	entries      map[string]*Meta
	entriesMutex sync.Mutex
}

func (m MemoryMetaStorage) ListPage(page int) ([]*Meta, error) {
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

func NewMemoryMetaStorage() MemoryMetaStorage {
	return MemoryMetaStorage{
		entries: make(map[string]*Meta, 0),
	}
}

func (m MemoryMetaStorage) FetchMeta(hash string) (*Meta, error) {
	meta, ok := m.entries[hash]

	if !ok {
		return nil, errors.New("file not found")
	}

	return meta, nil
}

func (m MemoryMetaStorage) StoreMeta(meta Meta) error {
	// m.entries[meta.Hash]

	// 	if ok {
	// 		return errors.New("file already exists")
	// 	}

	m.entriesMutex.Lock()
	m.entries[meta.Hash] = &meta
	m.entriesMutex.Unlock()

	return nil
}
