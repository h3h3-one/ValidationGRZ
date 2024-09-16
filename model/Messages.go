package model

import "fmt"

type Message struct {
	X     byte   // Координата X
	Y     byte   // Координата Y
	Color byte   // Значение цвета
	Text  string // Текст сообщения
}

func NewMessage() Message {
	return Message{
		X:     0,   // Устанавливаем X в 0 по умолчанию
		Y:     0,   // Устанавливаем Y в 0 по умолчанию
		Color: 255, // Устанавливаем цвет по умолчанию (например, белый)
		Text:  "",  // Пустой текст по умолчанию
	}
}

func NewMessageWithParams(x, y, color byte, text string) Message {
	return Message{
		X:     x,
		Y:     y,
		Color: color,
		Text:  text,
	}
}

func (m Message) Display() string {
	return fmt.Sprintf("Message: '%s' at (%d, %d) with color %d", m.Text, m.X, m.Y, m.Color)
}

func (m *Message) UpdateText(newText string) {
	m.Text = newText
}

func (m *Message) UpdateCoordinates(newX, newY byte) {
	m.X = newX
	m.Y = newY
}

func (m *Message) UpdateColor(newColor byte) {
	m.Color = newColor
}
