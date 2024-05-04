package model

import (
	"errors"
	"testing"
	"time"

	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TdTransaction struct {
	db.Datastore
	mock.Mock
}

func (td *TdTransaction) GetOrg(orgId string, userId string) (*types.Org, error) {
	org := &types.Org{
		Currency: "USD",
	}

	return org, nil
}

func (td *TdTransaction) GetPermissionedAccountIds(userId string, orgId string, tokenId string) ([]string, error) {
	return []string{"1", "2"}, nil
}

func (td *TdTransaction) GetAccountsByOrgId(orgId string) ([]*types.Account, error) {
	return []*types.Account{{Id: "1", Currency: "USD"}, {Id: "2"}}, nil
}

func (td *TdTransaction) InsertTransaction(transaction *types.Transaction) (err error) {
	return nil
}

func (td *TdTransaction) GetTransactionById(id string) (*types.Transaction, error) {
	args := td.Called(id)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (td *TdTransaction) UpdateTransaction(oldId string, transaction *types.Transaction) error {
	args := td.Called(oldId, transaction)
	return args.Error(0)
}

func (td *TdTransaction) GetOrgUserIds(id string) ([]string, error) {
	return []string{"1"}, nil
}

func TestCreateTransaction(t *testing.T) {
	tests := map[string]struct {
		err error
		tx  *types.Transaction
	}{
		"successful": {
			err: nil,
			tx: &types.Transaction{
				Id:          "1",
				OrgId:       "2",
				UserId:      "3",
				Date:        time.Now(),
				Inserted:    time.Now(),
				Updated:     time.Now(),
				Description: "description",
				Data:        "",
				Deleted:     false,
				Splits: []*types.Split{
					{TransactionId: "1", AccountId: "1", Amount: 1000, NativeAmount: 1000},
					{TransactionId: "1", AccountId: "2", Amount: -1000, NativeAmount: -1000},
				},
			},
		},
		"bad split amounts": {
			err: errors.New("splits must add up to 0"),
			tx: &types.Transaction{
				Id:          "1",
				OrgId:       "2",
				UserId:      "3",
				Date:        time.Now(),
				Inserted:    time.Now(),
				Updated:     time.Now(),
				Description: "description",
				Data:        "",
				Deleted:     false,
				Splits: []*types.Split{
					{TransactionId: "1", AccountId: "1", Amount: 1000, NativeAmount: 1000},
					{TransactionId: "1", AccountId: "2", Amount: -500, NativeAmount: -500},
				},
			},
		},
		"lacking permission": {
			err: errors.New("user does not have permission to access account 3"),
			tx: &types.Transaction{
				Id:          "1",
				OrgId:       "2",
				UserId:      "3",
				Date:        time.Now(),
				Inserted:    time.Now(),
				Updated:     time.Now(),
				Description: "description",
				Data:        "",
				Deleted:     false,
				Splits: []*types.Split{
					{TransactionId: "1", AccountId: "1", Amount: 1000, NativeAmount: 1000},
					{TransactionId: "1", AccountId: "3", Amount: -1000, NativeAmount: -1000},
				},
			},
		},
		"nativeAmount mismatch": {
			err: errors.New("nativeAmount must equal amount for native currency splits"),
			tx: &types.Transaction{
				Id:          "1",
				OrgId:       "2",
				UserId:      "3",
				Date:        time.Now(),
				Inserted:    time.Now(),
				Updated:     time.Now(),
				Description: "description",
				Data:        "",
				Deleted:     false,
				Splits: []*types.Split{
					{TransactionId: "1", AccountId: "1", Amount: 1000, NativeAmount: 500},
					{TransactionId: "1", AccountId: "2", Amount: -1000, NativeAmount: -500},
				},
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdTransaction{}
		model := NewModel(td, nil, types.Config{})

		err := model.CreateTransaction(test.tx)

		assert.Equal(t, err, test.err)
	}
}
