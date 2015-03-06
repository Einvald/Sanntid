package networkModule

//mangler broadcastMasterQueue og checkMasterAlive. Disse bør muligens flyttes til runElevator?
// Videre trenger vi å utvikle nettverksmodulen slik at heisene kan kommunisere mer direkte med hverandre. F.eks broaccaste bestillinger


import (
    	"fmt"
    	"net"
    	"os"
	)

type dataObject struct {
	newIp string
	masterQueue []string
	
MasterQueue:= [] string {}
IsMaster:= false
IsBackup := false
Isonnected := false
MyIp := string 			//Lovlig??

func runElevator(){
	ipBroadcast := "129.241.187.255"
	portIp := "20017" //Her velges port som egen IP-adresse skal sendes og leses fra
	portMasterQueue:= "20018" //Her velges port som Masterqueue skal sendes og leses fra
	recieveIpChan:= make(chan string,1024) 
	connectElevator(ipBroadcast string,portMasterQueue string)
	go readToUpdateQueue(ipBroadcast,portIp,recieveIpChan) //Trenger egentlig ikke ipBroadcast her fordi den leser kun fra egen port // obs 0
	go updateMasterQueue(recieveIpChan) //myIp Legges nå inn gjennom broadcastIP og updateM.Queue
	go broadcastIp(MyIp,ipBroadcast,ipPort) // Fungerer også som Imalive
	if IsMaster{
		broadcastMasterQueue(ipBroadcast, queuePort) 	//Denne må lages. Fungerer som imAlive

	}
	if IsBackup{
		checkMaster() // Denne må lages. Sjekker om Master broadcaster køen fortsatt. Hvis ikke blir man selv master og nestemann blir backup
	} 



}

func broadcastIp(myIp string, ipBroadcast string, port string){
	udpAddr, err := net.ResolveUDPAddr("udp",ipBroadcast+":"+port)
	if err != nil {
                fmt.Println("error resolving UDP address on ", portNumber)
                fmt.Println(err)
                os.Exit(1)
        }
	
	broadcastSocket, err := net.DialUDP("udp",nil, udpAddr)
	if err1 != nil {
                fmt.Println("error listening on UDP port ", portNumber)
                fmt.Println(err1)
                os.Exit(1)
	}

	sendingObject := dataObject{myIp,[]string{}}
	jsonFile := struct2json(sendingObject)
	broadcastSocket.Write(jsonFile)
			
}

func struct2json(packageToSend dataObject) [] byte {
	jsonObject, _ := json.Marshal(packageToSend)
	return jsonObject
}

func json2struct(jsonObject []byte) dataObject{
	structObject := dataObject{}
	json.Unmarshal(jsonObject, &structObject)  //Her kan det være noe som ikke stemmer helt
	return structObject
}


func updateMasterQueue(recieveIpChan chan string) { //Leser inn ny ip fra channel. lager en temp kø lik nåværende MasterQueue. Sjekker om ny ip ligger i køen. // Hvis ikke legges den til i lista. Hvis det er deg selv settes connected = tru
	for {
		allreadyInQueue := false
		newQueueObject:= <- recieveIpChan
		for i,element:= range Masterqueue{
			if element == newQueueObject{
				allreadyInQueue = true
			}
		}
		if newQueueObject == myIp{
			connected = true
			
		}
		if !allreadyInQueue{
			MasterQueue = append(MasterQueue,newQueueObject)	
		}
	}	
}
	

func readServerToUpdateQueue(ipAddress string, portNumber string,recieveIpChan Chan string){ 		
	bufferToRead := make([] byte, 1024)
	UDPadr, err:= net.ResolveUDPAddr("udp",ipAddress+":"+portNumber)

	if err != nil {
                fmt.Println("error resolving UDP address on ", portNumber)
                fmt.Println(err)
                os.Exit(1)
        }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
            fmt.Println("error listening on UDP port ", portNumber)
            fmt.Println(err)
            os.Exit(1)
            
	}
	
	for {
		n,UDPadr, err := readerSocket.ReadFromUDP(bufferToRead)
        
	 	if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            os.Exit(1)
            
        }
        
        fmt.Println("got message from ", UDPadr, " with n = ", n)

       	if n > 0 {
           	fmt.Println("printer melding vi leste over UDP",json2struct(bufferToRead[0:n]))  
            structObject := json2struct(bufferToRead[0:n])
            ip = structObject.newIp
            recieveIpChan <- ip // Deadlock?
        }
   	}

}

func broadcastMasterQueue(ipBroadcast string,queuePort string){

}

