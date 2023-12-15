package logger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/abinashphulkonwar/route-master/services"
)

type Logger struct {
	queue *services.Queue
	file  *os.File
}

type Log struct {
	Method  string
	Path    string
	Address string
	Scheme  string
	Name    string
	Time    time.Time
}

func NewLogger() *Logger {
	newQueue := services.NewQueue()
	file, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := &Logger{
		queue: newQueue,
		file:  file,
	}

	go logger.start()

	return logger
}

func (l *Logger) push(event *Log) {
	data, err := json.Marshal(event)
	if err != nil {
		println(err.Error())
		return
	}

	go l.queue.Enqueue([]byte(data))
}

func (l *Logger) Log(event *Log) {
	go l.push(event)
}

func (l *Logger) start() {
	defer l.file.Close()

	endLine := []byte("\n")
	logBuf := []byte{}

	for {

		isFound, list := l.queue.DequeueMany(50)
		if !isFound {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		for _, log := range *list {
			log = append(log, endLine...)

			logBuf = append(logBuf, log...)

		}

		l.file.Write(logBuf)
		println(string(logBuf))
		logBuf = []byte{}

	}
}
