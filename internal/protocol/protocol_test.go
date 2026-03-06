package protocol

import (
	"bytes"
	"testing"

	"paqet/internal/conf"
	"paqet/internal/tnet"
)

func TestProtoRoundTrip(t *testing.T) {
	addr, err := tnet.NewAddr("127.0.0.1:443")
	if err != nil {
		t.Fatalf("new addr: %v", err)
	}

	in := Proto{
		Type: PTCPF,
		Addr: addr,
		TCPF: []conf.TCPF{{SYN: true, ACK: true}, {FIN: true, NS: true}},
	}

	var buf bytes.Buffer
	if err := in.Write(&buf); err != nil {
		t.Fatalf("write: %v", err)
	}

	var out Proto
	if err := out.Read(&buf); err != nil {
		t.Fatalf("read: %v", err)
	}

	if out.Type != in.Type {
		t.Fatalf("type mismatch: got %d want %d", out.Type, in.Type)
	}
	if out.Addr == nil || out.Addr.String() != in.Addr.String() {
		t.Fatalf("addr mismatch: got %v want %v", out.Addr, in.Addr)
	}
	if len(out.TCPF) != len(in.TCPF) {
		t.Fatalf("tcpf len mismatch: got %d want %d", len(out.TCPF), len(in.TCPF))
	}
	for i := range in.TCPF {
		if out.TCPF[i] != in.TCPF[i] {
			t.Fatalf("tcpf[%d] mismatch: got %+v want %+v", i, out.TCPF[i], in.TCPF[i])
		}
	}
}

func TestProtoRejectsBadMagic(t *testing.T) {
	bad := []byte{0, 0, 1, byte(PPING), 0, 0, 0}
	var p Proto
	if err := p.Read(bytes.NewReader(bad)); err == nil {
		t.Fatal("expected error for invalid magic")
	}
}

func TestPacketTypeIDsV2(t *testing.T) {
	tests := []struct {
		name string
		got  PType
		want PType
	}{
		{name: "PPING", got: PPING, want: 0x11},
		{name: "PPONG", got: PPONG, want: 0x12},
		{name: "PTCPF", got: PTCPF, want: 0x21},
		{name: "PTCP", got: PTCP, want: 0x22},
		{name: "PUDP", got: PUDP, want: 0x23},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Fatalf("%s mismatch: got 0x%02x want 0x%02x", tt.name, tt.got, tt.want)
		}
	}
}
