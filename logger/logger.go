package logger

import "github.com/abinashphulkonwar/route-master/services"

type Logger struct {
	queue *services.Queue
}

func NewLogger() *Logger {
	newQueue := services.NewQueue()
	logger := &Logger{
		queue: newQueue,
	}

	go logger.start()

	return logger
}

func (l *Logger) Log(msg string) {
	go l.queue.Enqueue([]byte(msg))
}

func (l *Logger) start() {
	for {
		isFound, list := l.queue.DequeueMany(5)
		if isFound {

			for _, log := range *list {
				println(string(log))
			}

		}
	}
}
