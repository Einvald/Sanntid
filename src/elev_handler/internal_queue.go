package elev_handler

import(
	"fmt"
	)

var queueUp  = [4] int {};
var queueDown = [4] int {};
var queueInElev = [4] int {};

func nextDirection(currentFloor int, currentDirection int) int {
	switch currentDirection{
		case 1:
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor){
					return 1;
				}
			}
			for i := range queueDown{
				if ((queueDown[i] < currentFloor && queueDown[i]>=0) || (queueInElev[i] < currentFloor && queueInElev[i]>=0) || (queueUp[i] < currentFloor && queueUp[i] >= 0) ){
					return -1;
				}
			}
		case 0:
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor || queueDown[i] > currentFloor){
					return 1;
				}
				if ((queueDown[i] < currentFloor && queueDown[i]>=0) || (queueInElev[i] < currentFloor && queueInElev[i]>=0) || (queueUp[i] < currentFloor && queueUp[i] >= 0)){
					return -1;
				}
			}
		case -1:
			for i := range queueDown{
				if ((queueDown[i] < currentFloor && queueDown[i]>=0) || (queueInElev[i] < currentFloor && queueInElev[i]>=0)) {
					return -1;
				}
			}
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor){
					return 1;
				}
			}
	}
	return 0;
}

func CheckIfEmptyQueues() bool {
	for i := range queueUp{
		if (queueUp[i] >= 0 || queueDown[i] >= 0 || queueInElev[i] >= 0){
			return false;
			fmt.Println("Not empty Queues")
		}
	}
	fmt.Println("Empty Queues")
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
			queueUp[floor] = -1;
		case 1:
			queueDown[floor] = -1;
		case 2:
			queueInElev[floor] = -1;
	}
}

func checkIfFloorInQueue(floor int, currentDirection int) bool{
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
		if ((queueUp[i] == i && i>floor) || (queueInElev[i]==i && i>floor)){
			return true;
		}
	}
	return false;
}

func checkIfOrdersAtLowerFloors(floor int) bool{
	for i := range queueUp{
		if ((queueDown[i] == i && i<floor) || (queueInElev[i]==i && i<floor)){
			return true;
		}
	}
	return false;
}