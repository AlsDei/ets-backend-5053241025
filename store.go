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

// hitungPemakaianSlot menghitung total unit terpakai pada slot 1-8 untuk sumber daya dan tanggal tertentu,
// mengecualikan peminjaman dengan ID excludeID (berguna untuk PUT dan ketersediaan).
func (s *Store) hitungPemakaianSlot(sumberDayaID int, tanggal string, excludeID int) [9]int {
	var pemakaian [9]int
	for _, p := range s.peminjaman {
		if p.ID == excludeID {
			continue
		}
		if p.SumberDayaID == sumberDayaID && p.Tanggal == tanggal && p.Status == "aktif" {
			for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
				if slot >= 1 && slot <= 8 {
					pemakaian[slot] += p.Jumlah
				}
			}
		}
	}
	return pemakaian
}

// cekBentrok memeriksa apakah peminjaman baru/update melebihi stok yang ada pada rentang slot yang diminta.
// Mengembalikan true jika terjadi bentrok (kapasitas terlampaui).
func (s *Store) cekBentrok(sumberDayaID int, tanggal string, slotMulai, slotSelesai, jumlah int, stok int, excludeID int) bool {
	pemakaian := s.hitungPemakaianSlot(sumberDayaID, tanggal, excludeID)
	for slot := slotMulai; slot <= slotSelesai; slot++ {
		if slot >= 1 && slot <= 8 {
			if pemakaian[slot]+jumlah > stok {
				return true
			}
		}
	}
	return false
}
