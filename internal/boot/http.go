package boot

import (
	"log"
	"net/http"

	"go-tutorial-2020/internal/config"

	userData "go-tutorial-2020/internal/data/user"
	server "go-tutorial-2020/internal/delivery/http"
	userHandler "go-tutorial-2020/internal/delivery/http/user"
	userService "go-tutorial-2020/internal/service/user"
	"go-tutorial-2020/pkg/firebaseclient"

	"github.com/jmoiron/sqlx"

	"go-tutorial-2020/pkg/kafka"

	kConsumer "go-tutorial-2020/internal/delivery/kafka"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
func HTTP() error {
	var (
		fc  *firebaseclient.Client // Firebase initiation
		s   server.Server          // HTTP Server Object
		ud  userData.Data          // User domain data layer
		us  userService.Service    // User domain service layer
		uh  *userHandler.Handler   // User domain handler
		cfg *config.Config         // Configuration object
		k   *kafka.Kafka           // Kafka Producer
	)

	// Get configuration
	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg = config.Get()

	// Open MySQL DB Connection
	db, err := sqlx.Open("mysql", cfg.Database.Master)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	fc, err = firebaseclient.NewClient(cfg)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	k, err = kafka.New(cfg.Kafka.Username, cfg.Kafka.Password, cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("[KAFKA] Failed to initialize kafka producer: %v", err)
	}

	// User domain initialization
	ud = userData.New(db, fc)
	us = userService.New(ud, k)
	uh = userHandler.New(us)

	// Inject service used on handler
	s = server.Server{
		User: uh,
	}

	go kConsumer.New(us, k, cfg.Kafka.Subscriptions)
	// Error Handling
	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
