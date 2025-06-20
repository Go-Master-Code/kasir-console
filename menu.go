package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Go-Master-Code/kasir-console/models"
	"github.com/Go-Master-Code/kasir-console/views"

	"github.com/go-playground/validator/v10"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// global var
var reader = bufio.NewScanner(os.Stdin) //deklarasi bufio scanner
var db = OpenConnectionMaster()
var harga int
var stok int
var idKategori int
var idBarang int
var namaBarang string
var qty int
var idTransaksi string
var lanjut bool
var userName string
var password string
var levelUser int
var periode1 string
var periode1Full string
var periode2 string
var periode2Full string
var layout = "2006-01-02" // Layout untuk "YYYY-MM-DD"

// data date n time untuk laporan barang dan transaksi

// struct data untuk detil_produk
// coba add beberapa data map
type BarangDetilTransaksi struct {
	IdTransaksi string `gorm:"column:id_transaksi"`
	Id          int    `gorm:"column:id_barang"`
	Jumlah      int    `gorm:"column:jumlah"`
}

var brg []BarangDetilTransaksi

//--end of struct

// sementara user di set dulu rini
var userID string

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {
	validate = validator.New()

	TampilkanHalamanLogin()

	//TambahKategoriBarang("Makanan")

	fmt.Println("====================================================================")
	fmt.Println("============== SISTEM INFORMASI MINIMARKET v.1.0.0 =================")
	// fmt.Println("==    = = = =     =      =   =      =       ==")
	// fmt.Println("==       =       = =      = =      = =      ==")
	// fmt.Println("==       =      =   =      =      =   =     ==")
	// fmt.Println("==   =   =     = = = =     =     = = = =    ==")
	// fmt.Println("==   = = =    =       =    =    =       =   ==")
	if levelUser == 1 {
		views.PrintHeaderAdmin()
	} else {
		views.PrintHeaderUserBiasa()
	}

	//experiment
	MainMenu()

	//TambahBarang("2", "Floridina 450 ml", 4500, "2")
	//TampilkanBarang()

	/*validator
	myEmail := "joeybloggs.gmail.com"

	errs := validate.Var(myEmail, "required,email")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}
	*/
}

func MainMenu() {
	fmt.Print("Ketik nomor menu: ")
	if reader.Scan() {
		pil := reader.Text()
		err := validate.Var(pil, "required")

		if err != nil {
			fmt.Println("Nomor menu tidak valid!")
			MainMenu()
		} else {
			switch pil {
			case "0":
				break
			case "9":
				InputKategoriBarang()
				models.TampilkanBarangPerKategori(db, idKategori)
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "4":
				InputPeriodeLapTransaksi1()
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "6": //case bisa pakai tipe data string
				models.TampilkanBarang(db)
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "5":
				TampilkanInputBarang()
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "7":
				TampilkanUpdateBarang()
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "10":
				TampilkanUpdateStokBarang()
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "8":
				HapusStokBarang()
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "1":
				TampilkanInputTransaksi()
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "11":
				models.TampilkanBarangSedikit(db)
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "2":
				models.CetakLaporanStokBarang(db)
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			case "3":
				//CetakLaporanTransaksiSummary() menampilkan semua data transaksi
				//mencetak data transaksi hari ini
				models.CetakLaporanTransaksiSummaryToday(db)
				if levelUser == 1 {
					views.PrintHeaderAdmin()
				} else {
					views.PrintHeaderUserBiasa()
				}
				MainMenu()
			default:
				fmt.Println("Ketik nomor menu yang valid!")
				MainMenu()
			}
		}

	}
}

func OpenConnectionMaster() *gorm.DB { //Open connection isinya diambil dari web gorm
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dialect := "root:root@tcp(127.0.0.1:3306)/kasir?charset=utf8mb4&parseTime=True&loc=Local" //root:root artinya username=root, password=root
	db, err := gorm.Open(mysql.Open(dialect), &gorm.Config{
		//comment Logger di bawah agar perintah sql nya tidak muncul
		//Logger: logger.Default.LogMode(logger.Info), //ubah logger pada gorm: semua perintah SQL akan di log sebagai info

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

func InputPeriodeLapTransaksi1() {
	fmt.Print("Ketik tanggal mulai ('yyyy-mm-dd'): ")
	if reader.Scan() {
		periode1 = reader.Text()
		err := validate.Var(periode1, "required") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("Tanggal mulai tidak valid!")
			InputPeriodeLapTransaksi1()
		} else {
			_, err := time.Parse(layout, periode1)
			if err != nil {
				fmt.Println("Tanggal mulai tidak valid!")
				InputPeriodeLapTransaksi1()
			} else {
				//tambahkan jam spesifik
				periode1Full = periode1 + " 00:00:00"
				InputPeriodeLapTransaksi2()
			}
		}
	}
}

func InputPeriodeLapTransaksi2() {
	fmt.Print("Ketik tanggal akhir ('yyyy-mm-dd'): ")
	if reader.Scan() {
		periode2 = reader.Text()
		err := validate.Var(periode2, "required") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("Tanggal akhir tidak valid!")
			InputPeriodeLapTransaksi2()
		} else {
			_, err := time.Parse(layout, periode2)
			if err != nil {
				fmt.Println("Tanggal akhir tidak valid!")
				InputPeriodeLapTransaksi2()
			} else {
				//tambahkan jam spesifik
				periode2Full = periode2 + " 23:59:59"
				//fmt.Println("Tanggal mulai: " + periode1)
				//fmt.Println("Tanggal mulai: " + periode2)
				models.CetakLapTransaksiSummaryPeriode(db, periode1Full, periode2Full, periode1, periode2)
			}
		}
	}
}

func InputBarang() {
	fmt.Print("Ketik nama barang: ")
	if reader.Scan() {
		namaBarang = reader.Text()
		err := validate.Var(namaBarang, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("Nama barang tidak valid!")
			TampilkanInputBarang()
		}
	}
}

func TampilkanInputTransaksi() {
	fmt.Print("Ketik ID barang: ")
	if reader.Scan() {
		idBarangString := reader.Text()
		err := validate.Var(idBarangString, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("ID barang tidak valid!")
			TampilkanInputTransaksi()
		} else {
			barang := models.Barang{}
			_ = db.First(&barang, "id = ?", idBarangString) //ambil 1 row dengan ID tertentu
			if len(barang.ID) == 0 {                        //cek apakah row barang ada
				fmt.Println("ID barang tidak valid!")
				TampilkanInputTransaksi()
			} else {
				idBarang, _ = strconv.Atoi(barang.ID)
				InputDanCekStokBarang()
				//anonumous func tambah master transaksi
				/*
					//anonumouse func tambah master transaksi
					go func () {
						transaksi := models.Transaksi{
							UserId: userID,
						}

						//query: INSERT INTO `products` (`id`,`name`,`price`,`created_at`,`updated_at`) VALUES ('P001','Laptop ASUS',10250000,'2024-12-06 15:15:51.069','2024-12-06 15:15:51.069')
						err := db.Create(&transaksi).Error
						if err != nil {
							panic(err)
						}
					}()
				*/
				transaksi := models.Transaksi{
					UserId: userID,
				}

				//query: INSERT INTO `products` (`id`,`name`,`price`,`created_at`,`updated_at`) VALUES ('P001','Laptop ASUS',10250000,'2024-12-06 15:15:51.069','2024-12-06 15:15:51.069')
				err := db.Create(&transaksi).Error
				if err != nil {
					panic(err)
				}
				idTransaksi = transaksi.ID //simpan ID transaksi untuk add more detil_jual

				models.TambahDetilTransaksi(transaksi.ID, idBarang, qty)
				//SaveMasterDetilTransaksi()
				fmt.Println("Transaksi berhasil diinput!")

				//konfirmasi apakah mau tambah barang atau tidak

				lanjut = KonfirmasiLanjut() //panggil function KonfirmasiLanjut()

				//for di bawah sebenarnya berfungsi sebagai while
				for lanjut { //selama lanjut = true (dijawab y saat konfirmasi lanjut)
					InputMoreDetilTransaksi()
					//lanjut = KonfirmasiLanjut()
				}
			}
		}
	}
}

func InputMoreDetilTransaksi() {
	fmt.Print("Ketik ID barang: ")
	if reader.Scan() {
		idBarangString := reader.Text()
		err := validate.Var(idBarangString, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("ID barang tidak valid!")
			TampilkanInputTransaksi()
		} else {
			barang := models.Barang{}
			_ = db.First(&barang, "id = ?", idBarangString) //ambil 1 row dengan ID tertentu
			if len(barang.ID) == 0 {                        //cek apakah row barang ada
				fmt.Println("ID barang tidak valid!")
				TampilkanInputTransaksi()
			} else {
				idBarang, _ = strconv.Atoi(barang.ID)
				InputDanCekStokBarang()

				models.TambahDetilTransaksi(idTransaksi, idBarang, qty)
				//SaveMasterDetilTransaksi()
				fmt.Println("Transaksi berhasil diinput!")
				//InputMoreDetilTransaksi()

				lanjut = KonfirmasiLanjut() //panggil function KonfirmasiLanjut()

				//for di bawah sebenarnya berfungsi sebagai while
				for lanjut { //selama lanjut = true (dijawab y saat konfirmasi lanjut)
					InputMoreDetilTransaksi()
					//lanjut = KonfirmasiLanjut()
				}
			}
		}
	}
}

func KonfirmasiLanjut() bool {
	reader := bufio.NewScanner(os.Stdin)

	fmt.Print("Apakah ada barang lain (y/n): ")
	if reader.Scan() {
		input := reader.Text()
		//masukkan hasil input ke var
		if input == "y" {
			lanjut = true
		} else if input == "n" {
			//fmt.Println("Sudah tidak ada lagi barang belanjaan")
			models.SaveDetilTransaksi(db) //hanya insert detil transaksi, master sudah diinput sebelumnya
			lanjut = false

			models.UpdateStokBarangDetilTransaksi(db)
			models.CetakStruk(db, idTransaksi, userName) //full atribut detil produk

			//kosongkan kembali slice data detil barang
			brg = []BarangDetilTransaksi{}
		} else {
			fmt.Println("Input salah!")
			KonfirmasiLanjut()
		}
	}
	return lanjut
}

func InputDanCekStokBarang() {
	fmt.Print("Ketik jumlah barang: ")
	if reader.Scan() {
		qtyString, err := strconv.Atoi(reader.Text())

		if err != nil {
			fmt.Println("Stok harus berupa angka!")
			InputDanCekStokBarang() //recursive
		} else {
			err := validate.Var(qtyString, "required")
			if err != nil {
				fmt.Println("Stok barang tidak valid!")
				InputDanCekStokBarang()
			} else {
				//cek kesediaan stok barang
				var barang = models.Barang{}

				result := db.First(&barang, "id = ?", idBarang).Error
				//result := db.Model(&models.Barang{}).Where("id = ?", idBarang).Find(&barang).Error
				if result != nil {
					panic(result)
				}

				//*Cek stok yang tersedia dan qty yang diminta
				//fmt.Println(qtyString)
				//fmt.Println(barang.Stok)
				if barang.Stok < qtyString {
					fmt.Println("Stok barang tidak cukup!")
					InputDanCekStokBarang()
				} else {
					qty = qtyString
				}
			}
		}
	}
}

func TampilkanInputBarang() {
	fmt.Print("Ketik nama barang: ")
	if reader.Scan() {
		namaBarang = reader.Text()
		err := validate.Var(namaBarang, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("Nama barang tidak valid!")
			TampilkanInputBarang()
		} else {
			InputHargaBarang()
			InputKategoriBarang()
			models.TambahBarang(db, namaBarang, harga, idKategori)
			fmt.Println("Barang berhasil diinput!")
		}
	}
}

func TampilkanHalamanLogin() {
	fmt.Println("============Login============")
	fmt.Print("Ketik username: ")
	if reader.Scan() {
		userName = reader.Text()
		err := validate.Var(userName, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("Username tidak valid!")
			TampilkanHalamanLogin()
		} else {
			userID = userName
			TampilkanInputPassword()
		}
	}
}

func TampilkanInputPassword() {
	fmt.Print("Ketik password: ")
	if reader.Scan() {
		password = reader.Text()
		err := validate.Var(password, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("Password tidak valid!")
			TampilkanInputPassword()
		} else {
			ValidasiLogin(userID, password)
		}
	}
}

func ValidasiLogin(userName string, password string) {
	var user []models.User

	_ = db.Model(&models.User{}).Preload("KategoriUser").Where("id = ? and password = ?", userName, password).Find(&user)

	if len(user) < 1 {
		fmt.Println("Username / password salah!")
		TampilkanHalamanLogin()
	} else {
		userID = user[0].ID
		levelUser, _ = strconv.Atoi(user[0].IdLevel)
		fmt.Println("=============================")
		fmt.Println("Selamat datang " + userID + "!")
	}
}

func HapusStokBarang() {
	fmt.Print("Ketik ID barang: ")
	if reader.Scan() {
		idBarangString := reader.Text()
		err := validate.Var(idBarangString, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("ID barang tidak valid!")
			HapusStokBarang()
		} else {
			barang := models.Barang{}
			_ = db.First(&barang, "id = ?", idBarangString) //ambil 1 row dengan ID tertentu
			if len(barang.ID) == 0 {                        //cek apakah row barang ada
				fmt.Println("ID barang tidak valid!")
				HapusStokBarang()
			} else {
				barang.Stok = 0      //jadikan stok barang dengan id di atas = 0
				_ = db.Save(&barang) //update data ke database
				fmt.Println("Stok barang " + barang.NamaBarang + " telah dihapus!")
			}
		}
	}
}

func TampilkanUpdateBarang() {
	fmt.Print("Ketik ID barang: ")
	if reader.Scan() {
		idBarangString := reader.Text()
		err := validate.Var(idBarangString, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("ID barang tidak valid!")
			TampilkanUpdateBarang()
		} else {
			barang := models.Barang{}
			_ = db.First(&barang, "id = ?", idBarangString) //ambil 1 row dengan ID=1

			if len(barang.ID) == 0 {
				fmt.Println("ID barang tidak valid!")
				TampilkanUpdateBarang()
			} else {
				idBarang, _ = strconv.Atoi(idBarangString) //kirim parameter idBarang ke var global
				InputBarang()
				InputHargaBarang()
				InputKategoriBarang()
				InputStokBarang()
				models.UpdateBarang(db, idBarang, namaBarang, harga, stok, idKategori)
				fmt.Println("Barang berhasil diupdate!")
			}
		}
	}
}

func TampilkanUpdateStokBarang() {
	fmt.Print("Ketik ID barang: ")
	if reader.Scan() {
		idBarangString := reader.Text()
		err := validate.Var(idBarangString, "required,ascii") //harus diisi dan boleh alpha numeric

		if err != nil {
			fmt.Println("ID barang tidak valid!")
			TampilkanUpdateStokBarang()
		} else {
			barang := models.Barang{}
			//_ = db.First(&barang, "id = ?", idBarangString) //ambil 1 row dengan ID=1

			result := db.Model(&models.Barang{}).Preload("KategoriBarang").Where("id = ?", idBarangString).Find(&barang).Error
			if result != nil {
				panic(result)
			}

			if len(barang.ID) == 0 {
				fmt.Println("ID barang tidak valid!")
				TampilkanUpdateStokBarang()
			} else {
				// fmt.Println("Nama barang: " + barang.NamaBarang)
				fmt.Println("============================Data Barang=============================")
				fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", "ID", "Nama Barang", "Harga", "Stok", "Kategori")
				//fmt.Println("ID | Nama Barang | Harga | Stok | Kategori |")
				fmt.Println("====================================================================")
				fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", barang.ID, barang.NamaBarang, models.FormatAngka(barang.Harga), models.FormatAngka(barang.Stok), barang.KategoriBarang.NamaKategori)
				fmt.Println("====================================================================")

				idBarang, _ = strconv.Atoi(idBarangString) //kirim parameter idBarang ke var global
				InputStokBarang()

				barang.Stok += stok //update stok barang dijumlahkan dengan input user

				err := db.Save(&barang).Error
				if err != nil {
					panic(err)
				}

				fmt.Println("Stok barang " + barang.NamaBarang + " berhasil diupdate!")
			}
		}
	}
}

func InputHargaBarang() {
	fmt.Print("Ketik harga barang: ")
	if reader.Scan() {
		hargaString, err := strconv.Atoi(reader.Text())

		if err != nil {
			fmt.Println("Harga harus berupa angka!")
			InputHargaBarang() //recursive
		} else {
			err := validate.Var(hargaString, "required")
			if err != nil {
				fmt.Println("Harga barang tidak valid!")
				InputHargaBarang()
			} else {
				harga = hargaString
			}
		}
	}
}

func InputStokBarang() {
	fmt.Print("Ketik stok barang: ")
	if reader.Scan() {
		stokString, err := strconv.Atoi(reader.Text())

		if err != nil {
			fmt.Println("Stok harus berupa angka!")
			InputStokBarang() //recursive
		} else {
			err := validate.Var(stokString, "required")
			if err != nil {
				fmt.Println("Stok barang tidak valid!")
				InputStokBarang()
			} else {
				stok = stokString
			}
		}
	}
}

func InputKategoriBarang() {
	models.TampilkanKategoriBarang(db)
	fmt.Print("Ketik ID kategori: ")
	if reader.Scan() {
		kategoriString, err := strconv.Atoi(reader.Text())

		if err != nil {
			fmt.Println("ID kategori harus berupa angka!")
			InputKategoriBarang() //recursive
		} else {
			err := validate.Var(kategoriString, "required")
			if err != nil {
				fmt.Println("ID kategori tidak valid!")
				InputKategoriBarang()
			} else {
				//validasi idKategori barang ada atau tidak
				//var kategoriBarang []models.KategoriBarang

				/*
					var kategoriBarang []models.KategoriBarang
					errs := db.Find(&kategoriBarang).Error //select

					if errs != nil {
						panic(errs)
					}

					var ketemu bool = false
					//iterasi untuk setiap record kategoriBarang, cari ID yang sesuai
					for i := range kategoriBarang {
						if kategoriString == kategoriBarang[i].ID {
							idKategori = kategoriString //masukkan ID kategori ke var global
							//fmt.Println("Ketemu! " + "ID: " + strconv.Itoa(kategoriBarang[i].ID))
							ketemu = true
						}
					}

					if !ketemu {
						InputKategoriBarang()
					}*/

				var kategoriBarang []models.KategoriBarang
				_ = db.Select("id").Where("id =?", kategoriString).Find(&kategoriBarang).Error //select field tertentu
				//var di atas (harusnya err) diignore karena bisa saja datanya tidak ada (0 rows) -> ini tidak menghasilkan error
				if len(kategoriBarang) < 1 {
					fmt.Println("ID kategori tidak ada!")
					InputKategoriBarang()
				} else {
					idKategori = kategoriBarang[0].ID //hanya ambil nilai row pertama (hasil query) dalam slice kategoriBarang
				}
			}
		}
	}
}
