package concurrency

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type TestEntity struct {
	ID   uint
	Name string
	Version
}

func init() {
	var err error
	if DB, err = OpenTestConnection(); err != nil {
		log.Printf("failed to connect database, got error %v", err)
		os.Exit(1)
	} else {
		sqlDB, err := DB.DB()
		if err == nil {
			err = sqlDB.Ping()
		}

		if err != nil {
			log.Printf("failed to connect database, got error %v", err)
		}

		DB.AutoMigrate(&TestEntity{})
	}
}

func OpenTestConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.Logger = db.Logger.LogMode(logger.Info)
	return
}

func TestConcurrency(t *testing.T) {
	// test auto set if empty
	e := TestEntity{
		ID:   1,
		Name: "1",
	}
	err := DB.Create(&e).Error
	assert.NoError(t, err)
	ev1 := e.Version
	assert.True(t, ev1.Valid)
	assert.NotEmpty(t, ev1.String)
	// test not set if not empty
	e1 := TestEntity{
		ID:      2,
		Version: NewVersion(),
	}
	e1v := e1.Version.String
	err = DB.Create(&e1).Error
	assert.NoError(t, err)
	assert.Equal(t, e1v, e1.Version.String)

	// query for later error test
	var ec TestEntity
	err = DB.First(&ec, 1).Error
	assert.NoError(t, err)

	//test update
	tx := DB.Model(&e).Update("name", "2")
	assert.NoError(t, tx.Error)
	assert.Equal(t, int64(1), tx.RowsAffected)
	ev2 := e.Version
	assert.True(t, ev2.Valid)
	assert.NotEmpty(t, ev2.String)
	assert.NotEqual(t, ev1.String, ev2.String)

	affected := DB.Model(&ec).Update("name", "3").RowsAffected
	assert.Equal(t, int64(0), affected)
	//test error
	err = ConcurrentUpdate(DB.Model(&ec), "name", "3").Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdates(DB.Model(&ec), map[string]interface{}{"name": "3"}).Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdateColumn(DB.Model(&ec), "name", "3").Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdateColumns(DB.Model(&ec), map[string]interface{}{"name": "3"}).Error
	assert.ErrorIs(t, err, ErrConcurrent)

}
