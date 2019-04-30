package lunchordersdb

import (
	"encoding/binary"
	"encoding/json"

	"github.com/dgraph-io/badger"
)

// LunchDB represents a lunchorder database
type LunchDB interface {
	AddItem(table Table, item canSetID) error
	UpdateItem(table Table, item canGetID) error
	GetItem(table Table, id uint64, result interface{}) error
	IterateItems(table Table, process func(value []byte) error) error
	DeleteItem(table Table, id uint64) error
	Close()
}

type lunchDB struct {
	db *badger.DB
}

type canGetID interface {
	getID() uint64
}

type canSetID interface {
	setID(id uint64)
}

// Person represents a person ordering lunch
type Person struct {
	ID        uint64 `json:"id,omitempty"`
	Firstname string `json:"fn"`
	Lastname  string `json:"ln"`
}

func (p Person) getID() uint64 {
	return p.ID
}

func (p *Person) setID(id uint64) {
	p.ID = id
}

// Meal represents a meal that can be ordered
type Meal struct {
	ID     uint64 `json:"id,omitempty"`
	Desc   string `json:"desc"`
	Price  uint32 `json:"price"`
	Active bool   `json:"active,omitempty"`
}

func (m *Meal) getID() uint64 {
	return m.ID
}

func (m *Meal) setID(id uint64) {
	m.ID = id
}

// LunchOrder represents a lunch order
type LunchOrder struct {
	ID       uint64 `json:"id,omitempty"`
	Date     uint32 `json:"d"`
	MealID   uint16 `json:"m"`
	PersonID uint32 `json:"p"`
}

type Table struct {
	prefix      byte
	prefixBytes []byte
}

func NewTable(prefix byte) Table {
	return Table{prefix: prefix, prefixBytes: []byte{prefix}}
}

func (t Table) getKey(id uint64) []byte {
	keyBuf := make([]byte, 1+8)
	keyBuf[0] = t.prefix
	binary.BigEndian.PutUint64(keyBuf[1:], id)
	return keyBuf
}

// Open opens the lunchorder DB in the given directory
func Open(dbDir string) (db LunchDB, err error) {
	opts := badger.DefaultOptions
	opts.Dir = dbDir
	opts.ValueDir = dbDir
	badgerDb, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return lunchDB{db: badgerDb}, nil
}

func (db lunchDB) getNewID(prefix []byte) (uint64, error) {
	seq, err := db.db.GetSequence(prefix, 1000)
	if err != nil {
		return 0, err
	}
	defer seq.Release()
	num, err := seq.Next()
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (db lunchDB) GetItem(table Table, id uint64, result interface{}) error {
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(table.getKey(id))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &result)
		})
	})

	return err
}

func (db lunchDB) IterateItems(table Table, process func(value []byte) error) error {
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := make([]byte, 8)
		prefix[0] = table.prefix
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			item.Key()
			err := item.Value(func(v []byte) error {
				return process(v)
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (db lunchDB) AddItem(table Table, item canSetID) error {
	personID, err := db.getNewID(table.prefixBytes)
	if err != nil {
		return err
	}

	item.setID(personID)
	return db.db.Update(func(txn *badger.Txn) error {
		valueBuffer, _ := json.Marshal(item)
		return txn.Set(table.getKey(personID), valueBuffer)
	})
}

func (db lunchDB) UpdateItem(table Table, item canGetID) error {
	return db.db.Update(func(txn *badger.Txn) error {
		key := table.getKey(item.getID())
		_, err := txn.Get(key)
		if err != nil {
			return err
		}

		valueBuffer, _ := json.Marshal(item)
		return txn.Set(key, valueBuffer)
	})
}

func (db lunchDB) DeleteItem(table Table, id uint64) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(table.getKey(id))
	})
}

// Close closes the lunchorder DB
func (db lunchDB) Close() {
	db.db.Close()
}
