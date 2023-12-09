package main

import (
	"context"
	"custom_db/api"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Test")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx := context.Background()
	go api.NewAPIRunner().Go(ctx, wg)
	wg.Wait()
}
