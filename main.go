package main

import (
	"fmt"
	"github.com/er1c-zh/crdt-viz/view"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8080", &view.Server{}); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
