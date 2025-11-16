package order

import (
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"jihulab.com/yanlin/food-api/pkg/test"
)

type BaseTestSuite struct {
	suite.Suite
	dbName      string
	pool        *pgxpool.Pool
	testPool    *test.Pool
	redisClient *redis.Client
}

// this function executes before the test suite begins execution
func (s *BaseTestSuite) SetupSuite() {
	// set up a new test db pool
	s.testPool = test.NewPool()
	s.dbName = s.testPool.Name
	s.pool = s.testPool.Pool
	s.redisClient = test.NewRedisClient()
}

// this function executes after all tests executed
func (s *BaseTestSuite) TearDownSuite() {
	// drop the test db pool
	s.testPool.Drop()
}

// this function executes before each test case
func (s *BaseTestSuite) SetupTest() {
	// depending on use cases, we may want to set up a brand new test db for each test case, instead of
	// the whole test suite
}

// this function executes after each test case
func (s *BaseTestSuite) TearDownTest() {
}

func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}
