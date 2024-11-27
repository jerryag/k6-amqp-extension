import amqp from 'k6/x/amqp';  // Importa a extensão AMQP

export default function () {
    // Configurações de conexão AMQP
    const host = `${__ENV.HOST}`;        // Endereço do RabbitMQ
    const vhost = `${__ENV.VHOST}`;           // Vhost a ser utilizado
    const username = `${__ENV.USERNAME}`;             // Usuário padrão do RabbitMQ
    const password = `${__ENV.PASSWORD}`;             // Senha padrão do RabbitMQ
    const exchange = `${__ENV.EXCHANGE}`;     // Nome do exchange existente
    const routingKey = `${__ENV.ROUTING_KEY}`;// Routing key usada
    const message = `Mensagem do K6: ${__VU} - ${__ITER}`; // Mensagem com VU e Iteração

    // Publicar a mensagem no exchange
    amqp.publishToExchange(host, vhost, username, password, exchange, routingKey, message);

    console.log(`Mensagem enviada: ${message}`);
}
