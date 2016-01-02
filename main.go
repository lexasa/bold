package main

import (
	. "./model"
	"log"
	"github.com/boltdb/bolt"
	"time"
)

func main() {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // tag
	t := TagJSON{
		TagID:1,
		Status:"enabled",
		URL:"http://tag1url.com",
		CPM:6.5,
	}

	vt := t.Load()

	start := time.Now()

	for i := 1; i < 10; i++ {
		vt.TagID++
		vt.Store(db)
	}

	elapsed := time.Since(start)
	log.Print(elapsed / 10)

	zt, err := GetTagByID(db, vt.TagID)

	log.Print(vt)
	log.Print(zt,err)

	// account

	a := AccountJSON{
		AccountID: 1,
		Status:"enabled",
		Tag:"acc1",
		CPM:5.5,
	}

	va := a.Load()

	start = time.Now()

	for i := 1; i < 10; i++ {
		va.AccountID++
		va.Store(db)
	}

	elapsed = time.Since(start)
	log.Print(elapsed / 10)

	za, err := GetAccountByID(db, va.AccountID)
	log.Print(va)
	log.Print(za,err)

	za, err = GetAccountByTag(db, va.Tag)
	log.Print(va)
	log.Print(za,err)

	// domainset

	ds:= DomainsetJSON{
		DomainsetID:2,
		Name:"All",
		Status:"enabled",
		Action:"allow",
	}

	vds := ds.Load()

	start = time.Now()

	for i := 1; i < 10; i++ {
		vds.DomainsetID++
		vds.Store(db)
	}

	elapsed = time.Since(start)
	log.Print(elapsed / 10)

	zds, err := GetDomainsetByID(db, vds.DomainsetID)
	log.Print(vds)
	log.Print(zds,err)

}