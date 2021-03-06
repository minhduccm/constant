package wire

import (
	"encoding/hex"
	"encoding/json"

	"github.com/libp2p/go-libp2p-peer"
)

const (
	MaxPingPayload = 1000 // 1 1Kb
)

type MessagePing struct {
}

func (self MessagePing) MessageType() string {
	return CmdPing
}

func (self MessagePing) MaxPayloadLength(pver int) int {
	return MaxPingPayload
}

func (self MessagePing) JsonSerialize() ([]byte, error) {
	jsonBytes, err := json.Marshal(self)
	return jsonBytes, err
}

func (self MessagePing) JsonDeserialize(jsonStr string) error {
	jsonDecodeString, _ := hex.DecodeString(jsonStr)
	err := json.Unmarshal([]byte(jsonDecodeString), self)
	return err
}
func (self MessagePing) SetSenderID(senderID peer.ID) error {
	return nil
}
