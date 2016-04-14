// ----------------------------------
// --- Generated with go:generate ---
// ---        DO NOT EDIT         ---
// ----------------------------------

package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

func (ds *DS) Init() {
	Printers["Sys"] = ds.PrintSys
	Printers["Cloud"] = ds.PrintCloud
	Printers["Model"] = ds.PrintModel
	Printers["Service"] = ds.PrintService
	Printers["Engine"] = ds.PrintEngine
}

var Buckets = []string{
	"Sys",
	"Cloud",
	"Model",
	"Service",
	"Engine",
}

type Sys struct {
	*Record
	Version uint32
}

func NewSys(id string, version uint32) *Sys {
	return &Sys{
		&Record{
			id,
			time.Now().UTC().Unix(),
			0,
		},
		version,
	}
}

type Cloud struct {
	*Record
	ApplicationID string
	Size          int
}

func NewCloud(id string, applicationID string, size int) *Cloud {
	return &Cloud{
		&Record{
			id,
			time.Now().UTC().Unix(),
			0,
		},
		applicationID,
		size,
	}
}

type Model struct {
	*Record
	CloudName    string
	CloudAddress string
	Data         []byte
}

func NewModel(id string, cloudName string, cloudAddress string, data []byte) *Model {
	return &Model{
		&Record{
			id,
			time.Now().UTC().Unix(),
			0,
		},
		cloudName,
		cloudAddress,
		data,
	}
}

type Service struct {
	*Record
	Caption     string
	Description string
	Source      string
	Target      string
	IsBuilt     bool
}

func NewService(id string, caption string, description string, source string, target string, isBuilt bool) *Service {
	return &Service{
		&Record{
			id,
			time.Now().UTC().Unix(),
			0,
		},
		caption,
		description,
		source,
		target,
		isBuilt,
	}
}

type Engine struct {
	*Record
	Name string
}

func NewEngine(id string, name string) *Engine {
	return &Engine{
		&Record{
			id,
			time.Now().UTC().Unix(),
			0,
		},
		name,
	}
}

func (ds *DS) HasSys(id string) (bool, error) {
	has := false
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sys"))
		if b == nil {
			return fmt.Errorf("Bucket Sys does not exist.")
		}

		v := b.Get([]byte(id))
		if v != nil {
			has = true
		}
		return nil
	})
	return has, err
}

func (ds *DS) HasSyss(ids []string) (bool, error) {
	has := true
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sys"))
		if b == nil {
			return fmt.Errorf("Bucket Sys does not exist.")
		}

		for _, id := range ids {
			v := b.Get([]byte(id))
			if v == nil {
				has = false
				return nil
			}
		}
		return nil
	})
	return has, err
}

func readSyss(tx *bolt.Tx, ids []string) ([]*Sys, error) {
	objs := make([]*Sys, len(ids))
	b := tx.Bucket([]byte("Sys"))
	if b == nil {
		return nil, fmt.Errorf("Bucket Sys does not exist.")
	}

	for i, id := range ids {
		v := b.Get([]byte(id))
		if v == nil {
			continue
		}

		o, err := DecodeSys(v)
		if err != nil {
			return nil, err
		}

		objs[i] = o
	}

	return objs, nil
}

func (ds *DS) ReadSyss(ids []string) ([]*Sys, error) {
	var objs []*Sys
	err := ds.db.View(func(tx *bolt.Tx) error {
		os, err := readSyss(tx, ids)
		if err != nil {
			return err
		}
		objs = os
		return nil
	})
	return objs, err
}

func (ds *DS) ReadSys(id string) (*Sys, error) {
	var obj *Sys
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sys"))
		if b == nil {
			return fmt.Errorf("Bucket Sys does not exist.")
		}

		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}

		o, err := DecodeSys(v)
		if err != nil {
			return err
		}

		obj = o

		return nil
	})
	return obj, err
}

func (ds *DS) CreateSys(o *Sys) error {
	return ds.writeSys(o, true)
}

func (ds *DS) UpdateSys(o *Sys) error {
	return ds.writeSys(o, false)
}

func (ds *DS) writeSys(o *Sys, create bool) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sys"))
		if b == nil {
			return fmt.Errorf("Bucket Sys does not exist.")
		}

		o.ModifiedAt = time.Now().UTC().Unix()

		v, err := EncodeSys(o)
		if err != nil {
			return err
		}

		err = b.Put([]byte(o.ID), v)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Sys write failed: %v", err)
	}

	return nil
}

func (ds *DS) DeleteSys(id string) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sys"))
		if b == nil {
			return fmt.Errorf("Bucket Sys does not exist.")
		}
		return b.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("Sys delete failed: %v", err)
	}
	return nil
}

func (ds *DS) ListSys() ([]*Sys, error) {
	var objs []*Sys
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sys"))
		if b == nil {
			return fmt.Errorf("Bucket Sys does not exist.")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o, err := DecodeSys(v)
			if err != nil {
				return err
			}
			objs = append(objs, o)
		}
		return nil
	})
	return objs, err
}

func EncodeSys(o *Sys) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(o)
	if err != nil {
		return nil, fmt.Errorf("Sys encode failed: %v", err)
	}
	return b.Bytes(), nil
}

func DecodeSys(b []byte) (*Sys, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	var o Sys
	err := dec.Decode(&o)
	if err != nil {
		return nil, fmt.Errorf("Sys decode failed: %v", err)
	}
	return &o, nil
}

func (ds *DS) PrintSys(b []byte) (string, error) {
	o, err := DecodeSys(b)
	if err != nil {
		return "", err
	}

	js, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func (ds *DS) HasCloud(id string) (bool, error) {
	has := false
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Cloud"))
		if b == nil {
			return fmt.Errorf("Bucket Cloud does not exist.")
		}

		v := b.Get([]byte(id))
		if v != nil {
			has = true
		}
		return nil
	})
	return has, err
}

func (ds *DS) HasClouds(ids []string) (bool, error) {
	has := true
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Cloud"))
		if b == nil {
			return fmt.Errorf("Bucket Cloud does not exist.")
		}

		for _, id := range ids {
			v := b.Get([]byte(id))
			if v == nil {
				has = false
				return nil
			}
		}
		return nil
	})
	return has, err
}

func readClouds(tx *bolt.Tx, ids []string) ([]*Cloud, error) {
	objs := make([]*Cloud, len(ids))
	b := tx.Bucket([]byte("Cloud"))
	if b == nil {
		return nil, fmt.Errorf("Bucket Cloud does not exist.")
	}

	for i, id := range ids {
		v := b.Get([]byte(id))
		if v == nil {
			continue
		}

		o, err := DecodeCloud(v)
		if err != nil {
			return nil, err
		}

		objs[i] = o
	}

	return objs, nil
}

func (ds *DS) ReadClouds(ids []string) ([]*Cloud, error) {
	var objs []*Cloud
	err := ds.db.View(func(tx *bolt.Tx) error {
		os, err := readClouds(tx, ids)
		if err != nil {
			return err
		}
		objs = os
		return nil
	})
	return objs, err
}

func (ds *DS) ReadCloud(id string) (*Cloud, error) {
	var obj *Cloud
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Cloud"))
		if b == nil {
			return fmt.Errorf("Bucket Cloud does not exist.")
		}

		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}

		o, err := DecodeCloud(v)
		if err != nil {
			return err
		}

		obj = o

		return nil
	})
	return obj, err
}

func (ds *DS) CreateCloud(o *Cloud) error {
	return ds.writeCloud(o, true)
}

func (ds *DS) UpdateCloud(o *Cloud) error {
	return ds.writeCloud(o, false)
}

func (ds *DS) writeCloud(o *Cloud, create bool) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Cloud"))
		if b == nil {
			return fmt.Errorf("Bucket Cloud does not exist.")
		}

		o.ModifiedAt = time.Now().UTC().Unix()

		v, err := EncodeCloud(o)
		if err != nil {
			return err
		}

		err = b.Put([]byte(o.ID), v)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Cloud write failed: %v", err)
	}

	return nil
}

func (ds *DS) DeleteCloud(id string) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Cloud"))
		if b == nil {
			return fmt.Errorf("Bucket Cloud does not exist.")
		}
		return b.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("Cloud delete failed: %v", err)
	}
	return nil
}

func (ds *DS) ListCloud() ([]*Cloud, error) {
	var objs []*Cloud
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Cloud"))
		if b == nil {
			return fmt.Errorf("Bucket Cloud does not exist.")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o, err := DecodeCloud(v)
			if err != nil {
				return err
			}
			objs = append(objs, o)
		}
		return nil
	})
	return objs, err
}

func EncodeCloud(o *Cloud) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(o)
	if err != nil {
		return nil, fmt.Errorf("Cloud encode failed: %v", err)
	}
	return b.Bytes(), nil
}

func DecodeCloud(b []byte) (*Cloud, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	var o Cloud
	err := dec.Decode(&o)
	if err != nil {
		return nil, fmt.Errorf("Cloud decode failed: %v", err)
	}
	return &o, nil
}

func (ds *DS) PrintCloud(b []byte) (string, error) {
	o, err := DecodeCloud(b)
	if err != nil {
		return "", err
	}

	js, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func (ds *DS) HasModel(id string) (bool, error) {
	has := false
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Model"))
		if b == nil {
			return fmt.Errorf("Bucket Model does not exist.")
		}

		v := b.Get([]byte(id))
		if v != nil {
			has = true
		}
		return nil
	})
	return has, err
}

func (ds *DS) HasModels(ids []string) (bool, error) {
	has := true
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Model"))
		if b == nil {
			return fmt.Errorf("Bucket Model does not exist.")
		}

		for _, id := range ids {
			v := b.Get([]byte(id))
			if v == nil {
				has = false
				return nil
			}
		}
		return nil
	})
	return has, err
}

func readModels(tx *bolt.Tx, ids []string) ([]*Model, error) {
	objs := make([]*Model, len(ids))
	b := tx.Bucket([]byte("Model"))
	if b == nil {
		return nil, fmt.Errorf("Bucket Model does not exist.")
	}

	for i, id := range ids {
		v := b.Get([]byte(id))
		if v == nil {
			continue
		}

		o, err := DecodeModel(v)
		if err != nil {
			return nil, err
		}

		objs[i] = o
	}

	return objs, nil
}

func (ds *DS) ReadModels(ids []string) ([]*Model, error) {
	var objs []*Model
	err := ds.db.View(func(tx *bolt.Tx) error {
		os, err := readModels(tx, ids)
		if err != nil {
			return err
		}
		objs = os
		return nil
	})
	return objs, err
}

func (ds *DS) ReadModel(id string) (*Model, error) {
	var obj *Model
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Model"))
		if b == nil {
			return fmt.Errorf("Bucket Model does not exist.")
		}

		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}

		o, err := DecodeModel(v)
		if err != nil {
			return err
		}

		obj = o

		return nil
	})
	return obj, err
}

func (ds *DS) CreateModel(o *Model) error {
	return ds.writeModel(o, true)
}

func (ds *DS) UpdateModel(o *Model) error {
	return ds.writeModel(o, false)
}

func (ds *DS) writeModel(o *Model, create bool) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Model"))
		if b == nil {
			return fmt.Errorf("Bucket Model does not exist.")
		}

		o.ModifiedAt = time.Now().UTC().Unix()

		v, err := EncodeModel(o)
		if err != nil {
			return err
		}

		err = b.Put([]byte(o.ID), v)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Model write failed: %v", err)
	}

	return nil
}

func (ds *DS) DeleteModel(id string) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Model"))
		if b == nil {
			return fmt.Errorf("Bucket Model does not exist.")
		}
		return b.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("Model delete failed: %v", err)
	}
	return nil
}

func (ds *DS) ListModel() ([]*Model, error) {
	var objs []*Model
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Model"))
		if b == nil {
			return fmt.Errorf("Bucket Model does not exist.")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o, err := DecodeModel(v)
			if err != nil {
				return err
			}
			objs = append(objs, o)
		}
		return nil
	})
	return objs, err
}

func EncodeModel(o *Model) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(o)
	if err != nil {
		return nil, fmt.Errorf("Model encode failed: %v", err)
	}
	return b.Bytes(), nil
}

func DecodeModel(b []byte) (*Model, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	var o Model
	err := dec.Decode(&o)
	if err != nil {
		return nil, fmt.Errorf("Model decode failed: %v", err)
	}
	return &o, nil
}

func (ds *DS) PrintModel(b []byte) (string, error) {
	o, err := DecodeModel(b)
	if err != nil {
		return "", err
	}

	js, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func (ds *DS) HasService(id string) (bool, error) {
	has := false
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Service"))
		if b == nil {
			return fmt.Errorf("Bucket Service does not exist.")
		}

		v := b.Get([]byte(id))
		if v != nil {
			has = true
		}
		return nil
	})
	return has, err
}

func (ds *DS) HasServices(ids []string) (bool, error) {
	has := true
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Service"))
		if b == nil {
			return fmt.Errorf("Bucket Service does not exist.")
		}

		for _, id := range ids {
			v := b.Get([]byte(id))
			if v == nil {
				has = false
				return nil
			}
		}
		return nil
	})
	return has, err
}

func readServices(tx *bolt.Tx, ids []string) ([]*Service, error) {
	objs := make([]*Service, len(ids))
	b := tx.Bucket([]byte("Service"))
	if b == nil {
		return nil, fmt.Errorf("Bucket Service does not exist.")
	}

	for i, id := range ids {
		v := b.Get([]byte(id))
		if v == nil {
			continue
		}

		o, err := DecodeService(v)
		if err != nil {
			return nil, err
		}

		objs[i] = o
	}

	return objs, nil
}

func (ds *DS) ReadServices(ids []string) ([]*Service, error) {
	var objs []*Service
	err := ds.db.View(func(tx *bolt.Tx) error {
		os, err := readServices(tx, ids)
		if err != nil {
			return err
		}
		objs = os
		return nil
	})
	return objs, err
}

func (ds *DS) ReadService(id string) (*Service, error) {
	var obj *Service
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Service"))
		if b == nil {
			return fmt.Errorf("Bucket Service does not exist.")
		}

		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}

		o, err := DecodeService(v)
		if err != nil {
			return err
		}

		obj = o

		return nil
	})
	return obj, err
}

func (ds *DS) CreateService(o *Service) error {
	return ds.writeService(o, true)
}

func (ds *DS) UpdateService(o *Service) error {
	return ds.writeService(o, false)
}

func (ds *DS) writeService(o *Service, create bool) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Service"))
		if b == nil {
			return fmt.Errorf("Bucket Service does not exist.")
		}

		o.ModifiedAt = time.Now().UTC().Unix()

		v, err := EncodeService(o)
		if err != nil {
			return err
		}

		err = b.Put([]byte(o.ID), v)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Service write failed: %v", err)
	}

	return nil
}

func (ds *DS) DeleteService(id string) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Service"))
		if b == nil {
			return fmt.Errorf("Bucket Service does not exist.")
		}
		return b.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("Service delete failed: %v", err)
	}
	return nil
}

func (ds *DS) ListService() ([]*Service, error) {
	var objs []*Service
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Service"))
		if b == nil {
			return fmt.Errorf("Bucket Service does not exist.")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o, err := DecodeService(v)
			if err != nil {
				return err
			}
			objs = append(objs, o)
		}
		return nil
	})
	return objs, err
}

func EncodeService(o *Service) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(o)
	if err != nil {
		return nil, fmt.Errorf("Service encode failed: %v", err)
	}
	return b.Bytes(), nil
}

func DecodeService(b []byte) (*Service, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	var o Service
	err := dec.Decode(&o)
	if err != nil {
		return nil, fmt.Errorf("Service decode failed: %v", err)
	}
	return &o, nil
}

func (ds *DS) PrintService(b []byte) (string, error) {
	o, err := DecodeService(b)
	if err != nil {
		return "", err
	}

	js, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func (ds *DS) HasEngine(id string) (bool, error) {
	has := false
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Engine"))
		if b == nil {
			return fmt.Errorf("Bucket Engine does not exist.")
		}

		v := b.Get([]byte(id))
		if v != nil {
			has = true
		}
		return nil
	})
	return has, err
}

func (ds *DS) HasEngines(ids []string) (bool, error) {
	has := true
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Engine"))
		if b == nil {
			return fmt.Errorf("Bucket Engine does not exist.")
		}

		for _, id := range ids {
			v := b.Get([]byte(id))
			if v == nil {
				has = false
				return nil
			}
		}
		return nil
	})
	return has, err
}

func readEngines(tx *bolt.Tx, ids []string) ([]*Engine, error) {
	objs := make([]*Engine, len(ids))
	b := tx.Bucket([]byte("Engine"))
	if b == nil {
		return nil, fmt.Errorf("Bucket Engine does not exist.")
	}

	for i, id := range ids {
		v := b.Get([]byte(id))
		if v == nil {
			continue
		}

		o, err := DecodeEngine(v)
		if err != nil {
			return nil, err
		}

		objs[i] = o
	}

	return objs, nil
}

func (ds *DS) ReadEngines(ids []string) ([]*Engine, error) {
	var objs []*Engine
	err := ds.db.View(func(tx *bolt.Tx) error {
		os, err := readEngines(tx, ids)
		if err != nil {
			return err
		}
		objs = os
		return nil
	})
	return objs, err
}

func (ds *DS) ReadEngine(id string) (*Engine, error) {
	var obj *Engine
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Engine"))
		if b == nil {
			return fmt.Errorf("Bucket Engine does not exist.")
		}

		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}

		o, err := DecodeEngine(v)
		if err != nil {
			return err
		}

		obj = o

		return nil
	})
	return obj, err
}

func (ds *DS) CreateEngine(o *Engine) error {
	return ds.writeEngine(o, true)
}

func (ds *DS) UpdateEngine(o *Engine) error {
	return ds.writeEngine(o, false)
}

func (ds *DS) writeEngine(o *Engine, create bool) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Engine"))
		if b == nil {
			return fmt.Errorf("Bucket Engine does not exist.")
		}

		o.ModifiedAt = time.Now().UTC().Unix()

		v, err := EncodeEngine(o)
		if err != nil {
			return err
		}

		err = b.Put([]byte(o.ID), v)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Engine write failed: %v", err)
	}

	return nil
}

func (ds *DS) DeleteEngine(id string) error {
	err := ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Engine"))
		if b == nil {
			return fmt.Errorf("Bucket Engine does not exist.")
		}
		return b.Delete([]byte(id))
	})
	if err != nil {
		return fmt.Errorf("Engine delete failed: %v", err)
	}
	return nil
}

func (ds *DS) ListEngine() ([]*Engine, error) {
	var objs []*Engine
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Engine"))
		if b == nil {
			return fmt.Errorf("Bucket Engine does not exist.")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o, err := DecodeEngine(v)
			if err != nil {
				return err
			}
			objs = append(objs, o)
		}
		return nil
	})
	return objs, err
}

func EncodeEngine(o *Engine) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(o)
	if err != nil {
		return nil, fmt.Errorf("Engine encode failed: %v", err)
	}
	return b.Bytes(), nil
}

func DecodeEngine(b []byte) (*Engine, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	var o Engine
	err := dec.Decode(&o)
	if err != nil {
		return nil, fmt.Errorf("Engine decode failed: %v", err)
	}
	return &o, nil
}

func (ds *DS) PrintEngine(b []byte) (string, error) {
	o, err := DecodeEngine(b)
	if err != nil {
		return "", err
	}

	js, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}
