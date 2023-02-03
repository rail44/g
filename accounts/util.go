package accounts;

import (
	"fmt"
	"github.com/rail44/g/sqlc/generated"
	"strconv"
)


func mapToSubtype(tx sqlc.GetTransactionsRow) (interface{}, error) {
	if tx.MintID.Valid {
		amount, err := strconv.Atoi(tx.MintAmount.String)
		if err != nil {
			return nil, fmt.Errorf("parse amount as decimal: %w", err)
		}

		_type := MintTypeMint
		return Mint{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type}, nil
	}

	if tx.SpendID.Valid {
		amount, err := strconv.Atoi(tx.SpendAmount.String)
		if err != nil {
			return nil, fmt.Errorf("parse amount as decimal: %w", err)
		}

		_type := SpendTypeSpend
		return Spend{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type}, nil
	}

	if tx.TransferID.Valid {
		amount, err := strconv.Atoi(tx.TransferAmount.String)
		if err != nil {
			return nil, fmt.Errorf("parse amount as decimal: %w", err)
		}

		_type := TransferTypeTransfer
		recipient := int(tx.TransferRecipient.Int64)
		return Transfer{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type, Recipient: &recipient}, nil
	}
	return nil, fmt.Errorf("failed to determine tx type")
}
