# GoChat

GoChat é uma aplicação de chat inteligente que permite aos usuários fazer upload de seus próprios documentos e interagir com eles. Utilizando modelos de linguagem avançados, o GoChat analisa o conteúdo dos documentos e responde a perguntas, transformando seus arquivos em uma base de conhecimento interativa.

## Para o Usuário Final

Se você deseja apenas utilizar a aplicação, a maneira mais simples é através do Docker.

### Funcionalidades

*   **Autenticação de Usuário:** Crie sua conta e faça login para ter acesso privado aos seus documentos.
*   **Upload de Documentos:** Envie seus documentos em formato de texto.
*   **Chat Interativo:** Faça perguntas em linguagem natural e receba respostas baseadas no conteúdo dos seus documentos.

### Como Executar com Docker

Com o Docker e o Docker Compose instalados, basta executar o seguinte comando na raiz do projeto:

```bash
docker-compose up --build
```

A aplicação estará disponível em `http://localhost:8080`.

## Para Desenvolvedores

Esta seção fornece informações técnicas para quem deseja entender o funcionamento do projeto ou contribuir com o desenvolvimento.

### Tecnologias Utilizadas

*   **Backend:** Go
*   **Framework Web:** Gin
*   **Banco de Dados:** SQLite (com GORM como ORM)
*   **Containerização:** Docker

### Estrutura do Projeto

O projeto segue uma estrutura organizada para separar as responsabilidades:

-   `cmd/api`: Ponto de entrada da aplicação.
-   `internal/`: Contém toda a lógica de negócio.
    -   `config`: Gerenciamento de configurações.
    -   `database`: Conexão e migrações do banco de dados.
    -   `handlers`: Manipuladores de requisições HTTP.
    -   `middleware`: Middlewares para as rotas.
    -   `models`: Estruturas de dados do banco de dados.
    -   `services`: Lógica de negócio e integrações com serviços externos (como OpenAI).
-   `data/`: Armazenamento do banco de dados SQLite.
-   `Dockerfile`: Define a imagem Docker para a aplicação.
-   `docker-compose.yml`: Orquestra os contêineres da aplicação.

### API Endpoints

| Método | Rota          | Descrição                                     | Autenticação |
| ------ | ------------- | --------------------------------------------- | ------------ |
| POST   | /register     | Registra um novo usuário.                     | Não          |
| POST   | /login        | Autentica um usuário e retorna um token JWT.  | Não          |
| POST   | /documents    | Faz o upload de um novo documento.            | Sim          |
| GET    | /documents    | Lista os documentos do usuário autenticado.   | Sim          |
| POST   | /ask          | Envia uma pergunta sobre um documento.        | Sim          |
| GET    | /health       | Verifica o status da aplicação.               | Não          |

### Como as Perguntas Funcionam (Chat com Memória)

A funcionalidade de perguntas e respostas foi projetada para simular uma conversa real. O sistema armazena o histórico de perguntas e respostas de uma sessão, permitindo que você faça perguntas de acompanhamento que dependem do contexto anterior.

Para que o histórico funcione, é essencial enviar um identificador único para cada "usuário" ou "sessão" de chat através do campo `external_id` no corpo da requisição para o endpoint `/ask`. O GoChat usará esse ID para recuperar o contexto das perguntas e respostas anteriores daquela sessão específica.

Por exemplo, você pode perguntar "Quem foi o primeiro presidente do Brasil?" e, em seguida, na próxima requisição com o mesmo `external_id`, perguntar "E quantos anos ele tinha?". O sistema entenderá que "ele" se refere ao presidente mencionado na pergunta anterior.

### Exemplos de Requisições

**Observação:** Substitua `<seu_token_jwt>` pelo token JWT obtido no login e os caminhos de arquivo/IDs conforme necessário.

#### Registrar Novo Usuário

```bash
curl --location 'http://localhost:8080/register' \
--header 'Content-Type: application/json' \
--data '{
  "username": "testuser",
  "password": "password123"
}'
```

#### Login

```bash
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
  "username": "testuser",
  "password": "password123"
}'
```
*A resposta incluirá o token JWT para ser usado nas rotas autenticadas.*

#### Upload de Documento

Para que o chat possa responder às perguntas, você deve primeiro fazer o upload de um documento (`.md`) contendo a base de conhecimento (regras de negócio, informações, etc.).

```bash
curl --location 'http://localhost:8080/documents' \
--header 'Authorization: Bearer <seu_token_jwt>' \
--form 'file=@"/path/to/your/document.md"'
```

#### Listar Documentos

```bash
curl --location 'http://localhost:8080/documents' \
--header 'Authorization: Bearer <seu_token_jwt>'
```

#### Fazer uma Pergunta

Ao fazer uma pergunta, você deve fornecer o `document_id` (obtido na listagem ou no upload), a `question` e o `external_id` para identificar a sessão da conversa.

```bash
curl --location 'http://localhost:8080/ask' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <seu_token_jwt>' \
--data '{
  "document_id": "1",
  "question": "Qual o dia do Natal?",
  "external_id": "user_001"
}'
```

### Executando Localmente

Para executar o projeto em seu ambiente de desenvolvimento local:

1.  **Instale as dependências:**

    ```bash
    go mod tidy
    ```

2.  **Crie um arquivo `.env`** na raiz do projeto com as seguintes variáveis (utilize o `.env.example` como base):

    ```
    PORT=8080
    SECRET_KEY=sua_chave_secreta
    OPENAI_API_KEY=sua_chave_da_openai
    ```

3.  **Execute a aplicação:**

    ```bash
    go run cmd/api/main.go
    ```

A aplicação estará disponível em `http://localhost:8080`.