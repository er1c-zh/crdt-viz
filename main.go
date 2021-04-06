package main

import (
	"fmt"
	"github.com/er1c-zh/crdt-viz/CvRDT"
	"github.com/er1c-zh/crdt-viz/impl/g_counter"
	"math/rand"
	"time"
)

func main() {
	obj := CvRDT.NewObject(g_counter.Factory)
	obj.Init(CvRDT.Option{Count: 3})
	fmt.Printf("%v\n", obj.Query(0))
	obj.Update(0, CvRDT.UpdateArgs{})
	fmt.Printf("%v\n", obj.Query(0))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i ++ {
		obj.Update(rand.Intn(3), CvRDT.UpdateArgs{})
	}
	fmt.Printf("%v\n", obj.State())
}
