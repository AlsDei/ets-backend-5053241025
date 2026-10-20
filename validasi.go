package main

import (
	"strings"
	"time"
	"unicode/utf8"
)

// validateCreateSumberDaya memvalidasi request pembuatan sumber daya.
// Mengembalikan pesan error jika tidak valid, atau string kosong jika valid.
func validateCreateSumberDaya(req *CreateSumberDayaRequest) string {
	req.Nama = strings.TrimSpace(req.Nama)

	if req.Nama == "" {
		return "nama wajib diisi dan tidak boleh kosong"
	}

	if utf8.RuneCountInString(req.Nama) > 50 {
		return "nama maksimal 50 karakter"
	}

	if req.Jenis != "ruang" && req.Jenis != "alat" {
		return "jenis harus berupa 'ruang' atau 'alat'"
	}

	if req.Stok < 1 {
		return "stok minimal 1"
	}

	return ""
}

// validatePeminjaman memvalidasi payload request peminjaman.
// Memeriksa peminjam, format tanggal kalender, rentang slot 1-8, keberadaan sumber daya, dan jumlah <= stok.
func validatePeminjaman(req *PeminjamanRequest, sd *SumberDaya) string {
	req.Peminjam = strings.TrimSpace(req.Peminjam)
	if req.Peminjam == "" {
		return "peminjam wajib diisi dan tidak boleh kosong"
	}
	if utf8.RuneCountInString(req.Peminjam) > 50 {
		return "peminjam maksimal 50 karakter"
	}

	t, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil || t.Format("2006-01-02") != req.Tanggal {
		return "format tanggal harus YYYY-MM-DD dan merupakan tanggal valid"
	}

	if req.SlotMulai < 1 || req.SlotSelesai > 8 || req.SlotMulai > req.SlotSelesai {
		return "slot harus bernilai antara 1 sampai 8 dan slot_mulai harus <= slot_selesai"
	}

	if sd == nil {
		return "sumber daya tidak ditemukan"
	}

	if req.Jumlah < 1 || req.Jumlah > sd.Stok {
		return "jumlah harus antara 1 dan stok sumber daya"
	}

	return ""
}
