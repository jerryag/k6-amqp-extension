import amqp from 'k6/x/amqp';  // Importa a extensão AMQP

export default function () {
    // Configurações de conexão AMQP
    const host = 'localhost:5672';        // Endereço do RabbitMQ
    const vhost = 'test-vhost';           // Vhost a ser utilizado
    const username = 'guest';             // Usuário padrão do RabbitMQ
    const password = 'guest';             // Senha padrão do RabbitMQ
    const exchange = 'test-exchange';     // Nome do exchange existente
    const routingKey = 'test-routing-key';// Routing key usada
    const message = `Mensagem do K6: ${__VU} - ${__ITER}`; // Mensagem com VU e Iteração

    console.log("AMQP: ", amqp);

    // Publicar a mensagem no exchange
    //amqp.PublishToExchange(host, vhost, username, password, exchange, routingKey, message);

    console.log(`Mensagem enviada: ${message}`);
}
