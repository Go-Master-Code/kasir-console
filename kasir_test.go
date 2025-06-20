package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Go-Master-Code/kasir-console/models"

	"github.com/dustin/go-humanize"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB { //Open connection isinya diambil dari web gorm
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

var dbTest = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, dbTest)
}

// kategori barang
func TestAddKategoriBarang(t *testing.T) { //test insert data ke db
	//db.Exec untuk mengeksekusi raw SQL
	//keempat query di bawah akan tetap tereksekusi (4 rows added)
	err := db.Exec("insert into kategori_barang(id,nama_kategori) values(?,?)", "1", "Makanan").Error
	assert.Nil(t, err) //assert memastikan err kosong (tidak ada error)
}

func TestAddKategoriBarangSimple(t *testing.T) { //untuk memasukkan 1 row (data) ke db
	kategoribarang := models.KategoriBarang{ //masukkan data (single) pada struct
		NamaKategori: "Minuman",
	}

	response := db.Create(&kategoribarang)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

func TestTampilkanKategoriBarang(t *testing.T) {
	var kategoriBarang []models.KategoriBarang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	result := db.Find(&kategoriBarang)
	assert.Nil(t, result.Error)

	fmt.Println("=============Kategori Barang===============")
	fmt.Println("ID | Kategori |")
	fmt.Println("============================================")

	for i := range kategoriBarang {
		fmt.Println(strconv.Itoa(kategoriBarang[i].ID) + " | " + kategoriBarang[i].NamaKategori)
	}
	fmt.Println("============================================")
}

//--end of kategori barang

// user
func TestAddUser(t *testing.T) {
	user := models.User{ //masukkan data (single) pada struct
		//ID:         "1", auto_increment
		ID:       "Budi",
		IdLevel:  "2",
		Password: "rahasia",
	}
	response := db.Create(&user)
	assert.Nil(t, response.Error)
}

//end of user

// barang
func TestAddBarang(t *testing.T) {
	barang := models.Barang{ //masukkan data (single) pada struct
		//ID:         "1", auto_increment
		NamaBarang: "Chitato Indomie",
		Harga:      8500,
		//Stok:       0, sudah ada default value 0
		IdKategori: 1,
	}
	response := db.Create(&barang)
	assert.Nil(t, response.Error)
}

func TestTampilkanBarang(t *testing.T) {
	var barang []models.Barang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	result := db.Model(&models.Barang{}).Preload("KategoriBarang").Find(&barang)
	assert.Nil(t, result.Error)

	fmt.Println("=============Data Stok barang===============")
	fmt.Println("ID | Nama Barang | Harga | Stok | Kategori |")
	fmt.Println("============================================")

	for i := range barang {
		fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + strconv.Itoa(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori)
	}
	fmt.Println("============================================")
}

func TestTampilkanBarangSedikit(t *testing.T) {
	var barang []models.Barang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	result := db.Model(&models.Barang{}).Preload("KategoriBarang").Where("stok < ?", "10").Find(&barang)
	assert.Nil(t, result.Error)

	fmt.Println("=============Data Stok barang===============")
	fmt.Println("ID | Nama Barang | Harga | Stok | Kategori |")
	fmt.Println("============================================")

	for i := range barang {
		fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + strconv.Itoa(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori)
	}
	fmt.Println("============================================")
}

func TestUpdateBarang(t *testing.T) {
	barang := models.Barang{}
	result := db.First(&barang, "id = ?", "1") //ambil 1 row dengan ID=1

	assert.Nil(t, result.Error)

	barang.NamaBarang = "Piatos 200 gr"
	barang.Harga = 7500
	barang.Stok = 4
	barang.KategoriBarang.ID = 1

	result = db.Save(&barang) //update data ke database
	assert.Nil(t, result.Error)
}

func TestUpdateStok(t *testing.T) {
	barang := models.Barang{}
	result := db.First(&barang, "id = ?", "1") //ambil 1 row dengan ID=1

	assert.Nil(t, result.Error)

	barang.Stok = 4

	result = db.Save(&barang) //update data ke database
	assert.Nil(t, result.Error)
}

func TestValidatorString(t *testing.T) {
	nama := "Lays rasa tomat @50 g"
	var validate *validator.Validate = validator.New()
	err := validate.Var(nama, "required,ascii")
	if err != nil {
		panic(err)
	}
}

func TestValidatorTanggal(t *testing.T) {
	layout := "2006-01-02" // Layout untuk "YYYY-MM-DD"

	str := "2024-11-05"

	tgl, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	} else {
		fmt.Println(tgl)
	}
}

func TestValidasiIdKategori(t *testing.T) {
	var kategoriBarang []models.KategoriBarang
	result := db.Select("id").Where("id =?", "144").Find(&kategoriBarang) //select field tertentu
	assert.Nil(t, result.Error)

	if len(kategoriBarang) < 1 {
		fmt.Println("Data kategori tidak ada!")
	} else {
		fmt.Println(kategoriBarang[0].ID) //hanya ambil nilai row pertama (hasil query) dalam slice kategoriBarang
	}
}

func TestThousandSeparator(testing *testing.T) {
	num := 1000000

	// Mengonversi integer ke string
	numStr := strconv.FormatInt(int64(num), 10)

	// Menambahkan titik sebagai pemisah ribuan
	var result []string
	for i, v := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(v))
	}

	// Menampilkan hasil
	fmt.Println(strings.Join(result, ""))
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

func TestCurrencyRupiahDustin(t *testing.T) {
	amount := 15000

	fmt.Println("Format angka: " + models.FormatAngka(amount))
	// Memformat angka sebagai mata uang Rupiah
	fmt.Println("Format rupiah: " + models.FormatRupiah(amount))
}

func TestAddDataStruct(t *testing.T) {
	type Mahasiswa struct {
		Nama   string
		NIM    string
		Alamat string
	}

	mhs1 := Mahasiswa{
		Nama:   "Budi",
		NIM:    "12345678",
		Alamat: "Jalan Merdeka 1",
	}

	// Menambahkan data baru pada struct mhs1
	fmt.Println("Data Mahasiswa 1:")
	fmt.Println("Nama:", mhs1.Nama)
	fmt.Println("NIM:", mhs1.NIM)
	fmt.Println("Alamat:", mhs1.Alamat)

	// Misalkan kita ingin menambahkan data pada mahasiswa kedua
	mhs2 := Mahasiswa{
		Nama:   "Siti",
		NIM:    "87654321",
		Alamat: "Jalan Merdeka 2",
	}

	fmt.Println("\nData Mahasiswa 2:")
	fmt.Println("Nama:", mhs2.Nama)
	fmt.Println("NIM:", mhs2.NIM)
	fmt.Println("Alamat:", mhs2.Alamat)

	fmt.Println(mhs1)
	fmt.Println(mhs2)

	//Pengisian data struct secara langsung
	mhs := Mahasiswa{}

	// Mengisi data ke dalam struct
	mhs.Nama = "Andi"
	mhs.NIM = "11223344"
	mhs.Alamat = "Jalan Raya No. 5"

	fmt.Println("Data Mahasiswa:")
	fmt.Println("Nama:", mhs.Nama)
	fmt.Println("NIM:", mhs.NIM)
	fmt.Println("Alamat:", mhs.Alamat)
}

func TestAddDataSliceOfStruct(t *testing.T) {
	type Mahasiswa struct {
		Nama   string
		NIM    string
		Alamat string
	}

	var mhs []Mahasiswa

	mhs = append(mhs, Mahasiswa{"Bambang", "2011928273", "Jl. Kembar"})
	mhs = append(mhs, Mahasiswa{"Anissa", "2011198273", "Jl. Dewata"})

	fmt.Println(mhs)             //versi struct
	for _, result := range mhs { //versi single data
		fmt.Println(result.NIM + " " + result.Nama + " " + result.Alamat)
	}
}

func TestTambahMasterDetilTransaksi(t *testing.T) {
	//buata data product
	transaksi := models.Transaksi{
		UserId: "rini",
	}

	//query: INSERT INTO `products` (`id`,`name`,`price`,`created_at`,`updated_at`) VALUES ('P001','Laptop ASUS',10250000,'2024-12-06 15:15:51.069','2024-12-06 15:15:51.069')
	result := db.Create(&transaksi)
	assert.Nil(t, result.Error)

	//coba add beberapa data map
	type BarangDetilTransaksi struct {
		IdTransaksi string `gorm:"column:id_transaksi"`
		Id          string `gorm:"column:id_barang"`
		Jumlah      int    `gorm:"column:jumlah"`
	}

	var brg []BarangDetilTransaksi
	//tambah barang pada slice BarangDetilTransaksi untuk diinput sekaligus
	brg = append(brg, BarangDetilTransaksi{transaksi.ID, "2", 1})
	brg = append(brg, BarangDetilTransaksi{transaksi.ID, "3", 2})

	/*Skenario goroutine:
	Struct dan slice dibuat di method berbeda
	Di dalamnya ada perintah append juga
	Saat semua data dalam map terkumpul, dilakukan goroutine dengan perulangan for range size map tsb, map dipecah menjadi variable dalam setiap iterasi
	semua data di .Create secara single
	*/

	result = db.Table("detil_transaksi").Create(brg)
	assert.Nil(t, result.Error)

	//--coba dibuat goroutine
	/*
		//buat data baru pada table detil -> user 2 like product P001
		result = db.Table("user_like_product").Create(map[string]interface{}{
			"user_id":    "2",
			"product_id": "P001",
		})*/
	//query: INSERT INTO `user_like_product` (`product_id`,`user_id`) VALUES ('P001','2')
}

func TestTambahDetilTransaksi(t *testing.T) {
	//buat data baru pada table detil -> user 1 like product P001
	result := db.Table("detil_transaksi").Create(map[string]interface{}{
		"id_transaksi": "2",
		"id_barang":    "18",
		"jumlah":       "1",
	})
	//query: INSERT INTO `user_like_product` (`product_id`,`user_id`) VALUES ('P001','1')
	assert.Nil(t, result.Error)
}

// Coba insert data ke record map
// coba add beberapa data map
type BarangDetilTransaksiGlobal struct {
	IdTransaksi string `gorm:"column:id_transaksi"`
	Id          int    `gorm:"column:id_barang"`
	Jumlah      int    `gorm:"column:jumlah"`
}

var barang []BarangDetilTransaksiGlobal

func TambahRecordBarangDetilTransaksi(idTransaksi string, idBarang int, jumlah int) {
	//tambah barang pada slice BarangDetilTransaksi untuk diinput sekaligus
	barang = append(barang, BarangDetilTransaksiGlobal{idTransaksi, idBarang, jumlah})
	fmt.Println(barang)
}

func TestTambahRecordDetilTransaksi(t *testing.T) {
	var idTransaksi = "1"
	var idBarang = 2
	var jumlah = 77

	TambahRecordBarangDetilTransaksi(idTransaksi, idBarang, jumlah)

	idTransaksi = "1"
	idBarang = 5
	jumlah = 33
	TambahRecordBarangDetilTransaksi(idTransaksi, idBarang, jumlah)
}

// percobaan goroutine
type Barang struct {
	Id   int
	Nama string
	Qty  int
}

var barangs []Barang

var wg sync.WaitGroup //wait group for goroutine

func InsertDataMap() {
	for i := 0; i < 500; i++ {
		barangs = append(barangs, Barang{i, "Barang " + strconv.Itoa(i), i * 10})
		//fmt.Println(i)
	}
	fmt.Println()

}

func TestGoroutine(t *testing.T) {

	InsertDataMap() //masukkan 500 data ke map
	InsertDataBarang()
}

func InsertDataBarang() {
	for _, item := range barangs {
		wg.Add(1)
		go InsertGoRoutine(&wg, item.Id, item.Nama, item.Qty)
	}
	wg.Wait()                   //menunggu hingga semua go routine selesai
	time.Sleep(time.Second * 1) //boleh pakai / tidak karena sudah ada wg.Wait tidak masalah
	fmt.Println("All data shown")
}

func InsertGoRoutine(wg *sync.WaitGroup, id int, nama string, qty int) { //hasil adalah var baru dari struct barang
	defer wg.Done()
	fmt.Println(id, nama, qty) //simpan function insert ke db
}

func TestTampilkanStrukBelanja(t *testing.T) { //untuk cetak struk belanja berdasarkan data slice (setelah selesai input data detil_transaksi)
	var idTransaksi string = "1"
	var total int
	var totalItem int
	var subtotal int

	type DetilBarang struct {
		ID    string
		Nama  string
		Harga int
		Qty   int
	}

	var slice []DetilBarang

	slice = append(slice, DetilBarang{"1000", "Panadol merah", 4000, 2})
	slice = append(slice, DetilBarang{"2000", "Paramex otot", 2500, 6})
	slice = append(slice, DetilBarang{"3000", "Redoxon 50 mg", 12500, 3})
	fmt.Println(slice)

	fmt.Println("ID Transaksi: " + idTransaksi)
	fmt.Println("===============Struk belanja=================")
	fmt.Println("ID | Nama Barang | Harga | Qty | Subtotal")
	fmt.Println("=============================================")

	for i := range slice {
		subtotal = slice[i].Harga * slice[i].Qty
		fmt.Println(slice[i].ID + " | " + slice[i].Nama + " | " + models.FormatAngka(slice[i].Harga) + " | " + models.FormatAngka(slice[i].Qty) + " | " + models.FormatAngka(subtotal))
		total += subtotal
		totalItem++
	}
	fmt.Println("=========================================")
	fmt.Println("Total item puchased: " + models.FormatAngka(totalItem))
	fmt.Println("Total purchase: " + models.FormatRupiah(total))
	fmt.Println("=========================================")
}

// Tampilkan data barang full untuk disimpan di struk
func TestTampilDataBarangStruk(t *testing.T) {
	var barang models.Barang
	//preload table detil

	/*
		Query:
		SELECT * FROM `user_like_product` WHERE `user_like_product`.`product_id` = 'P001'
		SELECT * FROM `users` WHERE `users`.`id` IN ('1','2')
		SELECT * FROM `products` WHERE id = 'P001' ORDER BY `products`.`id` LIMIT 1
	*/

	//Joins memuat nama struct "JualBarang" yang ada di model barang!
	result := db.Preload("JualBarang").Select("barang.id, nama_barang").First(&barang, "id = ?", "1")
	assert.Nil(t, result.Error)
	fmt.Println(barang)
}

func TestPreloadManyToManyTransaksi(t *testing.T) { //Untuk Show Transaksi berdasarkan id_transaksi
	var total int
	var totalItem int
	var subtotal int

	var transaksi []models.Transaksi

	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int    `gorm:"column:id_transaksi"`
		IdBarang    int    `gorm:"column:id_barang"`
		Jumlah      int    `gorm:"column:jumlah"`
		NamaBarang  string `gorm:"column:nama_barang"`
		Harga       int    `gorm:"column:harga"`
		//NamaBarang  string
		//Harga       int
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("detil_transaksi.id_transaksi, detil_transaksi.id_barang, detil_transaksi.jumlah, barang.nama_barang, barang.harga").Joins("join barang on detil_transaksi.id_barang=barang.id").Where("id_transaksi = ?", "20").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("==Qty setiap row barang==")
	fmt.Println(detil)
	fmt.Println("ID Transaksi: " + "20") //idTransaksi ambil dari klausa where di atas
	fmt.Println("======================Struk belanja=====================")
	fmt.Printf("%-3s %-25s %-8s %-5s %-14s\n", "ID", "Nama Barang", "Harga", "Qty", "Subtotal")
	fmt.Println("========================================================")
	//pengecekan per row
	for _, row := range detil {
		subtotal = row.Harga * row.Jumlah
		fmt.Printf("%-3s %-25s %-8s %-5s %-14s\n", strconv.Itoa(row.IdBarang), row.NamaBarang, models.FormatAngka(row.Harga), strconv.Itoa(row.Jumlah), models.FormatAngka(subtotal))
		//fmt.Println( /*strconv.Itoa(row.IdTransaksi)+" | "+*/ strconv.Itoa(row.IdBarang) + " | " + row.NamaBarang + " | " + FormatAngka(row.Harga) + " | " + strconv.Itoa(row.Jumlah) + " | " + FormatAngka(subtotal))
		total += subtotal
		totalItem++
	}

	fmt.Println("========================================================")
	fmt.Printf("%-22s %-3s\n", "Total item purchased: ", models.FormatAngka(totalItem))
	//fmt.Println("Total item purchased: " + FormatAngka(totalItem))
	fmt.Printf("%-22s %-3s\n", "Total purchase: ", models.FormatRupiah(total))
	//fmt.Println("Total purchase: " + FormatRupiah(total))
	fmt.Println("========================================================")

	//Joins memuat nama struct "JualBarang" yang ada di model barang!
	result := db.Preload("BarangTerjual").First(&transaksi, "id_transaksi = ?", "20")
	assert.Nil(t, result.Error)
	//pastikan ada 3 row yang muncul (3 barang) pada id transaksi 20
	assert.Equal(t, 3, len(transaksi[0].BarangTerjual))

	fmt.Println(transaksi[0].BarangTerjual)
	fmt.Println("==Print per line barang==")
	for _, tx := range transaksi[0].BarangTerjual { //range diambila terhadap relasi transaksi ke tabel detilnya melalui struct BarangTerjual
		fmt.Println(tx)
		//fmt.Println(i, tx.NamaBarang, FormatAngka(tx.Harga))
	}

	fmt.Println("==Pengecekan dengan for==")
	for i := range len(transaksi[0].BarangTerjual) {
		//fmt.Println(i)
		fmt.Println(transaksi[0].BarangTerjual[i].NamaBarang) //tx menjadi single instance, tidak perlu pakai index [0]
	}

	//pengecekan tanpa for, transaksi masih berbentuk slice, cara printnya beda
	fmt.Println("==Pengecekan manual per line==")
	fmt.Println(transaksi[0].BarangTerjual[0].NamaBarang)
	fmt.Println(transaksi[0].BarangTerjual[1].NamaBarang)
	fmt.Println(transaksi[0].BarangTerjual[2].NamaBarang)
}

func TestCetakStrukRata(t *testing.T) {
	fmt.Println("=====================Contoh Struk=====================")
	fmt.Printf("%-5s %-20s %-10s %-15s %-12s\n", "No", "Nama Item", "Jumlah", "Harga Satuan", "Subtotal")
	fmt.Println("======================================================")
	fmt.Printf("%-5s %-20s %-10s %-15s %-12s\n", "1", "Beng beng mini", "4", "18000", "72000")
}

func TestUpdateStokBarang(t *testing.T) {
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int `gorm:"column:id_transaksi"`
		IdBarang    int `gorm:"column:id_barang"`
		Jumlah      int `gorm:"column:jumlah"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("id_transaksi, id_barang, jumlah").Where("id_transaksi = ?", "25").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(detil)

	//coba update stok per detil transaksi

	barang := models.Barang{}

	result := db.First(&barang, "id = ?", detil[1].IdBarang) //ambil 1 row dengan ID=1

	assert.Nil(t, result.Error)

	var stok int = barang.Stok - detil[1].Jumlah
	barang.Stok = stok

	result = db.Save(&barang) //update data ke database
	assert.Nil(t, result.Error)
}

//update tiap stok row barang berdasarkan detil_transaksi

var brgDetilTransaksi []BarangDetilTransaksi

func TambahListBarangDetilTrans() {
	//tambah barang pada slice BarangDetilTransaksi untuk diinput sekaligus
	brgDetilTransaksi = append(brgDetilTransaksi, BarangDetilTransaksi{"25", 11, 1})
	brgDetilTransaksi = append(brgDetilTransaksi, BarangDetilTransaksi{"25", 4, 5})
	fmt.Println(brg)
}
func TestSaveDetilTransaksi(t *testing.T) { //save semua row barang ke dalam tabel detil_transaksi
	TambahListBarangDetilTrans()
	fmt.Println(brgDetilTransaksi[0].Id)
	barang := models.Barang{}
	fmt.Println("Kurangi stok tiap barang")

	for i := range brgDetilTransaksi {
		barang.ID = strconv.Itoa(brgDetilTransaksi[i].Id)
		_ = db.First(&barang, "id = ?", barang.ID) //ambil 1 row dengan ID tertentu

		barang.Stok = barang.Stok - brgDetilTransaksi[i].Jumlah //update stok berdasarkan qty terjual

		_ = db.Save(&barang) //update data ke database
		fmt.Println("Stok barang " + barang.ID + " telah diupdate menjadi: " + strconv.Itoa(barang.Stok))
	}
}

func TestValidasiLogin(t *testing.T) {
	var user []models.User

	result := db.Model(&models.User{}).Preload("KategoriUser").Where("id = ? and password = ?", "admin", "admin").Find(&user)
	assert.Nil(t, result.Error)

	if len(user) < 1 {
		fmt.Println("Username / password salah!")
	} else {
		fmt.Println(user[0].ID, user[0].Password) //hanya ambil nilai row pertama (hasil query) dalam slice kategoriBarang
	}
}

func TestHashPassword(t *testing.T) {
	//password yang akan di hash
	password := []byte("rahasia")
	//menggunakan argon2 untuk menghasilkan hash
	salt := []byte("random_salt") //salting sangat penting untuk keamanan
	hash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	//mencetak hash yang dihasilkan
	fmt.Printf("Hash: %x/n", hash)
	fmt.Println(len(hash))

	//verifikasi jika hash yang dihasilkan cocok dengan password yang seharusnya
	newHash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	if string(hash) == string(newHash) { //konversi byte ke string
		fmt.Println("Password cocok")
	} else {
		fmt.Println("Password salah")
	}
}
