package elev_handler
import "time"


func DoorTimer(){
	time.Sleep(3000 * time.Millisecond);
	TimerOut();
}