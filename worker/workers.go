package worker

import (
	"encoding/json"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"log"
	"gitlab.com/salismt/microservice-pattern-user-service/model"
	"sync"
	"fmt"
	"gitlab.com/salismt/microservice-pattern-user-service/caches"
)

type Worker struct {
	cache *caches.Cache
	db    *sqlx.DB
	id    int
	queue string
}

func newWorker(id int, db *sqlx.DB, c *caches.Cache, queue string) Worker {
	return Worker{cache: c, db: db, id: id, queue: queue}
}

// run the queue and send data to database
func (w Worker) process(id int) {
	for {
		conn := w.cache.Pool.Get()
		var channel string
		var uuid int

		if reply, err := redigo.Values(conn.Do("BLPOP", w.queue, 30+id)); err == nil {

			if _, err := redigo.Scan(reply, &channel, &uuid); err != nil {
				w.cache.EnqueueValue(w.queue, uuid)
				continue
			}

			values, err := redigo.String(conn.Do("GET", uuid))
			if err != nil {
				w.cache.EnqueueValue(w.queue, uuid)
				continue
			}

			user := model.User{}
			if err := json.Unmarshal([]byte(values), &user); err != nil {
				w.cache.EnqueueValue(w.queue, uuid)
				continue
			}

			log.Println(user)
			if err := user.Create(w.db); err != nil {
				w.cache.EnqueueValue(w.queue, uuid)
				continue
			}

		} else if err != redigo.ErrNil {
			log.Fatal(err)
		}
		conn.Close()
	}
}

// creates number of workers for the queue to instantiate additionally
// initialize the workers asynchronously
func UsersToDB(numWorkers int, db *sqlx.DB, cache *caches.Cache, queue string) {
	cache.EnqueueValue("da", 123)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(i)
		go func(id int, db *sqlx.DB, c *caches.Cache, queue string) {
			fmt.Println(c)
			worker := newWorker(i, db, c, queue)
			worker.process(i)
			defer wg.Done()

		}(i, db, cache, queue)
	}
}