package orchestrator

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/badger/v3"
)

//type Pod struct {
//	ID      string `json:"id"`
//	Service string `json:"service"`
//	Status  string `json:"status"` // Pending | Running
//	IP      string `json:"ip"`
//	Image   string `json:"image"`
//	Port    int    `json:"port"`
//}

type Datastore struct {
	db *badger.DB
}

func NewDatastore(path string) *Datastore {
	opts := badger.DefaultOptions(path).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("failed to open BadgerDB: %v", err)
	}
	return &Datastore{db: db}
}

func (d *Datastore) Close() {
	d.db.Close()
}

func (d *Datastore) PutPod(ctx context.Context, pod Pod) error {
	data, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := PodKey(pod.Service, pod.ID)
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (d *Datastore) GetPod(ctx context.Context, id string) (*Pod, error) {
	var pod Pod
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &pod)
		})
	})
	if err != nil {
		return nil, err
	}
	return &pod, nil
}

func (d *Datastore) ListPods(ctx context.Context, service string) ([]Pod, error) {
	var pods []Pod

	prefix := []byte("pods/")
	if service != "" {
		prefix = []byte("pods/" + service + "/")
	}

	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var pod Pod
				if err := json.Unmarshal(val, &pod); err != nil {
					return err
				}
				pods = append(pods, pod)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return pods, err
}

func (d *Datastore) DeletePod(ctx context.Context, service, id string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(PodKey(service, id)))
	})
}
