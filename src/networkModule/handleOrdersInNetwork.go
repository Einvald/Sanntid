package networkModule

import (
    	"net"
    	"time"
   )

type OrderData struct {
	FromMaster bool
	Type MessageType
	Order ButtonOrder
	Cost int 
	Ip string
	Deadline int64	
}

type MessageType int 
	const(
	ORDER  = 0 + iota
	COST 
	ORDER_COMPLETE
	REQUEST_AUCTION
	)

type ButtonOrder struct {
	Floor int
	ButtonType int
}

const IP_BROADCAST = "129.241.187.255"
const BROADCAST_PORT = "31012"
const SLAVE_TO_MASTER_PORT = "31015"
const MASTER_TO_SLAVE_PORT = "31014"
const COST_CEILING = 100 

func handleOrdersInNetwork(isMasterChan chan bool, masterQueueChan chan []IpObject, unfinishedOrdersChan chan []OrderData, recievedMessage chan OrderData, recievedMessageToMaster chan OrderData, recievedOrderChan chan OrderData){
	unactiveElevatorChan := make(chan string, 1)
	recievedCostChan := make(chan OrderData,1024)  
	recievedOrderComplete := make(chan OrderData, 1024) 
	auctionLock := make(chan int, 1)
	select{
		case auctionLock <- 1:
		case <-time.After(100 * time.Millisecond):
	}
	go checkOrderDeadline(unactiveElevatorChan, isMasterChan, unfinishedOrdersChan, recievedOrderChan)
	go readOrderData(SLAVE_TO_MASTER_PORT, isMasterChan,recievedMessageToMaster,recievedMessage)
	isMaster := false
	for {
		isMaster = <- isMasterChan; isMasterChan <- isMaster
		if !isMaster{break}
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
			if !isInQueue(order,unfinishedOrdersChan){
				go auction(order, unactiveElevatorChan,masterQueueChan,unfinishedOrdersChan,recievedCostChan, auctionLock) 
			}
		case orderComplete := <- recievedOrderComplete:
			orderComplete.FromMaster = true 
			removeOrder(orderComplete,unfinishedOrdersChan)
			sendOrderData(IP_BROADCAST,MASTER_TO_SLAVE_PORT,orderComplete)
		}	
	}
}

func auction(newOrderData OrderData, unactiveElevatorChan chan string,masterQueueChan chan []IpObject, unfinishedOrdersChan chan [] OrderData, recievedCostChan chan OrderData, auctionLock chan int){
	<- auctionLock
	if !isInQueue(newOrderData, unfinishedOrdersChan){
		deadline := time.Now().UnixNano() / int64(time.Millisecond) + 800
		elevatorsInAuction := [] OrderData {} 
		newOrderData.Type = REQUEST_AUCTION
		newOrderData.FromMaster = true  
		sendOrderData(IP_BROADCAST,BROADCAST_PORT,newOrderData)  
		for {
			select{
			case cost := <- recievedCostChan:
				allreadyInList := false
				for _,element:= range elevatorsInAuction{
					if element == cost {			
						allreadyInList = true
						break
					}
				}
				if !allreadyInList && cost.Order == newOrderData.Order {
					elevatorsInAuction = append(elevatorsInAuction,cost)
				}
			}
			timeNow := time.Now().UnixNano() / int64(time.Millisecond)
			masterQueue := <- masterQueueChan; masterQueueChan <- masterQueue
			if len(elevatorsInAuction) == len(masterQueue) || timeNow > deadline  {
				break
			}
		}
		ip := "nil"
		select {
			case ip = <- unactiveElevatorChan:
			case <-time.After(50 * time.Millisecond):
		}
		lowestCost := COST_CEILING
		elevatorWithLowestCost := newOrderData
		elevatorWithLowestCost.Cost = COST_CEILING
		
		for _,element:= range elevatorsInAuction{

			if element.Cost < lowestCost && element.Ip != ip {
				lowestCost = element.Cost
				elevatorWithLowestCost = element
			}
		}
		addNewOrder(elevatorWithLowestCost,unfinishedOrdersChan)
		elevatorWithLowestCost.Type = ORDER
		elevatorWithLowestCost.FromMaster = true
		sendOrderData(elevatorWithLowestCost.Ip,MASTER_TO_SLAVE_PORT,elevatorWithLowestCost) 
	}
	auctionLock <- 1	
}

func handleOrdersFromMaster(orderDataFromMasterChan chan OrderData, isMasterChan chan bool,recievedMessage chan OrderData, recievedMessageToMaster chan OrderData,myIp string){
	go readOrderData(MASTER_TO_SLAVE_PORT,isMasterChan,recievedMessageToMaster,recievedMessage)	
	go readOrderData(BROADCAST_PORT,isMasterChan,recievedMessageToMaster,recievedMessage)
	for {
		select{
		case recievedData := <- recievedMessage:
			orderDataFromMasterChan <-recievedData
		default:
			time.Sleep(5 * time.Millisecond)
		}	
	}	
}

func isInQueue(newOrderData OrderData,unfinishedOrdersChan chan []OrderData) bool {
	newOrder := newOrderData.Order
	unfinishedOrders := <- unfinishedOrdersChan
	allreadyInList := false
	for _,element:= range unfinishedOrders{
		if element.Order == newOrder{			
			allreadyInList = true
			break
		}
	}
	unfinishedOrdersChan <- unfinishedOrders 
	return allreadyInList
}

func addNewOrder(newOrderData OrderData, unfinishedOrdersChan chan []OrderData){
	deadline := time.Now().UnixNano()/int64(time.Millisecond) + 20000
	allreadyInList := isInQueue (newOrderData,unfinishedOrdersChan)
	unfinishedOrders := <- unfinishedOrdersChan
	if !allreadyInList {
		newOrderData.Deadline = deadline
		unfinishedOrders = append(unfinishedOrders,newOrderData)
	}
	unfinishedOrdersChan <- unfinishedOrders
}

func removeOrder(orderComplete OrderData,unfinishedOrdersChan chan []OrderData){
	order2Remove := orderComplete.Order
	unfinishedOrders := <- unfinishedOrdersChan
	for i,element:= range unfinishedOrders{
		if element.Order == order2Remove{
			n := len (unfinishedOrders)			
			newUnfinishedOrders :=unfinishedOrders[0:i]
			newUnfinishedOrders = append(newUnfinishedOrders,unfinishedOrders[i+1:n]...)
			unfinishedOrders = newUnfinishedOrders
			break
		}
	}
	unfinishedOrdersChan <- unfinishedOrders 
}

func readOrderData(port string, isMasterChan chan bool,recievedMessageToMaster chan OrderData, recievedMessage chan OrderData){
	bufferToRead := make([] byte, 1024)
	UDPadr,_:= net.ResolveUDPAddr("udp",""+":"+port)
    readerSocket,_ := net.ListenUDP("udp",UDPadr)
    isMaster := false
	for {
		isMaster = <- isMasterChan; isMasterChan <- isMaster
		if (port == SLAVE_TO_MASTER_PORT && !isMaster){break}
		deadline := time.Now().Add(100*time.Millisecond)
		readerSocket.SetReadDeadline(deadline)
		n,_, _ := readerSocket.ReadFromUDP(bufferToRead)     
       	if n > 0 { 
            structObject := json2struct(bufferToRead,n)
            data := structObject.OrderMessage
            if data.FromMaster{
            	recievedMessage <- data
            }else if !data.FromMaster{
            	recievedMessageToMaster <- data
            } 
        }
        time.Sleep(30 * time.Millisecond)
   	}
   	readerSocket.Close()
}
	
func sendOrderData(ip string, port string, message OrderData){
	udpAddr,_ := net.ResolveUDPAddr("udp",ip+":"+port)
	broadcastSocket,_ := net.DialUDP("udp",nil, udpAddr)	
	deadline := time.Now().UnixNano() / int64(time.Millisecond) + 2
	for {
		timeNow := time.Now().UnixNano() / int64(time.Millisecond)
		if timeNow > deadline {
			break
		}
		sendingObject := DataObject{"",[]IpObject {},message,[]OrderData {} }
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)
		time.Sleep(1 * time.Millisecond)	
	}		
}

func sendOrderDataToMaster(orderDataToMasterChan chan OrderData, myIp string, masterQueueChan chan []IpObject){
	for {
		select {
		case sendingObject := <- orderDataToMasterChan:
			sendingObject.Ip = myIp
			masterQueue:= <- masterQueueChan ; masterQueueChan <- masterQueue
			if len(masterQueue)>0 {
				ip_master := masterQueue[0].Ip			
				sendOrderData(ip_master,SLAVE_TO_MASTER_PORT,sendingObject)
			}
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func checkOrderDeadline(unactiveElevatorChan chan string, isMasterChan chan bool, unfinishedOrdersChan chan []OrderData, recievedOrderChan chan OrderData){
	isMaster := false
	for {
		unfinishedOrders := <- unfinishedOrdersChan; unfinishedOrdersChan <- unfinishedOrders
		isMaster = <- isMasterChan; isMasterChan <- isMaster
		for _,element := range unfinishedOrders{
			if !isMaster{break}
			if element.Deadline < time.Now().UnixNano()/int64(time.Millisecond) {
				unactiveElevatorChan <- element.Ip
				removeOrder(element,unfinishedOrdersChan)
				recievedOrderChan <- element
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func allocateElevatorOrders(deadElevator IpObject,unfinishedOrdersChan chan []OrderData, recievedOrderChan chan OrderData ){   
	unfinishedOrders :=  <- unfinishedOrdersChan ; unfinishedOrdersChan <- unfinishedOrders
	for _,element := range unfinishedOrders{
		if element.Ip == deadElevator.Ip {
			element.Type = ORDER
			removeOrder(element, unfinishedOrdersChan)
			recievedOrderChan <- element
		}
	}	
}