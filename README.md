Deskripsi sistem:
Berikut adalah mesin kasir berbasis teks menggunakan bahasa Go (Golang).
Aplikasi ini akan digunakan di minimarket untuk mempermudah pengelolaan barang, transaksi, dan stok.
Berikut adalah fitur-fitur yang tersedia:

1. **Menu Inventory (CRUD)**: 
    - Tambah barang baru.
    - Lihat daftar barang beserta detailnya (kode barang, nama barang, harga, stok).
    - Edit data barang tertentu.
    - Hapus barang dari inventory.

2. **Menu Gudang**:
    - Cek stok barang per kategori.
    - Tambah stok ke barang tertentu.
    - Lihat barang dengan stok rendah.

3. **Menu Transaksi**:
    - Buat transaksi baru dengan memasukkan kode barang dan jumlah.
    - Hitung total harga secara otomatis.
    - Cetak struk transaksi sederhana.

4. **Menu Laporan**:
    - Tampilkan laporan penjualan harian.
    - Tampilkan total pendapatan dalam periode tertentu.
    - Simpan laporan ke file.

5. **Autentikasi Pengguna**:
    - Terdapat login untuk admin dan kasir.
    - Admin memiliki akses penuh ke semua fitur.
    - Kasir hanya dapat mengakses fitur transaksi dan laporan penjualan.
