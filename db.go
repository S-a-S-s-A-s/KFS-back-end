package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

type File struct {
	ID        int64
	Name      string
	Size      string
	Project   string
	UserName  string
	CreatTime string
	IsPublic  bool
}

type FileInfo struct {
	ID        int64
	FileId    int64
	Filer     File   `gorm:"foreignKey:FileId"`
	Operation string //
	Status    bool
	Operator  string
	CreatTime string
}

// Init init DB
func Initdb() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/khfs?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	DB.AutoMigrate(&File{}, &FileInfo{})
}

func GetQuery(role, name, project, userName string) (query string) {
	if name != "" {
		query = "name LIKE \"%" + name + "%\""
	}
	if project != "" {
		if query == "" {
			query = "project = \"" + project + "\""
		} else {
			query += " AND project = \"" + project + "\""
		}
	}
	if role == "vip" {
		query += " AND (user_name = \"" + userName + "\"" + " OR is_public = true)"
	}
	if role == "common" {
		query += " AND is_public = true"
	}
	return query
}
