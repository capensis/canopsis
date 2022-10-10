package main

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	var err error

	fmt.Printf("%v\n", errors.Is(err, mongo.ErrNoDocuments))
}
