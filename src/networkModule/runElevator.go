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
	
var MasterQueue = [] IpObject {} // Er det egentlig lurt at disse er public?
var IsMaster = false
var IsBackup = false
var IsConnected = false //trenger vi denne?
var MyIp string 			

func RunElevator(){
	ipBroadcast := "129.241.187.255"
	portIp := "20017" //Her velges port som egen IP-adresse skal sendes og leses fra
	portMasterQueue := "20019" //Her velges port som Masterqueue skal sendes og leses fra
	recieveIpChan := make(chan string,1024) 
	
	go readBroadcastToUpdateQueue(portIp,recieveIpChan) //Trenger egentlig ikke ipBroadcast her fordi den leser kun fra egen port // obs 0
	go updateMasterQueue(recieveIpChan) //myIp Legges nå inn gjennom broadcastIP og updateM.Queue
	go broadcastIp(ipBroadcast,portIp) // Fungerer også som Imalive
	
	
	if IsMaster{
		go broadcastMasterQueue(ipBroadcast, portMasterQueue)									//Denne må lages. Fungerer som imAlive
		go removeDeadElevators()



	}
	if IsBackup{
		fmt.Println("jeg er Backup")
		go checkMasterAlive(portMasterQueue) // Denne må lages. Sjekker om Master broadcaster køen fortsatt. Hvis ikke blir man selv master og nestemann blir backup
	} 
	time.Sleep(5000 * time.Millisecond)
	
	fmt.Println("jodle","Masterqueue =",MasterQueue,"IsMaster=",IsMaster,"IsBackup=",IsBackup)

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
		//fmt.Println("Er igang med å broadcaste")
		
		fmt.Println("å hei hvor det går")
		if err != nil {
		    fmt.Println("error listening on UDP port ", portIp)
		    fmt.Println(err)
		    os.Exit(1)
		}
		//fmt.Println("detter er min IP",MyIp)
		sendingObject := DataObject{MyIp,[]IpObject {}}
		//fmt.Println("Slik ser sendingObject ut:",sendingObject)
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)

		
		time.Sleep(50*time.Millisecond)	


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

//Leser inn ny ip fra channel. lager en temp kø lik nåværende MasterQueue. Sjekker om ny ip ligger i køen. // Hvis ikke legges den til i lista. Hvis det er deg selv settes connected = tru
func updateMasterQueue(recieveIpChan chan string) { 
	for {
		fmt.Println("er igang med å oppdatere MasterQ")
		allreadyInQueue := false
		select{
			case newIpObject:= <- recieveIpChan:
				fmt.Println("leste fra kanal")
				index := 0
				for i,element:= range MasterQueue{
					if element.Ip == newIpObject{
						allreadyInQueue = true
						index = i
						fmt.Println("den var allerede på plass den gitt")
						break
					}
				}
				if newIpObject == MyIp && !IsConnected{
					IsConnected = true
				
				}
				if allreadyInQueue && IsMaster{
					deadline := time.Now().UnixNano() / int64(time.Millisecond) + 1500
					MasterQueue[index].Deadline = deadline

				}
				if !allreadyInQueue{
					deadline := time.Now().UnixNano() / int64(time.Millisecond)  + 1500
					object := IpObject {newIpObject,deadline}
					MasterQueue = append(MasterQueue,object)	
				}

		default:
			fmt.Println("går i default")
			time.Sleep(50*time.Millisecond)	
		}
		
		
	}	
}

	
func readBroadcastToUpdateQueue(portIp string,recieveIpChan chan string) { 		
	bufferToRead := make([] byte, 1024)
	
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+portIp)
	
	if err != nil {
        fmt.Println("error resolving UDP address on ", portIp)
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("Har nå kommet til ListenUDP")
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
        fmt.Println("error listening on UDP port ", portIp)
        fmt.Println(err)
        os.Exit(1)
        
	}
	
	for {
		fmt.Println("er i gang med å lese fra Broadcast")
		n,UDPadr, err := readerSocket.ReadFromUDP(bufferToRead)
        
	 	if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            os.Exit(1)
            
        }
        
        fmt.Println("got message from ", UDPadr, " with n = ", n)

       	if n > 0 {
           	fmt.Println("printer melding vi leste over UDP",string(bufferToRead))  
            structObject := json2struct(bufferToRead,n)
            ip := structObject.NewIp
            fmt.Println("fikk noe med lengde større enn 0")
            select{
            	case recieveIpChan <- ip:
            		fmt.Println("har sendt noe til kanalen")
            	default:
            		time.Sleep(50*time.Millisecond)
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
		//fmt.Println("Er igang med å broadcaste")
		
		fmt.Println("er igang med å broadcaste MasterQueue")
		if err != nil {
		    fmt.Println("error listening on UDP port ", portMasterQueue)
		    fmt.Println(err)
		    os.Exit(1)
		}
		//fmt.Println("detter er min IP",MyIp)
		sendingObject := DataObject{"",MasterQueue}
		//fmt.Println("Slik ser sendingObject ut:",sendingObject)
		jsonFile := struct2json(sendingObject)
		broadcastSocket.Write(jsonFile)

		
		time.Sleep(50*time.Millisecond)	


	}	
}

func checkMasterAlive(portMasterQueue string){
	bufferToRead := make([] byte, 1024)
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+portMasterQueue)

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
	deadline := time.Now().Add(1500*time.Millisecond)
	readerSocket.SetReadDeadline(deadline)
	n,UDPadr,err := readerSocket.ReadFromUDP(bufferToRead[0:])
		
    readerSocket.Close()
 	if err != nil {
 		fmt.Println("Alle mann til pumpene, Master er død. Jeg tar over, follow my command.")
        IsMaster = true
        IsBackup = false      
    }
    
    fmt.Println("got message from ", UDPadr, " with n = ", n)

   	if n > 0 {
       	fmt.Println("Master er i kjempeform og jer er fortsatt backup. I køen hans ligger:",json2struct(bufferToRead,n))  
      		 
   }
   
 
  }
}

func removeDeadElevators(){ // Problemer her me at denne går så fort at de andre ikke kommer til?
	for{
		tempQueue := MasterQueue
		if len(tempQueue) > 0{
			for _,element:= range tempQueue{
				timeNow := time.Now().UnixNano() / int64(time.Millisecond)
				if timeNow > element.Deadline{
					fmt.Println("tiden gikk ut")
					fmt.Println("det har gått for lang tid siden vi hørte fra heisen med ip",element.Ip,"Den fjernes derfor fra Masterqueue")
					newMasterQueue := [] IpObject {}
					for _,element2:= range tempQueue{
						if element != element2{
							newMasterQueue =append(newMasterQueue,element2)
						}
					}
					MasterQueue = newMasterQueue
					break
				}
			}
		}else{
			time.Sleep(500 * time.Millisecond)
		}	
	}
}





