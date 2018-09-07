package client

import (
	"crypto/rand"

	"golang.org/x/crypto/curve25519"
)

const (
	SpendingKeyLength     = 32 // bytes
	ReceivingKeyLength    = 32
	SpendingAddressLength = 32
)

type SpendingKey [SpendingKeyLength]byte

// RandBits generates random bits and return as bytes; zero out redundant bits
func RandBits(n int) []byte {
	m := 1 + (n-1)/8
	b := make([]byte, m)
	rand.Read(b)

	if n%8 > 0 {
		b[m-1] &= ((1 << uint(n%8)) - 1)
	}
	return b
}

// RandSpendingKey generates a random SpendingKey
func RandSpendingKey() SpendingKey {
	b := RandBits(SpendingKeyLength*8 - 4)
	// b := make([]byte, SpendingKeyLength)
	// rand.Read(b)
	// b[SpendingKeyLength-1] &= 0x0F // First 4 bits are 0

	ask := *new(SpendingKey)
	copy(ask[:], b)
	return ask
}

type ReceivingKey [ReceivingKeyLength]byte

func GenReceivingKey(ask SpendingKey) ReceivingKey {
	data := PRF_addr_x(ask[:], 1)
	clamped := clampCurve25519(data)
	var skenc ReceivingKey
	copy(skenc[:], clamped)
	return skenc
}

func clampCurve25519(x []byte) []byte {
	x[0] &= 0xF8  // Clear bit 0, 1, 2 of first byte
	x[31] &= 0x7F // Clear bit 7 of last byte
	x[31] |= 0x40 // Set bit 6 of last byte
	return x
}

type SpendingAddress [SpendingAddressLength]byte

func GenSpendingAddress(ask SpendingKey) SpendingAddress {
	data := PRF_addr_x(ask[:], 0)
	var apk SpendingAddress
	copy(apk[:], data)
	return apk
}

type ViewingKey struct {
	Apk   SpendingAddress
	Skenc ReceivingKey
}

func GenViewingKey(ask SpendingKey) ViewingKey {
	var ivk ViewingKey
	ivk.Apk = GenSpendingAddress(ask)
	ivk.Skenc = GenReceivingKey(ask)
	return ivk
}

type TransmissionKey [32]byte

type PaymentAddress struct {
	Apk   SpendingAddress
	Pkenc TransmissionKey
}

func GenTransmissionKey(skenc ReceivingKey) TransmissionKey {
	// TODO: reduce copy
	var x, y [32]byte
	copy(y[:], skenc[:])
	curve25519.ScalarBaseMult(&x, &y)

	var pkenc TransmissionKey
	copy(pkenc[:], x[:])
	return pkenc
}

func GenPaymentAddress(ask SpendingKey) PaymentAddress {
	var addr PaymentAddress
	addr.Apk = GenSpendingAddress(ask)
	addr.Pkenc = GenTransmissionKey(GenReceivingKey(ask))
	return addr
}

// FullKey convenient struct storing all keys and addresses
type FullKey struct {
	Ask  SpendingKey
	Ivk  ViewingKey
	Addr PaymentAddress
}

// GenFullKey generates all needed keys from a single SpendingKey
func (ask SpendingKey) GenFullKey() FullKey {
	return FullKey{Ask: ask, Ivk: GenViewingKey(ask), Addr: GenPaymentAddress(ask)}
}