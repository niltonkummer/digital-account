package repository

import (
	"context"
	"digital-account/application/models"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type Transfer interface {
	Create(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)
	List(ctx context.Context, id int64) (models.Transfers, error)
}

type transferRepo struct {
	DB *gorm.DB
}

// Create executes a transfer process as a transaction
func (a *transferRepo) Create(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error) {
	err := a.DB.Transaction(func(tx *gorm.DB) (err error) {

		transferLock := &models.TransferLock{
			AccountOriginID:      transfer.AccountOriginID,
			AccountDestinationID: transfer.AccountDestinationID,
		}

		res := tx.WithContext(ctx).
			Create(transferLock)
		if res.Error != nil {
			return res.Error
		}
		defer func(tx *gorm.DB, transferLock *models.TransferLock, transfer *models.Transfer) {
			res := tx.WithContext(ctx).
				Where("account_origin_id", transfer.AccountOriginID).
				Where("account_destination_id", transfer.AccountDestinationID).
				Delete(&transferLock)
			err = res.Error
		}(tx, transferLock, transfer)

		var origin models.Account
		res = tx.WithContext(ctx).
			First(&origin, transfer.AccountOriginID)
		if res.Error != nil {
			return res.Error
		}

		var destination models.Account
		res = tx.WithContext(ctx).
			First(&destination, transfer.AccountDestinationID)
		if res.Error != nil {
			return res.Error
		}

		res = tx.WithContext(ctx).
			Model(&transfer.AccountOrigin).
			Clauses(clause.Returning{}).
			Where("id", transfer.AccountOriginID).
			Update("balance", origin.Balance.Sub(transfer.Amount))

		if res.Error != nil {
			return res.Error
		}

		res = tx.WithContext(ctx).
			Model(&transfer.AccountDestination).
			Clauses(clause.Returning{}).
			Where("id", transfer.AccountDestinationID).
			Update("balance", destination.Balance.Add(transfer.Amount))

		if res.Error != nil {
			return res.Error
		}

		res = tx.WithContext(ctx).
			Clauses(clause.Returning{}).
			Create(&transfer)
		if res.Error != nil {
			return res.Error
		}

		return nil
	})
	return transfer, err
}

// List retrieves a list of transfers in accord with filter
func (a *transferRepo) List(ctx context.Context, id int64) (transfer models.Transfers, err error) {

	res := a.DB.WithContext(ctx).Scopes(Paginate(ctx)).Where("account_origin_id = ?", id).Find(&transfer)
	err = res.Error
	return
}

func CreateTransferRepository(db *gorm.DB) Transfer {
	return &transferRepo{
		DB: db,
	}
}
