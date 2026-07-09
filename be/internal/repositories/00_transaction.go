package repositories

import (
	"gorm.io/gorm"
)

// TransactionManager handles transaction lifecycle.
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithinTransaction executes a function within a transaction.
// Automatically commits on success, rolls back on error or panic.
func (tm *TransactionManager) WithinTransaction(fn func(tx *gorm.DB) error) error {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw panic
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// WithinTransactionWithResult executes a function within a transaction and returns a result.
func (tm *TransactionManager) WithinTransactionWithResult(fn func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw panic
		}
	}()

	result, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}
