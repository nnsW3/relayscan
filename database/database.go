// Package database exposes the postgres database
package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IDatabaseService interface {
	SaveBidForSlot(relay string, slot uint64, parentHash, proposerPubkey string, respStatus uint64, respBid any, respError string, durationMs uint64) error

	GetDataAPILatestPayloadDelivered(relay string) (*PayloadDeliveredEntry, error)
	SaveDataAPIPayloadDelivered(entry *PayloadDeliveredEntry) error
}

type DatabaseService struct {
	DB *sqlx.DB
}

func NewDatabaseService(dsn string) (*DatabaseService, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.DB.SetMaxOpenConns(50)
	db.DB.SetMaxIdleConns(10)
	db.DB.SetConnMaxIdleTime(0)

	if os.Getenv("PRINT_SCHEMA") == "1" {
		fmt.Println(schema)
	}

	if os.Getenv("DB_DONT_APPLY_SCHEMA") == "" {
		_, err = db.Exec(schema)
		if err != nil {
			return nil, err
		}
	}

	return &DatabaseService{
		DB: db,
	}, nil
}

func (s *DatabaseService) Close() error {
	return s.DB.Close()
}

func (s *DatabaseService) SaveBidForSlot(relay string, slot uint64, parentHash, proposerPubkey string, respStatus uint64, respBid any, respError string, durationMs uint64) error {
	return nil
}

func (s *DatabaseService) SaveDataAPIPayloadDelivered(entry *PayloadDeliveredEntry) error {
	query := `INSERT INTO ` + TableDataAPIPayloadDelivered + `
		(relay, epoch, slot, parent_hash, block_hash, builder_pubkey, proposer_pubkey, proposer_fee_recipient, gas_limit, gas_used, value, num_tx, block_number) VALUES
		(:relay, :epoch, :slot, :parent_hash, :block_hash, :builder_pubkey, :proposer_pubkey, :proposer_fee_recipient, :gas_limit, :gas_used, :value, :num_tx, :block_number)
		ON CONFLICT DO NOTHING`
	_, err := s.DB.NamedExec(query, entry)
	return err
}

func (s *DatabaseService) GetDataAPILatestPayloadDelivered(relay string) (*PayloadDeliveredEntry, error) {
	entry := new(PayloadDeliveredEntry)
	query := `SELECT id, inserted_at, relay, epoch, slot, parent_hash, block_hash, builder_pubkey, proposer_pubkey, proposer_fee_recipient, gas_limit, gas_used, value, num_tx, block_number FROM ` + TableDataAPIPayloadDelivered + ` WHERE relay=$1 ORDER BY slot DESC LIMIT 1`
	err := s.DB.Get(entry, query, relay)
	return entry, err
}

// func (s *DatabaseService) SaveValidatorRegistration(registration types.SignedValidatorRegistration) error {
// 	entry := ValidatorRegistrationEntry{
// 		Pubkey:       registration.Message.Pubkey.String(),
// 		FeeRecipient: registration.Message.FeeRecipient.String(),
// 		Timestamp:    registration.Message.Timestamp,
// 		GasLimit:     registration.Message.GasLimit,
// 		Signature:    registration.Signature.String(),
// 	}

// 	query := `INSERT INTO ` + TableValidatorRegistration + `
// 		(pubkey, fee_recipient, timestamp, gas_limit, signature) VALUES
// 		(:pubkey, :fee_recipient, :timestamp, :gas_limit, :signature)
// 		ON CONFLICT (pubkey, fee_recipient) DO NOTHING;`
// 	_, err := s.DB.NamedExec(query, entry)
// 	return err
// }
