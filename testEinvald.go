package main
import(
	"fmt"
	)

var queueUp  = [4] int {-1, -1, -1, -1};
var queueDown = [4] int {-1, -1, -1, -1};
var queueInElev = [4] int {-1, -1, -1, -1};

func main() {
	AddToQueue(1, 0)
	AddToQueue(1, 0)
	//RemoveOrderFromQueue(1, 0)
	var currentDirection int
	var currentFloor = 0
	currentDirection = nextDirection(currentFloor, 0);
	fmt.Println("currentDirection: ", currentDirection)
	for i:= range queueUp{
		fmt.Println(queueUp[i])
	} 
}
//%HOMEDRIVE%%HOMEPATH%



func nextDirection(currentFloor int, currentDirection int) int {
	switch currentDirection{
		case 1:
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor){
					return 1;
				}
			}
			for i := range queueDown{
				if (queueDown[i] < currentFloor || queueInElev[i] < currentFloor){
					return -1;
				}
			}
		case 0:
			for i := range queueUp{
				if (queueUp[i] > currentFloor || queueInElev[i] > currentFloor){
					return 1;
				}
				if (queueDown[i] < currentFloor || queueInElev[i] < currentFloor){
					return -1;
				}
			}
		case -1:
			for i := range queueDown{
				if (queueDown[i] < currentFloor || queueInElev[i] < currentFloor){
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


func checkIfEmptyQueues() bool {
	for i := range queueUp{
		if (queueUp[i] >= 0 || queueDown[i] >= 0 || queueInElev[i] >= 0){
			return false;
		}
	}
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

func RemoveOrderFromQueue(floor int, buttonType int){
	switch buttonType{
		case 0:
			queueUp[floor] = -1;
		case 1:
			queueDown[floor] = -1;
		case 2:
			queueInElev[floor] = -1;
	}
}

/*
func RemoveOrderFromQueue(floor int, buttonType int){
	switch buttonType{
		case 0:
			for i:= range queueUp{
				if queueUp[i] == floor{
					fmt.Println("Slettet fra QueueUp: ", queueUp[i]);
					queueUp[i] = -1;
				}
			}
		case 1:
			for i:= range queueDown{
				if queueUp[i] == floor{
					queueDown[i] = -1;
				}
			}
		case 2:
			for i:= range queueInElev{
				if queueInElev[i] == floor{
					queueInElev[i] = -1;
				}
			}
	}
}


func AddToQueue(floor int, buttonType int){
	switch buttonType{
		case 0:
			addToQueueUp(floor)
		case 1:
			addToQueueDown(floor)
		case 2: 
			addToQueueInElev(floor)
	}
}

func addToQueueUp(floor int) int{
	for i:= range queueUp{
		if queueUp[i] == floor{
			fmt.Println("Ligger i koen fra for")
			return 0;
		} 
			
		
	}
	for i:= range queueUp{
		if queueUp[i] < 0 {
			queueUp[i] = floor;
			return 1;
		}
	}
	return 0;
}


func addToQueueDown(floor int) int{
	for i:= range queueDown{
		if queueDown[i] == floor{
			fmt.Println("Ligger i koen fra for")
			return 0;
		} 
			
		
	}
	for i:= range queueDown{
		if queueDown[i] < 0 {
			queueDown[i] = floor;
			return 1;
		}
	}
	return 0;
}	

func addToQueueInElev(floor int) int{
	for i:= range queueInElev{
		if queueInElev[i] == floor{
			fmt.Println("Ligger i koen fra for")
			return 0;
		} 
			
		
	}
	for i:= range queueInElev{
		if queueInElev[i] < 0 {
			queueInElev[i] = floor;
			return 1;
		}
	}
	return 0;
}
*/			

/*
func NewOrderInEmptyQueue {
	switch currentState{
		case RUN_UP:
		case RUN_DOWN:
		case IDLE:
		case DOOR_OPEN:
		case STOP:
		default:
	}
}
*/
/*
func CalculateDistanceToOrder(currentFloor int, orderedFloor int, buttonType int, currentDirection int){
	distance := 0
	switch buttonType{
		case 0:
			switch currentDirection{
				case 1:
					

		case 1:
	}
}

*/