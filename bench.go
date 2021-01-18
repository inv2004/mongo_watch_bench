package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Key                                    string
	Time                                   time.Time
	I0, I1, I2, I3, I4, I5, I6, I7, I8, I9 int64
	F0, F1, F2, F3, F4, F5, F6, F7, F8, F9 float64
}

const DBName = "testbench"
const DBColl = "test1"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func sender(wg *sync.WaitGroup, num int, init int) {
	defer wg.Done()

	client, err := mongo.NewClient(options.Client())
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		log.Panic(err)
	}

	defer client.Disconnect(ctx)

	coll := client.Database(DBName).Collection(DBColl)

	ticker := time.NewTicker(1000 * time.Millisecond)

	for range ticker.C {
		go func() {
			post := Post{}

			start := time.Now()

			for i := 0; i < init; i++ {
				post.Key = randSeq(6)
				post.Time = time.Now()

				_, err = coll.InsertOne(context.TODO(), post)
			}
			fmt.Printf("[sender-%d] %d in %s\n", num, init, time.Since(start))
		}()
	}

}

func reader(wg *sync.WaitGroup, num int) {
	defer wg.Done()

	client, err := mongo.NewClient(options.Client())
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		log.Panic(err)
	}

	defer client.Disconnect(ctx)

	coll := client.Database(DBName).Collection(DBColl)

	watch, err := coll.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		log.Panic(err)
	}
	defer watch.Close(context.TODO())

	watchCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// post := Post{}
	var data struct {
		FullDocument Post
	}

	maxDuration := time.Duration(0)

	go func(maxDuration *time.Duration) {
		for range time.NewTicker(1 * time.Second).C {
			fmt.Printf("[reader-%d] %s\n", num, maxDuration)
			*maxDuration = 0
		}
	}(&maxDuration)

	for watch.Next(watchCtx) {
		if err := watch.Decode(&data); err != nil {
			log.Panic(err)
		}
		diff := time.Since(data.FullDocument.Time)
		if diff > maxDuration {
			maxDuration = diff
		}
	}
}

func main() {
	rows := flag.Int("rows", 100, "number of rows to insert per sender")
	senders := flag.Int("s", 1, "amount of sender processes")
	readers := flag.Int("r", 1, "amount of readers processes (read all data without filtering)")

	flag.Parse()

	if flag.NFlag() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("rows = %d, senders = %d, readers = %d\n", *rows, *senders, *readers)

	var wg sync.WaitGroup

	for i := 0; i < *senders; i++ {
		wg.Add(1)
		go sender(&wg, i, *rows)
	}
	for i := 0; i < *readers; i++ {
		wg.Add(1)
		go reader(&wg, i)
	}

	wg.Wait()
}
