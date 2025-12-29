package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type PodRepository struct {
	db *badger.DB
}

func NewPodRepository(db *badger.DB) *PodRepository {
	return &PodRepository{db: db}
}

func (pr *PodRepository) PutPod(ctx context.Context, pod entities.Pod) error {
	data, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := podKey(pod.Service, pod.ID)
	return pr.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (pr *PodRepository) GetPod(ctx context.Context, id string) (*entities.Pod, error) {
	var pod entities.Pod
	err := pr.db.View(func(txn *badger.Txn) error {
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

func (pr *PodRepository) ListPodsByService(ctx context.Context, service string) ([]entities.Pod, error) {
	var pods []entities.Pod

	prefix := []byte("pods/")
	if service != "" {
		prefix = []byte("pods/" + service + "/")
	}

	err := pr.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var pod entities.Pod
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

func (pr *PodRepository) ListRunningPods(ctx context.Context, service string) ([]entities.Pod, error) {
	var pods []entities.Pod

	prefix := []byte("pods/")
	if service != "" {
		prefix = []byte("pods/" + service + "/")
	}

	err := pr.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var pod entities.Pod
				if err := json.Unmarshal(val, &pod); err != nil {
					return err
				}

				if pod.Status == entities.PodRunning {
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

func (pr *PodRepository) ListPodsGroupedByService(ctx context.Context) ([]entities.Pod, error) {
	var result []entities.Pod
	seenServices := make(map[string]bool)

	prefix := []byte("pods/")

	err := pr.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			// key format: pods/{service}/{podID}
			key := string(item.Key())
			parts := strings.Split(key, "/")
			if len(parts) < 3 {
				continue
			}

			service := parts[1]

			if seenServices[service] {
				continue
			}

			err := item.Value(func(val []byte) error {
				var pod entities.Pod
				if err := json.Unmarshal(val, &pod); err != nil {
					return err
				}

				result = append(result, pod)
				seenServices[service] = true
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err
}

func (pr *PodRepository) DeletePod(ctx context.Context, service, id string) error {
	return pr.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(podKey(service, id)))
	})
}

func (pr *PodRepository) DeleteAllPods(ctx context.Context) error {
	prefix := []byte("pods/")

	return pr.db.Update(func(txn *badger.Txn) error {
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

/*
PodKey design:

	pods/<service>/<podID>
	This allows:
	- List pods by service
	- Efficient deletes
*/
func podKey(service string, id string) string {
	return fmt.Sprintf("pods/%s/%s", service, id)
}
