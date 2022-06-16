package models

type TransferLock struct {
	AccountOriginID      int64 `gorm:"uniqueIndex:transfer_lock_uniq"`
	AccountDestinationID int64 `gorm:"uniqueIndex:transfer_lock_uniq"`
}
