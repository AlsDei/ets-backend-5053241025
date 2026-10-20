# PROSES

## 1. Persona dan masalahnya

Persona saya adalah petugas gedung. ketika ingin mencari tahu kapan suatu ruangan selesai digunakan, petuga gedung perlu membuka daftar peminjaman satu per satu. saya akan  membuat endpoint yang menampilkan kapan ruangan selesai digunakan, sehingga petugas gedung tidak perlu membuka daftar peminjaman satu per satu.

## 2. Rancangan struktur data
Storage {
  mu               sync.Mutex        // satu kunci untuk semua data bersama
  sumberDaya       []SumberDaya      // slice; id = urutan + 1, tidak pernah dihapus
  peminjaman       map[int]Peminjaman // kunci = id; mudah dihapus dan dicari
  idPeminjamanNext int               // mulai 1, naik tiap peminjaman baru yang BERHASIL
  idSumberDayaNext int               // mulai 1
}

SumberDaya  { id, nama, jenis, stok }
Peminjaman  { id, sumber_daya_id, peminjam, tanggal, slot_mulai, slot_selesai, jumlah, status }

## 3. Satu keputusan desain

Keputusan: satu fungsi yang menerima (sumber daya, tanggal, id yang dikecualikan) dan mengembalikan array jumlah unit terpakai untuk slot 1-8. Fungsi ini dipakai oleh cek bentrok (B5/B6) dan endpoint ketersediaan (Persona A).

Alternatif:
(1) menghitung terpisah di dua tempat; 
(2) menyimpan grid pemakaian map[sumber_daya+tanggal][8]int yang diperbarui setiap POST/PUT/DELETE.

Alasan memilih: satu sumber kebenaran, jadi ketersediaan tidak mungkin bilang "kosong" padahal POST ditolak 409. Parameter id yang dikecualikan membuat PUT tidak bentrok dengan dirinya sendiri. Untuk data di memori, memindai peminjaman per request sudah cukup cepat.

## 4. Bug atau kendala (minimal dua)

### Bug 1

- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2

- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request

- Apa yang diminta: Menyediakan endpoint `GET /ringkasan?tanggal=YYYY-MM-DD` yang mengembalikan ringkasan harian: jumlah peminjaman aktif, total unit slot terpakai (dihitung dari jumlah * durasi slot), dan daftar peminjaman per sumber daya yang mencakup semua sumber daya berurutan berdasarkan id naik (bernilai 0 jika tidak dipinjam).
- Bagian kode mana yang kamu ubah:
  1. Menambahkan struct `RingkasanSumberDaya` dan `RingkasanResponse` di `models.go`.
  2. Mendaftarkan route `GET /ringkasan` di `NewHandler()` pada `api.go`.
  3. Mengimplementasikan handler `handleGetRingkasan` di `api.go` dengan validasi query tanggal menggunakan helper `validateTanggal()`.
- Satu keputusan desain yang kamu ambil: Menggunakan map sementara `peminjamanPerSD[int]int` saat mengagregasi peminjaman aktif di tanggal tersebut, kemudian melakukan iterasi terhadap `s.store.sumberDaya` untuk menyusun `per_sumber_daya`. Dengan cara ini, semua sumber daya (termasuk yang tidak memiliki peminjaman pada hari itu) secara otomatis tercakup dengan nilai `0` dan urutan ID tetap terjamin ascending tanpa memerlukan kueri berulang.
