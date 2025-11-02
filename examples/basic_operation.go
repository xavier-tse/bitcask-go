package main

import (
	"fmt"
	"log"

	bitcask "github.com/xavier-tse/bitcask-go"
)

func main() {
	opt := bitcask.DefaultOptions
	opt.DirPath = "/tmp/bitcask-go"
	db, err := bitcask.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Put([]byte("name"), []byte("bitcask"))
	if err != nil {
		log.Fatal(err)
	}

	val, err := db.Get([]byte("name"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("val: ", string(val))

	err = db.Delete([]byte("name"))
	if err != nil {
		log.Fatal(err)
	}
}
