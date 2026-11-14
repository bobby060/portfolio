package config

import (
	"slices"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestConfig(t *testing.T) {
	tests := []struct {
		desc            string
		env             map[string]string
		valid           bool
		expectedSecrets []string
	}{
		{"Nothing set should fail", map[string]string{}, false, []string{}},

		{"DEV_MODE=false and ORIGIN_AUTH should pass", map[string]string{"DEV_MODE": "false", "ORIGIN_AUTH": "dummy"}, true, []string{"dummy"}},
		{"DEV_MODE=true and no ORIGIN_AUTH should pass", map[string]string{"DEV_MODE": "true", "ORIGIN_AUTH": ""}, true, []string{}},
		{"DEV_MODE=TRUE and no ORIGIN_AUTH should fail", map[string]string{"DEV_MODE": "TRUE", "ORIGIN_AUTH": "", "ORIGIN_AUTH_PREVIOUS": ""}, false, []string{}},
		{"ORIGIN_AUTH=a, PREV=B => OriginSecrets a, b", map[string]string{"DEV_MODE": "true", "ORIGIN_AUTH": "a", "ORIGIN_AUTH_PREVIOUS": "b"}, true, []string{"a", "b"}},
		{"Only ORIGIN_AUTH_PREVIOUS set fails", map[string]string{"DEV_MODE": "false", "ORIGIN_AUTH_PREVIOUS": "b"}, false, []string{"b"}},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {

			c, e := FromEnv(env(tc.env))
			if tc.valid && e != nil {
				t.Errorf("Test %s should not have failed! Env: \n%+v", tc.desc, tc.env)
			} else if !tc.valid {
				if e == nil {
					t.Errorf("Test %s should have failed! Env: \n%+v", tc.desc, tc.env)
				}
				return

			} else if !slices.Equal(c.OriginSecrets, tc.expectedSecrets) {
				t.Errorf("Expected OriginSecrests %+v are not equal to actual %+v for %s", tc.expectedSecrets, c.OriginSecrets, tc.desc)
			}
		})
	}
}
