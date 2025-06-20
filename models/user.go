package models

type User struct {
	ID           string       `gorm:"primary_key;column:id"`
	IdLevel      string       `gorm:"column:id_level"`
	Password     string       `gorm:"column:password"`
	KategoriUser KategoriUser `gorm:"foreignKey:id_level;references:id"`
	//Transaksi    []Transaksi  `gorm:"foreignKey:id;references:id"`
	//relasi one to many
	//1 user menangani banyak transaksi
}

func (u *User) TableName() string {
	return "user" //nama table pada db nya adalah user_logs
}
