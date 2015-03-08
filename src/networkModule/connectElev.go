package networkModule
																
import (
    	"fmt"
    	"net"
    	"os"
    	"time"
	)

// Den nyoppstartede heisen får MasterQueue fra Master, og gjør seg selv til Master, Backup eller ingen av delene. Man blir ikke selv lagt til i MasterQueue enda

func InitializeElevator() { 
	portMasterQueue := "20019"
	MyIp = getIpAddr()
	fmt.Println("2her kommer min IP",MyIp)
	isEmpty := setMasterQueue(portMasterQueue) 
	if isEmpty{
		fmt.Println("jeg er Master")
		IsMaster = true 
	}
	fmt.Println("lengden på MasterQueue er",len(MasterQueue))
	if len(MasterQueue)==1{
		IsBackup = true
		fmt.Println("jeg er Backup")
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
	
func setMasterQueue(portMasterQueue string) bool{	// må endre navna på portNumber osv i funksjonene nedover
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
	deadline := time.Now().Add(1500*time.Millisecond)
	readerSocket.SetReadDeadline(deadline)
	n,UDPadr,err := readerSocket.ReadFromUDP(bufferToRead[0:])
		
    readerSocket.Close()
 	if err != nil {
        return true
    }
    
    fmt.Println("got message from ", UDPadr, " with n = ", n)

   	if n > 0 {
       	fmt.Println("printer melding vi leste over UDP",json2struct(bufferToRead,n))  
        structObject := DataObject{} 
        structObject = json2struct(bufferToRead,n)
       	MasterQueue = structObject.MasterQueue
       	return false
       		 
   }
   readerSocket.Close()
  	return true
}






