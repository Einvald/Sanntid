package networkModule

						
						//mangler broadcastMasterQueue og checkMaster. Disse bør muligens flyttes til runElevator?
						// Videre trenger vi å utvikle nettverksmodulen slik at heisene kan kommunisere mer direkte med hverandre. F.eks broaccaste bestillinger
						// 
				
import (
    	"fmt"
    	"net"
    	"os"
	)
type dataObject struct {
	newIp string
	masterQueue []string
	
func runElevator()
{
	isMaster,isBackup := connectElevator() //noe sånn
}


func connectElevator() (bool,bool){ // denne gjør kanskje mer enn å bare connecte? Hva med å returenere isMaster og is Backup og så bruke disse i en annen kode/package
	ReceiveChan:= make(chan string,1024)
	MasterQueue:= [] string {}  
	ipBroadcast:= "129.241.187.255"
	ipPort := "20017" //Her velges port som egen IP-adresse skal sendes og leses fra
	queuePort:= "20018" //Her velges port som køen skal sendes og motas fra
	myIp:= getIpAddr()
	isMaster:= false
	isBackup := false
	initializeElevator(queuePort,myIp) //får køen fra Master, gjør seg selv til eventuell master eller Backup 
	go readFromServer(ipBroadcast,ipPport,0) //Trenger egentlig ikke ipBroadcast her fordi den leser kun fra egen port // obs 0
	go updateMasterQueue() //Her må vi passe på at ip ikke legges inn dobbelt. Tror det skal være ordnet. Legges nå inn gjennom broadcastIP og updateM.Queue
	broadcastIp(myIp,ipBroadcast,ipPort) // Hva hvis den ikke kommer frem?
	if Master{
		broadcastMasterQueue(ipBroadcast, queuePort) 	//Denne må lages. Fungerer som imAlive

	}
	if Backup{
		checkMaster() // Denne må lages. Sjekker om Master broadcaster køen fortsatt. Hvis ikke blir man selv master og nestemann blir backup
	}
}



func getIpAddr() string {
	addrs, err := net.InterfaceAddrs()
	var streng string
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				streng = ipnet.IP.String()
															

			}
		}
	} 
	
	return streng
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


func updateMasterQueue(MasterQueue []string) {
	newQueueObject:= <- RecieveChan
	if newQueueObject != myIp{
		MasterQueue = append(MasterQueue,newQueueObject)
	}	
}
	

func readFromServer(ipAddress string, portNumber string,updateQueue bool, getQueue bool) bool{ 		
	bufferToRead := make([] byte, 1024)
	UDPadr, err:= net.ResolveUDPAddr("udp",ipAddress+":"+portNumber)

	if err != nil {
                fmt.Println("error resolving UDP address on ", portNumber)
                fmt.Println(err)
                return
        }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
            fmt.Println("error listening on UDP port ", portNumber)
            fmt.Println(err)
            return
	}
	
	for {
		n,UDPadr, err := readerSocket.ReadFromUDP(bufferToRead)
        
	 	if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            return
        }
        
        fmt.Println("got message from ", UDPadr, " with n = ", n)

       	if n > 0 {
           	fmt.Println("printer melding vi leste over UDP",json2struct(bufferToRead[0:n]))  
            structObject := json2struct(bufferToRead[0:n])
            if updateQueue{
            	ip = structObject.newIp
            	RecieveChan <- ip // Deadlock?

            }
           	if getQueue{
           		MasterQueue = structObject.masterQueue
           		return false
           	} 
        }else{
        	if getQueue{									///Obs! er vi sikre på at bufferet er tomt når lista er tom?
        		MasterQueue = append(MasterQueue,myIp)
        		return true
        	}
        }
	}
	return false
}
	

func initializeElevator(queuePort string, myIp string){
	isEmpty := readFromServer(myIp,queuePort,false,true) // Her settes hvilken versjon av readFrom som skal brukes
	if isEmpty{
		Master = true 
	}
	if len(MasterQueue)==1{
		Backup = true
	}
}

func broadcastMasterQueue(ipBroadcast string,queuePort string){

}


