package main

import(
	"fmt"
	d "driver"
	elev "elev_handler"
	//"time"
	)

func main() {
	i := d.D_init()
	fmt.Println(i)
	elev.DoorTimer()
	fmt.Println("Timer done")
	fmt.Println(elev.STOP)
	
}
	
func checkForInput(){

}

func readFromNetwork() {
	
}