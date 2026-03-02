# Desafio Client-Server-API

## Descrição Server
### Arquivo server.go
Ao receber uma requisição em /cotacao, o server deve consumir a API de Cambio: https://economia.awesomeapi.com.br/json/last/USD-BRL.
    Timeout: O timeout maximo para chamar essa API externa deve ser de 200ms (usando o pacote context)
    Persistencia (Banco de Dados):

O servidor deve registrar cada cotação recebida em um banco de dados SQLite
	Timeout: O timeout m&aacute;ximo para persistir os dados no banco deve ser de 10ms (usando o pacote context).
	Resposta: O endpoint deve retornar o resultado da cotação em formato JSON para o cliente.
	Logs: Caso os timeouts (API ou Banco) sejam excedidos, o erro deve ser logado no console do servidor.


## Descrição Client
### Arquivo client.go
Deve solicitar a cotação ao endpoint /cotacao do servidor local.
	Timeout: O timeout máximo para receber o resultado do servidor deve ser de 300ms (usando o pacote context).
	
Processamento e Arquivo:
	O cliente deve receber apenas o valor atual do cambio (campo bid do JSON).
    Deve salvar a cotação em um arquivo chamado cotacao.txt.
	Formato do arquivo: Dólar: {valor}
	Logs: Caso o timeout de 300ms seja excedido, o erro deve ser logado no console do cliente.