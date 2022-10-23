package database

import (
	"github.com/flashbots/mev-boost-relay/common"
)

var (
	tableBase = common.GetEnv("DB_TABLE_PREFIX", "rsdev")

	TableSignedBuilderBid        = tableBase + "_signed_builder_bid"
	TableDataAPIPayloadDelivered = tableBase + "_data_api_payload_delivered"
	TableDataAPIBuilderBid       = tableBase + "_data_api_builder_bid"
	TableError                   = tableBase + "_error"
)

var schema = `
CREATE TABLE IF NOT EXISTS ` + TableSignedBuilderBid + ` (
	id          bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,

	relay        text NOT NULL,
	requested_at timestamp NOT NULL,
	received_at  timestamp NOT NULL,
	duration_ms	 bigint NOT NULL,

	slot            bigint NOT NULL,
	parent_hash     varchar(66) NOT NULL,
	proposer_pubkey	varchar(98) NOT NULL,

	pubkey 		varchar(98) NOT NULL,
	signature   text NOT NULL,

	value         NUMERIC(48, 0) NOT NULL,
	fee_recipient varchar(42) NOT NULL,
	block_hash    varchar(66) NOT NULL,
	block_number  bigint NOT NULL,
	gas_limit     bigint NOT NULL,
	gas_used      bigint NOT NULL,
	extra_data    text NOT NULL,

	epoch bigint NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ` + TableSignedBuilderBid + `_u_relay_slot_n_hashes_idx ON ` + TableSignedBuilderBid + `("relay", "slot", "parent_hash", "block_hash");
CREATE INDEX IF NOT EXISTS ` + TableSignedBuilderBid + `_insertedat_idx ON ` + TableSignedBuilderBid + `("inserted_at");
CREATE INDEX IF NOT EXISTS ` + TableSignedBuilderBid + `_slot_idx ON ` + TableSignedBuilderBid + `("slot");
CREATE INDEX IF NOT EXISTS ` + TableSignedBuilderBid + `_block_number_idx ON ` + TableSignedBuilderBid + `("block_number");
CREATE INDEX IF NOT EXISTS ` + TableSignedBuilderBid + `_value_idx ON ` + TableSignedBuilderBid + `("value");


CREATE TABLE IF NOT EXISTS ` + TableDataAPIPayloadDelivered + ` (
	id          bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,
	relay       text NOT NULL,

	epoch bigint NOT NULL,
	slot  bigint NOT NULL,

	parent_hash            varchar(66) NOT NULL,
	block_hash             varchar(66) NOT NULL,
	builder_pubkey         varchar(98) NOT NULL,
	proposer_pubkey        varchar(98) NOT NULL,
	proposer_fee_recipient varchar(42) NOT NULL,
	gas_limit              bigint NOT NULL,
	gas_used               bigint NOT NULL,
	value_claimed_wei      NUMERIC(48, 0) NOT NULL,
	value_claimed_eth      NUMERIC(16, 8) NOT NULL,
	num_tx                 int,
	block_number           bigint,

	value_check_ok              boolean, 		-- null means not yet checked
	value_check_method          text,  		    -- how value was checked (i.e. blockBalanceDiff)
	value_delivered_wei         NUMERIC(48, 0), -- actually delivered value
	value_delivered_eth         NUMERIC(16, 8), -- actually delivered value
	value_delivered_diff_wei    NUMERIC(48, 0), -- value_delivered - value_claimed
	value_delivered_diff_eth    NUMERIC(16, 8), -- value_delivered - value_claimed
	block_coinbase_addr		    varchar(42),    -- block coinbase address
	block_coinbase_is_proposer  boolean,        -- true if coinbase == proposerFeeRecipient
	coinbase_diff_wei           NUMERIC(48, 0), -- builder value difference
	coinbase_diff_eth           NUMERIC(16, 8), -- builder value difference
	found_onchain               boolean         -- whether this block was found on chain
);

CREATE UNIQUE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_u_relay_slot_blockhash_idx ON ` + TableDataAPIPayloadDelivered + `("relay", "slot", "parent_hash", "block_hash");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_insertedat_idx ON ` + TableDataAPIPayloadDelivered + `("inserted_at");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_slot_idx ON ` + TableDataAPIPayloadDelivered + `("slot");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_builder_pubkey_idx ON ` + TableDataAPIPayloadDelivered + `("builder_pubkey");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_block_number_idx ON ` + TableDataAPIPayloadDelivered + `("block_number");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_value_wei_idx ON ` + TableDataAPIPayloadDelivered + `("value_claimed_wei");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_value_eth_idx ON ` + TableDataAPIPayloadDelivered + `("value_claimed_eth");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIPayloadDelivered + `_valuecheck_ok_idx ON ` + TableDataAPIPayloadDelivered + `("value_check_ok");


CREATE TABLE IF NOT EXISTS ` + TableDataAPIBuilderBid + ` (
	id          bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	inserted_at timestamp NOT NULL default current_timestamp,
	relay       text NOT NULL,

	epoch bigint NOT NULL,
	slot  bigint NOT NULL,

	parent_hash            varchar(66) NOT NULL,
	block_hash             varchar(66) NOT NULL,
	builder_pubkey         varchar(98) NOT NULL,
	proposer_pubkey        varchar(98) NOT NULL,
	proposer_fee_recipient varchar(42) NOT NULL,
	gas_limit              bigint NOT NULL,
	gas_used               bigint NOT NULL,
	value                  NUMERIC(48, 0) NOT NULL,
	num_tx                 int,
	block_number           bigint,
	timestamp			   timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ` + TableDataAPIBuilderBid + `_unique_idx ON ` + TableDataAPIBuilderBid + `("relay", "slot", "builder_pubkey", "parent_hash", "block_hash");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIBuilderBid + `_insertedat_idx ON ` + TableDataAPIBuilderBid + `("inserted_at");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIBuilderBid + `_slot_idx ON ` + TableDataAPIBuilderBid + `("slot");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIBuilderBid + `_builder_pubkey_idx ON ` + TableDataAPIBuilderBid + `("builder_pubkey");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIBuilderBid + `_block_number_idx ON ` + TableDataAPIBuilderBid + `("block_number");
CREATE INDEX IF NOT EXISTS ` + TableDataAPIBuilderBid + `_value_idx ON ` + TableDataAPIBuilderBid + `("value");
`
