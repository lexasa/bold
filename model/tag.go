
package model

import (
	"github.com/boltdb/bolt"
	"encoding/binary"
)

type TagJSON struct {
	TagID  int
	Status string
	URL    string
	CPM    float64
}

type Tag struct {
	TagID  int
	Active bool
	URL    string
	CPM    float64
}

func (data *TagJSON) Load() *Tag {
	v := Tag{TagID:data.TagID, URL:data.URL, CPM:data.CPM}
	if data.Status == StatusEnabledTxt {
		v.Active = StatusEnabled
	}else {
		v.Active = StatusDisabled
	}
	return &v
}

func (tag *Tag) Store(db *bolt.DB) error {

	if err := db.Update(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)
		binary.PutVarint(idBuf, int64(tag.TagID))

		b, err := tx.CreateBucketIfNotExists([]byte("tags"))
		if err != nil {
			return err
		}

		vBuf, err := Marshal(tag);
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

func GetTagByID(db *bolt.DB, id int) (*Tag, error) {
	var ret Tag
	if err := db.View(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)

		binary.PutVarint(idBuf, int64(id))
		data:=tx.Bucket([]byte("tags")).Get(idBuf)

		if err := Unmarshal(data, &ret); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}
