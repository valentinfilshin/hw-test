package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = "." + strings.ToLower(domain)
	scanner := bufio.NewScanner(r)

	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		var user User
		if err := user.UnmarshalJSON(line); err != nil {
			return make(DomainStat), err
		}

		dogIndex := strings.LastIndex(user.Email, "@")
		if dogIndex < 0 {
			continue
		}
		email := strings.ToLower(user.Email[dogIndex+1:])
		if strings.HasSuffix(email, domain) {
			result[email]++
		}
		i++
	}

	return result, nil
}
