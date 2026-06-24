# go-clean-architecture

Projeto de estudo de Clean Architecture em Go, implementando um sistema de pedidos (orders) com três interfaces de comunicação: REST HTTP, gRPC e GraphQL.

## Conceitos de Clean Architecture

O projeto é organizado em camadas com dependências sempre apontando para o centro, nunca para fora.

```
cmd/
└── ordersystem/              # Ponto de entrada — monta dependências e sobe os servidores

internal/
├── entity/                   # Camada de domínio
│   ├── order.go              # Entidade Order com regras de negócio (validação, FinalPrice)
│   ├── order_repository.go   # Interface do repositório (contrato, sem implementação)
│   └── errors.go
│
├── usecase/                  # Casos de uso — orquestram o domínio
│   ├── create_order.go
│   └── list_orders.go
│
└── infra/                    # Adaptadores externos (implementações concretas)
    ├── database/             # Implementação MySQL do OrderRepository
    ├── grpc/                 # Handler e código gerado pelo protoc
    ├── graphql/              # Handler GraphQL gerado pelo gqlgen
    └── web/                  # Handler HTTP e webserver (chi)

configs/                      # Leitura de variáveis de ambiente via viper
sql/migrations/               # Migrations gerenciadas pelo golang-migrate (embeddadas no binário)
```

### Princípios aplicados

- **Dependency Rule**: `entity` não conhece `usecase`, que não conhece `infra`. A dependência flui sempre de fora para dentro.
- **Interface como contrato**: `OrderRepository` é uma interface definida em `entity`. Os casos de uso dependem dela, não da implementação MySQL.
- **Use Cases independentes de framework**: `CreateOrderUseCase` e `ListOrdersUseCase` não importam nada de HTTP, gRPC ou banco — são testáveis isoladamente.
- **Inversão de dependência no `main`**: o `main.go` é o único ponto que conhece todas as camadas, responsável por montar (wiring) as dependências concretas.

## Pré-requisitos

- Go 1.25+
- Docker e Docker Compose

## Executando a aplicação

### Com Docker

```bash
docker compose -f docker_compose.yml up --build
```

### Com Podman

```bash
podman-compose -f podman_compose.yml up -d
```

### Localmente

1. Suba apenas o banco de dados:

```bash
docker compose -f docker_compose.yml up mysql -d
```

2. Configure o arquivo `.env` na raiz (já existe um exemplo):

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=:8000
GRPC_SERVER_PORT=:50051
GRAPHQL_SERVER_PORT=:8080
```

3. Execute a aplicação:

```bash
go run cmd/ordersystem/main.go
```

As migrations são aplicadas automaticamente na inicialização.

## Portas

| Serviço        | Porta   |
|----------------|---------|
| HTTP REST      | `8000`  |
| gRPC           | `50051` |
| GraphQL        | `8080`  |

## Testando via HTTP (REST)

### Criar pedido

```bash
curl -s -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"ID": "1", "Number": "ORD-001", "Price": 100.0, "Tax": 10.0}' | jq
```

Resposta:

```json
{
  "id": "1",
  "number": "ORD-001",
  "price": 100,
  "tax": 10,
  "final_price": 110
}
```

### Listar pedidos

```bash
curl -s http://localhost:8000/order | jq
```

> Também é possível usar o arquivo `api/orders.http` com a extensão REST Client do VS Code.

## Testando via GraphQL

O playground interativo está disponível em `http://localhost:8080`.

### Criar pedido (Mutation)

```graphql
mutation {
  createOrder(
    id: "1"
    number: "ORD-001"
    price: 100.0
    tax: 10.0
  ) {
    id
    number
    price
    tax
    finalPrice
  }
}
```

### Listar pedidos (Query)

```graphql
query {
  listOrders(page: 1, limit: 10) {
    orders {
      id
      number
      price
      tax
      finalPrice
    }
    total
    page
    limit
  }
}
```

## Testando via gRPC

Instale o [Evans](https://github.com/ktr0731/evans) (cliente gRPC interativo):

```bash
go install github.com/ktr0731/evans@latest
```

Conecte ao servidor usando o arquivo `.proto`:

```bash
evans --proto internal/infra/grpc/proto/order.proto --host localhost --port 50051
```

Dentro do Evans:

```
package pb
service OrderService

call CreateOrder
call ListOrders
```

### Alternativa com grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Criar pedido
grpcurl -plaintext \
  -proto internal/infra/grpc/proto/order.proto \
  -d '{"id": "1", "number": "ORD-001", "price": 100.0, "tax": 10.0}' \
  localhost:50051 pb.OrderService/CreateOrder

# Listar pedidos
grpcurl -plaintext \
  -proto internal/infra/grpc/proto/order.proto \
  -d '{"page": 1, "limit": 10}' \
  localhost:50051 pb.OrderService/ListOrders
```

## Executando os testes

```bash
go test ./...
```
