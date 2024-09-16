package mqtt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

type ConfigurationModel struct {
	MqttClientIp               string      `json:"mqtt_client_ip"`
	MqttClientPort             int         `json:"mqtt_client_port"`
	MqttUsername               string      `json:"mqtt_username"`
	MqttPassword               string      `json:"mqtt_password"`
	CameraIdDeviceIdDictionary map[int]int `json:"camera_id_device_id_dictionary"`
}

type MessageIntegration struct {
	Grz       string `json:"grz"`
	CamNumber int    `json:"cam_number"`
}

type ConnectionDatabase struct {
}

func (c *ConnectionDatabase) connected() (*sql.DB, error) {
	dsn := "username:password@tcp(127.0.0.1:3306)/dbname"
	return sql.Open("mysql", dsn)
}

type BaseMqttClient struct {
	mqttConfig         string
	configurationModel ConfigurationModel
	eventType          string
	idPep              string
	connectionDB       ConnectionDatabase
}

func NewBaseMqttClient() (*BaseMqttClient, error) {
	return &BaseMqttClient{
		mqttConfig:   "ValidatedConfig.json",
		connectionDB: ConnectionDatabase{},
	}, nil
}

func (b *BaseMqttClient) getConnection() error {
	b.isNewFile(b.mqttConfig)

	file, err := ioutil.ReadFile(b.mqttConfig)
	if err != nil {
		return fmt.Errorf("ошибка чтения конфигурационного файла: %w", err)
	}

	if err := json.Unmarshal(file, &b.configurationModel); err != nil {
		return fmt.Errorf("ошибка десериализации конфигурации: %w", err)
	}

	// Логируем параметры подключения
	log.Printf("Создание подключения клиента... HOST_NAME = %s, PORT = %d, USERNAME = %s",
		b.configurationModel.MqttClientIp,
		b.configurationModel.MqttClientPort,
		b.configurationModel.MqttUsername)

	// Настройка клиента MQTT
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", b.configurationModel.MqttClientIp, b.configurationModel.MqttClientPort))
	opts.SetClientID(fmt.Sprintf("%s-Validation", getLocalHostName()))
	opts.SetUsername(b.configurationModel.MqttUsername)
	opts.SetPassword(b.configurationModel.MqttPassword)
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5 * time.Second)

	client := mqtt.NewClient(opts)

	// Подключение к брокеру
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("ошибка подключения к MQTT брокеру: %w", token.Error())
	}

	// Подписка на топик
	client.Subscribe("Parking/IntegratorCVS", 0, func(client mqtt.Client, msg mqtt.Message) {
		b.messageHandler(msg)
	})

	log.Println("Подписка на топик Parking/IntegratorCVS прошла успешно.")
	return nil
}

func (b *BaseMqttClient) messageHandler(msg mqtt.Message) {
	log.Printf("Получено сообщение! ТОПИК: %s СООБЩЕНИЕ: %s", msg.Topic(), msg.Payload())

	var messages MessageIntegration
	if err := json.Unmarshal(msg.Payload(), &messages); err != nil {
		log.Printf("Ошибка обработки сообщения: %s", err)
		return
	}

	grz := messages.Grz
	camNumber := messages.CamNumber

	b.eventType = "" // Обновляем значение, чтобы не кэшировалось
	b.execute(grz, camNumber)
}

// execute выполняет процедуру в базе данных.
func (b *BaseMqttClient) execute(grz string, camNumber int) {
	idDev, exists := b.configurationModel.CameraIdDeviceIdDictionary[camNumber]
	if !exists {
		log.Printf("Ошибка: Не найден ID устройства для камеры: %d", camNumber)
		return
	}

	connection, err := b.connectionDB.connected()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %s", err)
		return
	}
	defer connection.Close()

	log.Printf("Входящие параметры для процедуры: ID_DEV: %d ID_CARD: %s GRZ: %s", idDev, grz, grz)

	procedure := "{ call REGISTERPASS_HL(?,?,?) }"
	call, err := connection.Prepare(procedure)
	if err != nil {
		log.Printf("Ошибка подготовки процедуры: %s", err)
		return
	}
	defer call.Close()

	if _, err = call.Exec(idDev, grz, grz); err != nil {
		log.Printf("Ошибка выполнения процедуры: %s", err)
		return
	}

	log.Println("Успешное выполнение процедуры.")
}

func (b *BaseMqttClient) isNewFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Ошибка создания файла: %s", err)
			return
		}
		defer file.Close()

		ow := json.NewEncoder(file)
		ow.SetIndent("", "  ")
		if err := ow.Encode(ConfigurationModel{}); err != nil {
			log.Printf("Ошибка записи в файл конфигурации: %s", err)
		}

		log.Printf("Файл конфигурации успешно создан. Запустите программу заново. ПУТЬ: %s", filePath)
		os.Exit(0)
	}
}

func getLocalHostName() string {
	host, err := net.LookupHost("localhost")
	if err != nil {
		log.Printf("Ошибка получения имени хоста: %s", err)
		return "localhost"
	}
	return host[0]
}
