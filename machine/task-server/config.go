package task_server

import "github.com/RichardKnop/machinery/v1/config"

func LoadInMemConfig() *config.Config {
	return &config.Config{
		Broker:        "amqp://guest:guest@localhost:5672/",
		ResultBackend: "redis://127.0.0.1:6379",
		DefaultQueue:  "jakob_tasks",
		AMQP: &config.AMQPConfig{
			Exchange:      "jakob_exchange",
			ExchangeType:  "direct",
			BindingKey:    "jakob_task",
			PrefetchCount: 1,
		},
	}
}
