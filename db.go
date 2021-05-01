package concurrency

import "gorm.io/gorm"

func ConcurrentUpdate(db *gorm.DB, column string, value interface{}) (tx *gorm.DB) {
	tx = db.Update(column, value)
	if tx.RowsAffected == 0 {
		tx.AddError(ErrConcurrent)
	}
	return
}
func ConcurrentUpdates(db *gorm.DB, values interface{}) (tx *gorm.DB) {
	tx = db.Updates(values)
	if tx.RowsAffected == 0 {
		tx.AddError(ErrConcurrent)
	}
	return
}
func ConcurrentUpdateColumn(db *gorm.DB, column string, value interface{}) (tx *gorm.DB) {
	tx = db.UpdateColumn(column, value)
	if tx.RowsAffected == 0 {
		tx.AddError(ErrConcurrent)
	}
	return
}
func ConcurrentUpdateColumns(db *gorm.DB, values interface{}) (tx *gorm.DB) {
	tx = db.UpdateColumns(values)
	if tx.RowsAffected == 0 {
		tx.AddError(ErrConcurrent)
	}
	return
}
