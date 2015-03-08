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
var SetTimerChan = make(chan bool, 1)
var currentDirection = -1;
var currentFloor = 1

/*func readStopOrderChannel() int{
	select {
		case stopButtonOrder := 
	}
}
*/
var currentState State = RUN_DOWN;
//Må vurdere hvor nåværende state skal være lagret, men tenker at den skal være lagret her, også "last direction" kan være lurt å ha med
func FloorReached(floor int){
	fmt.Println("Inside FloorReached with current state: ");
	SetCurrentFloorLampChan <- floor;
	currentFloor = floor;
	switch currentState{
		case RUN_UP:
			if checkIfFloorInQueue(floor, currentDirection){
				SetMotorChan <- 0;
				currentState = DOOR_OPEN;
				SetTimerChan <- true;

			}
					
		case RUN_DOWN:
			if checkIfFloorInQueue(floor, currentDirection){
				SetMotorChan <- 0;
				currentState = DOOR_OPEN;
				fmt.Println("setter timer chan");
				SetTimerChan <- true;
			} 
		case DOOR_OPEN:
		case IDLE:
		case STOP:
			 
	}
}

func TimerOut() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:
				fmt.Println("Døren er lukket som den skulle")
				SetTimerChan <- false;
				currentDirection = nextDirection(currentFloor, currentDirection);
				SetMotorChan <- currentDirection;
				switch currentDirection{
				case 1:
					currentState = RUN_UP;
				case 0:
					currentState = IDLE;
				case -1: 
					currentState = RUN_DOWN;
				}
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
