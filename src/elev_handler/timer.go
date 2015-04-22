package elev_handler
import "time"


func DoorTimer(timerChan chan int64, CurrentDirection chan int, CurrentFloor chan int, CurrentState chan State){
	startTime := <- timerChan;

	for {
		select{
			case startTime = <- timerChan:
			case <- time.After(300 *time.Millisecond):
				if (time.Now().UnixNano()/int64(time.Millisecond) - startTime) > 3000{
					TimerOut(CurrentDirection, CurrentFloor, CurrentState);
					startTime = <- timerChan;
				}
		}
	}
}