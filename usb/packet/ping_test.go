package packet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const pingPacketHexDump = "10000000676e69700000000001000000"

func TestPingSerialization(t *testing.T) {
	assert.Equal(t, pingPacketHexDump, fmt.Sprintf("%x", NewPingPacketAsBytes()))
}
