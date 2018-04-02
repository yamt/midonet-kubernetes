package midonet

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	kubernetesSpaceUUIDString = "CAC60164-F74C-404A-AB39-3C1320124A17"
)

func idForKey(key string) uuid.UUID {
	kubernetes, err := uuid.Parse(kubernetesSpaceUUIDString)
	if err != nil {
		log.WithError(err).Fatal("kubernetesSpaceUUIDString")
	}
	return uuid.NewSHA1(kubernetes, []byte(key))
}

func subID(id uuid.UUID, s string) uuid.UUID {
	return uuid.NewSHA1(id, []byte(s))
}
