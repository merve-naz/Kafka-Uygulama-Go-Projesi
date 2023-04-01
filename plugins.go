package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	// cache
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/config"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/models"

	// in-app-cache
	"github.com/gin-contrib/cache/persistence"

	// s3
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewS3Session(conf *models.S3Config) *session.Session {
	// open s3 session
	creds := credentials.NewStaticCredentials(
		conf.AccessKey,
		conf.SecretKey,
		"",
	)
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(conf.Endpoint),
		Region:           aws.String(conf.Region),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Fatalln("error_create_s3_session: ", err)
	}
	return sess
}

func NewRedisCacheConnection(redisUrl string) (*redis.Client, context.Context) {
	// redis://username:password@host:port/db
	_, username, password, host, port, db := UrlStringToOptions(redisUrl)

	// convert db to int
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalln("INIT: redis connection database is not integer ", err)
	}

	rContext := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Username: username,
		Password: password,
		DB:       dbInt,
	})
	// ping redis for check connection
	_, err = rdb.Ping(rContext).Result()
	if err != nil {
		log.Fatalln("INIT: redis ping request failed ", err)
	}

	return rdb, rContext
}

func NewInAppCacheStore(defaultTTL time.Duration) *persistence.InMemoryStore {
	return persistence.NewInMemoryStore(defaultTTL)
}

func NewKafkaConsumerConnection(url string, cg string, topic string) *kafka.Reader {
	log.Println("INIT: kafka consumer connection is initializing...")

	// split all options (app_name, username, password, host, port, db)
	_, username, password, host, port, partitionVal := UrlStringToOptions(url)

	// partition is digit string. cast to int.
	partition, err := strconv.Atoi(partitionVal)
	if err != nil {
		log.Fatalln("INIT: kafka consumer partition is not digit string ", err)
	}

	// create kafka dialer
	dialer := &kafka.Dialer{
		Timeout:   time.Minute,
		DualStack: true,
		// SASLMechanism: mechanism,
	}

	// if env is "local" skip scram-auth.
	if config.C.App.Env != "local" {
		// scram.Mechanism
		mechanism, err := scram.Mechanism(scram.SHA512, username, password)
		if err != nil {
			log.Fatal("INIT: failed to kafka scram auth:", err)
		}

		dialer.SASLMechanism = mechanism
	}

	// create reader (consumer)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			host + ":" + port,
		},
		// kafka generic opts.
		GroupID:   cg,
		Topic:     topic,
		Dialer:    dialer,
		Partition: partition,
		// no wait for message bucket to be full.
		MinBytes: 1,   // 1 byte
		MaxBytes: 1e6, // 1MB
	})

	// log success
	log.Println("INIT: kafka consumer connection is initialized successfully.")

	return r
}

func NewKafkaProducerConnection(url string, cg string, topic string) *kafka.Writer {
	log.Println("INIT: kafka producer connection is initializing...")

	// split all options (app_name, username, password, host, port, db)
	_, username, password, host, port, _ := UrlStringToOptions(url)

	// create kafka dialer
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		// SASLMechanism: mechanism,
	}

	// if env is "local" skip scram-auth.
	if config.C.App.Env != "local" {
		// scram.Mechanism
		log.Println("INIT: kafka connection scram.Mechanism is initializing...")
		mechanism, err := scram.Mechanism(scram.SHA512, username, password)
		if err != nil {
			log.Fatal("INIT: kafka connection scram.Mechanism is failed: ", err)
		}

		dialer.SASLMechanism = mechanism
	}

	// create producer writer
	log.Println("INIT: kafka producer writer connection is initializing...")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{
			host + ":" + port,
		},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
		Async:    false,
	})

	// log success
	log.Println("INIT: kafka producer connection is initialized successfully.")

	return w
}

func UrlStringToOptions(url string) (string, string, string, string, string, string) {

	options := strings.Split(url, "://")

	// split protocol and info
	protocol := options[0]
	info := options[1]

	// split info to username, password, host, port, db
	infoOptions := strings.Split(info, "@")
	// split username and password
	usernamePassword := infoOptions[0]
	hostPortDb := infoOptions[1]

	usernamePasswordOptions := strings.Split(usernamePassword, ":")
	username := usernamePasswordOptions[0]
	password := usernamePasswordOptions[1]
	// split host, port and db
	hostPortDbOptions := strings.Split(hostPortDb, "/")
	hostPort := hostPortDbOptions[0]
	db := hostPortDbOptions[1]
	// split host and port
	hostPortOptions := strings.Split(hostPort, ":")
	host := hostPortOptions[0]
	port := hostPortOptions[1]

	return protocol, username, password, host, port, db
}
