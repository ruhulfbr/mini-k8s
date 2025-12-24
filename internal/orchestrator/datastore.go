package orchestrator

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/dgraph-io/badger/v3"
)

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

func (d *Datastore) ListPodsByService(ctx context.Context, service string) ([]Pod, error) {
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

func (d *Datastore) ListRunningPods(ctx context.Context, service string) ([]Pod, error) {
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

				if pod.Status == PodRunning {
					pods = append(pods, pod)
				}

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

func (d *Datastore) ListPodsGroupedByService(ctx context.Context) (map[string][]Pod, error) {
	podsByService := make(map[string][]Pod)

	prefix := []byte("pods/")

	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			key := string(item.Key())
			parts := strings.Split(key, "/")
			if len(parts) < 3 {
				continue
			}

			service := parts[1]

			err := item.Value(func(val []byte) error {
				var pod Pod
				if err := json.Unmarshal(val, &pod); err != nil {
					return err
				}

				podsByService[service] = append(podsByService[service], pod)
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	return podsByService, err
}

func (d *Datastore) DeletePod(ctx context.Context, service, id string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(PodKey(service, id)))
	})
}

func (d *Datastore) DeleteAllPods(ctx context.Context) error {
	prefix := []byte("pods/")

	return d.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if err := txn.Delete(it.Item().Key()); err != nil {
				return err
			}
		}
		return nil
	})
}
