package elev_handler

import(
	"fmt"
	)

type State int
const (
RUN_UP State = 0 + iota
RUN_DOWN
DOOR_OPEN
IDLE
STOP
)

/*
type FloorIndicatorLamp struct{
	floor int;
	turnOn bool;
}
*/

type ButtonLamp struct{
	floor int;
	buttonType int;
	turnOn bool;
}



var StopButtonLampChan = make(chan int);
var SetCurrentFloorLampChan = make(chan int);//FloorIndicatorLamp);
var SetButtonLampChan = make(chan ButtonLamp);
var SetMotorChan = make(chan int)
var currentDirection = 1;
var currentFloor = 0

/*func readStopOrderChannel() int{
	select {
		case stopButtonOrder := 
	}
}
*/
var currentState State = RUN_UP;
//Må vurdere hvor nåværende state skal være lagret, men tenker at den skal være lagret her, også "last direction" kan være lurt å ha med
func FloorReached(floor int){
	fmt.Println("Inside FloorReached with current state: ");

	
	switch currentState{
		case RUN_UP:
			SetCurrentFloorLampChan <- floor;
			if checkIfFloorInQueue(floor, currentDirection){
				SetMotorChan <- 0;
				currentState = DOOR_OPEN

			}
					
		case RUN_DOWN:
			SetCurrentFloorLampChan <- floor;
			if checkIfFloorInQueue(floor, currentDirection){
				SetMotorChan <- 0;
				currentState = DOOR_OPEN;
			} 
		case DOOR_OPEN:
			SetCurrentFloorLampChan <- floor;
		case IDLE:
			SetCurrentFloorLampChan <- floor; 
		case STOP:
			SetCurrentFloorLampChan <- floor; 
	}
}

func TimerOut() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:
		

		case STOP:

	}
}

func StopButton() {
	switch currentState{
		case RUN_UP:
			StopButtonLampChan <- 1;
			currentState = STOP;
		case RUN_DOWN:
			StopButtonLampChan <- 1;
			currentState = STOP;
		case DOOR_OPEN:
			StopButtonLampChan <- 1;
			currentState = STOP;
		case IDLE:
			StopButtonLampChan <- 1;
			currentState = STOP;

		case STOP:

	}
}
/*
func newOrderInEmptyQueue() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:

		case STOP:

	}
}

func newOrderToCurrentFloor() {currentState = STOP;
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:

		case STOP:

	}
}

func leftFloor() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:

		case STOP:

	}
}

func nextDirection() {

}
*/
