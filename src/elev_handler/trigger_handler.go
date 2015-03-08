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
	turnOn int;
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

func FloorReached(floor int){
	fmt.Println("Inside FloorReached with current state: ", currentState);
	SetCurrentFloorLampChan <- floor;
	currentFloor = floor;
	switch currentState{
		case RUN_UP:
			if checkIfFloorInQueue(floor, currentDirection){
				SetMotorChan <- 0;
				currentState = DOOR_OPEN;
				SetTimerChan <- true;
				SetButtonLampChan <- ButtonLamp {floor, 0, 0};
				SetButtonLampChan <- ButtonLamp {floor, 2, 0};
				removeOrderFromQueue(floor, 0); //MÅ SENDE MELDING OM AT KØ FOR ALLE HEISER ER FERDIG
				removeOrderFromQueue(floor, 2);
			}
					
		case RUN_DOWN:
			if checkIfFloorInQueue(floor, currentDirection){
				SetMotorChan <- 0;
				currentState = DOOR_OPEN;
				fmt.Println("setter timer chan");
				SetTimerChan <- true;
				removeOrderFromQueue(floor, 1); //MÅ SENDE MELDING OM AT KØ FOR ALLE HEISER ER FERDIG
				removeOrderFromQueue(floor, 2);
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
					if EmptyQueues(){
						currentState = IDLE;
					} else {
						SetTimerChan <- true;
					}
				case -1: 
					currentState = RUN_DOWN;
				}
		case IDLE:
		case STOP://Hvordan skal døren oppføre seg i stop-situasjoner? 
	}
}

func StopButton() {
	StopButtonLampChan <- 1;
	EmptyQueues();
	currentState = STOP;
	switch currentState{
		case RUN_UP:
			SetMotorChan <- 0;
			currentDirection = 0;
		case RUN_DOWN:
			SetMotorChan <- 0;
			currentDirection = 0;
		case DOOR_OPEN:
		case IDLE:
		case STOP:
	}
}

func newOrderInEmptyQueue() {
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
		case IDLE:
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
		case STOP:
	}
}

func newOrderToCurrentFloor() {currentState = STOP;
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
		case IDLE:
			currentState = DOOR_OPEN;
			SetTimerChan <- true;

		case STOP:
			SetTimerChan <- true;
	}
}
/*
func leftFloor() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:

		case STOP:

	}
}
*/


