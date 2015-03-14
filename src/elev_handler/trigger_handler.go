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
var currentDirectionLockChan = make(chan int, 1);
var currentFloorLockChan = make(chan int, 1);
var currentStateLockChan = make(chan int, 1);
var currentState State = RUN_DOWN;

func FloorReached(floor int){
	if floor<0{floor = floor*(-1)}
	fmt.Println("Inside FloorReached with current state: ", getCurrentState());
	SetCurrentFloorLampChan <- floor;
	setCurrentFloor(floor);
	switch getCurrentState(){
		case RUN_UP:
			if checkIfFloorInQueue(floor, getCurrentDirection()){
				
				setCurrentState(DOOR_OPEN);
				SetTimerChan <- true;
				SetButtonLampChan <- ButtonLamp {floor, 0, 0};
				SetButtonLampChan <- ButtonLamp {floor, 2, 0};
				removeOrderFromQueue(floor, 0); //MÅ SENDE MELDING OM AT KØ FOR ALLE HEISER ER FERDIG
				removeOrderFromQueue(floor, 2);
				if !checkIfOrdersAtHigherFloors(floor){
					removeOrderFromQueue(floor, 1);
					SetButtonLampChan <- ButtonLamp {floor, 1, 0};
				}
				printQueues();
			}
					
		case RUN_DOWN:
			if checkIfFloorInQueue(floor, getCurrentDirection()){

				setCurrentState(DOOR_OPEN);
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
	switch getCurrentState(){
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
				SetTimerChan <- false;
				prevDirection := getCurrentDirection();
				setCurrentDirection(nextDirection(getCurrentFloor(), getCurrentDirection()));
				SetMotorChan <- getCurrentDirection();
				switch currentDirection{
					case 1:
						setCurrentState(RUN_UP);
						removeOrderFromQueue(getCurrentFloor(), 2);
						SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 2, 0};
						if prevDirection != 1{
							removeOrderFromQueue(getCurrentFloor(), 0);
							SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 0, 0};
						}

					case 0:
								EmptyQueues();
								SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 0, 0};
								SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 1, 0};
								SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 2, 0};
								setCurrentState(IDLE);	
					case -1: 
						setCurrentState(RUN_DOWN);
						removeOrderFromQueue(getCurrentFloor(), 2);
						SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 2, 0};
						if prevDirection != -1{
							removeOrderFromQueue(getCurrentFloor(), 1);
							SetButtonLampChan <- ButtonLamp {getCurrentFloor(), 1, 0};
						}
				}
		case IDLE:
		case STOP:
	}
}

func StopButton() {
	StopButtonLampChan <- 1;
	//EmptyQueues();
	fmt.Println("STOP")
	
	switch getCurrentState(){
		case RUN_UP:
			SetMotorChan <- 0;
			setCurrentDirection(0);
			setCurrentState(STOP);
		case RUN_DOWN:
			SetMotorChan <- 0;
			setCurrentDirection(0);
			setCurrentState(STOP);
		case DOOR_OPEN:
			setCurrentState(STOP);
		case IDLE:
			setCurrentState(STOP);
		case STOP:
			StopButtonLampChan <- 0;
			setCurrentState(IDLE);
	}
	
}

func NewOrderInEmptyQueue() {
	fmt.Println("New Order In EMpty Queue")
	printQueues();
	switch getCurrentState(){
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
		case IDLE:
			fmt.Println("Case IDLE")
			setCurrentDirection(nextDirection(getCurrentFloor(), getCurrentDirection()));
			fmt.Println("Current direction is now set to: ", getCurrentDirection());
			SetMotorChan <- getCurrentDirection();
			switch getCurrentDirection(){
				case 1:
					setCurrentState(RUN_UP);

				case 0:
					EmptyQueues()
					SetTimerChan <- true;
					setCurrentState(DOOR_OPEN);
				case -1: 
					setCurrentState(RUN_DOWN);
			}
		case STOP:
	}
}

func NewOrderToCurrentFloor() {
	switch getCurrentState(){
		case RUN_UP:
		case RUN_DOWN:
		case DOOR_OPEN:
			SetTimerChan <- true;
		case IDLE:
			setCurrentState(DOOR_OPEN);
			SetTimerChan <- true;

		case STOP:
			SetTimerChan <- true;
	}
}

//Litt "stygge" funksjoner nedenfor
func CheckIfFloorInQueue(floor int) bool{
	return checkIfFloorInQueue(floor, getCurrentDirection());
}

func CheckIfCurrentFloor(floor int) bool {
	return (floor == getCurrentFloor());
}

func GetCostForOrder(floor int, buttonType int) int {
	return getCostForOrder(floor, buttonType, getCurrentDirection(), getCurrentFloor(), getCurrentState());
}

func InitializeChanLocks(){
	currentDirectionLockChan <- 1;
	currentFloorLockChan  <- 1;
	currentStateLockChan <- 1;

}
func setCurrentDirection(direction int){
	<- currentDirectionLockChan
	currentDirection = direction;
	currentDirectionLockChan <- 1;
}

func getCurrentDirection() int{
	<- currentDirectionLockChan
	temp := currentDirection;
	currentDirectionLockChan <- 1;
	return temp;
}

func setCurrentFloor(floor int){
	<-currentFloorLockChan;
	currentFloor = floor;
	currentFloorLockChan <- 1;
}

func getCurrentFloor() int{
	<- currentFloorLockChan ;
	temp := currentFloor;
	currentFloorLockChan  <- 1;
	return temp;
}

func setCurrentState(state State) {
	<- currentStateLockChan;
	currentState = state;
	currentStateLockChan <- 1;
}

func getCurrentState() State{
	<- currentStateLockChan;
	temp := currentState;
	currentStateLockChan <- 1;
	return temp;
}
