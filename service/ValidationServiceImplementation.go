package validationgrz

import (
	"context"
	"log"
	"sync"
)

type ValidationServiceImplementation struct {
	mqttServices MqttServices
}

func NewValidationServiceImplementation(mqttServices MqttServices) *ValidationServiceImplementation {
	return &ValidationServiceImplementation{
		mqttServices: mqttServices,
	}
}

func (v *ValidationServiceImplementation) GetConnectionMqttClient(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	err := v.mqttServices.GetConnection()
	if err != nil {
		log.Printf("Ошибка при подключении к MQTT: %v", err)
		return
	}
	log.Println("Соединение с MQTT установлено успешно.")
}

func (v *ValidationServiceImplementation) SubscribeToTopic(ctx context.Context, topic string, callback func(string), wg *sync.WaitGroup) {
	defer wg.Done()
	err := v.mqttServices.GetSubscribe(topic, callback)
	if err != nil {
		log.Printf("Ошибка подписки на топик %s: %v", topic, err)
		return
	}
	log.Printf("Подписка на топик %s успешна.", topic)
}

func (v *ValidationServiceImplementation) HandleMqttMessage(mqttMessage string, camNumber int) {
	err := v.mqttServices.ImplementQueryProcedure(mqttMessage, camNumber)
	if err != nil {
		log.Printf("Ошибка в обработке сообщения для камеры %d: %v", camNumber, err)
	}
}

func (v *ValidationServiceImplementation) PublishResult(camNumber int, eventType, grz string) {
	err := v.mqttServices.PublishResultProcedure(camNumber, eventType, grz)
	if err != nil {
		log.Printf("Ошибка публикации результата для камеры %d: %v", camNumber, err)
	} else {
		log.Printf("Результат для камеры %d успешно опубликован: событие %s, GRZ %s", camNumber, eventType, grz)
	}
}
