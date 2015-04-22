package elev_handler

import(
	"fmt"
	"math"
	)
var FinishedOrderChan = make (chan ButtonOrder, 1024);
var queueUp  = [4] int {};
var queueDown = [4] int {};
var queueInElev = [4] int {};

func nextDirection(currentFloor int, currentDirection int) int {
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

func CheckIfEmptyQueues() bool {
	for i := range queueUp{
		if (queueUp[i] >= 0 || queueDown[i] >= 0 || queueInElev[i] >= 0){return false;}
	}
	fmt.Println("Empty Queues true");
	return true;
}

func AddToQueue(floor int, buttonType int){
	switch buttonType{
		case 0:
			queueUp[floor] = floor;
		case 1:
			queueDown[floor] = floor;
		case 2: 
			queueInElev[floor] = floor;
	}
}

func removeOrderFromQueue(floor int, buttonType int){
	switch buttonType{
		case 0:
			if queueUp[floor] == floor{FinishedOrderChan <- ButtonOrder {floor, buttonType};}
			queueUp[floor] = -1;
		case 1:
			if queueDown[floor] == floor{FinishedOrderChan <- ButtonOrder {floor, buttonType};}
			queueDown[floor] = -1;
		case 2:
			queueInElev[floor] = -1;
	}
}

func CheckIfFloorInQueue(floor int, CurrentDirection chan int) bool{
	if floor < 0 {floor = floor*(-1);}
	currentDirection := <- CurrentDirection;
	CurrentDirection <- currentDirection;

	fmt.Println("CHECKIFFLOORINQUEUE: ", floor);
	switch currentDirection{
		case 1:
			if queueUp[floor] == floor || queueInElev[floor] == floor{return true;}
			if !checkIfOrdersAtHigherFloors(floor){
				if queueDown[floor] == floor{return true;} 
			}
		case -1:
			if queueDown[floor] == floor || queueInElev[floor] == floor{return true;}
			if !checkIfOrdersAtLowerFloors(floor){
				if queueUp[floor] == floor{return true;}
			}
	}
	return false;
}

func EmptyQueues(){
	for i:=range queueUp{
		queueUp[i] = -1;
		queueDown[i] = -1;
		queueInElev[i] = -1;
	}
}

func printQueues(){
	for i := range queueUp{
		fmt.Println(queueUp[i],"  ", queueDown[i],"  ", queueInElev[i]);
	}
}


func checkIfOrdersAtHigherFloors(floor int) bool{
	for i := range queueUp{
		if ((queueUp[i] == i && i>floor) || (queueInElev[i]==i && i>floor) || (queueDown[i]==i && i>floor)){return true;}
	}
	return false;
}

func checkIfOrdersAtLowerFloors(floor int) bool{
	for i := range queueUp{
		if ((queueDown[i] == i && i<floor) || (queueInElev[i]==i && i<floor) || (queueUp[i] == i && i<floor)){return true;}
	}
	return false;
}

func getQueue(buttonType int) [4] int{
	switch buttonType{
	case 0:
		return queueUp;
	case 1:
		return queueDown;
	case 2:
		return queueInElev
	}
	return queueUp;
}

func GetCostForOrder(floor int, buttonType int, CurrentDirection chan int, CurrentFloor chan int, CurrentState chan State) int {
	currentDirection := <- CurrentDirection; CurrentDirection <- currentDirection;
	currentFloor := <- CurrentFloor; CurrentFloor <- currentFloor;
	currentState := <- CurrentState; CurrentState <- currentState;
	cost := 0;
	if currentState== IDLE {cost += (int(math.Abs(float64(floor - currentFloor))*3));} else if currentState == DOOR_OPEN{cost+=3;}
	if currentState == DOOR_OPEN && ((currentDirection == 1 && buttonType == 0) || (currentDirection == -1 && buttonType == 1)) && floor == currentFloor {return 0;}
	switch currentDirection {
		case 1:
			switch buttonType {
				case 0:
					if currentFloor >= floor{cost+=15;} else{cost += (int(math.Abs(float64(floor - currentFloor))*3));}
				case 1:
					cost += 10;

			}
		case -1:
			switch buttonType {
				case 0:
					cost += 10;
				case 1:
					if currentFloor <= floor{cost+=15;} else{cost += (int(math.Abs(float64(floor - currentFloor))*3));}
			}
	}
	stops := amountOfOrdersInQueue(floor, buttonType, currentDirection, currentFloor);
	cost += stops*3;
	return cost;
}

func amountOfOrdersInQueue(floor int, buttonType int, currentDirection int, currentFloor int) int {
	stopCounter := 0;
	for i := range queueUp{
		if queueUp[i] == i {stopCounter += 1;}
		if queueDown[i] == i {stopCounter += 1;}
		if queueInElev[i] == i {stopCounter += 1;}
	}
	return stopCounter;
}

func checkIfOrdersInFloor(floor int) bool{
	if (queueUp[floor]==floor || queueDown[floor]==floor || queueInElev[floor]==floor){
		return true;
	} else {return false;}
}

/* Tidligere innmat i amountOfStopsBeforeFloor
	stopCounter := 0;
	floorIsOver := floor>currentFloor;
	switch currentDirection{
		case 1:
			switch buttonType{
				case 0:
					if floorIsOver{
						for i := range queueUp{
							if i == floor {break;}
							if ((queueUp[i] == i || queueInElev[i] == i) && i>currentFloor){stopCounter++;} 
						}
					} else{
						for i := range queueUp{
							if queueDown[i]==i || queueInElev[i] == i  {stopCounter++;}
							if queueUp[i] == i && !(i>floor && i<=currentFloor) && queueInElev[i] != i {stopCounter++;}

						}
					}
				case 1:
						for i := range queueUp{
							if ((queueUp[i] == i || queueInElev[i] == i) && i>currentFloor){stopCounter++;}
							if (queueDown[i] == i || (queueInElev[i] == i && i<currentFloor)) && i>floor {stopCounter++;} 
						}
			}
		case -1:
			switch buttonType{
				case 0:
						for i := range queueDown{
							if ((queueDown[i] == i || queueInElev[i] == i) && i<currentFloor){stopCounter++;}
							if (queueUp[i] == i || (queueInElev[i] == i && i>=currentFloor)) && i<floor  {stopCounter++;}
						}
				case 1:
					if !floorIsOver && floor != currentFloor{
						for i := range queueDown{
							if ((queueDown[i] == i || queueInElev[i] == i) && i<currentFloor && i>floor){stopCounter++;}
						}
					} else{
						for i := range queueUp{
							if (queueUp[i]==i || queueInElev[i] == i) {stopCounter++;}
							if queueDown[i] == i && !(i<floor && i>=currentFloor) && queueInElev[i] != i {stopCounter ++;} 
						}
			}
		}

	}
	fmt.Println("STOPAMOUNT = ", stopCounter);
	fmt.Println("STOPAMOUNT = ", stopCounter);
	*/