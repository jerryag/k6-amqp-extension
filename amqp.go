package amqp

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

const version = "0.0.1"

// Estrutura principal da extensão AMQP
type AMQP struct {
	Version string
}

func init() {
	a := AMQP {
		Version: version
	}

	// Registrar o módulo personalizado no k6
	modules.Register("k6/x/amqp", &a)
}

// Conectar ao RabbitMQ e obter utilizando as informações e credenciais recebidas
func (a *AMQP) Connect(host, vhost, username, password) {
	// Montar a URL de conexão AMQP com usuário e senha
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost)
	redactedURL := fmt.Sprintf("amqp://%s:*********@%s/%s", username, host, vhost)
	
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		fmt.Printf("Erro ao conectar ao RabbitMQ (%S): %v\n", err, redactedURL)
		return nil, err
	}

	return conn, nil
}

// Encerrar a conexão 
func (a *AMQP) Disconnect(conn) {
	conn.Close()
}

// Abrir e obter um canal na conexão fornecida
func (a *AMQP) OpenChannel(conn) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Erro ao abrir o canal: %v\n", err)
		return nil, err
	}

	return ch, nil
}

// Fechar o canal
func (a *AMQP) CloseChannel(ch) {
	ch.Close()
}

// Publica uma mensagem em um exchange específico do RabbitMQ usando a conexão fornecida
func (a *AMQP) PublishToExchange(host, vhost, username, password, exchange, routingKey, message string) {
	conn, err := Connect(host, vhost, username, password)
	if err != nil {
		return
	}

	PublishToExchange(conn, exchange, routingKey, message)
	
	defer Disconnect;
}

// Publica uma mensagem em um exchange específico do RabbitMQ usando a conexão fornecida
func (a *AMQP) PublishToExchange(conn, exchange, routingKey, message string) {
	ch, err := OpenChannel(conn)
	if err != nil {
		return
	}

	PublishToExchange(ch, exchange, routingKey, message)

	defer CloseChannel(ch)
}

// Publica uma mensagem em um exchange específico do RabbitMQ usando a conexão fornecida
func (a *AMQP) PublishToExchange(ch, exchange, routingKey, message string) {
	// Publicar a mensagem no exchange com a routing key especificada
	err = ch.Publish(
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
	if err != nil {
		fmt.Printf("Erro ao publicar a mensagem: %v\n", err)
		return
	}

	fmt.Printf("Mensagem publicada com sucesso no exchange '%s' com routing key '%s'\n", exchange, routingKey)
}
