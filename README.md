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


## Exemplo como agregar em seu repositorio 

```yml
name: Trigger llm

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  trigger-ia:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write

    steps:
      - name: Check out the repository
        uses: actions/checkout@v2
        with:
          fetch-depth: 0
          
      - name: Get changes in files
        id: capture_changes
        run: |
            files=$(git diff --name-only FETCH_HEAD HEAD)
          
            echo "Changed files: $files" 
          
            if [ -z "$files" ]; then
                echo "No files changed."
                echo "fileChanges={}" >> $GITHUB_ENV
                exit 0
            fi
          
            jsonOutput="{"
          
            for file in $files; do
              currentContent=$(cat "$file" || echo "Cannot read file")
              changes=$(git diff --unified=0 HEAD^ HEAD -- "$file" || echo "No changes detected")
              fileChangesJson=$(jq -c -n --arg current "$currentContent" --arg changes "$changes" '{current: $current, changes: $changes}')
              jsonOutput+="\"$file\": $fileChangesJson,"
            done
          
            jsonOutput="${jsonOutput%,}}"
            echo "Constructed JSON for fileChanges: $jsonOutput"
            printf "fileChanges=%s\n" "$jsonOutput" >> $GITHUB_ENV
        
      - name: Code review in pull request
        uses: ezequielbugnon/llm-code-review-action/@main
        with:
          go_version: '1.23'
          client_id: ${{ secrets.CLIENTID }}
          client_secret: ${{ secrets.CLIENTSECRET }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          files_changes: ${{ env.fileChanges }}
```

## Contribuições

Contribuições são bem-vindas. Se você deseja colaborar, por favor abra um pull request ou issue para discutir suas ideias.
