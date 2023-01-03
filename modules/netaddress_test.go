package modules

import (
	"net"
	"strings"
	"testing"
)

var (
	// Networks such as 10.0.0.x have been omitted from testing - behavior
	// for these networks is currently undefined.

	invalidAddrs = []string{
		// Garbage addresses
		"",
		"bar:bar:baz",
		"garbage:6146:616",
		// Missing host / port
		":",
		"111.111.111.111",
		"12.34.45.64",
		"[::2]",
		"::2",
		"bar",
		"hn.com",
		"世界",
		"bar:",
		"世界:",
		":bar",
		":世界",
		"localhost:",
		"[::1]:",
		// Invalid host / port chars
		"localhost:-",
		"[::1]:-",
		"bar:{}",
		"{}:123",
		" bar:123",
		"bar :123",
		"f oo:123",
		"bar: 123",
		"bar:123 ",
		"bar:1 23",
		"\x00:123",
		"bar:\x00",
		"世界:123",
		"bar:世界",
		"世:界",
		`":"`,
		// Unspecified address
		"[::]:bar",
		"0.0.0.0:bar",
		// Invalid hostnames
		"unqualifiedhost:123",
		"Yo-Amazon.we-are-really-happy-for-you.and-we-will-let-you-finish.but-ScPrime-is-the-best-cloud-storage-of-all-time.of-all-time-of-all-time-of-all-time-of-all-time-of-all-time.of-all-time-of-all-time-of-all-time-of-all-time-of-all-time.of-all-time-of-all-time:123",
		strings.Repeat("a", 64) + ".com:123",                       // 64 char long label too long.
		strings.Repeat(strings.Repeat("a", 62)+".", 4) + "co:123",  // 254 char long hostname too long.
		strings.Repeat(strings.Repeat("a", 62)+".", 4) + "co.:123", // 254 char long hostname with trailing dot too long.
		"-bar.bar:123",
		"bar-.bar:123",
		"bar.-bar:123",
		"bar.bar-:123",
		"bar-bar.-baz:123",
		"bar-bar.baz-:123",
		"bar.-bar.baz:123",
		"bar.bar-.baz:123",
		".:123",
		".bar.com:123",
		"bar.com..:123",
		// invalid port numbers
		"bar:0",
		"bar:65536",
		"bar:-100",
		"bar:1000000",
		"localhost:0",
		"[::1]:0",
	}
	validAddrs = []string{
		// Loopback address (valid in testing only, can't really test this well)
		"localhost:123",
		"127.0.0.1:123",
		"[::1]:123",
		// Valid addresses.
		"bar.com:1",
		"bar.com.:1",
		"a.b.c:1",
		"a.b.c.:1",
		"bar-bar.com:123",
		"FOO.com:1",
		"1bar.com:1",
		"tld.bar.com:1",
		"hn.com:8811",

		"[::2]:65535",
		"111.111.111.111:111",
		"12.34.45.64:7777",

		strings.Repeat("bar.", 63) + "f:123",  // 253 chars long
		strings.Repeat("bar.", 63) + "f.:123", // 254 chars long, 253 chars long without trailing dot

		strings.Repeat(strings.Repeat("a", 63)+".", 3) + "a:123", // 3x63 char length labels + 1x1 char length label without trailing dot
		strings.Repeat(strings.Repeat("a", 63)+".", 3) + ":123",  // 3x63 char length labels with trailing dot
	}
)

// TestHostPort tests the Host and Port methods of the NetAddress type.
func TestHostPort(t *testing.T) {
	t.Parallel()

	// Test valid addrs.
	for _, addr := range validAddrs {
		na := NetAddress(addr)
		host := na.Host()
		port := na.Port()
		expectedHost, expectedPort, err := net.SplitHostPort(addr)
		if err != nil {
			t.Fatal(err)
		}
		if host != expectedHost {
			t.Errorf("Host() returned unexpected host for NetAddress '%v': expected '%v', got '%v'", na, expectedHost, host)
		}
		if port != expectedPort {
			t.Errorf("Port() returned unexpected port for NetAddress '%v': expected '%v', got '%v'", na, expectedPort, port)
		}
	}

	// Test that Host / Port return "" when net.SplitHostPort errors
	na := NetAddress("::")
	host := na.Host()
	port := na.Port()
	if host != "" {
		t.Error("expected Host() to return blank for an un-splittable NetAddress, but it returned:", host)
	}
	if port != "" {
		t.Error("expected Port() to return blank for an un-splittable NetAddress, but it returned:", port)
	}
}

// TestIsLoopback tests the IsLoopback method of the NetAddress type.
func TestIsLoopback(t *testing.T) {
	t.Parallel()

	testSet := []struct {
		query           NetAddress
		desiredResponse bool
	}{
		// Localhost tests.
		{"localhost", false},
		{"localhost:1234", true},
		{"127.0.0.1", false},
		{"127.0.0.1:6723", true},
		{"::1", false},
		{"[::1]:7124", true},

		// Local network tests.
		{"10.0.0.0", false},
		{"10.0.0.0:1234", false},
		{"10.2.2.5", false},
		{"10.2.2.5:16432", false},
		{"10.255.255.255", false},
		{"10.255.255.255:16432", false},
		{"172.16.0.0", false},
		{"172.16.0.0:1234", false},
		{"172.26.2.5", false},
		{"172.26.2.5:16432", false},
		{"172.31.255.255", false},
		{"172.31.255.255:16432", false},
		{"192.168.0.0", false},
		{"192.168.0.0:1234", false},
		{"192.168.2.5", false},
		{"192.168.2.5:16432", false},
		{"192.168.255.255", false},
		{"192.168.255.255:16432", false},
		{"1234:0000:0000:0000:0000:0000:0000:0000", false},
		{"[1234:0000:0000:0000:0000:0000:0000:0000]:1234", false},
		{"fc00:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fc00:0000:0000:0000:0000:0000:0000:0000]:1234", false},
		{"fd00:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fd00:0000:0000:0000:0000:0000:0000:0000]:1234", false},
		{"fd30:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fd30:0000:0000:0000:0000:0000:0000:0000]:1234", false},
		{"fd00:0000:0030:0000:0000:0000:0000:0000", false},
		{"[fd00:0000:0030:0000:0000:0000:0000:0000]:1234", false},
		{"fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", false},
		{"[fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff]:1234", false},
		{"fe00:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fe00:0000:0000:0000:0000:0000:0000:0000]:1234", false},

		// Unspecified address tests.
		{"0.0.0.0:1234", false},
		{"[::]:1234", false},

		// Public name tests.
		{"hn.com", false},
		{"hn.com:8811", false},
		{"2.34.45.64", false},
		{"2.34.45.64:7777", false},
		{"12.34.45.64", false},
		{"12.34.45.64:7777", false},
		{"122.34.45.64", false},
		{"122.34.45.64:7777", false},
		{"197.34.45.64", false},
		{"197.34.45.64:7777", false},
		{"222.34.45.64", false},
		{"222.34.45.64:7777", false},

		// Garbage name tests.
		{"", false},
		{"garbage", false},
		{"garbage:6432", false},
		{"garbage:6146:616", false},
		{"::1:4646", false},
		{"[::1]", false},
	}
	for _, test := range testSet {
		if test.query.IsLoopback() != test.desiredResponse {
			t.Error("test failed:", test, test.query.IsLoopback())
		}
	}
}

// TestIsValid tests that IsValid only returns nil for valid addresses.
func TestIsValid(t *testing.T) {
	t.Parallel()

	for _, addr := range validAddrs {
		na := NetAddress(addr)
		if err := na.IsValid(); err != nil {
			t.Errorf("IsValid returned non-nil for valid NetAddress %q: %v", addr, err)
		}
	}
	for _, addr := range invalidAddrs {
		na := NetAddress(addr)
		if err := na.IsValid(); err == nil {
			t.Errorf("IsValid returned nil for an invalid NetAddress %q: %v", addr, err)
		}
	}
}

// TestIsLocal checks that the correct values are returned for all local IP
// addresses.
func TestIsLocal(t *testing.T) {
	t.Parallel()

	testSet := []struct {
		query           NetAddress
		desiredResponse bool
	}{
		// Localhost tests.
		{"localhost", false},
		{"localhost:1234", true},
		{"127.0.0.1", false},
		{"127.0.0.1:6723", true},
		{"::1", false},
		{"[::1]:7124", true},

		// Local network tests.
		{"10.0.0.0", false},
		{"10.0.0.0:1234", true},
		{"10.2.2.5", false},
		{"10.2.2.5:16432", true},
		{"10.255.255.255", false},
		{"10.255.255.255:16432", true},
		{"172.16.0.0", false},
		{"172.16.0.0:1234", true},
		{"172.26.2.5", false},
		{"172.26.2.5:16432", true},
		{"172.31.255.255", false},
		{"172.31.255.255:16432", true},
		{"192.168.0.0", false},
		{"192.168.0.0:1234", true},
		{"192.168.2.5", false},
		{"192.168.2.5:16432", true},
		{"192.168.255.255", false},
		{"192.168.255.255:16432", true},
		{"169.254.255.255", false},
		{"169.254.255.255:12345", true},
		{"169.254.25.55", false},
		{"169.254.25.55:12345", true},
		{"169.254.44.207", false},
		{"169.254.44.207:9982", true},
		{"169.254.3.139", false},
		{"169.254.3.139:9982", true},
		{"100.65.70.95", false},
		{"100.65.70.95:9982", true},
		{"1234:0000:0000:0000:0000:0000:0000:0000", false},
		{"[1234:0000:0000:0000:0000:0000:0000:0000]:1234", false},
		{"fc00:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fc00:0000:0000:0000:0000:0000:0000:0000]:1234", false},
		{"fd00:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fd00:0000:0000:0000:0000:0000:0000:0000]:1234", true},
		{"fd30:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fd30:0000:0000:0000:0000:0000:0000:0000]:1234", true},
		{"fd00:0000:0030:0000:0000:0000:0000:0000", false},
		{"[fd00:0000:0030:0000:0000:0000:0000:0000]:1234", true},
		{"fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", false},
		{"[fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff]:1234", true},
		{"fe00:0000:0000:0000:0000:0000:0000:0000", false},
		{"[fe00:0000:0000:0000:0000:0000:0000:0000]:1234", false},

		// Unspecified address tests.
		{"0.0.0.0:1234", false},
		{"[::]:1234", false},

		// Public name tests.
		{"hn.com", false},
		{"hn.com:8811", false},
		{"2.34.45.64", false},
		{"2.34.45.64:7777", false},
		{"12.34.45.64", false},
		{"12.34.45.64:7777", false},
		{"122.34.45.64", false},
		{"122.34.45.64:7777", false},
		{"197.34.45.64", false},
		{"197.34.45.64:7777", false},
		{"222.34.45.64", false},
		{"222.34.45.64:7777", false},

		// Garbage name tests.
		{"", false},
		{"garbage", false},
		{"garbage:6432", false},
		{"garbage:6146:616", false},
		{"::1:4646", false},
		{"[::1]", false},
	}
	for _, test := range testSet {
		if test.query.IsLocal() != test.desiredResponse {
			t.Error("test failed:", test, test.query.IsLocal())
		}
	}
}
