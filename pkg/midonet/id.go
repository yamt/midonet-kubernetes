package midonet

import (
	"crypto/sha256"
	"net"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	kubernetesSpaceUUIDString = "CAC60164-F74C-404A-AB39-3C1320124A17"
)

func IDForKey(key string) uuid.UUID {
	kubernetes, err := uuid.Parse(kubernetesSpaceUUIDString)
	if err != nil {
		log.WithError(err).Fatal("kubernetesSpaceUUIDString")
	}
	return uuid.NewSHA1(kubernetes, []byte(key))
}

func SubID(id uuid.UUID, s string) uuid.UUID {
	return uuid.NewSHA1(id, []byte(s))
}

func MACForKey(key string) net.HardwareAddr {
	hash := sha256.Sum256([]byte(key))
	// AC-CA-BA  Midokura Co., Ltd.
	addr := [6]byte{0xac, 0xca, 0xba, hash[0], hash[1], hash[2]}
	return net.HardwareAddr(addr[:])
}
