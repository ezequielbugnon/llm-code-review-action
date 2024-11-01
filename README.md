# Code Review com Stack Spot llm

Este projeto utiliza um script em Go para realizar análises de código e gerar comentários em pull requests no GitHub. A análise é realizada por meio de uma API de revisão de código.

## Dependências
Go 1.23 ou superior

## Uso

Este projeto foi desenvolvido para ser executado automaticamente através do GitHub Actions. Ele é acionado ao abrir ou sincronizar um pull request. O fluxo de trabalho realiza o seguinte:

- Checkout do código: Obtém o código do repositório.
- Configuração do Go: Configura o ambiente Go.
- Instalação de dependências: Executa go mod tidy para gerenciar as dependências.
- Execução do analisador: Executa o script principal (main.go) que realiza a revisão do código usando a API da StackSpot.
- Publicação da revisão: Publica o resultado da análise como um comentário no pull request.

## Variáveis de Ambiente

Certifique-se de configurar as seguintes variáveis de ambiente nos segredos do seu repositório GitHub:

- CLIENTID: Seu Client ID para a API da StackSpot.
- CLIENTSECRET: Seu Client Secret para a API da StackSpot.
- GITHUB_TOKEN: Token de acesso para publicar comentários no pull request.

## Contribuições

Contribuições são bem-vindas. Se você deseja colaborar, por favor abra um pull request ou issue para discutir suas ideias.