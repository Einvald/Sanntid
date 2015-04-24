package elev_handler
import "time"

func DoorTimer(timerChan chan int64, values CurrentElevValues, orderQueueChannels OrderQueueChannels, output OutputChans){
	startTime := <- timerChan;

	for {
		select{
			case startTime = <- timerChan:
			case <- time.After(300 *time.Millisecond):
				if (time.Now().UnixNano()/int64(time.Millisecond) - startTime) > 3000{
					TimerOut(values, orderQueueChannels, output);
					startTime = <- timerChan;
				}
		}
	}
}