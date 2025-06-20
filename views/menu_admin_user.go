package views

import "fmt"

func PrintHeaderAdmin() {
	fmt.Println("====================================================================")
	// fmt.Println("User: " + userID)
	PrintMenuTransaksi()
	PrintMenuLaporan()
	PrintMenuInventory()
	PrintMenuGudang()
	fmt.Println("0. Keluar")
	fmt.Println("====================================================================")
}

func PrintHeaderUserBiasa() {
	fmt.Println("====================================================================")
	// fmt.Println("User: " + userID)
	PrintMenuTransaksi()
	PrintMenuLaporan()
	fmt.Println("0. Keluar")
	fmt.Println("====================================================================")
}

func PrintMenuTransaksi() {
	fmt.Println("Menu Transkasi:")
	fmt.Println("1. Tambah transaksi")
	fmt.Println("")
}

func PrintMenuLaporan() {
	fmt.Println("Menu Laporan:")
	fmt.Println("2. Laporan stok barang")
	fmt.Println("3. Laporan penjualan harian")
	fmt.Println("4. Laporan penjualan per periode")
	fmt.Println("")
}

func PrintMenuInventory() {
	fmt.Println("Menu Inventory (CRUD):")
	fmt.Println("5. Tambah barang baru")
	fmt.Println("6. Lihat daftar barang")
	fmt.Println("7. Edit data barang")
	fmt.Println("8. Hapus stok barang")
	fmt.Println("")
}

func PrintMenuGudang() {
	fmt.Println("Menu Gudang:")
	fmt.Println("9. Cek stok barang per kategori")
	fmt.Println("10. Tambah stok barang")
	fmt.Println("11. Lihat barang stok rendah (<10)")
	fmt.Println("")
}
