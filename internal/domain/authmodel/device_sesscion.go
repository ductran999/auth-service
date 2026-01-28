package authmodel

import "time"

type DeviceSession struct {
	JTI       string
	AccountID string

	UserAgent string
	IP        string

	SignAt    time.Time
	ExpiresAt time.Time
}

func NewDeviceSession(jti, accountID string, signAt time.Time) *DeviceSession {
	return &DeviceSession{
		JTI:       jti,
		AccountID: accountID,
		SignAt:    signAt,
	}
}

func (ds *DeviceSession) AttachClientInfo(userAgent, clientIP string) {
	ds.UserAgent = userAgent
	ds.IP = clientIP
}

func (ds *DeviceSession) SetExpiry(exp time.Time) {
	ds.ExpiresAt = exp
}
