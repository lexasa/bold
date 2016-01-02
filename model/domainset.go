
package model

import (
	"github.com/boltdb/bolt"
	"encoding/binary"
)

type DomainsetJSON struct {
	DomainsetID int
	Name        string
	Status      string
	Action      string
}

type Domainset struct {
	DomainsetID int
	Active      bool
	Action      bool
}

func (data *DomainsetJSON) Load() *Domainset {
	v := Domainset{DomainsetID:data.DomainsetID}
	if data.Status == StatusEnabledTxt {
		v.Active = StatusEnabled
	}else {
		v.Active = StatusDisabled
	}
	if data.Action == ActionAllowTxt {
		v.Action = ActionAllow
	}else {
		v.Action = ActionDeny
	}
	return &v
}

func (domainset *Domainset) Store(db *bolt.DB) error {

	if err := db.Update(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)
		binary.PutVarint(idBuf, int64(domainset.DomainsetID))

		b, err := tx.CreateBucketIfNotExists([]byte("domainsets"))
		if err != nil {
			return err
		}

		vBuf, err := Marshal(domainset);
		if err != nil {
			return err
		}

		if err := b.Put(idBuf, vBuf); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func GetDomainsetByID(db *bolt.DB, id int) (*Domainset, error) {
	var ret Domainset
	if err := db.View(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)

		binary.PutVarint(idBuf, int64(id))
		data:=tx.Bucket([]byte("domainsets")).Get(idBuf)

		if err := Unmarshal(data, &ret); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}
