package main

import (
	"fmt"

	"github.com/oarkflow/xid"
	"github.com/oarkflow/xid/wuid"
)

func main() {
	fmt.Println(xid.New().String())
	fmt.Println(xid.New().String())
	fmt.Println(xid.New().String())
	fmt.Println(xid.New().String())

	fmt.Printf("%d\n", wuid.New())
	fmt.Printf("%d\n", wuid.New())
	fmt.Printf("%d\n", wuid.New())
	fmt.Printf("%d\n", wuid.New())
}
