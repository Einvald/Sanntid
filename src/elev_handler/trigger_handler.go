package elev_handler

type State int
const (
RUN_UP State = 0 + iota
RUN_DOWN
DOOR_OPEN
IDLE
STOP
)

/*

}
var currentFloor int
var currentState State = IDLE
//Må vurdere hvor nåværende state skal være lagret, men tenker at den skal være lagret her, også "last direction" kan være lurt å ha med
func floorReached(floor int){
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:

		case STOP:

	}
}

func timerOut() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

		case DOOR_OPEN:

		case IDLE:

		case STOP:

	}
}

func stopButton() {
	switch currentState{
		case RUN_UP:

		case RUN_DOWN:

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

		case STOP:

	}
}

func newOrderToCurrentFloor() {
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