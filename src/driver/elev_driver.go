package driver
import (
	"fmt"
	)
// All function which starts with D_ is driver-functions

const N_FLOORS int = 4
const N_BUTTONS int = 3


var lamp_channel_matrix = [N_FLOORS][N_BUTTONS] int{
{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}
var button_channel_matrix = [N_FLOORS][N_BUTTONS] int{
{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}


func Driver_init() int {
	if(!Io_init()){
		return 0
	}
	for i := 0; i < N_FLOORS; i++ {
		if (i != 0){
			Driver_set_button_lamp(1, i, 0);
		}
		if (i != N_FLOORS -1){
			Driver_set_button_lamp(0, i, 0); 		
		}
		Driver_set_button_lamp(2, i, 0);
	}
	Driver_set_stop_lamp(0);
	Driver_set_door_open_lamp(0)
	return 1

}

func Driver_set_motor_direction(direction int){
	fmt.Println("Setter motorretning i retning", direction);
	switch {
		case direction == 0:
			Io_write_analog(MOTOR, 0)
		case direction > 0:
			Io_clear_bit(MOTORDIR)
			Io_write_analog(MOTOR, 2800)
		case direction < 0:
			Io_set_bit(MOTORDIR)
			Io_write_analog(MOTOR, 2800)
	}
}

func Driver_set_door_open_lamp(value int){
	if (value >0){
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Driver_set_stop_lamp(value int){
		if (value >0){
		Io_set_bit(LIGHT_STOP)
	} else {
		Io_clear_bit(LIGHT_STOP)
	}
}

func Driver_set_floor_indicator(floor int){
	if( (floor>=0) && (floor<N_FLOORS)){
	
		if (floor & 0x02 != 0 ){
			Io_set_bit(LIGHT_FLOOR_IND1)
		} else {
			Io_clear_bit(LIGHT_FLOOR_IND1)
		}
	
		if (floor & 0x01 != 0){
			Io_set_bit(LIGHT_FLOOR_IND2)
		} else {
			Io_clear_bit(LIGHT_FLOOR_IND2)
		}
	}
}

func Driver_set_button_lamp(buttonType int, floor int, value int){
	if((floor>=0) && (floor<N_FLOORS) && !(buttonType == 0 && floor ==3) && !(buttonType==1 && floor ==0) && buttonType>=0 && buttonType<=2){ //assertions
		if(value != 0 ){
			Io_set_bit(lamp_channel_matrix[floor][buttonType])
		} else {
			Io_clear_bit(lamp_channel_matrix[floor][buttonType])
		}
	}
}

func Driver_get_floor_sensor_signal() int {
	switch{
		case Io_read_bit(SENSOR_FLOOR1) != 0:
			return 0
		case Io_read_bit(SENSOR_FLOOR2) != 0:
			return 1
		case Io_read_bit(SENSOR_FLOOR3) != 0:
			return 2
		case Io_read_bit(SENSOR_FLOOR4) != 0:
			return 3
		default:
			return -1
	}
}

func Driver_get_stop_signal() bool{
	return (Io_read_bit(STOP) != 0)
}

func Driver_get_button_signal(buttonType int, floor int) int{
	if((floor>=0) && (floor<N_FLOORS) && !(buttonType == 0 && floor ==3) && !(buttonType==1 && floor ==0) && buttonType>=0 && buttonType<=2){ //assertions
		if(Io_read_bit(button_channel_matrix[floor][buttonType])!= 0) {
			return 1;
		}else{
			return 0;
		}
	}
	return 0;
}







