package models

import "github.com/google/uuid"

type SmartToy struct {
	Uuid        uuid.UUID `db:"smart_toy_uuid"`
	FingerPrint string    `db:"fingerprint"`
}
