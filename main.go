package main

import(
	"fmt"
	d "driver"
	elev "elev_handler"
	"time"
	net "networkModule"
	)

const N_FLOORS int = 4
const N_BUTTONS int = 3
const ON = 1
const OFF = 0

/*
type ButtonLamp struct{
	floor int;
	buttonType int;
	turnOn int;
}
*/
func main() {
	Floor_sensor_channel := make(chan int, 1);
	Order_button_signal_channel := make(chan net.ButtonOrder);
	Order_data_to_master_channel := make(chan net.OrderData, 1024);
	Order_data_from_master_channel := make(chan net.OrderData, 1024);
	d.Driver_init()
	d.Driver_set_button_lamp(2, 2, 1);
	d.Driver_set_motor_direction(-1);
	elev.EmptyQueues();
	elev.AddToQueue(0, 2);
	go checkForInput(Floor_sensor_channel, Order_button_signal_channel)
	go handleInputChannels(Floor_sensor_channel, Order_button_signal_channel, Order_data_from_master_channel, Order_data_to_master_channel)
	//time.Sleep(100 * time.Millisecond);
	go handleElevatorCommands(Order_data_to_master_channel);
	elev.InitializeChanLocks();
	net.InitializeElevator()
	go net.RunElevator(Order_data_to_master_channel, Order_data_from_master_channel)
	
	deadChan := make(chan int);
	<- deadChan;
}

func handleInputChannels(Floor_sensor_channel chan int, Order_button_signal_channel chan net.ButtonOrder, Order_data_from_master_channel chan net.OrderData, Order_data_to_master_channel chan net.OrderData){
	for{
		select {
			case floor:= <-Floor_sensor_channel:
				elev.FloorReached(floor);
			case order:= <- Order_button_signal_channel:
				if order.ButtonType == 2 {d.Driver_set_button_lamp(order.ButtonType, order.Floor, ON);
					if elev.CheckIfCurrentFloor(order.Floor){
						elev.AddToQueue(order.Floor, order.ButtonType);
						elev.NewOrderToCurrentFloor();
					}else if elev.CheckIfEmptyQueues(){
						elev.AddToQueue(order.Floor, order.ButtonType)
						elev.NewOrderInEmptyQueue();
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
					if elev.CheckIfCurrentFloor(order.Floor){
						elev.AddToQueue(order.Floor, order.ButtonType);
						elev.NewOrderToCurrentFloor();
					}else if elev.CheckIfEmptyQueues(){
						elev.AddToQueue(order.Floor, order.ButtonType)
						elev.NewOrderInEmptyQueue();
					}else {
						elev.AddToQueue(order.Floor, order.ButtonType)
					}
				}else if orderData.Type == net.REQUEST_AUCTION{	
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, 1)
					cost := elev.GetCostForOrder(order.Floor, order.ButtonType);
					orderData := net.OrderData {false, net.COST, orderData.Order, cost, ""}
					Order_data_to_master_channel <- orderData;
				} else if orderData.Type == net.ORDER_COMPLETE {
					d.Driver_set_button_lamp(order.ButtonType, order.Floor, 0);
				}
			case <- time.After(25* time.Millisecond):
		}
	}
}


func checkForInput(Floor_sensor_channel chan int, Order_button_signal_channel chan net.ButtonOrder){
	floorSensored := 0;
	floorPushed := [N_FLOORS * N_BUTTONS] int{};
	for i := range floorPushed{floorPushed[i] = 0;}
	for {
		if d.Driver_get_floor_sensor_signal() != (-1) && floorSensored ==0 {
			fmt.Println("reached floor")
			if elev.CheckIfFloorInQueue(d.Driver_get_floor_sensor_signal()){
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
				if d.Driver_get_button_signal(buttonType, floorLevel) != 0 && floorPushed[(N_FLOORS-1)*floorLevel + buttonType] == 0{
					Order_button_signal_channel <- net.ButtonOrder{floorLevel, buttonType};
					floorPushed[(N_FLOORS-1)*floorLevel + buttonType] = 1;
				}
				if d.Driver_get_button_signal(buttonType, floorLevel) == 0 && floorPushed[(N_FLOORS-1)*floorLevel + buttonType] != 0{
					floorPushed[(N_FLOORS-1)*floorLevel + buttonType] = 0;
				}
			}
			
		}
		time.Sleep(30 * time.Millisecond)
	}
		
}

func handleElevatorCommands(Order_data_to_master_channel chan net.OrderData){
	timeChan := make(chan int64);
	go elev.DoorTimer(timeChan);
	for {
		select {
			case floor := <- elev.SetCurrentFloorLampChan:
				d.Driver_set_floor_indicator(floor);
			case direction := <- elev.SetMotorChan:
				fmt.Println("Endret motor retning til: ", direction)
				d.Driver_set_motor_direction(direction);
			case turnOn := <- elev.SetTimerChan:
					if turnOn{
						timeChan <- time.Now().UnixNano()/int64(time.Millisecond)
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
			//case <- time.After(25*time.Millisecond):
		}
	}
}