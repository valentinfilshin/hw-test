package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		var user User
		if err = user.UnmarshalJSON(line); err != nil {
			return
		}

		result[i] = user
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = "." + strings.ToLower(domain)

	for _, user := range u {
		dogIndex := strings.LastIndex(user.Email, "@")
		if dogIndex < 0 {
			continue
		}
		email := strings.ToLower(user.Email[dogIndex+1:])
		if strings.HasSuffix(email, domain) {
			result[email]++
		}
	}
	return result, nil
}
