package database

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/storage/redis"
)

type Database struct {
	redis redis.Redis
}

func Init(host string) (*Database, error) {
	s := new(Database)
	err := s.redis.Init(host)
	if err != nil {
		return nil, errors.E("failed to init database", errors.Params{"host": host})
	}
	return s, nil
}

func (s *Database) SaveData(entity, key string, value interface{}) error {
	return s.redis.AddHM(entity, key, value)
}

func (s *Database) GetData(entity, key string, r interface{}) error {
	return s.redis.GetHMValue(entity, key, &r)
}
