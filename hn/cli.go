package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/peterhellberg/hn"
)

var (
	numStories = 10
	offset     = 0

	client *hn.Client
	cache  []*hn.Item
)

func init() {
	var err error
	if len(os.Args) > 1 {
		numStories, err = strconv.Atoi(os.Args[1])
		e(err)
	}
	if len(os.Args) > 2 {
		offset, err = strconv.Atoi(os.Args[2])
		e(err)
	}
}

func main() {
	client = hn.DefaultClient

	ids, err := client.TopStories()
	e(err)

	// fetch data
	pending := numStories
	c := make(chan bool)
	cache = make([]*hn.Item, numStories, numStories)

	for i, id := range ids[offset : offset+numStories] {
		go fetch(i, id, cache, c)
	}

	for pending > 0 {
		<-c
		pending--
	}

	for i, item := range cache {
		fmt.Printf("%2d. \033[1m%s\033[0m\n    \033[2m%s\033[0m\n",
			i+offset, item.Title, item.URL)
	}
}

func fetch(i, id int, cache []*hn.Item, c chan bool) {
	var err error
	cache[i], err = client.Item(id)
	e(err)
	c <- true
}

func e(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
