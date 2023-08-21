package gcp

import (
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/api/iterator"
)

func TestAdd(t *testing.T) {
	ctx := context.Background()
	client := CreateClient(ctx)
	_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
		"first":  "Alan",
		"middle": "Mathison",
		"last":   "Turing",
		"born":   1912,
	})
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}

	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}
