package main

import (
	"fmt"

	"github.com/oarkflow/xid"
)

func main() {
	fmt.Println(xid.New().String())
	fmt.Println(xid.New().String())
	fmt.Println(xid.New().String())
	fmt.Println(xid.New().String())
}
