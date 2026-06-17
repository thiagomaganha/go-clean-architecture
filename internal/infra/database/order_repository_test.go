package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/thiagomaganha/go-clean-architecture/internal/entity"
	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
	_ "modernc.org/sqlite"
)

type OrderRepositorySuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositorySuite) SetupSuite() {
	db, err := sql.Open("sqlite", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, number varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositorySuite) TearDownSuite() {
	suite.Db.Close()
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositorySuite))
}

func (suite *OrderRepositorySuite) TestList() {
	ctx := context.Background()
	repo := NewOrderRepository(suite.Db)
	input := usecase.ListOrdersInput{Query: "%", Page: 1, Limit: 10}
	output, err := repo.List(ctx, input)
	suite.NoError(err)
	suite.Equal(0, len(output.Orders))
	suite.Equal(0, output.Total)

	order := entity.NewOrder("123", "ORD1", 10.0, 2.0)
	suite.NoError(err)
	err = repo.Save(ctx, order)
	suite.NoError(err)

	output, err = repo.List(ctx, input)
	suite.NoError(err)
	suite.Equal(1, len(output.Orders))
	suite.Equal(1, output.Total)
}
