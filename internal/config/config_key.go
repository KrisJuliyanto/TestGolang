package config

type (
	// Config ...
	Config struct {
		Server   ServerConfig   `yaml:"server"`
		Database DatabaseConfig `yaml:"database"`
		Firebase FirebaseConfig `yaml:"firebase"`
		Kafka    KafkaConfig    `yaml:"kafka"`
	}

	// ServerConfig ...
	ServerConfig struct {
		Port string `yaml:"port"`
	}

	// FirebaseConfig ...
	FirebaseConfig struct {
		ProjectID string `yaml:"projectID"`
	}

	// DatabaseConfig ...
	DatabaseConfig struct {
		Master string `yaml:"master"`
	}

	//KafkaConfig ...
	KafkaConfig struct {
		Username      string   `yaml:"username"`
		Password      string   `yaml:"password"`
		Brokers       []string `yaml:"brokers"`
		Subscriptions []string `yaml:"subscriptions"`
	}
)
