package util

import (
    "fmt"
    "time"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

func NewDB(dsn string) (db *gorm.DB,err error){
    db,err = gorm.Open("mysql", dsn)
    if err != nil {
        fmt.Println("mysql conn error"+err.Error())
        return
    }
    db.DB().SetMaxIdleConns(5)
    db.DB().SetConnMaxLifetime(time.Second*60)
    db.DB().SetMaxOpenConns(50)
    db.LogMode(true)
    return
}