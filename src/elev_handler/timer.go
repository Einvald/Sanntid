package elev_handler
import "time"


func DoorTimer(timeChan chan int64){
	startTime := <- timeChan;

	for {
		select{
			case startTime = <- timeChan:
			case <- time.After(300 *time.Millisecond):
				if (time.Now().UnixNano()/int64(time.Millisecond) - startTime) > 3000{
					TimerOut();
					startTime = <- timeChan;
				}
		}
	}
}