package config

import (
	"testing"
)

type Config struct {
	RabbitMQ RabbitMQConf `toml:"RabbitMQ"`
	Canopsis CanopsisConf `toml:"Canopsis"`
}

func (c *Config) Clone() interface{} {
	cloned := *c
	return &cloned
}

func TestOverride(t *testing.T) {
	replacementIntValue := int64(0)
	replacementStrValue := "replaced-string-value"
	replacementBoolValue := false
	replacementArrayValue := []interface{}{
		"item2",
		"item1",
	}

	config := &Config{
		Canopsis: CanopsisConf{
			Alarm: SectionAlarm{
				StealthyInterval:    255,
				DisplayNameScheme:   "intial-string-value",
				EnableLastEventDate: false,
			},
			Logger: SectionLogger{
				ConsoleWriter: ConsoleWriter{
					PartsOrder: []string{
						"item1",
						"item2",
					},
				},
			},
		},
	}

	overrideConf := map[string]interface{}{
		"Canopsis": map[string]interface{}{
			"alarm": map[string]interface{}{
				"StealthyInterval":    replacementIntValue,
				"DisplayNameScheme":   replacementStrValue,
				"EnableLastEventDate": replacementBoolValue,
			},
			"logger": map[string]interface{}{
				"console_writer": map[string]interface{}{
					"PartsOrder": replacementArrayValue,
				},
			},
		},
	}

	var err error
	var newPtr interface{}
	if newPtr, err = Override(config, overrideConf); err != nil {
		t.Fatalf("failed overriding config. err: %v", err)
	}
	config = newPtr.(*Config)

	if config.Canopsis.Alarm.StealthyInterval != int(replacementIntValue) {
		t.Errorf("expected value: %v but got %v", replacementIntValue, config.Canopsis.Alarm.StealthyInterval)
	}

	if config.Canopsis.Alarm.DisplayNameScheme != replacementStrValue {
		t.Errorf("expected value: %v but got %v", replacementStrValue, config.Canopsis.Alarm.DisplayNameScheme)
	}

	if config.Canopsis.Alarm.EnableLastEventDate != replacementBoolValue {
		t.Errorf("expected value: %v but got %v", replacementBoolValue, config.Canopsis.Alarm.EnableLastEventDate)
	}

	if config.Canopsis.Logger.ConsoleWriter.PartsOrder[0] != replacementArrayValue[0] {
		t.Errorf("expected value: %v but got %v",
			replacementArrayValue[0],
			config.Canopsis.Logger.ConsoleWriter.PartsOrder[0],
		)
	}
}

func TestOverride_Array(t *testing.T) {
	replacementValue := "exchange-3"

	config := &Config{
		RabbitMQ: RabbitMQConf{
			Exchanges: []Exchange{
				{Name: "exchange-1"},
				{Name: "exchange-2"},
			},
		},
	}
	initialLen := len(config.RabbitMQ.Exchanges)

	overrideConf := map[string]interface{}{
		"RabbitMQ": map[string]interface{}{
			"exchanges": []map[string]interface{}{
				{
					"name": replacementValue,
				},
			},
		},
	}

	var err error
	var newPtr interface{}
	if newPtr, err = Override(config, overrideConf); err != nil {
		t.Fatalf("failed overriding config. err: %v", err)
	}
	config = newPtr.(*Config)

	if initialLen+1 != len(config.RabbitMQ.Exchanges) {
		t.Errorf("expected value: %v but got %v", initialLen+1, len(config.RabbitMQ.Exchanges))
	}

	if config.RabbitMQ.Exchanges[initialLen].Name != replacementValue {
		t.Errorf("expected value: %v but got %v", replacementValue, config.RabbitMQ.Exchanges[initialLen].Name)
	}
}

func TestOverride_AbortWhenError(t *testing.T) {
	initialValue := "initial-string-value"

	config := &Config{
		RabbitMQ: RabbitMQConf{
			Queues: []Queue{
				{Name: "queue-1"},
				{Name: "queue-2"},
			},
		},
		Canopsis: CanopsisConf{
			Alarm: SectionAlarm{
				DisplayNameScheme: initialValue,
			},
		},
	}
	initialLen := len(config.RabbitMQ.Queues)

	overrideConf := map[string]interface{}{
		"RabbitMQ": map[string]interface{}{
			"queues": []map[string]interface{}{
				{
					"name": "queue-3",
				},
			},
		},
		"Canopsis": map[string]interface{}{
			"alarm": map[string]interface{}{
				"DisplayNameScheme": "replaced-string-value",
				"not-matching-tag":  "",
			},
		},
	}

	var err error
	var newPtr interface{}
	if newPtr, err = Override(config, overrideConf); err == nil {
		t.Errorf("failed overriding config. err: %v", err)
	}
	config = newPtr.(*Config)

	if config.Canopsis.Alarm.DisplayNameScheme != initialValue {
		t.Errorf("expected value: %v but got %v", initialValue, config.Canopsis.Alarm.DisplayNameScheme)
	}

	if initialLen != len(config.RabbitMQ.Queues) {
		t.Errorf("expected value: %v but got %v", initialLen, len(config.RabbitMQ.Queues))
	}
}
