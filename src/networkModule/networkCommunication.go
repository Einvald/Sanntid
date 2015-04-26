package networkModule

import (
    	"fmt"
    	"net"
    	"encoding/json"
    	"time"
    	"strings"
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
const PORT_MASTER_DATA = "31019"
const PORT_MY_IP = "31017"

func RunNetworkCommunication(orderDataToMasterChan chan OrderData, orderDataFromMasterChan chan OrderData){
	recieveIpChan := make(chan string,1024)
	isMasterChan := make(chan bool,1)
	isBackupChan := make(chan bool,1)
	becomeMasterChan := make(chan int,1)
	becomeBackupChan := make(chan int,1) 
	becomeSlaveChan := make(chan int,1)
	masterStartTime := make(chan int64,1 )
	masterQueueChan := make(chan [] IpObject ,1)
	unfinishedOrdersChan := make(chan [] OrderData,1)
	recievedMessageToMaster := make(chan OrderData, 1024)
	recievedMessage := make (chan OrderData,1024)
	recievedOrderChan :=make(chan OrderData,1024)
	masterQueueChan <- []IpObject{}
	unfinishedOrdersChan <- [] OrderData {}
	masterStartTime <- 0
	isMasterChan <- false
	isBackupChan <- false
	myIp := getIpAddr()
	initializeElevatorConnection(myIp,masterQueueChan,becomeMasterChan,becomeBackupChan)
	go updateMasterData(unfinishedOrdersChan, masterQueueChan,myIp,becomeMasterChan,becomeBackupChan,becomeSlaveChan,isMasterChan,isBackupChan,masterStartTime)
	go broadcastIp(myIp)
	go sendOrderDataToMaster(orderDataToMasterChan,myIp,masterQueueChan)
	go handleOrdersFromMaster(orderDataFromMasterChan, isMasterChan,recievedMessage,recievedMessageToMaster,myIp)
	for{
		masterQueue := <- masterQueueChan; masterQueueChan <- masterQueue
		fmt.Println("Masterqueue=",masterQueue)
		select{
			case <-becomeMasterChan:
				<-isMasterChan; isMasterChan <- true
				<- masterStartTime; masterStartTime <- time.Now().UnixNano()/int64(time.Millisecond)
				fmt.Println("Is Master")
				go listenForActiveElevators(recieveIpChan,isMasterChan)
				go updateElevators(masterQueueChan,recieveIpChan,isMasterChan,isBackupChan) 
				go broadcastMasterData(masterQueueChan,isMasterChan,unfinishedOrdersChan)					
				go handleDeadElevators(masterQueueChan,isMasterChan,myIp,unfinishedOrdersChan,recievedOrderChan)
				go handleOrdersInNetwork(isMasterChan,masterQueueChan,unfinishedOrdersChan, recievedMessage,recievedMessageToMaster, recievedOrderChan)			
				<-becomeSlaveChan
				<-isMasterChan; isMasterChan <- false
			case <-becomeBackupChan:
				<-isBackupChan; isBackupChan <- true
				fmt.Println("Is Backup")
				<- becomeMasterChan	
				<-isBackupChan; isBackupChan <- false
				becomeMasterChan <- 1	
		}
	}			
}

func broadcastIp(myIp string){
	udpAddr,_:= net.ResolveUDPAddr("udp",IP_BROADCAST+":"+PORT_MY_IP)
	broadcastSocket,_ := net.DialUDP("udp",nil, udpAddr)
	for {
		sendingObject := DataObject{myIp,[]IpObject {},OrderData {},[]OrderData {} }
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)
		time.Sleep(300 * time.Millisecond)
	}	
}

func struct2json(packageToSend DataObject) [] byte {
	jsonObject, _ := json.Marshal(packageToSend)
	return jsonObject
}

func json2struct(jsonObject []byte,n int) DataObject{
	structObject := DataObject{}
	json.Unmarshal(jsonObject[0:n], &structObject)  
	return structObject
}

func updateMasterData(unfinishedOrdersChan chan []OrderData,masterQueueChan chan []IpObject,myIp string,becomeMasterChan chan int,becomeBackupChan chan int, becomeSlaveChan chan int,isMasterChan chan bool,isBackupChan chan bool, masterStartTime chan int64){
	UDPadr, _:= net.ResolveUDPAddr("udp",""+":"+PORT_MASTER_DATA)   
    readerSocket ,_ := net.ListenUDP("udp",UDPadr)
	deadMaster := false
	twoMasters := false
	isMaster := false
	isBackup := false
	for{
		isMaster = <-isMasterChan ; isMasterChan <- isMaster
		isBackup = <-isBackupChan; isBackupChan <- isBackup
		bufferToRead := make([] byte, 1024)
		deadline := time.Now().Add(1000*time.Millisecond)
		readerSocket.SetReadDeadline(deadline)
		n,_,timeout := readerSocket.ReadFromUDP(bufferToRead[0:])
				    
	 	if timeout != nil {
	 		if isBackup{
	        	becomeMasterChan <- 1
	        	continue
	        }else if !isMaster {
	        	if deadMaster{
	        		becomeMasterChan <- 1
	        		deadMaster = false
	        	}else{deadMaster = true;}	
	        	continue
	        } 
	    }
	    timeAsMaster := <-masterStartTime; masterStartTime <- timeAsMaster
	    timeAsMaster = time.Now().UnixNano()/int64(time.Millisecond) - timeAsMaster
	    masterQueue := <- masterQueueChan; masterQueueChan <-masterQueue
	    if len(masterQueue) > 0 {
		   	ipMaster := masterQueue[0].Ip	   			
	   		if isMaster && ipMaster != myIp {
	   			myIpSplit := strings.Split(myIp,".")
	   			ipMasterSplit := strings.Split(ipMaster,".")
	   			if myIpSplit[3] < ipMasterSplit[3] && twoMasters &&timeAsMaster > 1500{ 
 	   				becomeSlaveChan <- 1
	   				twoMasters = false
	   				continue
	   			}else if timeAsMaster > 1500{
	   				twoMasters = true
	   				time.Sleep(600 *time.Millisecond)
	   			}
	   		}
		}   
	   	if n > 0 {
	   		structObject := json2struct(bufferToRead,n)
	   		if len(structObject.MasterQueue) >0  {
		   		ipMaster := structObject.MasterQueue[0].Ip	   			
		   		if isMaster && ipMaster != myIp {
		   			myIpSplit := strings.Split(myIp,".")
		   			ipMasterSplit := strings.Split(ipMaster,".")
		   				if myIpSplit[3] < ipMasterSplit[3] && timeAsMaster > 1500{
		   					<-masterQueueChan; masterQueueChan <-structObject.MasterQueue
		   				}
		   		}	
	   		}
	   		if !isMaster {
	   			newMasterQueue := structObject.MasterQueue
	   			<-unfinishedOrdersChan; unfinishedOrdersChan <- structObject.UnfinishedOrders
		       	<-masterQueueChan	       	
		       	if !isBackup && len(newMasterQueue) > 1 {
		       		if newMasterQueue[1].Ip == myIp{
		       			becomeBackupChan <- 1
		       		}
		    	}
		    	masterQueueChan <-newMasterQueue
		    }
	    }  	
	    time.Sleep(5 * time.Millisecond)	 
    }  	 					  			
}

func updateElevators(masterQueueChan chan []IpObject,recieveIpChan chan string, isMasterChan chan bool, isBackupChan chan bool) { 
	isMaster := false
	isBackup := false 
	for {
		isMaster = <- isMasterChan; isMasterChan <- isMaster
		isBackup = <- isBackupChan; isBackupChan <- isBackup
		if !isMaster{break} 
		allreadyInQueue := false
		newIpObject:= <- recieveIpChan
		index := 0
		masterQueue := <- masterQueueChan
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
		masterQueueChan <- masterQueue 		
		time.Sleep(5 * time.Millisecond)
	}	
}
	
func listenForActiveElevators(recieveIpChan chan string, isMasterChan chan bool) { 
	bufferToRead := make([] byte, 1024)	
	UDPadr, _:= net.ResolveUDPAddr("udp",""+":"+PORT_MY_IP)   
    readerSocket ,_ := net.ListenUDP("udp",UDPadr)
	isMaster := false
	for {
		isMaster = <- isMasterChan ; isMasterChan <- isMaster
		if !isMaster{break}
		n,_, _ := readerSocket.ReadFromUDP(bufferToRead) 
       	if n > 0 { 
            structObject := json2struct(bufferToRead,n)
            ip := structObject.NewIp
            select{
            	case recieveIpChan <- ip:
            		time.Sleep(5 * time.Millisecond)	
            	default:
            		time.Sleep(5*time.Millisecond)		
            }    
        }
   	}
   	readerSocket.Close()
}

func broadcastMasterData(masterQueueChan chan []IpObject, isMasterChan chan bool, unfinishedOrdersChan chan []OrderData){
	udpAddr, _ := net.ResolveUDPAddr("udp",IP_BROADCAST+":"+PORT_MASTER_DATA)
	broadcastSocket, _ := net.DialUDP("udp",nil, udpAddr)
	isMaster := false
	for {
		isMaster = <- isMasterChan; isMasterChan <-isMaster
		if !isMaster{break} 
		select{
		case masterQueue:= <- masterQueueChan:			
			unfinishedOrders := <- unfinishedOrdersChan; unfinishedOrdersChan <- unfinishedOrders 
			sendingObject := DataObject{"",masterQueue,OrderData{},unfinishedOrders}
			masterQueueChan <- masterQueue
			jsonFile := struct2json(sendingObject)
			broadcastSocket.Write(jsonFile)
			time.Sleep(5 * time.Millisecond)
		default:
			time.Sleep(5 * time.Millisecond)					
		}
	}
	broadcastSocket.Close()	
}

func handleDeadElevators(masterQueueChan chan []IpObject, isMasterChan chan bool, myIp string,unfinishedOrdersChan chan []OrderData, recievedOrderChan chan OrderData){
	isMaster := false 
	for{
		isMaster = <- isMasterChan; isMasterChan <-isMaster
		if !isMaster{break} 
		masterQueue:= <- masterQueueChan
		n := len(masterQueue)
		if n > 0{
			for i,element:= range masterQueue{
				timeNow := time.Now().UnixNano() / int64(time.Millisecond)
				if timeNow > element.Deadline && element.Ip != myIp{
					masterQueueChan <- masterQueue
					removeElevator(masterQueueChan,i,n)
					masterQueue = <- masterQueueChan
					allocateElevatorOrders(element,unfinishedOrdersChan,recievedOrderChan)	
					break
				}
			}
		}
		masterQueueChan <- masterQueue
		time.Sleep(25 * time.Millisecond)		
 	}
}

func removeElevator(masterQueueChan chan []IpObject, deadIndex int,lengthMasterQueue int){
	masterQueue := <- masterQueueChan
	newMasterQueue :=masterQueue[0:deadIndex]
	newMasterQueue = append(newMasterQueue,masterQueue[deadIndex+1:lengthMasterQueue]...)
	fmt.Println("Remove dead elevator, new master queue = ",newMasterQueue)
	masterQueueChan <- newMasterQueue
}
