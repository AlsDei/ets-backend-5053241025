package main

// SumberDaya merepresentasikan ruangan atau alat yang dapat dipinjam.
type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

// CreateSumberDayaRequest merepresentasikan payload body untuk POST /sumber-daya.
type CreateSumberDayaRequest struct {
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

// Peminjaman merepresentasikan data peminjaman ruangan atau alat.
type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

// PeminjamanRequest merepresentasikan payload body untuk POST /peminjaman dan PUT /peminjaman/{id}.
type PeminjamanRequest struct {
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
}

// Persona A: Petugas Gedung
// SlotKetersediaan merepresentasikan sisa unit dan ketersediaan pada satu slot (1-8).
type SlotKetersediaan struct {
	Slot     int  `json:"slot"`
	Sisa     int  `json:"sisa"`
	Tersedia bool `json:"tersedia"`
}

// KetersediaanResponse merepresentasikan payload respons untuk GET /sumber-daya/{id}/ketersediaan.
type KetersediaanResponse struct {
	SumberDayaID int                `json:"sumber_daya_id"`
	Tanggal      string             `json:"tanggal"`
	Slot         []SlotKetersediaan `json:"slot"`
}

// Change Request A: Dashboard Ringkasan Harian
// RingkasanSumberDaya merepresentasikan ringkasan jumlah peminjaman per sumber daya pada tanggal tertentu.
type RingkasanSumberDaya struct {
	SumberDayaID     int    `json:"sumber_daya_id"`
	Nama             string `json:"nama"`
	JumlahPeminjaman int    `json:"jumlah_peminjaman"`
}

// RingkasanResponse merepresentasikan payload respons untuk GET /ringkasan?tanggal=YYYY-MM-DD.
type RingkasanResponse struct {
	Tanggal          string                `json:"tanggal"`
	JumlahPeminjaman int                   `json:"jumlah_peminjaman"`
	TotalUnitSlot    int                   `json:"total_unit_slot"`
	PerSumberDaya    []RingkasanSumberDaya `json:"per_sumber_daya"`
}

// ErrorResponse merepresentasikan format respons standar untuk error.
type ErrorResponse struct {
	Error string `json:"error"`
}
