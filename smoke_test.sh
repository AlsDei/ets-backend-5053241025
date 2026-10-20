#!/usr/bin/env bash
# Tes fitur khas Persona A: GET /sumber-daya/{id}/ketersediaan?tanggal=
# Jalankan di server BARU (data kosong):  bash persona_a_test.sh
# Butuh: bash, curl, python3. Alat bantu uji, bukan bagian dari repo.

BASE=${BASE:-http://localhost:8080}
PYBIN=$(command -v python3 || command -v python)
if [ -z "$PYBIN" ]; then echo "Python tidak ditemukan (butuh python3 atau python)"; exit 1; fi
pass=0; fail=0
TMP=$(mktemp); PY=$(mktemp --suffix=.py)

cat > "$PY" << 'PYEOF'
import sys, json
rid, tgl, want = int(sys.argv[1]), sys.argv[2], [int(x) for x in sys.argv[3].split(",")]
try:
    d = json.load(sys.stdin)
except Exception:
    print("respons bukan JSON valid"); sys.exit(1)
errs = []
if d.get("sumber_daya_id") != rid: errs.append("sumber_daya_id salah: %r" % d.get("sumber_daya_id"))
if d.get("tanggal") != tgl: errs.append("tanggal salah: %r" % d.get("tanggal"))
s = d.get("slot")
if not isinstance(s, list) or len(s) != 8:
    errs.append("field slot harus array berisi 8 elemen")
else:
    if [x.get("slot") for x in s] != list(range(1, 9)): errs.append("nomor slot harus 1..8 berurutan")
    got = [x.get("sisa") for x in s]
    if got != want: errs.append("sisa %s != harapan %s" % (got, want))
    for x in s:
        if x.get("tersedia") != ((x.get("sisa") or 0) > 0):
            errs.append("tersedia tidak konsisten di slot %s" % x.get("slot")); break
print("; ".join(errs)); sys.exit(1 if errs else 0)
PYEOF

check() { # name want METHOD path [body]
  local name=$1 want=$2 method=$3 path=$4 body=$5 got
  if [ -n "$body" ]; then
    got=$(curl -s -o "$TMP" -w '%{http_code}' -X "$method" -H 'Content-Type: application/json' -d "$body" "$BASE$path")
  else
    got=$(curl -s -o "$TMP" -w '%{http_code}' -X "$method" "$BASE$path")
  fi
  if [ "$got" = "$want" ]; then pass=$((pass+1)); echo "  ok    $name"
  else fail=$((fail+1)); echo "  FAIL  $name (harapan $want, dapat $got) -> $(head -c 120 "$TMP")"; fi
}

avail() { # name rid tanggal "sisa1,...,sisa8"
  local name=$1 rid=$2 tgl=$3 want=$4 msg
  msg=$(curl -s "$BASE/sumber-daya/$rid/ketersediaan?tanggal=$tgl" | "$PYBIN" "$PY" "$rid" "$tgl" "$want")
  if [ $? -eq 0 ]; then pass=$((pass+1)); echo "  ok    $name"
  else fail=$((fail+1)); echo "  FAIL  $name -> $msg"; fi
}

P() { echo "{\"sumber_daya_id\":$1,\"peminjam\":\"$2\",\"tanggal\":\"$3\",\"slot_mulai\":$4,\"slot_selesai\":$5,\"jumlah\":$6}"; }

echo "== Persiapan =="
check "buat alat Proyektor stok 3 (id 1)" 201 POST /sumber-daya '{"nama":"Proyektor","jenis":"alat","stok":3}'
check "buat ruang Aula stok 1 (id 2)"     201 POST /sumber-daya '{"nama":"Aula A","jenis":"ruang","stok":1}'

echo "== Urutan pemeriksaan & validasi =="
check "id bukan angka"                        400 GET '/sumber-daya/abc/ketersediaan?tanggal=2026-11-02'
check "id bukan angka, tanpa tanggal (400)"   400 GET '/sumber-daya/abc/ketersediaan'
check "sumber daya tidak ada"                 404 GET '/sumber-daya/999/ketersediaan?tanggal=2026-11-02'
check "tidak ada + tanggal kosong -> 404 dulu" 404 GET '/sumber-daya/999/ketersediaan'
check "tidak ada + tanggal jelek -> 404 dulu"  404 GET '/sumber-daya/999/ketersediaan?tanggal=xxx'
check "tanggal tidak diberikan"               400 GET '/sumber-daya/1/ketersediaan'
check "tanggal kosong (?tanggal=)"            400 GET '/sumber-daya/1/ketersediaan?tanggal='
check "tanggal 30 Februari"                   400 GET '/sumber-daya/1/ketersediaan?tanggal=2026-02-30'
check "tanggal format salah"                  400 GET '/sumber-daya/1/ketersediaan?tanggal=02-11-2026'
check "29 Feb 2028 diterima"                  200 GET '/sumber-daya/1/ketersediaan?tanggal=2028-02-29'

echo "== Content-Type & bentuk respons =="
ct=$(curl -s -D - -o /dev/null "$BASE/sumber-daya/1/ketersediaan?tanggal=2026-11-02" | grep -i '^content-type:' | tr -d '\r')
if echo "$ct" | grep -qi 'application/json'; then pass=$((pass+1)); echo "  ok    Content-Type application/json"
else fail=$((fail+1)); echo "  FAIL  Content-Type bukan application/json ($ct)"; fi
avail "belum ada peminjaman: 8 slot, sisa=stok"  1 2026-11-02 "3,3,3,3,3,3,3,3"
avail "ruang stok 1: semua sisa 1"               2 2026-11-02 "1,1,1,1,1,1,1,1"

echo "== Mencerminkan peminjaman =="
check "POST alat 2 unit, slot 3-4"        201 POST /peminjaman "$(P 1 Rina 2026-11-02 3 4 2)"
avail "slot 3-4 sisa 1"                   1 2026-11-02 "3,3,1,1,3,3,3,3"
check "POST alat 1 unit, slot 4-5"        201 POST /peminjaman "$(P 1 Budi 2026-11-02 4 5 1)"
avail "slot 4 sisa 0, slot 5 sisa 2"      1 2026-11-02 "3,3,1,0,2,3,3,3"
check "POST 409 (slot 4 sudah penuh)"     409 POST /peminjaman "$(P 1 Ani 2026-11-02 4 4 1)"
avail "peminjaman gagal tidak mengubah"   1 2026-11-02 "3,3,1,0,2,3,3,3"
avail "tanggal lain tidak terpengaruh"    1 2026-11-03 "3,3,3,3,3,3,3,3"
avail "sumber daya lain tidak terpengaruh" 2 2026-11-02 "1,1,1,1,1,1,1,1"

echo "== Mencerminkan PUT dan DELETE =="
check "PUT peminjaman 1 -> slot 1-1, 1 unit" 200 PUT /peminjaman/1 "$(P 1 Rina 2026-11-02 1 1 1)"
avail "setelah PUT"                       1 2026-11-02 "2,3,3,2,2,3,3,3"
check "PUT ke tanggal lain"               200 PUT /peminjaman/1 "$(P 1 Rina 2026-11-03 1 1 1)"
avail "tanggal lama pulih"                1 2026-11-02 "3,3,3,2,2,3,3,3"
avail "tanggal baru terisi"               1 2026-11-03 "2,3,3,3,3,3,3,3"
check "DELETE peminjaman 2"               204 DELETE /peminjaman/2
avail "setelah DELETE, 2026-11-02 kosong" 1 2026-11-02 "3,3,3,3,3,3,3,3"

echo "== Kasus penuh 8 slot =="
check "POST 3 unit slot 1-8"              201 POST /peminjaman "$(P 1 Citra 2026-12-01 1 8 3)"
avail "semua slot sisa 0 (tersedia=false)" 1 2026-12-01 "0,0,0,0,0,0,0,0"
check "DELETE peminjaman penuh"           204 DELETE /peminjaman/3
avail "pulih setelah dihapus"             1 2026-12-01 "3,3,3,3,3,3,3,3"

rm -f "$TMP" "$PY"
echo; echo "Hasil: $pass lulus, $fail gagal"
[ "$fail" -eq 0 ]
