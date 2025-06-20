create table kategori_barang (
    id int(2) not null auto_increment,
    nama_kategori varchar(100) not NULL,
    primary key(id)
)engine=InnoDB;

create table barang (
    id int(6) not null auto_increment,
    nama_barang varchar(100) not null,
    harga int not null,
    stok int not null,
    id_kategori int(2) not null,
    primary key(id),
    foreign key (id_kategori) references kategori_barang(id)
)engine=InnoDB;

create table kategori_user (
    id char(1) not null,
    level_user varchar(10) not null,
    primary key (id)
)engine=InnoDB;

create table user (
    id varchar(50) not null,
    password varchar(50) not null,
    id_level char(1) not null,
    primary key(id),
    foreign key (id_level) references kategori_user(id)
)engine=InnoDB;

create table transaksi (
    id_transaksi int(11) not null auto_increment,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    user_id varchar(50) not null,
    primary key(id_transaksi),
	foreign key(user_id) references user(id)
)engine=INNODB;

create table detil_transaksi (
	id_transaksi int(11) not null,
	id_barang int(6) not null,
    jumlah int(6) not null,
	primary key(id_transaksi, id_barang),
	foreign key (id_barang) references barang(id),
	foreign key (id_transaksi) references transaksi(id_transaksi)
)engine=INNODB;