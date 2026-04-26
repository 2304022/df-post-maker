package usecase

import (
	"fmt"
	"net"
	"os"
	"time"
)

type TestUseCase struct{}

func NewTestUseCase() *TestUseCase {
	return &TestUseCase{}
}

func (u *TestUseCase) Ping() error {
	ip := os.Getenv("PING_TARGET_IP")
	if ip == "" {
		ip = "93.90.233.146"
	}

	conn, err := net.DialTimeout("ip4:icmp", ip, 3*time.Second)
	if err != nil {
		// Fallback: try TCP port 80 if ICMP is not allowed
		conn2, err2 := net.DialTimeout("tcp", fmt.Sprintf("%s:80", ip), 3*time.Second)
		if err2 != nil {
			return fmt.Errorf("host unreachable: %w", err2)
		}
		conn2.Close()
		return nil
	}
	conn.Close()
	return nil
}
