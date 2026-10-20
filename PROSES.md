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

(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
