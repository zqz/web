package filedb

import "errors"

type MemoryMetaStorage struct {
	entries map[string]*Meta
}

func NewMemoryMetaStorage() MemoryMetaStorage {
	return MemoryMetaStorage{
		entries: make(map[string]*Meta),
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

	m.entries[meta.Hash] = &meta

	return nil
}
