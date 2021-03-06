package cloudflare

import (
	"fmt"
	"net"
	"strings"
)

// validateRecordType ensures that the cloudflare record type is valid
func validateRecordType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validTypes := map[string]struct{}{
		"A":     {},
		"AAAA":  {},
		"CNAME": {},
		"TXT":   {},
		"SRV":   {},
		"LOC":   {},
		"MX":    {},
		"NS":    {},
		"SPF":   {},
	}

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			`%q contains an invalid type %q. Valid types are "A", "AAAA", "CNAME", "TXT", "SRV", "LOC", "MX", "NS" or "SPF"`, k, value))
	}
	return
}

// validateRecordName ensures that based on supplied record type, the name content matches
// Currently only validates A and AAAA types
func validateRecordName(t string, value string) error {
	switch t {
	case "A":
		// Must be ipv4 addr
		addr := net.ParseIP(value)
		if addr == nil || !strings.Contains(value, ".") {
			return fmt.Errorf("A record must be a valid IPv4 address, got: %q", value)
		}
	case "AAAA":
		// Must be ipv6 addr
		addr := net.ParseIP(value)
		if addr == nil || !strings.Contains(value, ":") {
			return fmt.Errorf("AAAA record must be a valid IPv6 address, got: %q", value)
		}
	}

	return nil
}

func validatePageRuleStatus(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validStatuses := map[string]struct{}{
		"active": {},
		"paused": {},
	}

	if _, ok := validStatuses[value]; !ok {
		errors = append(errors, fmt.Errorf(
			`%q contains an invalid status %q. Valid statuses are "active" or "paused"`, k, value))
	}
	return
}

func assertIsOneOf(setting string, acceptables []interface{}, value interface{}) error {
	for _, acceptable := range acceptables {
		if value == acceptable {
			return nil
		}
	}
	return fmt.Errorf("%q %q invalid: must be one of %q", setting, value, acceptables)
}

func validateCacheLevel(v interface{}, k string) (ws []string, errors []error) {
	if err := assertIsOneOf("Cache level", []interface{}{"bypass", "basic", "simplified", "aggressive", "cache_everything"}, v.(string)); err != nil {
		errors = append(errors, err)
	}
	return
}

func validateForwardStatusCode(v interface{}, k string) (ws []string, errors []error) {
	if err := assertIsOneOf("Fowarding status code", []interface{}{301, 302}, v.(int)); err != nil {
		errors = append(errors, err)
	}
	return
}

func validateIsTrue(v interface{}, k string) (ws []string, errors []error) {
	if !v.(bool) {
		// We can't just ignore this, since if the action is *not set* by the
		// user it will appear as `false` too.
		errors = append(errors, fmt.Errorf("Action %q has no further setting; `true` is the only valid option.", k))
	}
	return
}

func validateOnOff(v interface{}, k string) (ws []string, errors []error) {
	if err := assertIsOneOf(k, []interface{}{"on", "off"}, v.(string)); err != nil {
		errors = append(errors, err)
	}
	return
}

func validateRocketLoader(v interface{}, k string) (ws []string, errors []error) {
	if err := assertIsOneOf("Rocket loader", []interface{}{"off", "manual", "automatic"}, v.(string)); err != nil {
		errors = append(errors, err)
	}
	return
}

func validateSecurityLevel(v interface{}, k string) (ws []string, errors []error) {
	if err := assertIsOneOf("Security level", []interface{}{"essentially_off", "low", "medium", "high", "under_attack"}, v.(string)); err != nil {
		errors = append(errors, err)
	}
	return
}

func validateSSL(v interface{}, k string) (ws []string, errors []error) {
	if err := assertIsOneOf("SSL mode", []interface{}{"off", "flexible", "full", "strict"}, v.(string)); err != nil {
		errors = append(errors, err)
	}
	return
}

func validateTTL(v interface{}, k string) (ws []string, errors []error) {
	if ttl, maxTTL := v.(int), 31536000; ttl > maxTTL {
		errors = append(errors, fmt.Errorf("Cache TTL of %q too long: max value is %q", ttl, maxTTL))
	}
	return
}
