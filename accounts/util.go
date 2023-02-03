package accounts

import (
	"fmt"
	"github.com/rail44/g/sqlc/generated"
	"strconv"
)

func mapToSubtype(entity sqlc.GetTransactionsRow) (interface{}, error) {
	accountId := int(entity.AccountID)
	if entity.MintID.Valid {
		amount, err := strconv.Atoi(entity.MintAmount.String)
		if err != nil {
			return nil, fmt.Errorf("parse amount as decimal: %w", err)
		}

		return Mint{Account: accountId, Amount: amount, InsertedAt: entity.InsertedAt, Type: MintTypeMint}, nil
	}

	if entity.SpendID.Valid {
		amount, err := strconv.Atoi(entity.SpendAmount.String)
		if err != nil {
			return nil, fmt.Errorf("parse amount as decimal: %w", err)
		}

		return Spend{Account: accountId, Amount: amount, InsertedAt: entity.InsertedAt, Type: SpendTypeSpend}, nil
	}

	if entity.TransferID.Valid {
		amount, err := strconv.Atoi(entity.TransferAmount.String)
		if err != nil {
			return nil, fmt.Errorf("parse amount as decimal: %w", err)
		}

		return Transfer{Account: accountId, Amount: amount, InsertedAt: entity.InsertedAt, Type: TransferTypeTransfer, Recipient: int(entity.TransferRecipient.Int64)}, nil
	}
	return nil, fmt.Errorf("failed to determine entity type")
}
