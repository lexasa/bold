package model

import (
	"github.com/boltdb/bolt"
	"encoding/binary"
)

type AccountJSON struct {
	AccountID int
	Status    string
	Tag       string
	CPM       float64
}

type Account struct {
	AccountID int
	Active    bool
	Tag       string
	CPM       float64
}

func (data *AccountJSON) Load() *Account {
	v := Account{AccountID:data.AccountID, Tag:data.Tag, CPM:data.CPM}
	if data.Status == StatusEnabledTxt {
		v.Active = StatusEnabled
	}else {
		v.Active = StatusDisabled
	}
	return &v
}

func (Account *Account) Store(db *bolt.DB) error {

	if err := db.Update(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)
		binary.PutVarint(idBuf, int64(Account.AccountID))

		b, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		if err != nil {
			return err
		}

		vBuf, err := Marshal(Account);
		if err != nil {
			return err
		}

		if err := b.Put(idBuf, vBuf); err != nil {
			return err
		}

		if err := b.Put([]byte(Account.Tag), vBuf); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func GetAccountByID(db *bolt.DB, id int) (*Account, error) {
	var ret Account
	if err := db.View(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)

		binary.PutVarint(idBuf, int64(id))
		data := tx.Bucket([]byte("accounts")).Get(idBuf)

		if err := Unmarshal(data, &ret); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetAccountByTag(db *bolt.DB, tag string) (*Account, error) {
	var ret Account
	if err := db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("accounts")).Get([]byte(tag))

		if err := Unmarshal(data, &ret); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}