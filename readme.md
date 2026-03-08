# Desafio Client-Server-API

Componentes do Sistema Você precisará entregar dois arquivos principais:

client.go
server.go

## Descrição Server
### Arquivo server.go
Requisitos Técnicos: Server.go O servidor HTTP deve operar na porta 8080 e expor o endpoint /cotacao.

Consumo de API Externa: 
- Ao receber uma requisição em /cotacao, o server deve consumir a API de Câmbio: https://economia.awesomeapi.com.br/json/last/USD-BRL.
Timeout: 
- O timeout máximo para chamar essa API externa deve ser de 200ms (usando o pacote context).

Persistência (Banco de Dados): O servidor deve registrar cada cotação recebida em um banco de dados SQLite.
Timeout: 
- O timeout máximo para persistir os dados no banco deve ser de 10ms (usando o pacote context).
Resposta: 
- O endpoint deve retornar o resultado da cotação em formato JSON para o cliente.
Logs: 
- Caso os timeouts (API ou Banco) sejam excedidos, o erro deve ser logado no console do servidor.

## Descrição Client
### Arquivo client.go
Requisitos Técnicos: Client.go O cliente deve realizar uma requisição HTTP ao server.go.

Requisição: 
- Deve solicitar a cotação ao endpoint /cotacao do servidor local.

Timeout: 
- O timeout máximo para receber o resultado do servidor deve ser de 300ms (usando o pacote context).

Processamento e Arquivo: 
- O cliente deve receber apenas o valor atual do câmbio (campo bid do JSON).
- Deve salvar a cotação em um arquivo chamado cotacao.txt.
- Formato do arquivo: Dólar: {valor}

## Execução
### Via Compose
O projeto tem um compose onde ele irá preparar um sqlite, executar o servidor do server.go e ao final, executar o script do client.go

### Rodar a cada arquivo
Você tambem pode rodar o projeto server.go primeiro e depois de iniciar o mesmo, rodar o script client.go, ambos os projetos vao criar:
- server.go: Deve criar a pasta data/cotacao.db
- client.go: Deve criar na raiz do projeto o cotacao.txt

Logs:
- Caso o timeout de 300ms seja excedido, o erro deve ser logado no console do cliente.
