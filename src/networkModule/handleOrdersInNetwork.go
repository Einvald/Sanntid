package networkModule


type OrderData struct 
	Type MessageType
	Cost int 
	Order buttonOrder
	OrderComplete bool
	Ip string

type MessageType int
	const(
	ORDER messageType = 0 + iota
	COST 
	ORDER_COMPLETE
	REQUEST_AUCTION
	)

var unfinishedOrders = [] orderData{} 
var recievedMessageToMaster = make(chan OrderData, 1024)
var recievedMessage = make (chan OrderData,1024)
var AuctionResultChan = make(chan OrderData,1)
var recievedCostChan = make(chan OrderData,1024)
var recievedOrderChan =make(chan buttonOrder,1024)
var recievedOrderComplete = make(chan buttonOrder, 1024)
var auctionLock = make(chan int, 1)
func handleOrdersInNetwork(){
	for {
		select{
		case recievedData := <- recievedMessageToMaster:
			switch recievedData.Type
				case COST:
					recievedCostChan <- recievedData
				case ORDER:
					recievedOrderChan <- recievedData.Order 
				case ORDER_COMPLETE:
					recievedOrderComplete <- recievedData.Order

		case order := <- recivedOrderChan:
			// Sjekk om finnes i liste. Legg til hvis ikke
		case ORDER_COMPLETE:
			//Fjern fra liste
			// Sett bit på channel
		}
	}

}

func handleOrdersFromMaster(LEGG TIL EN CHANNEL){
	for {
		select{
		case recievedData := <- recievedMessage:
			switch recievedData.Type{
				case ORDER:
					//Legges på channel som leses av main

				case REQUEST_AUCTION:
					// Legges på en annen channel som leses av main
			}
			
	}	}


}
//Funksjoner master bruker		
func removeFinishedOrder(){

}



func requestAuction(){}

func DecideAuctionWinner(){
	//Sjekker beste kosten av det som mottas
	//Sjekk at alle har svart
	//Send endelig bestilling til valgt heis
	// Lagre i liste 
	// Legg til timeout på 0.2 sek ish på hvor lenge du vil vente på svar

}
 
// Lag funksjonalitet for døde heiser, hvordan dens bestillinger skal fordeles osv.

func broadcastOrder()

//Funksjoner som alle bruker
func readOrderData(){
	//Hvis til Master legges det på masterkanalen, hvis ikke på den andre kanalen

}

func sendOrderData(){
	// Bruk TCP, da vet vi at alt blir sendt og mottatt
	
}
func readMessageFromMaster(){

}


func struct2json(packageToSend OrderData) [] byte {
	jsonObject, _ := json.Marshal(packageToSend)
	return jsonObject
}

func json2struct(jsonObject []byte,n int) OrderData{
	structObject := OrderData{}
	json.Unmarshal(jsonObject[0:n], &structObject)  
	return structObject
}

//Spør om lys skal lyse i alle på alle heispaneler hvis man bestemmer for eksempel ned fra 2. etasje?

