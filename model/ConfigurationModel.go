package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Message struct {
	Byte1 byte   // Первый байт сообщения
	Byte2 byte   // Второй байт сообщения
	Byte3 byte   // Третий байт сообщения
	Text  string // Текст сообщения
}

type ConfigurationModel struct {
	MqttUsername               string               // Имя пользователя для MQTT
	MqttPassword               string               // Пароль для MQTT
	MqttClientId               string               // Идентификатор клиента для MQTT
	MqttClientIp               string               // IP-адрес клиента для MQTT
	MqttClientPort             int                  // Порт клиента для MQTT
	DatabaseLogin              string               // Логин для базы данных
	DatabasePassword           string               // Пароль для базы данных
	DatabasePath               string               // Путь к базе данных
	DatabaseIp                 string               // IP-адрес базы данных
	DatabasePort               int                  // Порт базы данных
	CameraIdDeviceIdDictionary map[int]int          // Словарь идентификаторов камер и устройств
	StringDictionary           map[string][]Message // Словарь строковых ключей и массива сообщений
}

func NewConfigurationModel() *ConfigurationModel {
	return &ConfigurationModel{
		MqttUsername:     "admin",
		MqttPassword:     "333",
		MqttClientId:     "Validation",
		MqttClientIp:     "194.87.237.67",
		MqttClientPort:   1883,
		DatabaseLogin:    "SYSDBA",
		DatabasePassword: "temp",
		DatabasePath:     `C:\Program Files (x86)\CardSoft\DuoSE\Access\ShieldPro_rest.gdb`,
		DatabaseIp:       "127.0.0.1",
		DatabasePort:     3050,
		CameraIdDeviceIdDictionary: map[int]int{
			1: 365, // Пример сопоставления идентификаторов камеры и устройства
		},
		StringDictionary: map[string][]Message{
			"65": {{Byte1: 0x09, Byte2: 0x00, Byte3: 0x02, Text: "Недействительная карточка"}},
			"50": {{Byte1: 0x09, Byte2: 0x00, Byte3: 0x02, Text: "Действительная карточка"}},
			"46": {{Byte1: 0x09, Byte2: 0x00, Byte3: 0x02, Text: "Неизвестная карточка"}},
		},
	}
}

func LoadConfiguration(filePath string) (*ConfigurationModel, error) {
	var config ConfigurationModel

	// Читаем содержимое файла
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Декодируем JSON в структуру ConfigurationModel
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, errors.New("неверный формат конфигурационного файла")
	}

	return &config, nil
}

func (c *ConfigurationModel) SaveConfiguration(filePath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Записываем данные в файл
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}
