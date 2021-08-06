package service

import (
	"github.com/MaxFuhrich/containerNetworkExample/entities"
	"time"
)

//Creates RequestTime-Object for the database
func GetTime() entities.RequestTime {
	t := time.Now()
	return entities.RequestTime{Time: t.Format(time.RFC1123)}
}
