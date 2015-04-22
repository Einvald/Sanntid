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
	Order ButtonOrder
	Cost int 
	Ip string
	

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



// Trenger funksjonalitet til: Hvis ikke bestillingen mottas av Master må den tas selv
// Diskuter deadlocks knyttet til unfinishedOrdersLock
//Diskuter kommentar linje 147 runElevator

const IP_BROADCAST = "129.241.187.255"
const BROADCAST_PORT = "20012"
const SLAVE_TO_MASTER_PORT = "20013"
const MASTER_TO_SLAVE_PORT = "20014"
const COST_CEILING = 100 
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

	unfinishedOrdersLock <- 1
	auctionLock <- 1
	go readOrderData(SLAVE_TO_MASTER_PORT)
	for {
		select{
		case recievedData := <- recievedMessageToMaster:
			//fmt.Println("INNE I HANDLEORDERSINNETWORK har mottat type: ", recievedData.Type)
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
				fmt.Println("RecievedOrderChan er nå blitt lest ", order)
				go auction(order,IP_BROADCAST,BROADCAST_PORT,MASTER_TO_SLAVE_PORT) 

			}
			// Sjekk om finnes i liste. Legg til hvis ikke
		case orderComplete := <- recievedOrderComplete:
			<- unfinishedOrdersLock
			orderComplete.FromMaster = true 
			removeOrder(orderComplete)
			sendOrderData(IP_BROADCAST,MASTER_TO_SLAVE_PORT,orderComplete)
			unfinishedOrdersLock <- 1 
		}
		
	}

}

func auction(newOrderData OrderData, IP_BROADCAST string, BROADCAST_PORT string, portSomSlaverLeserFra string){
	<- auctionLock
	if !isInQueue(newOrderData){
		fmt.Println("AUCTION I GANG")
		deadline := time.Now().UnixNano() / int64(time.Millisecond) + 800
		elevatorsInAuction := [] OrderData {} 
		newOrderData.Type = REQUEST_AUCTION
		newOrderData.FromMaster = true  
		sendOrderData(IP_BROADCAST,BROADCAST_PORT,newOrderData)  
		for {

			select{
			case cost := <- recievedCostChan:
				fmt.Println("COST IS RECIEVED! ", cost)
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
			if len(elevatorsInAuction) == len(masterQueue) || timeNow > deadline  {


				break
			}
		}
		lowestCost := COST_CEILING
		elevatorWithLowestCost := OrderData {}
		elevatorWithLowestCost.Cost = COST_CEILING
		fmt.Println("AUCTIONCHECK", elevatorWithLowestCost)
		for _,element:= range elevatorsInAuction{
			if element.Cost < lowestCost{
				lowestCost = element.Cost
				elevatorWithLowestCost = element

			}
			fmt.Println("COSTSJEKKEN på følgende element: ", element)
		
		}
		//<- unfinishedOrdersLock 
		fmt.Println("LOwEST COST ER NÅ FUNNET TIL Å VÆRE: ", lowestCost)
		addNewOrder(elevatorWithLowestCost)
		//unfinishedOrdersLock <- 1
		elevatorWithLowestCost.Type = ORDER
		elevatorWithLowestCost.FromMaster = true
		fmt.Println("ELEVATOR WITH LOWEST COST SIN IP ER FØLGENDE: ", elevatorWithLowestCost.Ip)
		sendOrderData(elevatorWithLowestCost.Ip,portSomSlaverLeserFra,elevatorWithLowestCost) // Porten her må bestemmes eksternt!!
	}
	auctionLock <- 1	
}

func handleOrdersFromMaster(Order_data_from_master_chan chan OrderData){
	go readOrderData(MASTER_TO_SLAVE_PORT)
	go readOrderData(BROADCAST_PORT)
	for {
		select{
		case recievedData := <- recievedMessage:
			Order_data_from_master_chan <-recievedData
		default:
			time.Sleep(5 * time.Millisecond)

		}
			
	}	
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
	unfinishedOrdersLock <- 1
	return allreadyInList

}
//Funksjoner master bruker	
func addNewOrder(newOrderData OrderData){ 
	allreadyInList := isInQueue (newOrderData)
	if !allreadyInList {
		unfinishedOrders = append(unfinishedOrders,newOrderData)
	}
}


func removeOrder(orderComplete OrderData){
	order2Remove := orderComplete.Order
	fmt.Println("UNFINISHED ORDERS LISTEN: ", unfinishedOrders)
	for i,element:= range unfinishedOrders{
		if element.Order == order2Remove{
			n := len (unfinishedOrders)			
			newUnfinishedOrders :=unfinishedOrders[0:i]
			newUnfinishedOrders = append(newUnfinishedOrders,unfinishedOrders[i+1:n]...)
			unfinishedOrders = newUnfinishedOrders
			break
		}
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
            	fmt.Println("Data fra master lest, melding: ", data)
            
            }else if !data.FromMaster{
            	recievedMessageToMaster <- data
            	fmt.Println("RECIEVEDMESSAGE TO MASTER CHANNEL IS SET")
            } 
            
        }
        time.Sleep(30 * time.Millisecond)
   	}
}
	//Hvis til Master legges det på masterkanalen, hvis ikke på den andre kanalen


func sendOrderData(ip string, port string, message OrderData){	//Denne brukes både til å sende enkle meldinger, men også broadcasting
	fmt.Println("SENDER ER SATT I GANG, FÅ GANG OG MOTTA! ", ip, " ", port, " ", message)
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
		sendingObject := DataObject{"",[]IpObject {},message,[]OrderData {} }
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)
		time.Sleep(1 * time.Millisecond)	

	}
	// Vi må sende flere ganger og hvis en bestilling ikke mottas av master skal heisen ta den selv.
	
}

func sendOrderDataToMaster(Order_data_to_master_chan chan OrderData){
	for {
		select {
		case sendingObject := <- Order_data_to_master_chan:
			fmt.Println("Inne i send Order to Master")
			sendingObject.Ip = myIp
			<- masterQueueLock 
			if len(masterQueue)>0 {
				ip_master := masterQueue[0].Ip
				
				sendOrderData(ip_master,SLAVE_TO_MASTER_PORT,sendingObject)
				fmt.Println("Sendt order data to master")
			}
			masterQueueLock <-1
		default:
			time.Sleep(5 * time.Millisecond)

		
		}
	}
}


//Spør om lys skal lyse i alle på alle heispaneler hvis man bestemmer for eksempel ned fra 2. etasje?

