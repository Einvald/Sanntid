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
	Floor int;
	ButtonType int;
	TurnOn int;
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
				
				currentState = DOOR_OPEN;
				SetTimerChan <- true;
				SetButtonLampChan <- ButtonLamp {floor, 0, 0};
				SetButtonLampChan <- ButtonLamp {floor, 2, 0};
				removeOrderFromQueue(floor, 0); //MÅ SENDE MELDING OM AT KØ FOR ALLE HEISER ER FERDIG
				removeOrderFromQueue(floor, 2);
				if !checkIfOrdersAtHigherFloors(floor){
					removeOrderFromQueue(floor, 0);
					SetButtonLampChan <- ButtonLamp {floor, 0, 0};
				}
				printQueues();
			}
					
		case RUN_DOWN:
			if checkIfFloorInQueue(floor, currentDirection){

				currentState = DOOR_OPEN;
				SetTimerChan <- true;
				SetButtonLampChan <- ButtonLamp {floor, 1, 0};
				SetButtonLampChan <- ButtonLamp {floor, 2, 0};
				if !checkIfOrdersAtLowerFloors(floor){
					removeOrderFromQueue(floor, 0);
					SetButtonLampChan <- ButtonLamp {floor, 0, 0};
				}
				removeOrderFromQueue(floor, 1); //MÅ SENDE MELDING OM AT KØ FOR ALLE HEISER ER FERDIG
				removeOrderFromQueue(floor, 2);
				printQueues();
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
				prevDirection := currentDirection;
				currentDirection = nextDirection(currentFloor, currentDirection);
				SetMotorChan <- currentDirection;
				switch currentDirection{
					case 1:
						currentState = RUN_UP;
						removeOrderFromQueue(currentFloor, 2);
						SetButtonLampChan <- ButtonLamp {currentFloor, 2, 0};
						if prevDirection != 1{
							removeOrderFromQueue(currentFloor, 0);
							SetButtonLampChan <- ButtonLamp {currentFloor, 0, 0};
						}

					case 0:
						if CheckIfEmptyQueues(){
							fmt.Println("Settes i IDLE")
							currentState = IDLE;
						} else {
							SetTimerChan <- true;
							EmptyQueues();
							SetButtonLampChan <- ButtonLamp {currentFloor, 0, 0};
							SetButtonLampChan <- ButtonLamp {currentFloor, 1, 0};
							SetButtonLampChan <- ButtonLamp {currentFloor, 2, 0};
						}
					case -1: 
						currentState = RUN_DOWN;
						removeOrderFromQueue(currentFloor, 2);
						SetButtonLampChan <- ButtonLamp {currentFloor, 2, 0};
						if prevDirection != -1{
							removeOrderFromQueue(currentFloor, 1);
							SetButtonLampChan <- ButtonLamp {currentFloor, 1, 0};
						}
				}
		case IDLE:
		case STOP:
	}
}

func StopButton() {
	StopButtonLampChan <- 1;
	//EmptyQueues();
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
			StopButtonLampChan <- 0;
			currentState = IDLE;
	}
}

func NewOrderInEmptyQueue() {
	fmt.Println("New Order In EMpty Queue")
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
					EmptyQueues()
					SetTimerChan <- true;
					currentState = DOOR_OPEN;
				case -1: 
					currentState = RUN_DOWN;
			}
		case STOP:
	}
}

func NewOrderToCurrentFloor() {
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
			SetTimerChan <- true;
		case IDLE:
			currentState = DOOR_OPEN;
			SetTimerChan <- true;

		case STOP:
			SetTimerChan <- true;
	}
}

func CheckIfFloorInQueue(floor int) bool{
	return checkIfFloorInQueue(floor, currentDirection);
}

func CheckIfCurrentFloor(floor int) bool {
	return (floor == currentFloor);
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


