package machine

var running bool = false
var speed int = 2

var quitChannel chan int = make(chan int)
var speedChannel chan int = make(chan int, 1)

func (m Machine) Start() {
	go m.stepper(quitChannel, speedChannel, speed)
	running = true
}

func (m Machine) Stop() {
	running = false
	quitChannel <- 1

	for len(speedChannel) > 0 {
		<-speedChannel
	}
}

func (m Machine) IsRunning() bool {
	return running
}

func (m Machine) GetSpeed() int {
	return speed
}

func (m Machine) SetSpeed(kHz int) {
	speed = kHz
	if running {
		speedChannel <- speed
	}
}
