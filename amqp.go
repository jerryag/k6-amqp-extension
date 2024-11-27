package amqp

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

const version = "v0.0.1"

// Estrutura principal da extensão AMQP
type AMQP struct{
	version  string
}

func init() {
	oAmqp := AMQP{
		Version:     version,
	}

	// Registrar o módulo personalizado no k6
	modules.Register("k6/x/amqp", &oAmqp)
}

// PublishToExchange publica uma mensagem em um exchange específico do RabbitMQ usando autenticação
func (a *AMQP) PublishToExchange(host, vhost, username, password, exchange, routingKey, message string) {
	// Montar a URL de conexão AMQP com usuário e senha
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost)
	redactedURL := fmt.Sprintf("amqp://%s:*********@%s/%s", username, host, vhost)

	// Conectar ao RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		fmt.Printf("Erro ao conectar ao RabbitMQ (%S): %v\n", err, redactedURL)
		return
	}
	defer conn.Close()

	// Abrir um canal
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Erro ao abrir o canal: %v\n", err)
		return
	}
	defer ch.Close()

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
