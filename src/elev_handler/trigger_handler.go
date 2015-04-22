package elev_handler

import(
	//"fmt"
	)

type State int
const (
RUN_UP State = 0 + iota
RUN_DOWN
DOOR_OPEN
IDLE
)

type ButtonOrder struct{
	Floor int;
	ButtonType int; 
}
type ButtonLamp struct{
	Floor int;
	ButtonType int;
	TurnOn int;
}


var SetCurrentFloorLampChan = make(chan int);
var SetButtonLampChan = make(chan ButtonLamp);
var SetMotorChan = make(chan int)
var SetTimerChan = make(chan bool, 1)

func FloorReached(floor int, CurrentDirection chan int, CurrentFloor chan int, CurrentState chan State){
	if floor<0{floor = floor*(-1)}
	SetCurrentFloorLampChan <- floor;
	<- CurrentFloor; CurrentFloor <- floor;
	currentState := <- CurrentState; CurrentState <- currentState;

	switch currentState{
		case RUN_UP:
			if CheckIfFloorInQueue(floor, CurrentDirection){ 
				
				<-CurrentState; CurrentState <- DOOR_OPEN;
				SetTimerChan <- true;
				SetButtonLampChan <- ButtonLamp {floor, 2, 0};
				removeOrderFromQueue(floor, 0); 
				removeOrderFromQueue(floor, 2);
				if !checkIfOrdersAtHigherFloors(floor){
					removeOrderFromQueue(floor, 1);
				}
				printQueues();
			}
					
		case RUN_DOWN:
			if CheckIfFloorInQueue(floor, CurrentDirection){ 

				<-CurrentState; CurrentState <- DOOR_OPEN;
				SetTimerChan <- true;
				SetButtonLampChan <- ButtonLamp {floor, 2, 0};
				if !checkIfOrdersAtLowerFloors(floor){
					removeOrderFromQueue(floor, 0);
				}
				removeOrderFromQueue(floor, 1);
				removeOrderFromQueue(floor, 2);
				printQueues();
			} 
		case DOOR_OPEN:
		case IDLE:

			 
	}
}

func TimerOut(CurrentDirection chan int, CurrentFloor chan int, CurrentState chan State) {
	currentState := <- CurrentState; CurrentState <- currentState;
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
				SetTimerChan <- false;
				currentFloor := <- CurrentFloor;
				CurrentFloor <- currentFloor;
				prevDirection := <- CurrentDirection;
				currentDirection := nextDirection(currentFloor, prevDirection); 
				CurrentDirection <- currentDirection;
				SetMotorChan <- currentDirection;
				switch currentDirection{
					case 1:
						<-CurrentState; CurrentState <- RUN_UP;
						removeOrderFromQueue(currentFloor, 2);
						SetButtonLampChan <- ButtonLamp {currentFloor, 2, 0};
						if prevDirection != 1{
							removeOrderFromQueue(currentFloor, 0);
						}

					case 0:
								removeOrderFromQueue(currentFloor, 0)
								removeOrderFromQueue(currentFloor, 1)
								removeOrderFromQueue(currentFloor, 2)
								SetButtonLampChan <- ButtonLamp {currentFloor, 2, 0};
								<-CurrentState; CurrentState <- IDLE;
					case -1: 
						<-CurrentState; CurrentState <- RUN_DOWN;
						removeOrderFromQueue(currentFloor, 2);
						SetButtonLampChan <- ButtonLamp {currentFloor, 2, 0};
						if prevDirection != -1{
							removeOrderFromQueue(currentFloor, 1);
						}
				}
		case IDLE:
	}
}

func NewOrderInEmptyQueue(CurrentDirection chan int, CurrentFloor chan int, CurrentState chan State) {
	printQueues();
	currentState := <- CurrentState; CurrentState <- currentState;
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
		case IDLE:
			currentFloor := <- CurrentFloor; CurrentFloor <- currentFloor;
			currentDirection := <-CurrentDirection;
			currentDirection = nextDirection(currentFloor, currentDirection);
			CurrentDirection <- currentDirection;
			SetMotorChan <- currentDirection;
			switch currentDirection{
				case 1:
					<-CurrentState; CurrentState <- RUN_UP;

				case 0:
					EmptyQueues()
					SetTimerChan <- true;
					<-CurrentState; CurrentState <- DOOR_OPEN;
				case -1: 
					<-CurrentState; CurrentState <- RUN_DOWN;
			}
	}
}

func NewOrderToCurrentFloor(CurrentState chan State) {
	currentState := <- CurrentState; CurrentState <- currentState;
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
			SetTimerChan <- true;
		case IDLE:
			<- CurrentState; CurrentState <- DOOR_OPEN; 
			SetTimerChan <- true;
	}
}

func CheckIfCurrentFloor(floor int, CurrentFloor chan int) bool {
	currentFloor := <- CurrentFloor; CurrentFloor <- currentFloor;
	return (floor == currentFloor);
}