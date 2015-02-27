package elev_handler
import "time"


func DoorTimer() int{
	select {
		case <- time.After(3 * time.Second):
			return 1
	}
	return 1
}