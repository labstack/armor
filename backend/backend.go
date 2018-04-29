package main

import (
	"fmt"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/labstack/gommon/log"
)

func main() {
	Start()
}

func Start() {
	// Register
	consul.Register()
	etcd.Register()
	zookeeper.Register()

	client := "localhost:8500"

	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.CONSUL, // or "consul"
		[]string{client},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("Cannot create store consul")
	}

	key := "armor"
	err = kv.Put(key, []byte("bar"), &store.WriteOptions{IsDir: true})
	if err != nil {
		fmt.Errorf("Something went wrong when initializing key: %v", key)
	}

	stopCh := make(chan struct{})
	events, err := kv.WatchTree(key, stopCh)

	for {
		select {
		case pairs := <-events:
			// Do something with events
			for _, pair := range pairs {
				fmt.Printf("value changed on key %v: new value=%s\n", key, pair.Value)
			}
		}
	}

	// key := "foo"
	// err = kv.Put(key, []byte("bar"), nil)
	// if err != nil {
	// 	fmt.Errorf("Error trying to put value at key: %v", key)
	// }

	// pair, err := kv.Get(key)
	// if err != nil {
	// 	fmt.Errorf("Error trying accessing value at key: %v", key)
	// }

	// err = kv.Delete(key)
	// if err != nil {
	// 	fmt.Errorf("Error trying to delete key %v", key)
	// }

	// log.Info("value: ", string(pair.Value))
}
