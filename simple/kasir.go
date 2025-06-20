package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

// global var
var reader = bufio.NewScanner(os.Stdin)
var barang string
var harga int
var jumlah int
var total int
var err error

/*
type Transaksi struct {
	barang string
	harga  int
	jumlah int
}

var trans Transaksi
*/

// ini harusnya method main
func Trenggalek() {
	fmt.Println("Selamat datang di toko palugada!")
	Start()

	var lanjut bool

	lanjut = KonfirmasiLanjut() //panggil function KonfirmasiLanjut()

	//for di bawah sebenarnya berfungsi sebagai while
	for lanjut { //selama lanjut = true (dijawab y saat konfirmasi lanjut)
		Start()
		lanjut = KonfirmasiLanjut()
	}

	time.Sleep(1 * time.Second)

}

func Start() {
	InputBarang()
	InputHargaBarang()
	InputJumlahBarang()
	HitungTotal(harga, jumlah)
	/*
		trans = Transaksi{
			barang: barang,
			harga:  harga,
			jumlah: jumlah,
		}
		fmt.Println(trans)
	*/
}

// func untuk input barang
func InputBarang() {
	fmt.Print("Masukkan nama barang: ")
	if reader.Scan() {
		input := reader.Text()
		//masukkan hasil input ke var
		if input == "" {
			fmt.Println("Nama barang harus diisi!")
			InputBarang() //recursive
		} else {
			barang = input
		}
	}
}

func InputHargaBarang() {
	fmt.Print("Masukkan harga barang: ")
	if reader.Scan() {
		inputHarga := reader.Text()
		//validasi jika string kosong
		if inputHarga == "" {
			fmt.Println("Harga harus diisi!")
			InputHargaBarang() //recursive
		} else {
			//coba konversi ke angka
			harga, err = strconv.Atoi(inputHarga)
			if err != nil {
				fmt.Println("Harga harus berupa angka!")
				InputHargaBarang() //recursive
			}
		}
	}
}

func InputJumlahBarang() {
	fmt.Print("Masukkan jumlah barang: ")
	if reader.Scan() {
		inputJumlah := reader.Text()
		//validasi jika string kosong
		if inputJumlah == "" {
			fmt.Println("Jumlah harus diisi!")
			InputJumlahBarang() //recursive
		} else {
			//coba konversi ke angka
			jumlah, err = strconv.Atoi(inputJumlah)
			if err != nil {
				fmt.Println("Jumlah harus berupa angka!")
				InputJumlahBarang() //recursive
			}
		}
	}
}

func PesanAkhir() {
	fmt.Println("Terima kasih telah berbelanja di toko kami!")
}

// func experimental
func KonfirmasiLanjut() bool {
	var lanjut bool
	reader := bufio.NewScanner(os.Stdin)

	fmt.Print("Apakah ada barang lain (y/n): ")
	if reader.Scan() {
		input := reader.Text()
		//masukkan hasil input ke var
		if input == "y" {
			lanjut = true
		} else {
			defer PesanAkhir()
			lanjut = false
		}
	}
	return lanjut
}

func HitungTotal(harga int, jumlah int) {
	total = harga * jumlah //total global var

	fmt.Println("==========Struk==========")
	fmt.Println("Nama barang: " + barang)
	fmt.Println("Harga barang: Rp" + strconv.Itoa(harga) + ",00")
	fmt.Println("Jumlah barang: " + strconv.Itoa(jumlah))
	fmt.Println("Total belanja: Rp" + strconv.Itoa(total) + ",00")
	fmt.Println("=========================")
}

/* experimental code
package main

import (
        "fmt"
)

func main() {
    fmt.Println("Enter your Age: ")

    var userAge int
    for true {
        _, err := fmt.Scanf("%d", &userAge)
        if err == nil {
            break
        }
        fmt.Println("Not a valid age - try again")
        var dump string
        fmt.Scanln(&dump)
    }

    fmt.Println(userAge)
}
*/
