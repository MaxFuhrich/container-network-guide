package service

import (
	"github.com/MaxFuhrich/containerNetworkExample/entities"
	"time"
)

func GetTime() entities.RequestTime {
	t := time.Now()
	return entities.RequestTime{Time: t.Format(time.RFC1123)}
}
