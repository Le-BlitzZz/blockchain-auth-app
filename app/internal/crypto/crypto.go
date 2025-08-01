package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateNonce() (string, error) {
	b := make([]byte, 16) // 128 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// VerifySignature checks that `signatureHex` is a valid EIP‑191
// (“personal_sign”) signature of `message` by `address`.
//
// Arguments
//   - address:       Ethereum address users claim to control (0x‑hex, any case)
//   - message:       Exact message the client signed (must include your nonce)
//   - signatureHex:  65‑byte r||s||v signature as returned by MetaMask et al.
//
// Returns nil on success, or an error describing why verification failed.
func VerifySignature(address, message, signatureHex string) error {
	// Decode hex signature
	sig, err := hexutil.Decode(signatureHex)
	if err != nil {
		return fmt.Errorf("invalid signature hex: %w", err)
	}
	if len(sig) != 65 {
		return fmt.Errorf("invalid signature length: got %d, want 65", len(sig))
	}

	// Normalise the recovery ID (v) to {0,1}
	switch v := sig[64]; v {
	case 27, 28: // MetaMask / parity
		sig[64] = v - 27
	case 0, 1: // geth
		// already fine
	default:
		return fmt.Errorf("invalid recovery ID v=%d", v)
	}

	// Hash the message exactly like personal_sign
	msgHash := accounts.TextHash([]byte(message)) // []byte, 32 bytes long

	// Recover the public key, then the sender address
	pubKey, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		return fmt.Errorf("cannot recover public key: %w", err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Compare with claimed address (case‑insensitive)
	if recoveredAddr != common.HexToAddress(address) {
		return fmt.Errorf("signature valid, but for %s—not %s",
			recoveredAddr.Hex(), common.HexToAddress(address).Hex())
	}

	return nil
}
