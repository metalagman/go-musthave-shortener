package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"shortener/internal/app/logger"
)

var (
	ErrEmptyTrustedNetwork = errors.New("empty trusted network")
)

func TrustedNetwork(cidr string) func(next http.Handler) http.Handler {
	const component = "TrustedNetwork.Middleware"

	l := logger.Global().Component(component)

	enabled := true
	trusted, err := parseTrustedNetwork(cidr)
	if err != nil {
		l.Error().Err(err).Msg("Trusted network config error")
		enabled = false
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			l := logger.Ctx(ctx).Component(component)

			if !enabled {
				l.Error().Msg("Trusted network is missing")
				http.Error(w, "", http.StatusForbidden)
				return
			}

			ip := r.Header.Get("X-Real-IP")
			addr := net.ParseIP(ip)
			if addr == nil {
				l.Error().Str("x-real-ip", ip).Msg("X-Real-IP header parse error")
				http.Error(w, "", http.StatusForbidden)
				return
			}

			if !trusted.Contains(addr) {
				l.Error().Str("x-real-ip", addr.String()).Msg("Access denied")
				http.Error(w, "", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseTrustedNetwork(s string) (*net.IPNet, error) {
	if s == "" {
		return nil, ErrEmptyTrustedNetwork
	}

	_, n, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("trusted network parse: %w", err)
	}

	return n, nil
}
