# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Setup Kerangka Dasar Proyek
- Prompt: "bantu aku setup kerangka dasar di repo aku"
- Dipakai: `models.go`, `store.go`, `helpers.go`, dan `main.go`.
- Diubah: Menyesuaikan struct dan tag JSON dengan kontrak PinjamKampus serta Persona A, dan menyiapkan routing standar `http.NewServeMux` Go 1.22.

## B1: Kelola Sumber Daya
- Prompt: "Langkah berikutnya: B1 (POST /sumber-daya) ... dan GET /sumber-daya"
- Dipakai: `validateCreateSumberDaya` di `validasi.go`, `handleCreateSumberDaya`, `handleListSumberDaya`, dan `handleGetSumberDaya` di `api.go`.
- Diubah: Menghitung panjang nama dengan `utf8.RuneCountInString` (karakter, bukan byte), memastikan slice kosong mengembalikan `[]`, serta memvalidasi `parseID` sebelum mencari data.

## B2 & B3: Buat dan Lihat Peminjaman
- Prompt: "Validasi peminjaman di validasi.go, POST /peminjaman, GET /peminjaman dengan filter, dan GET /peminjaman/{id}"
- Dipakai: `validatePeminjaman` di `validasi.go`, `findSumberDaya` di `store.go`, serta handler peminjaman di `api.go`.
- Diubah: Menambahkan validasi kalender dengan `time.Parse`, filter query opsional (`sumber_daya_id`, `tanggal`, `status`), dan memastikan hasil list selalu urut ID naik.

## Fitur Khas Persona A: Ketersediaan Sumber Daya (GET /sumber-daya/{id}/ketersediaan)
- Prompt: "selesaikan khas persona dahulu"
- Dipakai: `validateTanggal` di `validasi.go`, `hitungPemakaianSlot` di `store.go`, dan `handleKetersediaanSumberDaya` di `api.go`.
- Diubah: Menerapkan urutan validasi ketat (id 400 -> sumber daya 404 -> tanggal 400) dan menghitung sisa unit untuk slot 1–8 dari peminjaman berstatus aktif.

## B5: Deteksi Bentrok Kapasitas Slot
- Prompt: "lanjut kerjakan b5"
- Dipakai: Method `cekBentrok` di `store.go` dan integrasi di `handleCreatePeminjaman`.
- Diubah: Menggunakan `hitungPemakaianSlot` sebagai sumber kebenaran tunggal untuk menolak request dengan status 409 jika unit terpakai + unit baru melebihi stok pada rentang slot.

## B6 & B7: Ubah/Hapus Peminjaman dan Keamanan Konkurensi
- Prompt: "coba sekarang bantu aku menerjakan b6-b7"
- Dipakai: `handleUpdatePeminjaman` dan `handleDeletePeminjaman` di `api.go`.
- Diubah: Pada PUT, mengecualikan ID peminjaman itu sendiri dari cek bentrok agar tidak bentrok dengan diri sendiri, menjaga `id` dan `status` tetap, serta membungkus semua operasi data dalam `s.store.mu.Lock()` untuk mencegah race condition.

## Change Request A: Dashboard Ringkasan Harian (GET /ringkasan)
- Prompt: "implementasi request change CR-A"
- Dipakai: Struct ringkasan di `models.go` dan `handleGetRingkasan` di `api.go`.
- Diubah: Menghitung `total_unit_slot` dengan rumus `jumlah * durasi slot`, serta menyertakan semua sumber daya terurut ID naik (dengan nilai 0 bila belum ada peminjaman).
