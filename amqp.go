package amqp

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

type (
	AMQP struct{}
)

func init() {
	modules.Register("k6/x/amqp", new(AMQP))
}

// Conectar ao RabbitMQ e obter utilizando as informações e credenciais recebidas
// x
func (a *AMQP) Connect(host, vhost, username, password string) (*amqp.Connection, error) {
	// Montar a URL de conexão AMQP com usuário e senha
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost)
	redactedURL := fmt.Sprintf("amqp://%s:*********@%s/%s", username, host, vhost)
	
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		fmt.Printf("Erro ao conectar ao RabbitMQ (%S): %v\n", err, redactedURL)
		return conn, err
	}

	return conn, nil
}

// Encerrar a conexão 
func (a *AMQP) Disconnect(conn amqp.Connection) {
	conn.Close()
}

// Abrir e obter um canal na conexão fornecida
func (a *AMQP) OpenChannel(conn amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Erro ao abrir o canal: %v\n", err)
		return ch, err
	}

	return ch, nil
}

// Fechar o canal
func (a *AMQP) CloseChannel(ch amqp.Channel) {
	ch.Close()
}

// Publica uma mensagem em um exchange específico do RabbitMQ usando as informações fornecidas
func (a *AMQP) Publish(host, vhost, username, password, exchange, routingKey, message string) error {
	conn, err := a.Connect(host, vhost, username, password)
	if err != nil {
		return err
	}

	err = a.PublishUsingConn(*conn, exchange, routingKey, message)
	
	defer a.Disconnect(*conn);

	return err
}

// Publica uma mensagem em um exchange específico do RabbitMQ usando a conexão e informações fornecida
func (a *AMQP) PublishUsingConn(conn amqp.Connection, exchange, routingKey, message string) error {
	ch, err := a.OpenChannel(conn)
	if err != nil {
		return err
	}

	err = a.PublishUsingChannel(*ch, exchange, routingKey, message)

	defer a.CloseChannel(*ch)

	return err
}

// Publica uma mensagem em um exchange específico do RabbitMQ usando a conexão fornecida
func (a *AMQP) PublishUsingChannel(ch amqp.Channel, exchange, routingKey, message string) error {
	// Publicar a mensagem no exchange com a routing key especificada
	return ch.Publish(
		exchange,   // Nome do exchange
		routingKey, // Routing key
		false,      // Mandatory (setando para false, se não houver um binding para a rk, a mensagem será descartada silenciosamente)
		false,      // Immediate (false, significa que a mensagem será enfileirada mesmo se não houver algum consumidor ativo para a(s) fila(s) bindadas)
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Timestamp:   time.Now(),
		},
	)
}
