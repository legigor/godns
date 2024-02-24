package dns

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type DNSFlags struct {
	QR     uint16
	Opcode uint16
	AA     uint16
	TC     uint16
	RD     uint16
	RA     uint16
	Z      uint16
	Rcode  uint16
}

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// Define DNSType as an enumeration of DNS record types.
type DNSType uint16

const (
	TypeA     DNSType = 1
	TypeNS    DNSType = 2
	TypeCNAME DNSType = 5
	TypeMX    DNSType = 15
	TypeTXT   DNSType = 16
	TypeAAAA  DNSType = 28
	// Add other types as needed.
)

// Define DNSClass as an enumeration of DNS classes.
type DNSClass uint16

const (
	ClassIN DNSClass = 1 // Internet
	// Add other classes as needed, but most uses will be ClassIN.
)

type DNSQuestion struct {
	Name  string
	Type  DNSType
	Class DNSClass
}

func Test_parsing(t *testing.T) {
	data := []byte("V\r\x01\x00\x00\x01\x00\x00\x00\x00\x00\x01\x03www\x06google\x03com\x00\x00\x01\x00\x01\x00\x00)\x04\xd0\x00\x00\x00\x00\x00\x00")
	buf := bytes.NewReader(data)

	var header DNSHeader
	err := binary.Read(buf, binary.BigEndian, &header)
	require.NoError(t, err)

	t.Logf("HEADER:\t%+v\n", header)

	flags := DNSFlags{
		QR:     (header.Flags >> 15) & 0x1,
		Opcode: (header.Flags >> 11) & 0xF,
		AA:     (header.Flags >> 10) & 0x1,
		TC:     (header.Flags >> 9) & 0x1,
		RD:     (header.Flags >> 8) & 0x1,
		RA:     (header.Flags >> 7) & 0x1,
		Z:      (header.Flags >> 4) & 0x7,
		Rcode:  header.Flags & 0xF,
	}

	t.Logf("FLAGS:\t%+v\n", flags)

	questions, err := parseDNSQuestions(data, header.QDCount)
	require.NoError(t, err)

	t.Logf("Question:\t%+v\n", questions)
}

func parseDNSQuestions(data []byte, qdCount uint16) ([]DNSQuestion, error) {
	var questions []DNSQuestion
	offset := 12 // Header length

	for i := 0; i < int(qdCount); i++ {
		var question DNSQuestion
		end := 0

		// Extract the name
		for {
			length := int(data[offset])
			if length == 0 {
				offset++ // Move past the null byte
				break
			}
			if end != 0 { // Add dot if not the first part of the name
				question.Name += "."
			}
			end = offset + length + 1
			question.Name += string(data[offset+1 : end])
			offset = end
		}

		// Check for sufficient data for Type and Class
		if len(data[offset:]) < 4 {
			return nil, fmt.Errorf("incomplete data for question at index %d", i)
		}

		// Extract Type and Class
		question.Type = DNSType(binary.BigEndian.Uint16(data[offset : offset+2]))
		question.Class = DNSClass(binary.BigEndian.Uint16(data[offset+2 : offset+4]))
		offset += 4 // Move past Type and Class

		questions = append(questions, question)
	}

	return questions, nil
}

func Ignored_test_server(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()

	srv := NewServer(logger, ctx)
	addr, err := srv.Start()
	require.NoError(t, err)

	logger.Info("server: " + addr.String())

	client, err := net.DialUDP("udp", nil, addr)
	require.NoError(t, err)

	err = client.SetDeadline(time.Now().Add(5 * time.Second))
	require.NoError(t, err)

	message := "Hello, UDP server!"

	_, err = client.Write([]byte(message))
	require.NoError(t, err)

	buffer := make([]byte, 1024)
	n, _, err := client.ReadFromUDP(buffer)
	require.NoError(t, err)

	response := string(buffer[:n])
	logger.Info("response: " + response)

	// DNS
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return client, nil
		},
	}

	ips, err := r.LookupIP(context.Background(), "ip", "www.google.com")
	require.NoError(t, err)

	for _, ip := range ips {
		logger.Info("Resolved", "ip", ip)
	}
}
