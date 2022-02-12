package concurrency

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type TestEntity struct {
	ID   uint
	Name string
	Version
	Relates []TestEntity2
}

type TestEntity2 struct {
	ID           uint
	TestEntityID uint
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

		DB.AutoMigrate(&TestEntity{}, &TestEntity2{})
	}
}

func OpenTestConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	return db.Debug(), err
}

func TestAutoSetIfEmpty(t *testing.T) {
	// test auto set if empty
	e := TestEntity{
		ID:   1,
		Name: "1",
		Relates: []TestEntity2{
			{ID: 1},
		},
	}
	err := DB.Create(&e).Error
	assert.NoError(t, err)
	assert.True(t, e.Version.Valid)
	assert.NotEmpty(t, e.Version.String)

}

//func TestNotSetIfPresent(t *testing.T) {
//	// test auto set if empty
//	e := TestEntity{
//		ID:      2,
//		Name:    "2",
//		Version: NewVersion(),
//		Relates: []TestEntity2{
//			{ID: 2},
//			{ID: 3, Version: NewVersion()},
//		},
//	}
//	ev := e.Version.String
//	err := DB.Create(&e).Error
//	assert.NoError(t, err)
//	assert.True(t, e.Version.Valid)
//	assert.NotEmpty(t, e.Version.String)
//	assert.Equal(t, ev, e.Version.String)
//}

func TestConcurrency(t *testing.T) {
	// test auto set if empty
	e := TestEntity{
		ID:   3,
		Name: "3",
	}
	err := DB.Create(&e).Error
	assert.NoError(t, err)

	// query for later error test
	var ec TestEntity
	err = DB.First(&ec, "id", 3).Error
	assert.NoError(t, err)
	assert.Equal(t, e.ID, ec.ID)

	//test update
	tx := DB.Model(&e).Update("name", "33")
	assert.NoError(t, tx.Error)
	assert.Equal(t, int64(1), tx.RowsAffected)
	assert.Equal(t, e.Name, "33")

	//version should be updated
	assert.True(t, e.Version.Valid)
	assert.NotEmpty(t, e.Version.String)
	assert.NotEqual(t, e.Version.String, ec.Version.String)

	//error version
	affected := DB.Model(&ec).Update("name", "33").RowsAffected
	assert.Equal(t, int64(0), affected)

	//test error
	err = ConcurrentUpdate(DB.Model(&ec), "name", "33").Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdates(DB.Model(&ec), map[string]interface{}{"name": "3"}).Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdates(DB.Model(&ec), ec).Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdateColumn(DB.Model(&ec), "name", "33").Error
	assert.ErrorIs(t, err, ErrConcurrent)

	err = ConcurrentUpdateColumns(DB.Model(&ec), map[string]interface{}{"name": "33"}).Error
	assert.ErrorIs(t, err, ErrConcurrent)

}
