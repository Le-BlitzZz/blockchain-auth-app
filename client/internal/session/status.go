package session

type Status int

const (
	Started Status = iota
	PendingWallet
	DeclinedWallet
	PendingSignature
	DeclinedSignature
	Verified
	Gone
)
