package session

type Status int

const (
	Started Status = iota
	PendingWallet
	WalletConnected
	PendingSignature
	Gone
)
