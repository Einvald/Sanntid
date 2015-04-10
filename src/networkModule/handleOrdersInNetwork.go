package networkModule

import (
    	"fmt"
    	"net"
    	"os"
    	"time"

   )




type OrderData struct {
	FromMaster bool
	Type MessageType
	Cost int 
	Order int //buttonOrder
	OrderComplete bool // trenger vell ikke denne?
	Ip string
}

type MessageType int 
	const(
	ORDER  = 0 + iota
	COST 
	ORDER_COMPLETE
	REQUEST_AUCTION
	)


//Obs!! Hvis master dør må backup også overta unfinishedOrderskøen

var unfinishedOrders = [] OrderData{}
var recievedMessageToMaster = make(chan OrderData, 1024)
var recievedMessage = make (chan OrderData,1024)
var AuctionResultChan = make(chan OrderData,1)
var recievedCostChan = make(chan OrderData,1024)
var recievedOrderChan =make(chan OrderData,1024)  
var recievedOrderComplete = make(chan OrderData, 1024) 
var auctionLock = make(chan int, 1)
var unfinishedOrdersLock = make(chan int, 1)
func handleOrdersInNetwork(){
	const IP_BROADCAST = "129.241.187.255"
	const BROADCAST_PORT = "20012" 
	portSomSlaverLeserFra := "20013" // Her må vi gjøre endringer
	unfinishedOrdersLock <- 1
	auctionLock <- 1
	for {
		select{
		case recievedData := <- recievedMessageToMaster:
			switch recievedData.Type {
				case COST:
					recievedCostChan <- recievedData
				case ORDER:
					recievedOrderChan <- recievedData 
				case ORDER_COMPLETE:
					recievedOrderComplete <- recievedData
			}
		case order := <- recievedOrderChan:
			if !isInQueue(order){
				go auction(order,IP_BROADCAST,BROADCAST_PORT,portSomSlaverLeserFra)
				
			}
			
			// Sjekk om finnes i liste. Legg til hvis ikke
		case orderComplete := <- recievedOrderComplete:
			removeFinishedOrder(orderComplete)
			// Sett bit på channel
		}
	}

}

func auction(newOrderData OrderData, IP_BROADCAST string, BROADCAST_PORT string, portSomSlaverLeserFra string){
	<- auctionLock
	deadline := time.Now().UnixNano() / int64(time.Millisecond) + 500
	elevatorsInAuction := [] OrderData {} 
	//Sett deadline
	newOrderData.Type = 3  
	sendOrderData(IP_BROADCAST,BROADCAST_PORT,newOrderData)  
	for {
		select{
		case cost := <- recievedCostChan:
			allreadyInList := false
			for _,element:= range elevatorsInAuction{
				if element == cost{			
					allreadyInList = true
					break
				}
			}
			if !allreadyInList{
				elevatorsInAuction = append(elevatorsInAuction,cost)
			}
		}
		timeNow := time.Now().UnixNano() / int64(time.Millisecond)
		if len(elevatorsInAuction) == len(masterQueue) || timeNow > deadline  {	
			break
		}
	}
	lowestCost := 100000000  //Denne må kanskje settes høyere?
	elevatorWithLowestCost := elevatorsInAuction[0]
	for _,element:= range elevatorsInAuction{
		if element.Cost < lowestCost{
			lowestCost = element.Cost
			elevatorWithLowestCost = element
		}

	}
	addNewOrder(elevatorWithLowestCost)
	elevatorWithLowestCost.Type = 0
	sendOrderData(elevatorWithLowestCost.Ip,portSomSlaverLeserFra,elevatorWithLowestCost) // Porten her må bestemmes eksternt!!
	auctionLock <- 1	
}

func handleOrdersFromMaster(){ //Her må vi ta en eller to channels som input
	for {
		select{
		case recievedData := <- recievedMessage:
			switch recievedData.Type{
				case ORDER:
					//Legges på channel som leses av main

				case REQUEST_AUCTION:
					// Legges på en annen channel som leses av main
			}
			
	}	}


}

func isInQueue(newOrderData OrderData) bool {
	newOrder := newOrderData.Order
	<- unfinishedOrdersLock
	allreadyInList := false
	for _,element:= range unfinishedOrders{
		if element.Order == newOrder{			
			allreadyInList = true
			break
		}
	}
	return allreadyInList

}
//Funksjoner master bruker	
func addNewOrder(newOrderData OrderData){
	<- unfinishedOrdersLock 
	allreadyInList := isInQueue (newOrderData)
	if !allreadyInList {
		unfinishedOrders = append(unfinishedOrders,newOrderData)
	}
	unfinishedOrdersLock <- 1
}



func removeFinishedOrder(orderComplete OrderData){
	order2Remove := orderComplete.Order
	<- unfinishedOrdersLock
	for i,element:= range unfinishedOrders{
		if element.Order == order2Remove{
			n := len (unfinishedOrders)			
			newUnfinishedOrders :=unfinishedOrders[0:i]
			newUnfinishedOrders = append(newUnfinishedOrders,unfinishedOrders[i+1:n]...)
			unfinishedOrders = newUnfinishedOrders
			break
		}
	unfinishedOrdersLock <- 1
	}

}

// Lag funksjonalitet for døde heiser, hvordan dens bestillinger skal fordeles osv.


//Funksjoner som alle bruker
func readOrderData(port string){  			
	bufferToRead := make([] byte, 1024)
	
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+port)
	
	if err != nil {
        fmt.Println("error resolving UDP address on ", port)
        fmt.Println(err)
        os.Exit(1)
    }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
        fmt.Println("error listening on UDP port ", port)
        fmt.Println(err)
        os.Exit(1)
	}
	for {
		n,_, err := readerSocket.ReadFromUDP(bufferToRead)  
	 	if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            os.Exit(1)     
        }
        
       	if n > 0 { 
            structObject := json2struct(bufferToRead,n)
            data := structObject.OrderMessage
            if data.FromMaster{
            	recievedMessage <- data
            
            }else if !data.FromMaster{
            	recievedMessageToMaster <- data
            } 
            
        }

   	}
}
	//Hvis til Master legges det på masterkanalen, hvis ikke på den andre kanalen


func sendOrderData(ip string, port string, message OrderData){	//Denne brukes både til å sende enkle meldinger, men også broadcasting
	udpAddr, err := net.ResolveUDPAddr("udp",ip+":"+port)
	if err != nil {
		fmt.Println("error resolving UDP address on ", port)
		fmt.Println(err)
		os.Exit(1)
	}
	broadcastSocket, err := net.DialUDP("udp",nil, udpAddr)
	if err != nil {
		    fmt.Println("error listening on UDP port ", port)
		    fmt.Println(err)
		    os.Exit(1)
	}
	deadline := time.Now().UnixNano() / int64(time.Millisecond) + 5
	for {
		timeNow := time.Now().UnixNano() / int64(time.Millisecond)
		if timeNow > deadline {
			break
		}
		sendingObject := DataObject{"",[]IpObject {},message}
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)
		time.Sleep(1 * time.Millisecond)	

	}
	// Vi må sende flere ganger og hvis en bestilling ikke mottas av master skal heisen ta den selv.
	
}



//Spør om lys skal lyse i alle på alle heispaneler hvis man bestemmer for eksempel ned fra 2. etasje?
