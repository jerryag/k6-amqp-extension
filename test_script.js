import amqp from 'k6/x/amqp';  // Importa a extensão AMQP

// Configurações de conexão AMQP
const host = `${__ENV.HOST}`;              // Endereço do RabbitMQ
const vhost = `${__ENV.VHOST}`;            // Vhost a ser utilizado
const username = `${__ENV.USERNAME}`;      // Usuário padrão do RabbitMQ
const password = `${__ENV.PASSWORD}`;      // Senha padrão do RabbitMQ
const exchange = `${__ENV.EXCHANGE}`;      // Nome do exchange existente
const routingKey = `${__ENV.ROUTING_KEY}`; // Routing key usada

// Cenário de execução
const scenarios = {
    scenario1 : {
      executor: 'constant-vus',
      vus: 1,
      duration: '5s',
      exec: 'executeScenario1'
    }
}

// Options de execução do k6 - Definição da dinâmica de execucao dos cenários de testes e outros parâmetros
export let options = { scenarios }

// Conectar com o RabbitMQ e abrir o Channel
const conn = amqp.connect(host, vhost, username, password);
const ch = amqp.openChannel(conn);

// Executores dos cenários
export function executeScenario1() {
    const message = `Mensagem do K6: ${__VU} - ${__ITER}`; 
    amqp.publishUsingChannel(ch, exchange, routingKey, message);
}
