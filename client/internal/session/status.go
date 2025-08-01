package session

type Status int

const (
	Started Status = iota
	PendingWallet
	PendingSignature
	DeclinedSignature
	Verified
	Gone
)
