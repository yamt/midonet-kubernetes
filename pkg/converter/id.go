package converter

import (
	"crypto/sha256"
	"net"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	kubernetesSpaceUUIDString    = "CAC60164-F74C-404A-AB39-3C1320124A17"
	midonetTenantSpaceUUIDString = "3978567E-91C4-465C-A0D1-67575F6B4C7F"
)

func idForString(spaceStr string, key string) uuid.UUID {
	space, err := uuid.Parse(spaceStr)
	if err != nil {
		log.WithError(err).Fatal("space")
	}
	return uuid.NewSHA1(space, []byte(key))
}

func IDForTenant(tenant string) uuid.UUID {
	return idForString(midonetTenantSpaceUUIDString, tenant)
}

func IDForKey(key string) uuid.UUID {
	return idForString(kubernetesSpaceUUIDString, key)
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
