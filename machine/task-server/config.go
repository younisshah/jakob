package task_server

import "github.com/RichardKnop/machinery/v1/config"

/**
*  Created by Galileo on 20/6/17.
 */

func LoadInMemConfig() *config.Config {
	return &config.Config{
		Broker:             "amqp://guest:guest@localhost:5672/",
		ResultBackend:      "redis://127.0.0.1:6379",
		MaxWorkerInstances: 10,
		DefaultQueue:       "jakob_tasks",
		AMQP: &config.AMQPConfig{
			Exchange:      "jakob_exchange",
			ExchangeType:  "direct",
			BindingKey:    "jakob_task",
			PrefetchCount: 1,
		},
	}
}
