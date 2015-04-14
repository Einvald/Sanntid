package main

import (
		network "networkModule"
		
    	//"fmt"
    	//"net"
    	//"os"
	)


	



func main(){
	Order_data_to_master_chan := make(chan network.OrderData,1024)
	Order_data_from_master_chan := make (chan network.OrderData,1024)
	network.InitializeElevator()
	network.RunElevator(Order_data_to_master_chan,Order_data_from_master_chan)

}



