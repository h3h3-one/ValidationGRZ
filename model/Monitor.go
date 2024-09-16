package model

import "fmt"

type Message struct {
	X     byte   // Координата X
	Y     byte   // Координата Y
	Color byte   // Значение цвета
	Text  string // Текст сообщения
}

type Monitor struct {
	CamNumber int       // Номер камеры
	Messages  []Message // Список сообщений, связанных с монитором
}

func NewMonitor() *Monitor {
	return &Monitor{
		CamNumber: 0,
		Messages:  []Message{}, // Инициализируем пустой срез сообщений
	}
}

func (monitor *Monitor) AddMessage(message Message) {
	monitor.Messages = append(monitor.Messages, message)
}

func NewMonitorWithParameters(camNumber int, messages []Message) *Monitor {
	return &Monitor{
		CamNumber: camNumber,
		Messages:  messages,
	}
}

func (monitor *Monitor) GetMessageCount() int {
	return len(monitor.Messages)
}

func (monitor *Monitor) GetMessages() []Message {
	return monitor.Messages
}

func (monitor *Monitor) ClearMessages() {
	monitor.Messages = []Message{}
}

func (monitor *Monitor) DisplayMessages() string {
	var result string
	for _, msg := range monitor.Messages {
		result += fmt.Sprintf("Камера %d: '%s' по координатам (%d, %d) с цветом %d\n", monitor.CamNumber, msg.Text, msg.X, msg.Y, msg.Color)
	}
	return result
}
