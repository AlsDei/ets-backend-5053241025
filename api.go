package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
)

// Server membungkus instance Store untuk handler HTTP.
type Server struct {
	store *Store
}

// NewHandler membuat Store baru dan menginisialisasi router http.ServeMux dengan semua route.
func NewHandler() http.Handler {
	s := &Server{
		store: NewStore(),
	}

	mux := http.NewServeMux()

	// Endpoints Sumber Daya (B1)
	mux.HandleFunc("POST /sumber-daya", s.handleCreateSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.handleListSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.handleGetSumberDaya)

	// Endpoint Fitur Khas Persona A (Petugas Gedung: Ketersediaan)
	mux.HandleFunc("GET /sumber-daya/{id}/ketersediaan", s.handleKetersediaanSumberDaya)

	// Endpoints Peminjaman (B2, B3, B6)
	mux.HandleFunc("POST /peminjaman", s.handleCreatePeminjaman)
	mux.HandleFunc("GET /peminjaman", s.handleListPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.handleGetPeminjaman)
	mux.HandleFunc("PUT /peminjaman/{id}", s.handleUpdatePeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", s.handleDeletePeminjaman)

	return mux
}

// B1: POST /sumber-daya
func (s *Server) handleCreateSumberDaya(w http.ResponseWriter, r *http.Request) {
	var req CreateSumberDayaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body request tidak valid")
		return
	}

	if errMsg := validateCreateSumberDaya(&req); errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	sumberDaya := SumberDaya{
		ID:    s.store.idSumberDayaNext,
		Nama:  req.Nama,
		Jenis: req.Jenis,
		Stok:  req.Stok,
	}

	s.store.sumberDaya = append(s.store.sumberDaya, sumberDaya)
	s.store.idSumberDayaNext++

	writeJSON(w, http.StatusCreated, sumberDaya)
}

// B1: GET /sumber-daya (daftar semua sumber daya, urut ID naik, empty array [] jika kosong)
func (s *Server) handleListSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	list := make([]SumberDaya, len(s.store.sumberDaya))
	copy(list, s.store.sumberDaya)

	writeJSON(w, http.StatusOK, list)
}

// B1: GET /sumber-daya/{id} (detail satu sumber daya, 400 jika id bukan angka, 404 jika tidak ditemukan)
func (s *Server) handleGetSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka positif")
		return
	}

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	for _, sd := range s.store.sumberDaya {
		if sd.ID == id {
			writeJSON(w, http.StatusOK, sd)
			return
		}
	}

	writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
}

func (s *Server) handleKetersediaanSumberDaya(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "belum diimplementasi")
}

// B2: POST /peminjaman (membuat peminjaman, status awal 'aktif')
func (s *Server) handleCreatePeminjaman(w http.ResponseWriter, r *http.Request) {
	var req PeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body request tidak valid")
		return
	}

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	sd := s.store.findSumberDaya(req.SumberDayaID)
	if errMsg := validatePeminjaman(&req, sd); errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	peminjaman := Peminjaman{
		ID:           s.store.idPeminjamanNext,
		SumberDayaID: req.SumberDayaID,
		Peminjam:     req.Peminjam,
		Tanggal:      req.Tanggal,
		SlotMulai:    req.SlotMulai,
		SlotSelesai:  req.SlotSelesai,
		Jumlah:       req.Jumlah,
		Status:       "aktif",
	}

	s.store.peminjaman[peminjaman.ID] = peminjaman
	s.store.idPeminjamanNext++

	writeJSON(w, http.StatusCreated, peminjaman)
}

// B3: GET /peminjaman (daftar peminjaman dengan filter opsional ?sumber_daya_id=, ?tanggal=, ?status=, urut id naik)
func (s *Server) handleListPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var filterSDID *int
	if val := q.Get("sumber_daya_id"); val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		filterSDID = &id
	}

	filterTanggal := q.Get("tanggal")
	filterStatus := q.Get("status")

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	list := make([]Peminjaman, 0)
	for _, p := range s.store.peminjaman {
		if filterSDID != nil && p.SumberDayaID != *filterSDID {
			continue
		}
		if filterTanggal != "" && p.Tanggal != filterTanggal {
			continue
		}
		if filterStatus != "" && p.Status != filterStatus {
			continue
		}
		list = append(list, p)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	writeJSON(w, http.StatusOK, list)
}

// B3: GET /peminjaman/{id} (detail satu peminjaman)
func (s *Server) handleGetPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka positif")
		return
	}

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	p, exists := s.store.peminjaman[id]
	if !exists {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handleUpdatePeminjaman(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "belum diimplementasi")
}

func (s *Server) handleDeletePeminjaman(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "belum diimplementasi")
}
