package networkModule


// Videre trenger vi å utvikle nettverksmodulen slik at heisene kan kommunisere mer direkte med hverandre. F.eks broaccaste bestillinger

//Mtp goRoutines: Vi kaller heller connectElevator fra main og så kjører vi Runelevator som en goroutine slik at dette gjøres for allti.d



import (
    	"fmt"
    	"net"
    	"os"
    	"encoding/json"
    	"time"
	)

type DataObject struct {
	NewIp string
	MasterQueue [] IpObject
	OrderMessage OrderData
	UnfinishedOrders [] OrderData
}

type IpObject struct{
	Ip string
	Deadline int64
}


var masterQueue = [] IpObject {} 
var isMaster = false
var isBackup = false
var myIp string 			
var masterQueueLock = make(chan int, 1);

func RunElevator(Order_data_to_master_chan chan OrderData, Order_data_from_master_chan chan OrderData){
	portIp := "20017" //Her velges port som egen IP-adresse skal sendes og leses fra
	portMasterData := "20019" //Her velges port som Masterqueue skal sendes og leses fra
	recieveIpChan := make(chan string,1024) 
	isMasterChan := make(chan bool,1)
	isBackupChan := make(chan bool,1)
	isOffline := make(chan bool,1)
	go updateMasterData(portMasterData,isMasterChan, isBackupChan)
	go broadcastIp(IP_BROADCAST,portIp) // Fungerer også som Imalive
	go sendOrderDataToMaster(Order_data_to_master_chan)
	go handleOrdersFromMaster(Order_data_from_master_chan)
	for{
		if isMaster{
			fmt.Println("Er nå Master")
			fmt.Println("Masterqueue =",masterQueue,"isMaster=",isMaster,"isBackup=",isBackup)
			go listenForActiveElevators(portIp,recieveIpChan)
			go updateElevators(recieveIpChan) //myIp Legges nå inn gjennom broadcastIP og updateM.Queue
			go broadcastMasterData(IP_BROADCAST, portMasterData)									//Denne må lages. Fungerer som imAlive
			go handleDeadElevators()
			go handleOrdersInNetwork()
			<- isOffline
			for _,element:= range unfinishedOrders{
				element.Type = ORDER
				Order_data_from_master_chan <- element
			
			time.Sleep(3000 * time.Millisecond)
			fmt.Println("MAsterqueu1:",masterQueue)
						
			break
			}
		}else if isBackup {
			fmt.Println("Masterqueue =",masterQueue,"isMaster=",isMaster,"isBackup=",isBackup)
			fmt.Println("jeg er Backup")
			isMaster = <- isMasterChan	
			isBackup = false
			
					
		}else{
			fmt.Println("Masterqueue =",masterQueue,"isMaster=",isMaster,"isBackup=",isBackup)
			fmt.Println("Jeg er bare slave")
			//select {
				//case isBackup <-isBackupChan:
				//case isMaster <-isMasterChan:  			
			
			//} 
		}
	}			
}

func broadcastIp(IP_BROADCAST string, portIp string){
	udpAddr, err := net.ResolveUDPAddr("udp",IP_BROADCAST+":"+portIp)
	if err != nil {
		fmt.Println("error resolving UDP address on ", portIp)
		fmt.Println(err)
		os.Exit(1)
	}
	broadcastSocket, err := net.DialUDP("udp",nil, udpAddr)
	for {
		
		if err != nil {
		    fmt.Println("error listening on UDP port ", portIp)
		    fmt.Println(err)
		    os.Exit(1)
		}
		sendingObject := DataObject{myIp,[]IpObject {},OrderData {},[]OrderData {} }
		//fmt.Println("Printer nå sendingObject:",sendingObject)
		jsonFile := struct2json(sendingObject)
		//test := json2struct(jsonFile,100)
		//fmt.Println("Printer nå tilbakekonvertert shit:",test)
		broadcastSocket.Write(jsonFile)
		time.Sleep(300 * time.Millisecond)
		// HER BØR VI KANSKJE SLEEPE


	}	
}


func struct2json(packageToSend DataObject) [] byte {
	jsonObject, _ := json.Marshal(packageToSend)
	return jsonObject
}

func json2struct(jsonObject []byte,n int) DataObject{
	structObject := DataObject{}
	json.Unmarshal(jsonObject[0:n], &structObject)  //Her kan det være noe som ikke stemmer helt
	return structObject
}

func updateMasterData(portMasterData string,isMasterChan chan bool,isBackupChan chan bool){
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+portMasterData) //muligens "" istedet for myIp
	if err != nil {
                fmt.Println("error resolving UDP address on ", portMasterData)
                fmt.Println(err)
                os.Exit(1)
    }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
            fmt.Println("error listening on UDP port ", portMasterData)
            fmt.Println(err)
            os.Exit(1)
	}
	deadMaster := false
	for{
		
		bufferToRead := make([] byte, 1024)
		deadline := time.Now().Add(1000*time.Millisecond)
		readerSocket.SetReadDeadline(deadline)
		n,_,timeout := readerSocket.ReadFromUDP(bufferToRead[0:])
				    
	 	if timeout != nil {
	 		if isBackup{
	 			fmt.Println("Alle mann til pumpene, Master er død. Jeg tar over, follow my command.")
	        	isMasterChan <- true
	        }else if !isMaster {
	        	if deadMaster{
	        		isMasterChan <- true

	        	}else{deadMaster = true; continue;}
	        	
	        } 
	    }
	    
	 
	    <- masterQueueLock   
	   	if n > 0 && !isMaster {		
	       	structObject := json2struct(bufferToRead,n)
	       	masterQueue = structObject.MasterQueue
	       	//<- unfinishedOrdersLock
	       	unfinishedOrders = structObject.UnfinishedOrders // Funksjonalitet lagt til
	       	//unfinishedOrdersLock <- 1
	       	}
	       	if !isBackup && len(masterQueue) > 1 {
	       		if masterQueue[1].Ip == myIp{
	       			isBackupChan <- true
	       			fmt.Println("Nå er jeg backup")
	       		}
	       	}
	    masterQueueLock <- 1   	
	       	 	   		 
    }
  	 	
  	 	//time.Sleep(5 * time.Millisecond)			
	  			
}

//Leser inn ny ip fra channel. lager en temp kø lik nåværende MasterQueue. Sjekker om ny ip ligger i køen. // Hvis ikke legges den til i lista.
func updateElevators(recieveIpChan chan string) { 
	for {
		allreadyInQueue := false
		select{
			case newIpObject:= <- recieveIpChan:
				index := 0
				<- masterQueueLock
				for i,element:= range masterQueue{
					if element.Ip == newIpObject{
						allreadyInQueue = true
						index = i						
						break
					}
				}
				
				
				if allreadyInQueue && isMaster{
					deadline := time.Now().UnixNano() / int64(time.Millisecond) + 2000
					masterQueue[index].Deadline = deadline

					

				}
				if !allreadyInQueue{
					deadline := time.Now().UnixNano() / int64(time.Millisecond)  + 2000
					object := IpObject {newIpObject,deadline}
					masterQueue = append(masterQueue,object)	
				}
				masterQueueLock <- 1
		default:
			//time.Sleep(5 * time.Millisecond)	
			
				
		}
		
		
	}	
}
	
func listenForActiveElevators(portIp string,recieveIpChan chan string) { 

	bufferToRead := make([] byte, 1024)
	
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+portIp)
	
	if err != nil {
        fmt.Println("error resolving UDP address on ", portIp)
        fmt.Println(err)
        os.Exit(1)
    }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
        fmt.Println("error listening on UDP port ", portIp)
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
           // fmt.Println("Printer nå mottatt pakke over UDP",structObject)
            ip := structObject.NewIp
            //fmt.Println("ip = ",ip)
            select{
            	case recieveIpChan <- ip:
            		time.Sleep(5 * time.Millisecond)	
            	default:
            		time.Sleep(5*time.Millisecond)		
            }    
        }

   	}
}

func broadcastMasterData(IP_BROADCAST string,portMasterData string){
	udpAddr, err := net.ResolveUDPAddr("udp",IP_BROADCAST+":"+portMasterData)
	if err != nil {
		fmt.Println("error resolving UDP address on ", portMasterData)
		fmt.Println(err)
		os.Exit(1)
	}
	broadcastSocket, err := net.DialUDP("udp",nil, udpAddr)
	for {
		select{
		case <- masterQueueLock:			
			if err != nil {
			    fmt.Println("error listening on UDP port ", portMasterData)
			    fmt.Println(err)
			    os.Exit(1)
			}
			//<-unfinishedOrdersLock 
			sendingObject := DataObject{"",masterQueue,OrderData{},unfinishedOrders}
			//unfinishedOrdersLock <- 1 													//Her kan vi muligens få en deadlock. Diskuter dette!
			masterQueueLock <- 1
			jsonFile := struct2json(sendingObject)
			broadcastSocket.Write(jsonFile)
			time.Sleep(30 * time.Millisecond)
		default:
			time.Sleep(5 * time.Millisecond)					

		}
			


	}	
}

func handleDeadElevators(){ 
	for{
		<- masterQueueLock
		n := len(masterQueue)
		if n > 0{
			for i,element:= range masterQueue{
				timeNow := time.Now().UnixNano() / int64(time.Millisecond)
				if timeNow > element.Deadline && element.Ip != myIp{
					removeElevator(i,n)
					allocateElevatorOrders(element)	
					break
				}
			}
		}
		masterQueueLock <- 1
		time.Sleep(5 * time.Millisecond)		
 	}
}

func removeElevator(deadIndex int,lengthMasterQueue int){
	fmt.Println("fjerner død heis")
	newMasterQueue :=masterQueue[0:deadIndex]
	newMasterQueue = append(newMasterQueue,masterQueue[deadIndex+1:lengthMasterQueue]...)
	masterQueue = newMasterQueue
}


func allocateElevatorOrders(deadElevator IpObject){   // Dette må vi diskutere
	<-unfinishedOrdersLock
	for _,element := range unfinishedOrders{
		if element.Ip == deadElevator.Ip {

			element.Type = ORDER
			removeOrder(element)
			recievedOrderChan <- element
			// Her må man starte en ny separat auksjon av alle bestillinger som ikke er tatt fra den døde heisen.
		}

	}
	unfinishedOrdersLock <- 1

}

