package main


import (
    	"fmt"
    	"time"
	)

var MasterQueue = [] IpObject {}
var allreadyInQueue = false
type DataObject struct {
	NewIp string
	MasterQueue []string
}

type IpObject struct{
	Ip string
	LastUpdated int64
}

type MessageType int 
	const(
	ORDER  = 0 + iota
	COST 
	ORDER_COMPLETE
	REQUEST_AUCTION
	)


func main() {
	deadline :=time.Now().UnixNano() / int64(time.Millisecond)

	testObject := IpObject{"123.123",deadline}
	deadline2 := time.Now().UnixNano() / int64(time.Millisecond)
	diff := deadline2 - deadline
	tull := IpObject{"123.124",deadline}
	
	ko := [] IpObject {testObject,tull}
	MasterQueue = ko
	//yolo := MasterQueue[1].LastUpdated
	fmt.Println(diff)
	kjau := []int {1,2,3,4,5,6}
	per := kjau [0:3]
	
	per = append(per,kjau[1:3]...)

	fmt.Println(per)
	test()
	//time.Sleep(10000 * time.Millisecond)
	

}

func test() {
	//go test2()
	fmt.Println ("cost",ORDER)
	fmt.Println("test er ferdig")
}

func test2 () {
	i := 0
	for{
		fmt.Println("jeg klarer fortsatt å kjore",i)
		i ++
	}

}

