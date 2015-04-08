package networkModule

//mangler broadcastMasterQueue og checkMasterAlive. Disse bør muligens flyttes til runElevator
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

func RunElevator(){
	ipBroadcast := "129.241.187.255"
	portIp := "20017" //Her velges port som egen IP-adresse skal sendes og leses fra
	portMasterQueue := "20019" //Her velges port som Masterqueue skal sendes og leses fra
	recieveIpChan := make(chan string,1024) 
	isMasterChan := make(chan bool,1)
	isBackupChan := make(chan bool,1)
	go updateMasterQueue(portMasterQueue,isMasterChan, isBackupChan)
	go broadcastIp(ipBroadcast,portIp) // Fungerer også som Imalive
	go handleOrdersFromMaster()
	for{
		if isMaster{
			fmt.Println("Er nå Master")
			fmt.Println("Masterqueue =",masterQueue,"isMaster=",isMaster,"isBackup=",isBackup)
			go listenForActiveElevators(portIp,recieveIpChan)
			go updateElevators(recieveIpChan) //myIp Legges nå inn gjennom broadcastIP og updateM.Queue
			go broadcastMasterQueue(ipBroadcast, portMasterQueue)									//Denne må lages. Fungerer som imAlive
			go removeDeadElevators()
			go handleOrdersInNetwork()
			deadChan := make(chan int)
			<-deadChan

		}else if isBackup {
			fmt.Println("Masterqueue =",masterQueue,"isMaster=",isMaster,"isBackup=",isBackup)
			fmt.Println("jeg er Backup")
			isMaster = <- isMasterChan	
			isBackup = false
			
					
		}else{
			fmt.Println("Masterqueue =",masterQueue,"isMaster=",isMaster,"isBackup=",isBackup)
			fmt.Println("Jeg er bare slave")
			isBackup = <-isBackupChan
			
			
		} 
	}

	
	
	

}

func broadcastIp(ipBroadcast string, portIp string){
	udpAddr, err := net.ResolveUDPAddr("udp",ipBroadcast+":"+portIp)
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
		sendingObject := DataObject{myIp,[]IpObject {}}
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)

		
		


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

func updateMasterQueue(portMasterQueue string,isMasterChan chan bool,isBackupChan chan bool){
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+portMasterQueue) //muligens "" istedet for myIp

	if err != nil {
                fmt.Println("error resolving UDP address on ", portMasterQueue)
                fmt.Println(err)
                os.Exit(1)
    }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
            fmt.Println("error listening on UDP port ", portMasterQueue)
            fmt.Println(err)
            os.Exit(1)
	}
	for{
		bufferToRead := make([] byte, 1024)
		deadline := time.Now().Add(1500*time.Millisecond)
		readerSocket.SetReadDeadline(deadline)
		n,_,err := readerSocket.ReadFromUDP(bufferToRead[0:])
			
	    
	 	if err != nil && isBackup {
	 		fmt.Println("Alle mann til pumpene, Master er død. Jeg tar over, follow my command.")
	        isMasterChan <- true
	                
	    }
	    
	    select{
	    case <- masterQueueLock:
		   	if n > 0 && !isMaster {
		   		//fmt.Println(n)
		       	structObject := json2struct(bufferToRead,n)
		       	
		       	masterQueue = structObject.MasterQueue
		       	if !isBackup && len(masterQueue) > 1 {
		       		if masterQueue[1].Ip == myIp{
		       			isBackupChan <- true
		       			fmt.Println("Nå er jeg backup")
		       		}
		       	}
		       	
		       	 	   		 
	    	}
	  	 	masterQueueLock <- 1
	  	 	time.Sleep(5 * time.Millisecond)			
	  	default:
	  		time.Sleep(5 * time.Millisecond)		
	  	}
   
	}
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
            ip := structObject.NewIp
            select{
            	case recieveIpChan <- ip:
            	default:
            		time.Sleep(5*time.Millisecond)		
            }    
        }

   	}
}

func broadcastMasterQueue(ipBroadcast string,portMasterQueue string){
	udpAddr, err := net.ResolveUDPAddr("udp",ipBroadcast+":"+portMasterQueue)
	if err != nil {
		fmt.Println("error resolving UDP address on ", portMasterQueue)
		fmt.Println(err)
		os.Exit(1)
	}
	broadcastSocket, err := net.DialUDP("udp",nil, udpAddr)
	for {
		select{
		case <- masterQueueLock:
			if err != nil {
			    fmt.Println("error listening on UDP port ", portMasterQueue)
			    fmt.Println(err)
			    os.Exit(1)
			}
			sendingObject := DataObject{"",masterQueue}
			masterQueueLock <- 1
			jsonFile := struct2json(sendingObject)
			broadcastSocket.Write(jsonFile)
		default:
			time.Sleep(5 * time.Millisecond)					

		}
			


	}	
}

func removeDeadElevators(){ 
	for{
		<- masterQueueLock
		n := len(masterQueue)
		if n > 0{
			for i,element:= range masterQueue{
				timeNow := time.Now().UnixNano() / int64(time.Millisecond)
				if timeNow > element.Deadline && element.Ip != myIp{
					fmt.Println("fjerner død heis")
					newMasterQueue :=masterQueue[0:i]
					newMasterQueue = append(newMasterQueue,masterQueue[i+1:n]...)
					masterQueue = newMasterQueue
					break
				}
			}
		}
		masterQueueLock <- 1
		time.Sleep(5 * time.Millisecond)		
 	}
}


