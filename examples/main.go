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
	
	uid := wuid.NewWUID("test")
	fmt.Printf("%#016x\n", uid.Next())
	fmt.Printf("%#016x\n", uid.Next())
	fmt.Printf("%#016x\n", uid.Next())
	fmt.Printf("%#016x\n", uid.Next())
}
