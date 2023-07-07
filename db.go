package main

import (
	basicLog "log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
)

type Sense struct {
	Usage     uint16
	Solar     uint16
	Day       uint16
	UpdatedAt time.Time
}

var db *gorm.DB

func dbInit(migrate bool) {
	var err error
	newLogger := gLog.New(
		basicLog.New(os.Stdout, "\r\n", basicLog.LstdFlags), // io writer
		gLog.Config{
			SlowThreshold: 5 * time.Second, // Slow SQL threshold
			Colorful:      true,            // Disable color
		},
	)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("PG_URI"),
		PreferSimpleProtocol: true, // necessary when pg pooling
	}), &gorm.Config{Logger: newLogger})
	check(err)

	if migrate {
		log.Info().Msg("Migrating")
		db.AutoMigrate(&Sense{})
	}
}

func write(sense Sense) {
	log.Info().Msg("Inserting Sense Data")
	res := db.Create(&sense) // pass pointer of data to Create
	check(res.Error)
	log.Info().Msg("inserted " + strconv.Itoa(int(res.RowsAffected)) + " row")
}

func clearPrevWeek(day uint16) {
	log.Print("deleting all data for day ", day)
	res := db.Where("day = ?", day).Delete(&Sense{})
	check(res.Error)
	log.Info().Msg("deleted " + strconv.Itoa(int(res.RowsAffected)) + " rows")
}
