package main

import(
	"fmt"
	d "driver"
	elev "elev_handler"
	"time"
	)

const N_FLOORS int = 4
const N_BUTTONS int = 3
const ON = 1
const OFF = 0

type buttonOrder struct{
	floor int;
	buttonType int;
}

/*
type ButtonLamp struct{
	floor int;
	buttonType int;
	turnOn int;
}
*/

func main() {
	Floor_sensor_channel := make(chan int, 1);
	Stop_button_signal_channel := make(chan bool,1);
	Order_button_signal_channel := make(chan buttonOrder);
	i := d.Driver_init()
	d.Driver_set_button_lamp(2, 1, 1);
	d.Driver_set_motor_direction(-1);
	elev.EmptyQueues();
	elev.AddToQueue(1, 2);
	fmt.Println(i)
	fmt.Println(elev.STOP)
	go checkForInput(Floor_sensor_channel, Stop_button_signal_channel, Order_button_signal_channel)
	go readFromInputChannels(Floor_sensor_channel, Stop_button_signal_channel, Order_button_signal_channel)
	time.Sleep(100 * time.Millisecond);
	go handleElevatorCommands();
	
	deadChan := make(chan int);
	<- deadChan;
}

func readFromInputChannels(Floor_sensor_channel chan int, Stop_button_signal_channel chan bool, Order_button_signal_channel chan buttonOrder){
	for{
		select {

			case floor:= <-Floor_sensor_channel:
				
				fmt.Println("FloorReached read from channel")
				elev.FloorReached(floor);

			case stop_sensor := <- Stop_button_signal_channel:
				elev.StopButton();
				fmt.Println("stop button pushed", stop_sensor);
			case order:= <- Order_button_signal_channel:
				
				d.Driver_set_button_lamp(order.buttonType, order.floor, ON);
				if elev.CheckIfCurrentFloor(order.floor){
					elev.AddToQueue(order.floor, order.buttonType);
					elev.NewOrderToCurrentFloor();
				}else if elev.CheckIfEmptyQueues(){
					elev.AddToQueue(order.floor, order.buttonType)
					elev.NewOrderInEmptyQueue();
				}else {
					elev.AddToQueue(order.floor, order.buttonType)
				}

				
			case <- time.After(10* time.Millisecond):
		}
	}
}

func checkForInput(Floor_sensor_channel chan int, Stop_button_signal_channel chan bool, Order_button_signal_channel chan buttonOrder){
	floorSensored := 0;
	stopPushed := 0;
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
		if d.Driver_get_stop_signal() && stopPushed == 0{
			Stop_button_signal_channel <- true;
			stopPushed = 1;
		}
		if !d.Driver_get_stop_signal() && stopPushed != 0{
			stopPushed = 0;
		}
		for floorLevel := 0; floorLevel<4 ; floorLevel++ {
			for buttonType := 0; buttonType<3; buttonType++{
				if d.Driver_get_button_signal(buttonType, floorLevel) != 0 && floorPushed[(N_FLOORS-1)*floorLevel + buttonType] == 0{
					Order_button_signal_channel <- buttonOrder{floorLevel, buttonType};
					floorPushed[(N_FLOORS-1)*floorLevel + buttonType] = 1;
				}
				if d.Driver_get_button_signal(buttonType, floorLevel) == 0 && floorPushed[(N_FLOORS-1)*floorLevel + buttonType] != 0{
					floorPushed[(N_FLOORS-1)*floorLevel + buttonType] = 0;
				}
			}

		}

	}
		
}

func handleElevatorCommands(){
	timeChan := make(chan int64);
	go elev.DoorTimer(timeChan);
	for {
		select {
		case stopButtonLamp := <- elev.StopButtonLampChan:
			d.Driver_set_stop_lamp(stopButtonLamp);
		case floor := <- elev.SetCurrentFloorLampChan:
			fmt.Println("Setter floor indicator lyset")
			d.Driver_set_floor_indicator(floor);
		case direction := <- elev.SetMotorChan:
			fmt.Println("Endret motor retning til: ", direction)
			d.Driver_set_motor_direction(direction);
		case turnOn := <- elev.SetTimerChan:
				
				if turnOn{
					fmt.Println("Lest fra SetTimerChan")
					timeChan <- time.Now().UnixNano()/int64(time.Millisecond)
					//go elev.DoorTimer(); //OBS!!! ER DET GREIT Å KJØRE EGEN TRÅD HER?
					d.Driver_set_door_open_lamp(1);
				} else{d.Driver_set_door_open_lamp(0);}
		case buttonOrder  := <- elev.SetButtonLampChan:
			order := elev.ButtonLamp {0, 0, 0};
			order = buttonOrder
			d.Driver_set_button_lamp(order.ButtonType, order.Floor, order.TurnOn);
		case <- time.After(10*time.Millisecond):

		}
	}
}

/*func readFromNetwork() {
	
}
*/