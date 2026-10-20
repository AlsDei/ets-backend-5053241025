# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Setup Kerangka Dasar Proyek
- Prompt: "bantu aku setup kerangka dasar di repo aku."
- Dipakai: `models.go`, `store.go`, `helpers.go`, dan kerangka router `main.go`.
- Diubah: Menyesuaikan tag JSON dan struct agar selaras dengan spesifikasi PinjamKampus serta Persona A.

## B1: Buat Sumber Daya (POST /sumber-daya)
- Prompt: "Sekarang bantu aku buat B1 (POST /sumber-daya). Urutannya: Decode body dengan json.NewDecoder. Kalau error balas 400. strings.TrimSpace pada nama. Validasi nama tidak kosong dan maksimal 50 karakter (utf8.RuneCountInString), jenis ruang atau alat, stok minimal 1 di validasi.go. s.store.mu.Lock() defer Unlock(), beri id dari idSumberDayaNext, tambahkan ke slice, naikkan counter. Balas 201 dengan objek SumberDaya."
- Dipakai: Fungsi `validateCreateSumberDaya` di `validasi.go` dan handler `handleCreateSumberDaya` di `main.go`.
- Diubah: Memastikan penghitungan panjang string menggunakan `utf8.RuneCountInString` agar sesuai ketentuan (per karakter, bukan per byte), dan mengunci mutex sebelum mengakses serta menginkremen `idSumberDayaNext` dan slice `sumberDaya` untuk mencegah race condition.

## B1: Lihat Sumber Daya (GET /sumber-daya dan GET /sumber-daya/{id})
- Prompt: "bisakah kamu kerjakan GET /sumber-daya dan GET /sumber-daya/{id}. setelah kamu kerjakan langsung isi ailog"
- Dipakai: `handleListSumberDaya` dan `handleGetSumberDaya` di `api.go`.
- Diubah: Memastikan pengembalian list kosong berupa `[]` (bukan `null`) dan slice di-copy di dalam critical section mutex (`s.store.mu.Lock()`). Pada `handleGetSumberDaya`, validasi ID dilakukan lebih dahulu menggunakan `parseID` (jika bukan angka positif menghasilkan 400), lalu mencari ke slice sumber daya, dan mengembalikan 404 jika tidak ditemukan.

## B2 & B3: Validasi, Buat Peminjaman, dan Lihat Peminjaman
- Prompt: "Validasi peminjaman di validasi.go (tanggal, slot, sumber_daya_id harus ada, jumlah ≤ stok). POST /peminjaman (belum perlu cek bentrok, itu B5). GET /peminjaman dengan filter, lalu GET /peminjaman/{id}."
- Dipakai: Fungsi `validatePeminjaman` di `validasi.go`, method `findSumberDaya` di `store.go`, serta handler `handleCreatePeminjaman`, `handleListPeminjaman`, dan `handleGetPeminjaman` di `api.go`.
- Diubah: Memvalidasi tanggal menggunakan `time.Parse("2006-01-02", ...)` dan verifikasi kecocokan kalender, mengecek rentang slot 1-8 serta `slot_mulai <= slot_selesai`, mengecek eksistensi sumber daya dan batas stok (status 400). Pada `GET /peminjaman`, parsing filter query opsional (`sumber_daya_id`, `tanggal`, `status`), mengembalikan 400 jika `sumber_daya_id` bukan angka, serta menyortir hasil ascending berdasarkan ID dan menjamin inisialisasi slice kosong mengembalikan `[]`.
