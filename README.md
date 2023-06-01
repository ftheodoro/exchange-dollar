# Exchange Dollar
Desafio da primeira fase do curso golang expert(https://goexpert.fullcycle.com.br/curso/)

Basicamente era criar um cliente e um Servidor.
Servidor serviria como um proxy para pegar a cotação do dolar e cada requisição é salva no banco de dados sqlite ,timeout de request é de 200ms e de salvar no banco é de 10ms.
Cliente consume os dados do servidor e tem que salvar em um arquivo .txt tempo de request do cliente para o servidor é de 300ms.
