package blockchain_test

import (
	"testing"

	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

type MockRepostiory struct {
	loadFunc func(aggregate midgard.Aggregate, aggregateID string) error
	saveFunc func(aggregateID string, events ...midgard.Event) error
}

func (m MockRepostiory) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return m.loadFunc(aggregate, aggregateID)
}

func (m MockRepostiory) Save(aggregateID string, events ...midgard.Event) error {
	return m.saveFunc(aggregateID, events...)
}

func (MockRepostiory) Close() {}

func TestCreateGenesisBlock(t *testing.T) {
	//given

	tests := map[string]struct {
		input struct {
			ConfigFilePath string
		}
		output blockchain.Block
		err    error
	}{
		"success create genesisBlock": {

			input: struct {
				ConfigFilePath string
			}{
				ConfigFilePath: "./GenesisBlockConfig.json",
			},

			output: &blockchain.DefaultBlock{
				PrevSeal:  make([]byte, 0),
				Height:    uint64(0),
				TxList:    make([]blockchain.Transaction, 0),
				TxSeal:    make([][]byte, 0),
				Timestamp: (time.Now()).Round(100 * time.Millisecond),
				Creator:   make([]byte, 0),
			},

			err: nil,
		},

		"fail create genesisBlock: wrong file path": {

			input: struct {
				ConfigFilePath string
			}{
				ConfigFilePath: "./WrongBlockConfig.json",
			},

			output: nil,

			err: blockchain.ErrSetConfig,
		},
	}

	timeStamp := (time.Now()).Round(100 * time.Millisecond)
	prevSeal := make([]byte, 0)
	txSeal := make([][]byte, 0)
	creator := make([]byte, 0)
	validator := blockchain.DefaultValidator{}
	Seal, err := validator.BuildSeal(timeStamp, prevSeal, txSeal, creator)

	if err != nil {
		log.Println(err.Error())
	}

	repo := MockRepostiory{}

	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, string(Seal), aggregateID)
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}

	eventstore.InitForMock(repo)
	defer eventstore.Close()

	GenesisFilePath := "./GenesisBlockConfig.json"

	defer os.Remove(GenesisFilePath)

	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)

	err = ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigJson, 0644)

	if err != nil {
		log.Println(err.Error())
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		GenesisBlock, err := blockchain.CreateGenesisBlock(test.input.ConfigFilePath)

		//then
		assert.Equal(t, test.err, err)

		if err != nil {
			assert.Equal(t, test.output, GenesisBlock)
			continue
		}

		assert.Equal(t, test.output.GetPrevSeal(), GenesisBlock.GetPrevSeal())
		assert.Equal(t, test.output.GetHeight(), GenesisBlock.GetHeight())
		assert.Equal(t, test.output.GetTxList(), GenesisBlock.GetTxList())
		assert.Equal(t, test.output.GetTxSeal(), GenesisBlock.GetTxSeal())
		assert.Equal(t, test.output.GetTimestamp(), GenesisBlock.GetTimestamp())
		assert.Equal(t, test.output.GetCreator(), GenesisBlock.GetCreator())

	}

}

func TestCreateProposedBlock(t *testing.T) {

	//given

	tests := map[string]struct {
		input struct {
			prevSeal []byte
			height   uint64
			txList   []blockchain.Transaction
			creator  []byte
		}
		output blockchain.Block
		err    error
	}{
		"success create proposed block": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				creator  []byte
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{},
				},
				creator: []byte("junksound"),
			},

			output: &blockchain.DefaultBlock{
				PrevSeal: []byte("prevseal"),
				Height:   1,
				TxList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{},
				},
				Timestamp: (time.Now()).Round(100 * time.Millisecond),
				Creator:   []byte("junksound"),
			},

			err: nil,
		},

		"fail case1: without transaction": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				creator  []byte
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList:   nil,
				creator:  []byte("junksound"),
			},

			output: nil,

			err: blockchain.ErrBuildingTxSeal,
		},

		"fail case2: without prevseal or creator": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				creator  []byte
			}{
				prevSeal: nil,
				height:   1,
				txList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{},
				},
				creator: nil,
			},

			output: nil,

			err: blockchain.ErrBuildingSeal,
		},
	}

	timeStamp := (time.Now()).Round(100 * time.Millisecond)
	prevSeal := []byte("prevseal")
	txList := []blockchain.Transaction{
		&blockchain.DefaultTransaction{},
	}
	creator := []byte("junksound")

	validator := blockchain.DefaultValidator{}
	txSeal, err := validator.BuildTxSeal(txList)
	Seal, err := validator.BuildSeal(timeStamp, prevSeal, txSeal, creator)

	if err != nil {
		log.Println(err.Error())
	}

	repo := MockRepostiory{}

	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, string(Seal), aggregateID)
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}

	eventstore.InitForMock(repo)
	defer eventstore.Close()

	for testName, test := range tests {

		t.Logf("Running test case %s", testName)

		//when
		ProposedBlock, err := blockchain.CreateProposedBlock(
			test.input.prevSeal,
			test.input.height,
			test.input.txList,
			test.input.creator,
		)

		//then
		assert.Equal(t, test.err, err)

		if err != nil {
			assert.Equal(t, test.output, ProposedBlock)
			continue
		}

		assert.Equal(t, test.output.GetPrevSeal(), ProposedBlock.GetPrevSeal())
		assert.Equal(t, test.output.GetHeight(), ProposedBlock.GetHeight())
		assert.Equal(t, test.output.GetTxList(), ProposedBlock.GetTxList())
		assert.Equal(t, test.output.GetTimestamp(), ProposedBlock.GetTimestamp())
		assert.Equal(t, test.output.GetCreator(), ProposedBlock.GetCreator())
	}

}
