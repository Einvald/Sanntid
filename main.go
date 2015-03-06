package main

import(
	"fmt"
	d "driver"
	elev "elev_handler"
	"time"
	)

const N_FLOORS int = 4
const N_BUTTONS int = 3

type buttonOrder struct{
	floor int;
	buttonType int;
}

func main() {
	Floor_sensor_channel := make(chan int, 1);
	Stop_button_signal_channel := make(chan bool,1);
	Order_button_signal_channel := make(chan buttonOrder);
	i := d.Driver_init()
	//d.Driver_set_button_lamp(0, 0, 1);
	fmt.Println(i)
	fmt.Println(elev.STOP)
	go checkForInput(Floor_sensor_channel, Stop_button_signal_channel, Order_button_signal_channel)
	go readFromChannels(Floor_sensor_channel, Stop_button_signal_channel, Order_button_signal_channel)
	time.Sleep(100 * time.Millisecond)
	Floor_sensor_channel <- d.Driver_get_floor_sensor_signal()
	deadChan := make(chan int);
	<- deadChan;
	
}

func readFromChannels(Floor_sensor_channel chan int, Stop_button_signal_channel chan bool, Order_button_signal_channel chan buttonOrder){
	for{
		select {

			case floor_sensor:= <-Floor_sensor_channel:
				elev.floorReached(floor_sensor);
			case stop_sensor := <- Stop_button_signal_channel:
				fmt.Println("stop button pushed", stop_sensor);
			case floor_order:= <- Order_button_signal_channel:
				fmt.Println("bestilt", floor_order.floor, floor_order.buttonType);
			case <- time.After(400* time.Millisecond):
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

/*func readFromNetwork() {
	
}
*/