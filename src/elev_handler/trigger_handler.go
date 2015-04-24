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

type CurrentElevValues struct{
	Direction chan int;
	Floor chan int; 
	State chan State;
}

type OutputChans struct{
	FloorLamp chan int;
	Button chan ButtonLamp;
	Motor chan int;
	Timer chan bool;
	FinishedOrder chan ButtonOrder;
}

func FloorReached(floor int,  values CurrentElevValues, queues OrderQueueChannels, output OutputChans){
	if floor<0{floor = floor*(-1)}
	output.FloorLamp <- floor;
	<- values.Floor; values.Floor <- floor;
	currentState := <- values.State; values.State <- currentState;
	switch currentState{
		case RUN_UP:
			if CheckIfFloorInQueue(floor, values.Direction, queues){ 
				<-values.State; values.State <- DOOR_OPEN;
				output.Timer <- true;
				output.Button <- ButtonLamp {floor, 2, 0};
				removeOrderFromQueue(floor, 0, queues, output.FinishedOrder); 
				removeOrderFromQueue(floor, 2, queues, output.FinishedOrder);
				if !checkIfOrdersAtHigherFloors(floor, queues){
					removeOrderFromQueue(floor, 1, queues, output.FinishedOrder);
				}
				printQueues(queues);
			}
					
		case RUN_DOWN:
			if CheckIfFloorInQueue(floor, values.Direction, queues){ 

				<-values.State; values.State <- DOOR_OPEN;
				output.Timer <- true;
				output.Button <- ButtonLamp {floor, 2, 0};
				if !checkIfOrdersAtLowerFloors(floor, queues){
					removeOrderFromQueue(floor, 0, queues, output.FinishedOrder);
				}
				removeOrderFromQueue(floor, 1, queues, output.FinishedOrder);
				removeOrderFromQueue(floor, 2, queues, output.FinishedOrder);
				printQueues(queues);
			} 
		case DOOR_OPEN:
		case IDLE:
	}
}

func TimerOut(values CurrentElevValues, queues OrderQueueChannels, output OutputChans) {
	currentState := <- values.State; values.State <- currentState;
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
				output.Timer <- false;
				currentFloor := <- values.Floor;
				values.Floor <- currentFloor;
				prevDirection := <- values.Direction;
				currentDirection := nextDirection(currentFloor, prevDirection, queues); 
				values.Direction <- currentDirection;
				output.Motor <- currentDirection;
				switch currentDirection{
					case 1:
						<-values.State; values.State <- RUN_UP;
						removeOrderFromQueue(currentFloor, 2, queues, output.FinishedOrder);
						output.Button <- ButtonLamp {currentFloor, 2, 0};
						if prevDirection != 1{
							removeOrderFromQueue(currentFloor, 0, queues, output.FinishedOrder);
						}
					case 0:
								for i := 0; i<3; i++ {removeOrderFromQueue(currentFloor, i, queues, output.FinishedOrder)}
								output.Button <- ButtonLamp {currentFloor, 2, 0};
								<-values.State; values.State <- IDLE;
					case -1: 
						<-values.State; values.State <- RUN_DOWN;
						removeOrderFromQueue(currentFloor, 2, queues, output.FinishedOrder);
						output.Button <- ButtonLamp {currentFloor, 2, 0};
						if prevDirection != -1{
							removeOrderFromQueue(currentFloor, 1, queues, output.FinishedOrder);
						}
				}
		case IDLE:
	}
}

func NewOrderInEmptyQueue(values CurrentElevValues, queues OrderQueueChannels, output OutputChans) {
	printQueues(queues);
	currentState := <- values.State; values.State <- currentState;
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
		case IDLE:
			currentFloor := <- values.Floor; values.Floor <- currentFloor;
			currentDirection := <-values.Direction;
			currentDirection = nextDirection(currentFloor, currentDirection, queues);
			values.Direction <- currentDirection;
			output.Motor <- currentDirection;
			switch currentDirection{
				case 1:
					<-values.State; values.State <- RUN_UP;

				case 0:
					for i := 0; i<3; i++ {removeOrderFromQueue(currentFloor, i, queues, output.FinishedOrder)}
					output.Timer <- true;
					<-values.State; values.State <- DOOR_OPEN;
				case -1: 
					<-values.State; values.State <- RUN_DOWN;
			}
	}
}

func NewOrderToCurrentFloor(CurrentState chan State, SetTimerChan chan bool) {
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

func CheckIfCurrentFloor(floor int, CurrentFloor chan int, queues OrderQueueChannels) bool {
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	currentFloor := <- CurrentFloor; CurrentFloor <- currentFloor;
	return (floor == currentFloor);
}