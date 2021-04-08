package main

import (
	"fmt"
	"github.com/er1c-zh/crdt-viz/CvRDT"
	"github.com/er1c-zh/crdt-viz/impl/g_counter"
	"github.com/er1c-zh/crdt-viz/network"
	"math/rand"
	"time"
)

func main() {
	cloud := network.NewCloud(network.DefaultOption())
	obj := CvRDT.NewObject(g_counter.Factory)
	obj.Init(cloud, CvRDT.Option{Count: 3})
	fmt.Printf("%v\n", obj.Query(0))
	obj.Update(0, CvRDT.UpdateArgs{})
	fmt.Printf("%v\n", obj.Query(0))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i ++ {
		obj.Update(rand.Intn(3), CvRDT.UpdateArgs{})
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("%v\n", obj.State())

	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %v\n", i, obj.Query(i))
	}
}
