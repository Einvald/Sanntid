package main

import(
	"fmt"
	d "driver"
	elev "elev_handler"
	"time"
	net "networkModule"
	)

const ON = 1
const OFF = 0

func main() {
	init := d.Driver_init()
	fmt.Println(init)
	Floor_sensor_channel := make(chan int, 1);
	Order_button_signal_channel := make(chan net.ButtonOrder);
	Order_data_to_master_channel := make(chan net.OrderData, 1024);
	Order_data_from_master_channel := make(chan net.OrderData, 1024);
	CurrentDirection := make(chan int, 1);
	CurrentFloor := make(chan int, 1);
	CurrentState := make(chan elev.State, 1);
	timerChan := make(chan int64);

	//d.Driver_init()
	//d.Driver_set_button_lamp(2, 2, 1);
	//d.Driver_set_motor_direction(-1);
	//elev.EmptyQueues();
	//elev.AddToQueue(0, 2);
	//initializeSystem(CurrentDirection, CurrentFloor, CurrentState);
	//go checkForInput(Floor_sensor_channel, Order_button_signal_channel, CurrentDirection, CurrentFloor, CurrentState)
	//go handleInput(Floor_sensor_channel, Order_button_signal_channel, Order_data_from_master_channel, Order_data_to_master_channel, CurrentDirection, CurrentFloor, CurrentState)
	//time.Sleep(100 * time.Millisecond);
	//go handleElevatorCommands(Order_data_to_masannter_channel, timerChan);
	//elev.InitializeChanLocks();
	//net.InitializeElevator()
	go checkForInput(Floor_sensor_channel, Order_button_signal_channel, CurrentDirection, CurrentFloor, CurrentState)
	go handleInput(Floor_sensor_channel, Order_button_signal_channel, Order_data_from_master_channel, Order_data_to_master_channel, CurrentDirection, CurrentFloor, CurrentState)
	go handleElevatorCommands(Order_data_to_master_channel, timerChan);
	go elev.DoorTimer(timerChan, CurrentDirection, CurrentFloor, CurrentState);
	
	
	//arrayChan := make(chan [4] int, 100)
	initializeSystem(CurrentDirection, CurrentFloor, CurrentState);
	

	go net.RunElevator(Order_data_to_master_channel, Order_data_from_master_channel)
	
	deadChan := make(chan int);
	<- deadChan;
}

func initializeSystem(CurrentDirection chan int, CurrentFloor chan int, CurrentState chan elev.State){
	CurrentState <- elev.RUN_DOWN;
	CurrentFloor <- 0;
	CurrentDirection <- -1;

	d.Driver_set_motor_direction(-1);
	elev.EmptyQueues();
	elev.AddToQueue(0, 2);
	d.Driver_init()
	d.Driver_set_button_lamp(2, 0, 1);
	net.InitializeElevator()
}

func handleInput(Floor_sensor_channel chan int, Order_button_signal_channel chan net.ButtonOrder, Order_data_from_master_channel chan net.OrderData, Order_data_to_master_channel chan net.OrderData, CurrentDirection chan int, CurrentFloor chan int, CurrentState chan elev.State){
	for{
		select {
			case floor:= <-Floor_sensor_channel:
				elev.FloorReached(floor, CurrentDirection, CurrentFloor, CurrentState);
			case order:= <- Order_button_signal_channel:
				if order.ButtonType == 2 {d.Driver_set_button_lamp(order.ButtonType, order.Floor, ON);
					if elev.CheckIfCurrentFloor(order.Floor, CurrentFloor){
						elev.AddToQueue(order.Floor, order.ButtonType);
						elev.NewOrderToCurrentFloor(CurrentState);
					}else if elev.CheckIfEmptyQueues(){
						elev.AddToQueue(order.Floor, order.ButtonType)
						elev.NewOrderInEmptyQueue(CurrentDirection, CurrentFloor, CurrentState);
					}else {
						elev.AddToQueue(order.Floor, order.ButtonType)
					}
				} else {
					orderData := net.OrderData {false, net.ORDER, order, 0, " "}
					Order_data_to_master_channel<-orderData;
				}
			case orderData := <- Order_data_from_master_channel:
				order := orderData.Order;
				if orderData.Type == net.ORDER{
					if elev.CheckIfCurrentFloor(order.Floor, CurrentFloor){
						elev.AddToQueue(order.Floor, order.ButtonType);
						elev.NewOrderToCurrentFloor(CurrentState);
					}else if elev.CheckIfEmptyQueues(){
						elev.AddToQueue(order.Floor, order.ButtonType)
						elev.NewOrderInEmptyQueue(CurrentDirection, CurrentFloor, CurrentState);
					}else {
						elev.AddToQueue(order.Floor, order.ButtonType)
					}
				}else if orderData.Type == net.REQUEST_AUCTION{	
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, 1)
					cost := elev.GetCostForOrder(order.Floor, order.ButtonType, CurrentDirection, CurrentFloor, CurrentState);
					orderData := net.OrderData {false, net.COST, orderData.Order, cost, ""}
					Order_data_to_master_channel <- orderData;
				} else if orderData.Type == net.ORDER_COMPLETE {
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, 0);
				}
			case <- time.After(25* time.Millisecond):
		}
	}
}


func checkForInput(Floor_sensor_channel chan int, Order_button_signal_channel chan net.ButtonOrder, CurrentDirection chan int, CurrentFloor chan int, CurrentState chan elev.State){
	floorSensored := 0;
	floorPushed := [elev.N_FLOORS * elev.N_BUTTONS] int{};
	for i := range floorPushed{floorPushed[i] = 0;}
	for {
		if d.Driver_get_floor_sensor_signal() != (-1) && floorSensored ==0 {
			fmt.Println("reached floor")
			if elev.CheckIfFloorInQueue(d.Driver_get_floor_sensor_signal(), CurrentDirection ){
					d.Driver_set_motor_direction(0);
			}
			floorSensored = 1;
			Floor_sensor_channel <- d.Driver_get_floor_sensor_signal();
		}
		if d.Driver_get_floor_sensor_signal() == (-1) && floorSensored !=0 {
			floorSensored = 0;
		}
		for floorLevel := 0; floorLevel<4 ; floorLevel++ {
			for buttonType := 0; buttonType<3; buttonType++{
				if d.Driver_get_button_signal(buttonType, floorLevel) != 0 && floorPushed[(elev.N_FLOORS-1)*floorLevel + buttonType] == 0{
					Order_button_signal_channel <- net.ButtonOrder{floorLevel, buttonType};
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

func handleElevatorCommands(Order_data_to_master_channel chan net.OrderData, timerChan chan int64){
	
	for {
		select {
			case floor := <- elev.SetCurrentFloorLampChan:
				d.Driver_set_floor_indicator(floor);
			case direction := <- elev.SetMotorChan:
				fmt.Println("Endret motor retning til: ", direction)
				d.Driver_set_motor_direction(direction);
			case turnOn := <- elev.SetTimerChan:
					if turnOn{
						timerChan <- time.Now().UnixNano()/int64(time.Millisecond)
						d.Driver_set_door_open_lamp(ON);
					} else{d.Driver_set_door_open_lamp(OFF);}
			case buttonOrder  := <- elev.SetButtonLampChan:
				order := elev.ButtonLamp {0, 0, 0};
				order = buttonOrder
				if order.ButtonType == 2{d.Driver_set_button_lamp(order.ButtonType, order.Floor, order.TurnOn);}
			case finishedOrder := <- elev.FinishedOrderChan:
				orderComplete := net.ButtonOrder{finishedOrder.Floor, finishedOrder.ButtonType}
				orderData := net.OrderData{false, net.ORDER_COMPLETE, orderComplete, 0, " "}
				Order_data_to_master_channel <- orderData;
		}
	}
}