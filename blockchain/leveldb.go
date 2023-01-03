package blockchain

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/qwertyqq2/sl2blockchain/crypto"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbname = "sl2"
)

type LevelDB struct {
	db *sql.DB
}

func NewLevelDb(filename string) (*LevelDB, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	file.Close()
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
	CREATE TABLE Blockchain (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Hash VARCHAR(44) UNIQUE,
		Block TEXT
	);
	`)
	if err != nil {
		return nil, err
	}
	return &LevelDB{
		db: db,
	}, nil
}

func loadLevel(db *sql.DB) *LevelDB {
	return &LevelDB{
		db: db,
	}
}

func loadBlockchain(filename string) (*Blockchain, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	l := loadLevel(db)
	bc := &Blockchain{
		levelDb: l,
	}
	return bc, nil
}

func (l *LevelDB) size() uint64 {
	var index uint64
	row := l.db.QueryRow("SELECT Id FROM Blockchain ORDER BY Id DESC")
	row.Scan(&index)
	fmt.Println("size", index)
	return index
}

func (l *LevelDB) insertBlock(hash []byte, block string) error {
	_, err := l.db.Exec("INSERT INTO Blockchain (Hash, Block) VALUES ($1, $2)",
		hash,
		block,
	)
	fmt.Println("insert block", block, "with hash ", string(hash))
	return err
}

func (l *LevelDB) balance(address string, size uint64) (uint64, error) {
	var (
		sblock  string
		block   *Block
		balance uint64
	)
	rows, err := l.db.Query("SELECT Block FROM Blockchain WHERE Id <= $1 ORDER BY Id DESC", size)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&sblock)
		block, err = DeserializeBlock(sblock)
		if err != nil {
			return 0, err
		}
		if value, ok := block.Mapping[address]; ok {
			balance = value
			break
		}
	}
	return balance, nil
}

func (l *LevelDB) lastBlock() *Block {
	var block Block
	row := l.db.QueryRow("SELECT * FROM Blockchain ORDER BY Id DESC")
	row.Scan(&block)
	return &block
}

func (l *LevelDB) idByHash(hash []byte) (uint64, error) {
	var idscan uint64
	row := l.db.QueryRow("select id from Blockchain where Hash = $1", crypto.Base64Encode(hash))
	err := row.Scan(&idscan)
	if err != nil {
		return 0, err
	}
	return idscan, nil
}

func (l *LevelDB) blockByHash(hash []byte) (*Block, error) {
	var pBlock string
	row := l.db.QueryRow("Select Block from Blockchain where Hash=$1", hash)
	row.Scan(pBlock)
	b, err := DeserializeBlock(pBlock)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (l *LevelDB) getBlocks() ([]*Block, error) {
	rows, err := l.db.Query("Select Block from Blockchain")
	if err != nil {
		return nil, err
	}
	blocks := make([]*Block, 0)
	var bs string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&bs)
		b, err := DeserializeBlock(bs)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}
