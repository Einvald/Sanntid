package networkModule
																
import (
    	"net"
    	"time"
    	"fmt"
	)

func initializeElevatorConnection(myIp string,masterQueueChan chan []IpObject, becomeMasterChan chan int,becomeBackupChan chan int) {
	isEmpty := setMasterQueue(PORT_MASTER_DATA,masterQueueChan)
	if isEmpty{
		becomeMasterChan <- 1 
	}
	masterQueue := <- masterQueueChan; 
	if len(masterQueue)==1{
		becomeBackupChan <- 1
	}
	masterQueueChan <- masterQueue 
	fmt.Println("The elevator with ip:",myIp,"has now been initialized")
}
func getIpAddr() string {
	addrs, _ := net.InterfaceAddrs()
	var ipString string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipString = ipnet.IP.String()															
			}
		}
	} 	
	return ipString
}
	
func setMasterQueue(PORT_MASTER_DATA string,masterQueueChan chan []IpObject ) bool{	
	bufferToRead := make([] byte, 1024)
	UDPadr,_:= net.ResolveUDPAddr("udp",""+":"+PORT_MASTER_DATA)
    readerSocket ,_ := net.ListenUDP("udp",UDPadr)
	deadline := time.Now().Add(1500*time.Millisecond)
	readerSocket.SetReadDeadline(deadline)
	n,UDPadr,timeout := readerSocket.ReadFromUDP(bufferToRead[0:])	
    readerSocket.Close()
 	if timeout != nil {
        return true
    }
   	if n > 0 {  
        structObject := DataObject{} 
        structObject = json2struct(bufferToRead,n)
        <- masterQueueChan
       	masterQueue := structObject.MasterQueue
       	masterQueueChan <- masterQueue
       	return false       		 
   }
   readerSocket.Close()
  	return true
}