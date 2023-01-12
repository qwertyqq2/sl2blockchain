package blockchain

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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

func existLevel(filename string) (bool, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return false, ErrNotFileExist
	}
	query := fmt.Sprintf("SELECT Block FROM Blockchain")
	row := db.QueryRow(query)
	var tmp interface{}
	err = row.Scan(&tmp)
	if tmp == nil {
		return false, ErrNotRows
	}
	return true, nil
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
	return index
}

func (l *LevelDB) insertBlock(hash []byte, block string) error {
	_, err := l.db.Exec("INSERT INTO Blockchain (Hash, Block) VALUES ($1, $2)",
		Base64Encode(hash),
		block,
	)
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
	var bs string
	row := l.db.QueryRow("SELECT Block FROM Blockchain ORDER BY Id DESC")
	err := row.Scan(&bs)
	if err != nil {
		fmt.Println(err)
	}
	block, err := DeserializeBlock(bs)
	if err != nil {
		log.Fatal(err)
	}
	return block
}

func (l *LevelDB) idByHash(hash []byte) (uint64, error) {
	var idscan uint64
	row := l.db.QueryRow("Select Id from Blockchain where Hash=$1", Base64Encode(hash))
	err := row.Scan(&idscan)
	if err != nil {
		return 0, err
	}
	return idscan, nil
}

func (l *LevelDB) blockByHash(hash []byte) (*Block, error) {
	var pBlock string
	row := l.db.QueryRow("Select Block from Blockchain where Hash=$1", Base64Encode(hash))
	err := row.Scan(&pBlock)
	if err != nil {
		return nil, err
	}
	b, err := DeserializeBlock(pBlock)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (l *LevelDB) blockById(id uint64) (*Block, error) {
	var pBlock string
	row := l.db.QueryRow("Select Block from Blockchain where Id=$1", id)
	err := row.Scan(&pBlock)
	if err != nil {
		return nil, err
	}
	b, err := DeserializeBlock(pBlock)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (l *LevelDB) getBlocksFromHash(hash []byte) ([]*Block, error) {
	idx, err := l.idByHash(hash)
	if err != nil {
		return nil, err
	}
	curId := l.size()
	blocks := make([]*Block, 0)
	for i := idx + 1; i <= curId; i++ {
		b, err := l.blockById(i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (l *LevelDB) getBlockAfter(hash []byte) (*Block, error) {
	idx, err := l.idByHash(hash)
	if err != nil {
		return nil, err
	}
	if idx == l.size() {
		return nil, ErrIsLastBlock
	}
	return l.blockById(idx + 1)
}

func (l *LevelDB) getBlocks() ([]string, error) {
	rows, err := l.db.Query("Select Block from Blockchain")
	if err != nil {
		return nil, err
	}
	blocks := make([]string, 0)
	var bs string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&bs)
		blocks = append(blocks, bs)
	}
	return blocks, nil
}
