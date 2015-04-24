package elev_handler

import(
	"fmt"
	"math"
	)

type OrderQueueChannels struct{
	Up chan [N_FLOORS] int;
	Down chan [N_FLOORS] int;
	InElev chan [N_FLOORS] int;
}

const N_FLOORS int = 4
const N_BUTTONS int = 3

func nextDirection(currentFloor int, currentDirection int, queues OrderQueueChannels) int {
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	switch currentDirection{
		case 1:
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor || queueDown[i]> currentFloor){return 1;}
			}
			for i := range queueDown{
				if ((queueDown[i] < currentFloor && queueDown[i]>=0) || (queueInElev[i] < currentFloor && queueInElev[i]>=0) || (queueUp[i] < currentFloor && queueUp[i] >= 0) ){return -1;}
			}
		case 0:
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor || queueDown[i] > currentFloor){return 1;}
				if ((queueDown[i] < currentFloor && queueDown[i]>=0) || (queueInElev[i] < currentFloor && queueInElev[i]>=0) || (queueUp[i] < currentFloor && queueUp[i] >= 0)){return -1;}
			}
		case -1:
			for i := range queueDown{
				if ((queueDown[i] < currentFloor && queueDown[i]>=0) || (queueInElev[i] < currentFloor && queueInElev[i]>=0) || (queueUp[i] < currentFloor && queueUp[i] >= 0)) {return -1;}
			}
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor || queueDown[i]>currentFloor){return 1;}
			}
	}
	return 0;
}

func CheckIfEmptyQueues(queues OrderQueueChannels) bool {
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	for i := range queueUp{
		if (queueUp[i] >= 0 || queueDown[i] >= 0 || queueInElev[i] >= 0){return false;}
	}; return true;
}

func AddToQueue(floor int, buttonType int, queues OrderQueueChannels){
	switch buttonType{
		case 0:
			queueUp := <- queues.Up
			queueUp[floor] = floor;
			queues.Up <- queueUp
		case 1:
			queueDown := <- queues.Down
			queueDown[floor] = floor;
			queues.Down <- queueDown;
		case 2: 
			queueInElev := <- queues.InElev
			queueInElev[floor] = floor;
			queues.InElev <- queueInElev 
	}
}

func removeOrderFromQueue(floor int, buttonType int, queues OrderQueueChannels, finishedOrderChan chan ButtonOrder){
	switch buttonType{
		case 0:
			queueUp := <- queues.Up
			if queueUp[floor] == floor{finishedOrderChan <- ButtonOrder {floor, buttonType};}
			queueUp[floor] = -1;
			queues.Up <- queueUp
		case 1:
			queueDown := <- queues.Down
			if queueDown[floor] == floor{finishedOrderChan <- ButtonOrder {floor, buttonType};}
			queueDown[floor] = -1;
			queues.Down <- queueDown;
		case 2:
			queueInElev := <- queues.InElev
			queueInElev[floor] = -1;
			queues.InElev <- queueInElev 
	}
}

func CheckIfFloorInQueue(floor int, CurrentDirection chan int, queues OrderQueueChannels) bool{
	if floor < 0 {floor = floor*(-1);}
	currentDirection := <- CurrentDirection;
	CurrentDirection <- currentDirection;
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	switch currentDirection{
		case 1:
			if queueUp[floor] == floor || queueInElev[floor] == floor{return true;}
			if !checkIfOrdersAtHigherFloors(floor, queues){
				if queueDown[floor] == floor{return true;} 
			}
		case -1:
			if queueDown[floor] == floor || queueInElev[floor] == floor{return true;}
			if !checkIfOrdersAtLowerFloors(floor, queues){
				if queueUp[floor] == floor{return true;}
			}
	}
	return false;
}

func InitializeQueues(queues OrderQueueChannels){
	queueUp := [N_FLOORS] int {}; queueDown := [N_FLOORS] int {}; queueInElev := [N_FLOORS] int {};
	for i:=range queueUp{
		queueUp[i] = -1;
		queueDown[i] = -1;
		queueInElev[i] = -1;
	}
	queues.Up <- queueUp; queues.Down <- queueDown; queues.InElev <- queueInElev;
}

func printQueues(queues OrderQueueChannels){
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	for i := range queueUp{
		fmt.Println(queueUp[i],"  ", queueDown[i],"  ", queueInElev[i]);
	}
}

func checkIfOrdersAtHigherFloors(floor int, queues OrderQueueChannels) bool{
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	for i := range queueUp{
		if ((queueUp[i] == i && i>floor) || (queueInElev[i]==i && i>floor) || (queueDown[i]==i && i>floor)){return true;}
	}; return false;
}

func checkIfOrdersAtLowerFloors(floor int, queues OrderQueueChannels) bool{
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	for i := range queueUp{
		if ((queueDown[i] == i && i<floor) || (queueInElev[i]==i && i<floor) || (queueUp[i] == i && i<floor)){return true;}
	}; return false;
}

func getQueue(buttonType int, queues OrderQueueChannels) [4] int{
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	switch buttonType{
		case 0:
			return queueUp;
		case 1:
			return queueDown;
		case 2:
			return queueInElev
	}; return queueUp;
}

func GetCostForOrder(floor int, buttonType int, values CurrentElevValues, queues OrderQueueChannels) int {
	currentDirection := <- values.Direction; values.Direction <- currentDirection;
	currentFloor := <- values.Floor; values.Floor <- currentFloor;
	currentState := <- values.State; values.State <- currentState;
	cost := 0;
	if currentState== IDLE {cost += (int(math.Abs(float64(floor - currentFloor))*3));} else if currentState == DOOR_OPEN{cost+=3;}
	if currentState == DOOR_OPEN && ((currentDirection == 1 && buttonType == 0) || (currentDirection == -1 && buttonType == 1)) && floor == currentFloor {return 0;}
	switch currentDirection {
		case 1:
			switch buttonType {
				case 0:
					if currentFloor >= floor{cost+=8;} else{cost += (int(math.Abs(float64(floor - currentFloor))*3));}
				case 1:
					cost += 5;
			}
		case -1:
			switch buttonType {
				case 0:
					cost += 5;
				case 1:
					if currentFloor <= floor{cost+=8;} else{cost += (int(math.Abs(float64(floor - currentFloor))*3));}
			}
	}
	stops := amountOfOrdersInQueue(floor, buttonType, currentDirection, currentFloor, queues);
	cost += stops*3;
	return cost;
}
func amountOfOrdersInQueue(floor int, buttonType int, currentDirection int, currentFloor int, queues OrderQueueChannels) int {
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	stopCounter := 0;
	for i := range queueUp{
		if queueUp[i] == i {stopCounter += 1;}
		if queueDown[i] == i {stopCounter += 1;}
		if queueInElev[i] == i {stopCounter += 1;}
	}
	return stopCounter;
}
func checkIfOrdersInFloor(floor int, queues OrderQueueChannels) bool{
	queueUp := <- queues.Up; queues.Up <- queueUp; 
	queueDown := <- queues.Down; queues.Down <- queueDown;
	queueInElev := <- queues.InElev; queues.InElev <- queueInElev;
	if (queueUp[floor]==floor || queueDown[floor]==floor || queueInElev[floor]==floor){return true;} else {return false;}
}