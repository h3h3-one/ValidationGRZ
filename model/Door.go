package model

import (
	"fmt"
)

type Door struct {
	CameraNumber int  // Номер камеры, связанной с этой дверью
	IsOpen       bool // Состояние двери: открыта или закрыта
	AccessLevel  int  // Уровень доступа, необходимый для открытия двери
}

func NewDoor(cameraNumber int, accessLevel int) *Door {
	return &Door{
		CameraNumber: cameraNumber,
		IsOpen:       false, // Изначально дверь закрыта
		AccessLevel:  accessLevel,
	}
}

func (d *Door) OpenDoor(userAccessLevel int) error {
	if d.IsOpen {
		return fmt.Errorf("дверь уже открыта")
	}
	if userAccessLevel < d.AccessLevel {
		return fmt.Errorf("недостаточный уровень доступа для открытия двери")
	}

	d.IsOpen = true
	return nil
}

func (d *Door) CloseDoor() error {
	if !d.IsOpen {
		return fmt.Errorf("дверь уже закрыта")
	}

	d.IsOpen = false
	return nil
}

func (d *Door) GetDoorStatus() string {
	if d.IsOpen {
		return "Дверь открыта"
	}
	return "Дверь закрыта"
}
