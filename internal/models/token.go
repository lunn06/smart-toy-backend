package models

import (
	"encoding/binary"
	"time"
)

type JwtToken struct {
	UserId       int       `db:"user_id"`
	CreationTime time.Time `db:"creation_time"`
	FingerPrint  string    `db:"fingerprint"`
}

func (jwtToken JwtToken) MarshalBinary() ([]byte, error) {
	bytes := make([]byte, 0, 12)

	bytes = binary.BigEndian.AppendUint32(bytes, uint32(jwtToken.UserId))
	bytes = binary.BigEndian.AppendUint64(bytes, uint64(jwtToken.CreationTime.Unix()))
	
	bytes = append(bytes, []byte(jwtToken.FingerPrint)...)

	return bytes, nil
}

func (jwtToken *JwtToken) UnmarshalBinary(data []byte) error {
	jwtToken.UserId = int(binary.BigEndian.Uint32(data[:4]))
	jwtToken.CreationTime = time.Unix(int64(binary.BigEndian.Uint64(data[4:12])), 0)
	jwtToken.FingerPrint = string(data[12:])

	return nil
}