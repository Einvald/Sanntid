package main
import(
	d "driver"
	elev "elev_handler"
	"time"
	net "networkModule"
	)

const ON = 1
const OFF = 0

func main() {
	d.Driver_init()
	orderDataToMasterChannel := make(chan net.OrderData, 1024);
	orderDataFromMasterChannel := make(chan net.OrderData, 1024);
	runInternalElevator(orderDataToMasterChannel, orderDataFromMasterChannel)
	go net.RunNetworkCommunication(orderDataToMasterChannel, orderDataFromMasterChannel)	
	deadChan := make(chan int);
	<- deadChan;
}

func handleInput(floorSensorChannel chan int, orderButtonSignalChannel chan net.ButtonOrder, orderDataFromMasterChannel chan net.OrderData, orderDataToMasterChannel chan net.OrderData, currentElevValues elev.CurrentElevValues, orderQueueChannels elev.OrderQueueChannels, outputCommands elev.OutputChans) {
	for{
		select {
			case floor:= <-floorSensorChannel:
				elev.FloorReached(floor, currentElevValues, orderQueueChannels, outputCommands);
			case order:= <- orderButtonSignalChannel:
				if order.ButtonType == 2 {
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, ON);
					if elev.CheckIfCurrentFloor(order.Floor, currentElevValues.Floor, orderQueueChannels){
						elev.AddToQueue(order.Floor, order.ButtonType, orderQueueChannels);
						elev.NewOrderToCurrentFloor(currentElevValues.State, outputCommands.Timer);
					}else if elev.CheckIfEmptyQueues(orderQueueChannels){
						elev.AddToQueue(order.Floor, order.ButtonType, orderQueueChannels)
						elev.NewOrderInEmptyQueue(currentElevValues, orderQueueChannels, outputCommands);
					}else {
						elev.AddToQueue(order.Floor, order.ButtonType, orderQueueChannels)
					}
				} else {
					orderData := net.OrderData {false, net.ORDER, order, 0, " ", 0}
					orderDataToMasterChannel<-orderData;
				}
			case orderData := <- orderDataFromMasterChannel:
				order := orderData.Order;
				if orderData.Type == net.ORDER{
					if elev.CheckIfCurrentFloor(order.Floor, currentElevValues.Floor, orderQueueChannels){
						elev.AddToQueue(order.Floor, order.ButtonType, orderQueueChannels);
						elev.NewOrderToCurrentFloor(currentElevValues.State, outputCommands.Timer);
					}else if elev.CheckIfEmptyQueues(orderQueueChannels){
						elev.AddToQueue(order.Floor, order.ButtonType, orderQueueChannels)
						elev.NewOrderInEmptyQueue(currentElevValues, orderQueueChannels, outputCommands);
					}else {
						elev.AddToQueue(order.Floor, order.ButtonType, orderQueueChannels)
					}
				}else if orderData.Type == net.REQUEST_AUCTION{	
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, 1)
					cost := elev.GetCostForOrder(order.Floor, order.ButtonType, currentElevValues, orderQueueChannels);
					orderData := net.OrderData {false, net.COST, orderData.Order, cost, "", 0}
					orderDataToMasterChannel <- orderData;
				} else if orderData.Type == net.ORDER_COMPLETE {
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, 0);
				}
			case <- time.After(25* time.Millisecond):
		}
	}
}

func checkForInput(floorSensorChannel chan int, orderButtonSignalChannel chan net.ButtonOrder, currentDirection chan int, orderQueueChannels elev.OrderQueueChannels){
	floorSensored := 0;
	floorPushed := [elev.N_FLOORS * elev.N_BUTTONS] int{};
	for i := range floorPushed{floorPushed[i] = 0;}
	for {
		if d.Driver_get_floor_sensor_signal() != (-1) && floorSensored ==0 {
			if elev.CheckIfFloorInQueue(d.Driver_get_floor_sensor_signal(), currentDirection, orderQueueChannels){
					d.Driver_set_motor_direction(0);
			}
			floorSensored = 1;
			floorSensorChannel <- d.Driver_get_floor_sensor_signal();
		}
		if d.Driver_get_floor_sensor_signal() == (-1) && floorSensored !=0 {
			floorSensored = 0;
		}
		for floorLevel := 0; floorLevel<4 ; floorLevel++ {
			for buttonType := 0; buttonType<3; buttonType++{
				if d.Driver_get_button_signal(buttonType, floorLevel) != 0 && floorPushed[(elev.N_FLOORS-1)*floorLevel + buttonType] == 0{
					orderButtonSignalChannel <- net.ButtonOrder{floorLevel, buttonType};
					floorPushed[(elev.N_FLOORS-1)*floorLevel + buttonType] = 1;
				}
				if d.Driver_get_button_signal(buttonType, floorLevel) == 0 && floorPushed[(elev.N_FLOORS-1)*floorLevel + buttonType] != 0{
					floorPushed[(elev.N_FLOORS-1)*floorLevel + buttonType] = 0;
				}
			}	
		}
		time.Sleep(30 * time.Millisecond)
	}
}

func handleElevatorCommands(orderDataToMasterChannel chan net.OrderData, timerChan chan int64, outputCommands elev.OutputChans){
	for {
		select {
			case floor := <- outputCommands.FloorLamp:
				d.Driver_set_floor_indicator(floor);
			case direction := <- outputCommands.Motor:
				d.Driver_set_motor_direction(direction);
			case turnOn := <- outputCommands.Timer:
					if turnOn{
						timerChan <- time.Now().UnixNano()/int64(time.Millisecond)
						d.Driver_set_door_open_lamp(ON);
					} else{d.Driver_set_door_open_lamp(OFF);}
			case buttonOrder  := <- outputCommands.Button:
				order := elev.ButtonLamp {};
				order = buttonOrder
				if order.ButtonType == 2{d.Driver_set_button_lamp(order.ButtonType, order.Floor, order.TurnOn);}
			case finishedOrder := <- outputCommands.FinishedOrder:
				orderComplete := net.ButtonOrder{finishedOrder.Floor, finishedOrder.ButtonType}
				orderData := net.OrderData{false, net.ORDER_COMPLETE, orderComplete, 0, " ", 0}
				orderDataToMasterChannel <- orderData;
		}
	}
}

func runInternalElevator(orderDataToMasterChannel chan net.OrderData, orderDataFromMasterChannel chan net.OrderData){
	floorSensorChannel := make(chan int, 1);
	orderButtonSignalChannel := make(chan net.ButtonOrder);
	currentDirection := make(chan int, 1);
	currentFloor := make(chan int, 1);
	currentState := make(chan elev.State, 1);
	currentElevValues := elev.CurrentElevValues {currentDirection, currentFloor, currentState}
	setCurrentFloorLampChan := make(chan int);
	setButtonLampChan := make(chan elev.ButtonLamp);
	setMotorChan := make(chan int);
	setTimerChan := make(chan bool, 1);
	finishedOrderChan := make (chan elev.ButtonOrder, 1024);
	outputCommands := elev.OutputChans {setCurrentFloorLampChan, setButtonLampChan, setMotorChan, setTimerChan, finishedOrderChan}
	timerChan := make(chan int64);
	queueUpChan := make(chan [elev.N_FLOORS] int, 1)
	queueDownChan := make(chan [elev.N_FLOORS] int, 1)
	queueInElevChan := make(chan [elev.N_FLOORS] int, 1)
	orderQueueChannels := elev.OrderQueueChannels {queueUpChan, queueDownChan, queueInElevChan}
	go checkForInput(floorSensorChannel, orderButtonSignalChannel, currentDirection, orderQueueChannels)
	go handleInput(floorSensorChannel, orderButtonSignalChannel, orderDataFromMasterChannel, orderDataToMasterChannel, currentElevValues, orderQueueChannels, outputCommands)
	go handleElevatorCommands(orderDataToMasterChannel, timerChan, outputCommands);
	go elev.DoorTimer(timerChan, currentElevValues, orderQueueChannels, outputCommands);
	initializeSystem(currentElevValues, orderQueueChannels);
}

func initializeSystem(currentElevValues elev.CurrentElevValues, orderQueueChannels elev.OrderQueueChannels){
	currentElevValues.State <- elev.RUN_DOWN;
	currentElevValues.Floor <- 0;
	currentElevValues.Direction <- -1;
	d.Driver_set_motor_direction(-1);
	elev.InitializeQueues(orderQueueChannels)
	elev.AddToQueue(0, 2, orderQueueChannels);
	d.Driver_set_button_lamp(2, 0, 1);
}