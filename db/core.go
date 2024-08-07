package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/chas3air/CarrotDB/models"
)

type DB struct {
	Values      map[any]any              `json:"values"`
	TempKey     map[any]models.TimedItem `json:"tempKey"`
	RecoverTime int                      `json:"recover_time"`
}

var mut sync.Mutex

func Init() DB {
	return DB{Values: make(map[any]any, 10), TempKey: make(map[any]models.TimedItem, 10), RecoverTime: 5}
}

func (db *DB) SetLifetime(seconds int) {
	db.RecoverTime = seconds
}

func (db *DB) GetItem(key any) (any, error) {
	mut.Lock()
	defer mut.Unlock()

	item, ok := db.Values[key]
	if ok {
		return item, nil
	} else {
		return nil, errors.New("key is undefined")
	}
}

func (db *DB) SetItem(key any, value any) {
	mut.Lock()
	defer mut.Unlock()

	db.Values[key] = value
}

func (db *DB) SetTempItem(key, value any, seconds int) {
	mut.Lock()
	defer mut.Unlock()

	db.Values[key] = value
	db.TempKey[key] = models.TimedItem{time.Now(), seconds}
}

func (db *DB) DeleteItem(key any) {
	mut.Lock()
	defer mut.Unlock()

	delete(db.Values, key)
}

func (db *DB) SaveAll(path string) error {
	mut.Lock()
	defer mut.Unlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := json.Marshal(db.Values)
	if err != nil {
		return err
	}

	_, err = file.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) LoadAll(path string) error {
	mut.Lock()
	defer mut.Unlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &db.Values)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Clear() {
	mut.Lock()
	defer mut.Unlock()
	db.Values = map[any]any{}
}

func (db *DB) GetItems() map[any]any {
	mut.Lock()
	defer mut.Unlock()

	return db.Values
}

func (db *DB) DbCleaner() {
	for {
		for i, v := range db.TempKey {
			if time.Now().Sub(v.Borntime) > time.Duration(v.Lifetime) {
				delete(db.Values, i)
				delete(db.TempKey, i)
			}
			time.Sleep(time.Second * time.Duration(db.RecoverTime))
		}
	}
}
