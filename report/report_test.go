package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Go-Master-Code/kasir-console/models"

	"github.com/dustin/go-humanize"
	"github.com/jung-kurt/gofpdf"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnectionReport() *gorm.DB { //Open connection isinya diambil dari web gorm
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dialect := "root:root@tcp(127.0.0.1:3306)/kasir?charset=utf8mb4&parseTime=True&loc=Local" //root:root artinya username=root, password=root
	db, err := gorm.Open(mysql.Open(dialect), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //ubah logger pada gorm: semua perintah SQL akan di log sebagai info

		//SkipDefaultTransaction: true,
		//PrepareStmt: true,
	})

	if err != nil {
		panic(err)
	}

	//connection pool
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute) //maksimal digunakan 30 mins
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  //max waktu connection nganggur

	return db
}

var db = OpenConnectionReport()

var dbTest = OpenConnectionReport()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, dbTest)
}

func TestReport(t *testing.T) {
	var barang []models.Barang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	err := db.Model(&models.Barang{}).Preload("KategoriBarang").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	type countryType struct {
		nameStr, capitalStr, areaStr, popStr string
	}

	var countries []countryType
	countries = append(countries, countryType{"Indonesia", "Jakarta", "Asia", "250000000"})
	countries = append(countries, countryType{"Vietnam", "Hanoi", "Asia", "148756000"})

	for _, c := range countries {
		fmt.Println(c.nameStr, c.capitalStr, c.areaStr, c.popStr)
	}

	header := []string{"ID", "Barang", "Harga", "Stok"}

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Stok Minimarket", "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	pdf.SetFont("Arial", "B", 12)
	left := (210.0 - 4*40) / 2
	pdf.SetX(left)

	for _, str := range header {
		pdf.CellFormat(40, 10, str, "1", 0, "C,M", false, 0, "")
	}

	//fmt.Println(countryList[0].nameStr)

	pdf.Ln(-1)

	//tampilkan list barang
	for _, b := range barang {
		pdf.SetX(left)
		pdf.CellFormat(40, 8, b.ID, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(40, 8, b.NamaBarang, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(40, 8, strconv.Itoa(b.Harga), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(40, 8, strconv.Itoa(b.Stok), "1", 0, "R,M", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(-1)
	//-->end of list barang

	for _, c := range countries {
		pdf.SetX(left)
		pdf.CellFormat(40, 8, c.nameStr, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(40, 8, c.capitalStr, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(40, 8, c.areaStr, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(40, 8, c.popStr, "1", 0, "R,M", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 12)
	//pdf.Cell(40, 8, "End of the report")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("hello.pdf")
	if errs != nil {
		panic(errs)
	}
}

func TestReportBarang(t *testing.T) {
	// 	Header
	header := []string{"ID", "Nama Barang", "Kategori", "Harga", "Stok"}

	// Column widths
	w := []float64{12.0, 88.0, 30.0, 22.0, 18.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Stok Minimarket", "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)
	var barang []models.Barang

	//select * from barang + barang.nama_kategori
	err := db.Model(&models.Barang{}).Preload("KategoriBarang").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, b := range barang {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, b.ID, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, b.NamaBarang, "1", 0, "L,M", false, 0, "")
		pdf.CellFormat(w[2], 8, b.KategoriBarang.NamaKategori, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngkaTest(b.Harga), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[4], 8, FormatAngkaTest(b.Stok), "1", 0, "R,M", false, 0, "")

		pdf.Ln(-1)
	}
	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("barang.pdf")
	if errs != nil {
		panic(errs)
	}
}

func FormatRupiahTest(amount int) string {
	// Menggunakan humanize.Comma untuk memformat angka dengan koma sebagai pemisah ribuan
	formatted := humanize.Comma(int64(amount))
	// Ganti koma dengan titik untuk format Rupiah
	return "Rp " + strings.ReplaceAll(formatted, ",", ".")
}

func FormatAngkaTest(amount int) string {
	// Menggunakan humanize.Comma untuk memformat angka dengan koma sebagai pemisah ribuan
	formatted := humanize.Comma(int64(amount))
	// Ganti koma dengan titik untuk format Rupiah
	return strings.ReplaceAll(formatted, ",", ".")
}

func TestLaporanTransaksiAll(t *testing.T) {
	var subtotal, total, totalItem int

	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"` //gorm tag untuk autocreatetime
		IdBarang    int       `gorm:"column:id_barang"`
		Jumlah      int       `gorm:"column:jumlah"`
		NamaBarang  string    `gorm:"column:nama_barang"`
		Harga       int       `gorm:"column:harga"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("created_at, detil_transaksi.id_transaksi, detil_transaksi.id_barang, detil_transaksi.jumlah, nama_barang, harga").Joins("join barang on detil_transaksi.id_barang=barang.id"). /*.Where("id_transaksi = ?", "20")*/ Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Order("detil_transaksi.id_transaksi asc").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("=======================Laporan Data Transaksi=======================")

	//pengecekan per row
	for _, row := range detil {
		subtotal = row.Harga * row.Jumlah
		fmt.Printf("%-4s %-11s %-30s %-8s %-5s %-14s\n", strconv.Itoa(row.IdTransaksi), row.CreatedAt.Format("2006-01-02"), row.NamaBarang, FormatAngkaTest(row.Harga), strconv.Itoa(row.Jumlah), FormatAngkaTest(subtotal))
		total += subtotal
		totalItem++
	}
	fmt.Println("=======================Rekap Laporan Transaksi======================")
	fmt.Println("Total transaksi: " + strconv.Itoa(totalItem))
	fmt.Println("Total nilai transaksi: Rp" + FormatAngkaTest(total))
	fmt.Println("====================================================================")
}

func TestLaporanTransaksiSummary(t *testing.T) {
	var total, totalItem int

	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("======================Laporan Transaksi Summary=====================")

	fmt.Printf("%-4s %-11s %-5s %-14s\n", "ID", "Tanggal", "#Item(s)", "Total")
	//pengecekan per row
	for _, row := range detil {
		fmt.Printf("%-4s %-11s %-8s %-14s\n", strconv.Itoa(row.IdTransaksi), row.CreatedAt.Format("2006-01-02"), strconv.Itoa(row.Item), FormatAngkaTest(row.Subtotal))
		total += row.Subtotal
		totalItem++
	}
	fmt.Println("======================Summary Laporan Transaksi=====================")
	fmt.Println("Total transaksi: " + strconv.Itoa(totalItem))
	fmt.Println("Total nilai transaksi: Rp" + FormatAngkaTest(total))
	fmt.Println("====================================================================")
}

func TestReportTransaksiSummary(t *testing.T) {
	var total, totalItem int

	// 	Header
	header := []string{"ID", "Tanggal", "#Item(s)", "Total"}

	// Column widths
	w := []float64{12.0, 28.0, 20.0, 27.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Penjualan Minimarket", "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)

	//deklarasi struct dan query
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, row := range detil {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, strconv.Itoa(row.IdTransaksi), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, row.CreatedAt.Format("2006-01-02"), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[2], 8, FormatAngkaTest(row.Item), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngkaTest(row.Subtotal), "1", 0, "R,M", false, 0, "")

		//hitung total item dan subtotal secara incremental
		total += row.Subtotal
		totalItem++

		pdf.Ln(-1)
	}
	//Summary

	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//harus selalu di set agar bisa nyambung dengan row sebelumnya
	pdf.SetX(left)
	pdf.CellFormat(60, 8, "Total Nilai Penjualan: ", "1", 0, "L,M", false, 0, "")
	pdf.CellFormat(27, 8, FormatAngkaTest(total), "1", 0, "R,M", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("laporan_transaksi.pdf")
	if errs != nil {
		panic(errs)
	}
}

func TestReportTransaksiSummaryToday(t *testing.T) {
	var total, totalItem int
	today := time.Now().Format("2006-01-02")
	todayString := time.Now().Format("2 Jan 2006")

	// 	Header
	header := []string{"ID", "Tanggal", "#Item(s)", "Total"}

	// Column widths
	w := []float64{12.0, 28.0, 20.0, 27.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Transaksi Tanggal "+todayString, "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)

	//deklarasi struct dan query
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Where("created_at like ?", today+"%").Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, row := range detil {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, strconv.Itoa(row.IdTransaksi), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, row.CreatedAt.Format("2006-01-02"), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[2], 8, FormatAngkaTest(row.Item), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngkaTest(row.Subtotal), "1", 0, "R,M", false, 0, "")

		//hitung total item dan subtotal secara incremental
		total += row.Subtotal
		totalItem++

		pdf.Ln(-1)
	}
	//Summary

	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//harus selalu di set agar bisa nyambung dengan row sebelumnya
	pdf.SetX(left)
	pdf.CellFormat(60, 8, "Total Nilai Penjualan: ", "1", 0, "L,M", false, 0, "")
	pdf.CellFormat(27, 8, FormatAngkaTest(total), "1", 0, "R,M", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("laporan_transaksi.pdf")
	if errs != nil {
		panic(errs)
	}
}

func TestReportTransaksiSummaryPeriode(t *testing.T) {
	var total, totalItem int

	// 	Header
	header := []string{"ID", "Tanggal", "#Item(s)", "Total"}

	// Column widths
	w := []float64{12.0, 28.0, 20.0, 27.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Transaksi Dari "+"tanggal sekarang"+"Sampai "+"tanggal sampai", "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)

	//deklarasi struct dan query
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Where("created_at between ? and ?", "2024-12-10", "2024-12-16").Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, row := range detil {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, strconv.Itoa(row.IdTransaksi), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, row.CreatedAt.Format("2006-01-02"), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[2], 8, FormatAngkaTest(row.Item), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngkaTest(row.Subtotal), "1", 0, "R,M", false, 0, "")

		//hitung total item dan subtotal secara incremental
		total += row.Subtotal
		totalItem++

		pdf.Ln(-1)
	}
	//Summary

	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//harus selalu di set agar bisa nyambung dengan row sebelumnya
	pdf.SetX(left)
	pdf.CellFormat(60, 8, "Total Nilai Penjualan: ", "1", 0, "L,M", false, 0, "")
	pdf.CellFormat(27, 8, FormatAngkaTest(total), "1", 0, "R,M", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("laporan_transaksi_per_periode.pdf")
	if errs != nil {
		panic(errs)
	}
}
