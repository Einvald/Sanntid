package networkModule
																
import (
    	"fmt"
    	"net"
    	"os"
	)



func connectElevator(ipBroadcast string, portMasterQueue string) { //Her settes MasterQueue, isMaster,isBackup og connected
	
	MyIp:= getIpAddr()
	initializeElevator(ipBroadcast,portMasterQueue) //får køen fra Master, gjør seg selv til eventuell master eller Backup 
	
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
	
func setMasterQueue(ipBroadcast string,portMasterQueue string) bool{	// må endre navna på portNumber osv i funksjonene nedover
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
       
       	MasterQueue =append(MasterQueue,structObject.masterQueue)
       	return false
       		
       	} 
    }else{
    	return true
    }
}

func initializeElevator(ipBroadcast string,portMasterQueue string){

	isEmpty := setMasterQueue(ipBroadcast,portMasterQueue) 
	if isEmpty{
		IsMaster = true 
	}
	if len(MasterQueue)==1{
		IsBackup = true
	}
}




