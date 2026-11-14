// Package config reads the app's settings from the environment once, at startup.
package config

import "errors"

// Config is everything main needs from the environment.
type Config struct {
	Addr       string
	ContentDir string // empty = embedded content
	DevMode    bool
	// OriginSecrets are the accepted X-Origin-Auth values, current first, then
	// previous, so the secret can be rotated without downtime (DESIGN.md §8.1).
	OriginSecrets []string
}

// FromEnv builds a Config. getenv is os.Getenv in main and a map lookup in tests.
// It fails closed: outside DEV_MODE a missing origin secret is an error, never "auth off".
func FromEnv(getenv func(string) string) (Config, error) {
	c := Config{
		Addr:       ":8080",
		ContentDir: getenv("CONTENT_DIR"),
		DevMode:    getenv("DEV_MODE") == "true", // exactly "true"; "1", "yes", "TRUE" all leave it off
	}
	for _, k := range []string{"ORIGIN_AUTH", "ORIGIN_AUTH_PREVIOUS"} {
		if v := getenv(k); v != "" {
			c.OriginSecrets = append(c.OriginSecrets, v)
		}
	}
	if !c.DevMode && getenv("ORIGIN_AUTH") == "" {
		return Config{}, errors.New("ORIGIN_AUTH is required unless DEV_MODE=true")
	}
	return c, nil
}
