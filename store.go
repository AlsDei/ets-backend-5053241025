package main

import "sync"

// Store menyimpan data dalam memori dan mengelola ID otomatis secara thread-safe.
type Store struct {
	mu               sync.Mutex
	sumberDaya       []SumberDaya
	peminjaman       map[int]Peminjaman
	idSumberDayaNext int
	idPeminjamanNext int
}

// NewStore membuat instance Store baru dengan state awal kosong dan ID dimulai dari 1.
func NewStore() *Store {
	return &Store{
		sumberDaya:       make([]SumberDaya, 0),
		peminjaman:       make(map[int]Peminjaman),
		idSumberDayaNext: 1,
		idPeminjamanNext: 1,
	}
}

// findSumberDaya mencari sumber daya berdasarkan ID. Pemanggil harus memegang kunci mu jika perlu thread-safety.
func (s *Store) findSumberDaya(id int) *SumberDaya {
	for i := range s.sumberDaya {
		if s.sumberDaya[i].ID == id {
			return &s.sumberDaya[i]
		}
	}
	return nil
}
