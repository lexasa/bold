package main

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/json"
	"fmt"
	"encoding/binary"
	"encoding/gob"
	"bytes"
//"github.com/google/flatbuffers/go"
)



type Account struct {
	AccountID int
	Status    string
	Tag       string
	CPM       float64
}

type Cap struct {
	CapID             int
	TagID             int
	AccountID         int
	CapValue          int
	TimeUnit          string
	TimeValue         int
	TimeInterval      int
	SessionDefinition string
}

type Permission struct {
	PermissionID int
	AccountID    int
	TagID        int
	Action       string
}

type Domainset struct {
	DomainsetID int
	Name        string
	Status      string
	Action      string
}

type Domain struct {
	DomainID    int
	DomainsetID int
	Domain      string
}

type TagDomainset struct {
	TagDomainsetID int
	TagID          int
	DomainsetID    int
}

type Geo struct {
	GeoID   int
	TagID   int
	Action  string
	Country string
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	obj := [](interface{}){}

	obj = append(obj,
		&TagJSON{
			TagID:1,
			Status:"enabled",
			URL:"http://tag1url.com",
			CPM:6.5,
		})

	obj = append(obj,
		&TagJSON{
			TagID:2,
			Status:"enabled",
			URL:"http://tag2url.com",
			CPM:7.5,
		})

	obj = append(obj,
		&Domainset{
			DomainsetID:1,
			Name:"All",
			Status:"enabled",
			Action:"block",
		})

	obj = append(obj,
		&Domainset{
			DomainsetID:2,
			Name:"All",
			Status:"enabled",
			Action:"allow",
		})

	obj = append(obj,
		&Domain{
			DomainID: 1,
			DomainsetID:1,
			Domain:"*",
		})

	obj = append(obj,
		&Domain{
			DomainID: 2,
			DomainsetID:2,
			Domain:"metacafe.com",
		})

	obj = append(obj,
		&TagDomainset{
			TagDomainsetID: 1,
			TagID:1,
			DomainsetID:1,
		})

	obj = append(obj,
		&TagDomainset{
			TagDomainsetID: 2,
			TagID:1,
			DomainsetID:2,
		})

	obj = append(obj,
		&TagDomainset{
			TagDomainsetID: 1,
			TagID:2,
			DomainsetID:1,
		})

	obj = append(obj,
		&TagDomainset{
			TagDomainsetID: 2,
			TagID:2,
			DomainsetID:2,
		})

	obj = append(obj,
		&Account{
			AccountID: 1,
			Status:"enabled",
			Tag:"acc1",
			CPM:5.5,
		})

	obj = append(obj,
		&Account{
			AccountID: 2,
			Status:"enabled",
			Tag:"acc2",
			CPM:4.5,
		})

	obj = append(obj,
		&Account{
			AccountID: 2,
			Status:"enabled",
			Tag:"acc3",
			CPM:4.5,
		})

	obj = append(obj,
		&Permission{
			AccountID: 1,
			TagID:1,
			Action:"allow",
		})

	obj = append(obj,
		&Permission{
			AccountID: 2,
			TagID:2,
			Action:"allow",
		})

	obj = append(obj,
		&Permission{
			AccountID: 3,
			TagID:1,
			Action:"allow",
		})

	obj = append(obj,
		&Permission{
			AccountID: 3,
			TagID:2,
			Action:"allow",
		})

	obj = append(obj,
		&Cap{
			CapID             :1,
			TagID             :1,
			AccountID         :0,
			CapValue          :1000,
			TimeUnit          :"all",
			TimeValue         :1,
			TimeInterval      :0,
			SessionDefinition :"",
		})

	obj = append(obj,
		&Cap{
			CapID             :2,
			TagID             :1,
			AccountID         :0,
			CapValue          :10,
			TimeUnit          :"hour",
			TimeValue         :1,
			TimeInterval      :1,
			SessionDefinition :"ip,session",
		})


	for _, o := range obj {
		st, _ := json.Marshal(o)
		switch o.(type) {
		case *TagJSON:



			o.(*Tag).Store(db)
			fmt.Println(string(st))
			break
		}
	}

	defer db.Close()
}


func (tag *Tag) Store(db *bolt.DB) error {

	if err := db.Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		idBuf := make([]byte, 8)
		binary.PutVarint(idBuf, int64(tag.TagID))

		b, err := tx.CreateBucketIfNotExists([]byte("tags"))
		if err != nil {
			return err
		}

		tag_s := tag.Storage()

		vBuf, err := Marshal(tag_s);
		if err != nil {
			return err
		}

		if err := b.Put(idBuf, vBuf); err != nil {
			return err
		}

		// Retrieve the key back from the database and verify it.
		value := b.Get(idBuf)

		v := TagBin{}
		err = Unmarshal(value, &v)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("The value of '%v' was: %v\n", idBuf, v)

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (db *bolt.DB) GetTag(id int) (*TagBin)
