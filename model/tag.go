package model

import (
	"github.com/boltdb/bolt"
	"encoding/binary"
)

type GeoCode [2]byte
type GeoMap map[GeoCode]bool

type TagJSON struct {
	TagID  int
	Status string
	URL    string
	CPM    float64
}

type Tag struct {
	TagID     int
	Active    bool
	URL       string
	CPM       float64
	Domainset []int
	Geo       GeoMap
}

type GeoJSON struct {
	GeoID   int
	TagID   int
	Action  string
	Country string
}

func (data *TagJSON) Load() *Tag {
	v := Tag{TagID:data.TagID, URL:data.URL, CPM:data.CPM}

	v.Geo = make(GeoMap)

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
		data := tx.Bucket([]byte("tags")).Get(idBuf)

		if err := Unmarshal(data, &ret); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (tag *Tag) modifyTagGeo(db *bolt.DB, geoj *GeoJSON, doAdd bool) error {

	if err := db.Update(func(tx *bolt.Tx) error {
		idBuf := make([]byte, 8)
		binary.PutVarint(idBuf, int64(tag.TagID))
		b := tx.Bucket([]byte("tags"))
		data := b.Get(idBuf)

		var geo GeoCode
		if (geoj.Country == "*") {
			geo = [2]byte{0, 0}
		} else {
			copy(geo[:], []byte(geoj.Country))
		}

		var action bool
		if (geoj.Action == ActionAllowTxt) {
			action = true
		}else {
			action = false
		}

		var ret Tag
		if err := Unmarshal(data, &ret); err != nil {
			return err
		}

		if doAdd {
			ret.Geo[geo] = action
		}else {
			delete(ret.Geo, geo)
		}

		vBuf, err := Marshal(ret);
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

func (data *GeoJSON) Add(db *bolt.DB, t *Tag) error {

	return t.modifyTagGeo(db, data, true)

}

func (data *GeoJSON) Remove(db *bolt.DB, t *Tag) error {

	return t.modifyTagGeo(db, data, false)

}