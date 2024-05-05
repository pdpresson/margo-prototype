package clients

import (
	"context"
	"errors"
	"fmt"
	"gitops_pushservice/appConfig"
	"time"

	"github.com/go-logr/logr"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MQClient struct {
	log             *logr.Logger
	appConfig       *appConfig.AppConfig
	queueName       string
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	resendDelay    = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("client is shutting down")
)

// init will initialize channel & declare queue
func (client *MQClient) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	err = ch.Confirm(false)

	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		client.queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	client.changeChannel(ch)
	client.isReady = true
	client.log.Info("Setup!")

	return nil
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (client *MQClient) handleReconnect(addr string) {
	for {
		client.isReady = false
		client.log.Info("Attempting to connect")

		conn, err := client.connect(addr)

		if err != nil {
			client.log.Error(err, "Failed to connect. Retrying...")

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (client *MQClient) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.log.Info("Connected!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (client *MQClient) handleReInit(conn *amqp.Connection) bool {
	for {
		client.isReady = false

		err := client.init(conn)

		if err != nil {
			client.log.Error(err, "Failed to initialize channel. Retrying...")

			select {
			case <-client.done:
				return true
			case <-client.notifyConnClose:
				client.log.Info("Connection closed. Reconnecting...")
				return false
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.log.Info("Connection closed. Reconnecting...")
			return false
		case <-client.notifyChanClose:
			client.log.Info("Channel closed. Re-running init...")
		}
	}
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (client *MQClient) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (client *MQClient) changeChannel(channel *amqp.Channel) {
	client.channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.channel.NotifyClose(client.notifyChanClose)
	client.channel.NotifyPublish(client.notifyConfirm)
}

// Push will push data onto the queue, and wait for a confirm.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (client *MQClient) Push(data []byte) error {
	if !client.isReady {
		return errors.New("failed to push: not connected")
	}
	for {
		err := client.UnsafePush(data)
		if err != nil {
			client.log.Error(err, "Push failed. Retrying...")
			select {
			case <-client.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		confirm := <-client.notifyConfirm
		if confirm.Ack {
			client.log.Info(fmt.Sprintf("Push confirmed [%d]!", confirm.DeliveryTag))
			return nil
		}
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (client *MQClient) UnsafePush(data []byte) error {
	if !client.isReady {
		return errNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return client.channel.PublishWithContext(
		ctx,
		"",
		client.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

// Consume will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (client *MQClient) Consume() (<-chan amqp.Delivery, error) {
	if !client.isReady {
		return nil, errNotConnected
	}

	if err := client.channel.Qos(
		1,
		0,
		false,
	); err != nil {
		return nil, err
	}

	return client.channel.Consume(
		client.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

// Close will cleanly shut down the channel and connection.
func (client *MQClient) Close() error {
	if !client.isReady {
		return errAlreadyClosed
	}
	close(client.done)
	err := client.channel.Close()
	if err != nil {
		return err
	}
	err = client.connection.Close()
	if err != nil {
		return err
	}

	client.isReady = false
	return nil
}

func NewMQClient(l *logr.Logger, c *appConfig.AppConfig, queueName string) *MQClient {
	client := MQClient{
		log:       l,
		appConfig: c,
		queueName: queueName,
		done:      make(chan bool),
	}
	go client.handleReconnect(c.MQAddress)
	return &client
}
