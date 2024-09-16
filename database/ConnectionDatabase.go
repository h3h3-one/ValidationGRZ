package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/nakagami/firebirdsql"
)

type ConfigurationModel struct {
	DatabaseIp       string `json:"databaseIp"`       // IP-адрес базы данных
	DatabasePort     string `json:"databasePort"`     // Порт базы данных
	DatabasePath     string `json:"databasePath"`     // Путь к базе данных
	DatabaseLogin    string `json:"databaseLogin"`    // Логин для доступа к базе данных
	DatabasePassword string `json:"databasePassword"` // Пароль для доступа к базе данных
}

type ConnectionDatabase struct {
	count              int                // Счетчик попыток подключения
	mapper             *json.Decoder      // Декодер для чтения JSON
	mqttConfig         *os.File           // Файл конфигурации MQTT
	configurationModel ConfigurationModel // Модель конфигурации базы данных
	connectionDB       *sql.DB            // Соединение с базой данных
}

// NewConnectionDaabase создает новый экземпляр ConnectionDatabase и загружает конфигурацию.
func NewConnectionDatabase() (*ConnectionDatabase, error) {
	db := &ConnectionDatabase{
		count:  0,
		mapper: json.NewDecoder(nil),
	}

	// Открываем конфигурационный файл
	var err error
	db.mqttConfig, err = os.Open("ValidatedConfig.json")
	if err != nil {
		return nil, err
	}
	defer db.mqttConfig.Close()

	// Декодируем конфигурацию из JSON в структуру
	if err := db.mapper.Decode(&db.configurationModel); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return db, nil
}

func (db *ConnectionDatabase) Connected() *sql.DB {
	db.count++
	log.Printf("Подключение к базе данных... Попытка: %d", db.count)

	if db.connectionDB == nil {
		var err error
		// Формируем строку подключения (DSN)
		dsn := fmt.Sprintf("firebirdsql://%s:%s/%s?encoding=WIN1251",
			db.configurationModel.DatabaseIp,
			db.configurationModel.DatabasePort,
			db.configurationModel.DatabasePath)

		// Открываем соединение с базой данных
		db.connectionDB, err = sql.Open("firebirdsql", dsn)
		if err != nil {
			log.Fatalf("Ошибка подключения к базе данных: %v", err)
		}

		// Логируем информацию о подключении
		log.Printf("Информация о подключении к базе данных. LOGIN: %s, PASSWORD: %s, IP: %s, PORT: %s, ПУТЬ: %s",
			db.configurationModel.DatabaseLogin,
			db.configurationModel.DatabasePassword,
			db.configurationModel.DatabaseIp,
			db.configurationModel.DatabasePort,
			db.configurationModel.DatabasePath)
	}

	// Проверяем, удалось ли подключиться к базе данных
	if err := db.connectionDB.Ping(); err == nil {
		log.Println("Соединение с базой данных произошло успешно!")
		return db.connectionDB
	} else {
		log.Printf("Ошибка: %v", err) // Логируем ошибку, если не удалось подключиться
	}

	return nil // Возвращаем nil, если соединение не удалось
}
