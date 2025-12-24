// thx for reference: https://altlinux.space/alt-atomic/atomic-installer/src/branch/main/lib/i18n.go
package lib

import (
	"os"
	"strings"

	gcore "github.com/diamondburned/gotk4/pkg/core/glib"
)

func InitLocales() {
	gcore.InitI18n("herbarium", "/usr/share/locale")
}

func GetSystemLocale() string {
	var locale string
	if v := os.Getenv("LC_ALL"); v != "" {
		locale = stripLocaleEncoding(v)
	} else if v := os.Getenv("LC_MESSAGES"); v != "" {
		locale = stripLocaleEncoding(v)
	} else {
		locale = stripLocaleEncoding(os.Getenv("LANG"))
	}

	return locale
}

func stripLocaleEncoding(locale string) string {
	if before, _, ok := strings.Cut(locale, "."); ok {
		return before
	}
	return locale
}

func T_(messageID string) string {
	return gcore.Local(messageID)
}
