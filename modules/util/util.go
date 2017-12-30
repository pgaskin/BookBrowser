package util

import (
	"net"
	"os"
	"strings"
)

// StringBetween gets the string in between two other strings, and returns an empty string if not found. It returns the first match.
func StringBetween(str, start, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}

// StringAfter gets the string after another.
func StringAfter(str, start string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	return str[s:]
}

// FixString fixes some issues with strings in metadata.
func FixString(s string) string {
	return strings.Map(func(in rune) rune {
		switch in {
		case '“', '‹', '”', '›':
			return '"'
		case '‘', '’':
			return '\''
		}
		return in
	}, s)
}

// GetIP gets the preferred outbound ip of this machine.
func GetIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Exists checks whether a path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
