package networkModule
																
import (
    	"fmt"
    	"net"
    	"os"
    	"time"
	)

// Den nyoppstartede heisen får MasterQueue fra Master, og gjør seg selv til Master, Backup eller ingen av delene. Man blir ikke selv lagt til i MasterQueue enda

func InitializeElevator() {
	masterQueueLock <- 1 
	myIp = getIpAddr()
	fmt.Println("Dette er min IP",myIp)
	isEmpty := setMasterQueue(PORT_MASTER_DATA)
	<- masterQueueLock 
	if isEmpty{
		fmt.Println("jeg er Master")
		isMaster = true 
	}
	if len(masterQueue)==1{
		isBackup = true
		fmt.Println("jeg er Backup. Backup =",isBackup)
	}
	
	masterQueueLock <- 1

	fmt.Println("Heisen med Ip:",myIp,"er nå initialisert. Følgende variabler er satt","isMaster=",isMaster,"isBackup=",isBackup,"MasterQueue=",masterQueue)
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
	
func setMasterQueue(PORT_MASTER_DATA string) bool{	// må endre navna på portNumber osv i funksjonene nedover
	bufferToRead := make([] byte, 1024)
	UDPadr, err:= net.ResolveUDPAddr("udp",""+":"+PORT_MASTER_DATA)
	
	if err != nil {
                fmt.Println("error resolving UDP address on ", PORT_MASTER_DATA)
                fmt.Println(err)
                os.Exit(1)
        }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
            fmt.Println("error listening on UDP port ", PORT_MASTER_DATA)
            fmt.Println(err)
            os.Exit(1)
	}
	deadline := time.Now().Add(1500*time.Millisecond)
	readerSocket.SetReadDeadline(deadline)
	n,UDPadr,err := readerSocket.ReadFromUDP(bufferToRead[0:])
		
    readerSocket.Close()
 	if err != nil {
        return true
    }
   	if n > 0 {
       	//fmt.Println("printer melding vi leste over UDP",json2struct(bufferToRead,n))  
        structObject := DataObject{} 
        structObject = json2struct(bufferToRead,n)
        <- masterQueueLock
       	masterQueue = structObject.MasterQueue
       	masterQueueLock <- 1
       	return false
       		 
   }
   readerSocket.Close()
  	return true
}






