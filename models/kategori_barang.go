package models

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type KategoriBarang struct {
	ID           int      `gorm:"primary_key;column:id;autoIncrement"`
	NamaKategori string   `gorm:"primary_key;column:nama_kategori"`
	Barang       []Barang `gorm:"foreignKey:id;references:id"`
	//relasi many to one terhadap barang
	//ID kategori barang yang sama dapat dimiliki beberapa barang
}

func (k *KategoriBarang) TableName() string {
	return "kategori_barang" //nama table pada db
}

// kategori barang
func TambahKategoriBarang(db *gorm.DB, kategori string) { //untuk memasukkan 1 row (data) ke db
	kategoribarang := KategoriBarang{ //masukkan data (single) pada struct
		NamaKategori: kategori, //ID tidak didefinisikan karena autoincrement
	}

	err := db.Create(&kategoribarang).Error
	if err != nil {
		panic(err)
	}
}

func TampilkanKategoriBarang(db *gorm.DB) {
	var kategoriBarang []KategoriBarang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	err := db.Find(&kategoriBarang).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("=============Kategori Barang=================")
	fmt.Println("ID | Kategori |")
	fmt.Println("=============================================")

	for i := range kategoriBarang {
		fmt.Println(strconv.Itoa(kategoriBarang[i].ID) + " | " + kategoriBarang[i].NamaKategori)
	}
	fmt.Println("=============================================")
}
