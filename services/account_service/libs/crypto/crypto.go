package crypto

import (
	"strings"
	"crypto/md5"
)

const (
	CRYPT_SETTING = "$P$BwQZDcQaNU9zAOF.6MOUdEhz9X68fL1"
)

func CryptPrivate(pw, setting string) string {
	const itoa64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Setting = "$P$BwQZDcQaNU9zAOF.6MOUdEhz9X68fL1"
	var outp = "*0"
	var count_log2 uint
	count_log2 = uint(strings.Index(itoa64, string(setting[3])))
	if count_log2 < 7 || count_log2 > 30 {
		return outp
	}
	count := 1 << count_log2
	salt := setting[4:12]
	if len(salt) != 8 {
		return outp
	}
	hasher := md5.New()
	hasher.Write([]byte(salt + pw))
	hx := hasher.Sum(nil)
	for count != 0 {
		hasher := md5.New()
		hasher.Write([]byte(string(hx) + pw))
		hx = hasher.Sum(nil)
		count -= 1
	}
	return setting[:12] + encode64(hx, 16)
}

func encode64(inp []byte, count int) string {
	const itoa64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var outp string
	cur := 0
	for cur < count {
		value := uint(inp[cur])
		cur += 1
		outp += string(itoa64[value&0x3f])
		if cur < count {
			value |= (uint(inp[cur]) << 8)
		}
		outp += string(itoa64[(value>>6)&0x3f])

		if cur >= count {
			break
		}
		cur += 1
		if cur < count {
			value |= (uint(inp[cur]) << 16)
		}
		outp += string(itoa64[(value>>12)&0x3f])
		if cur >= count {
			break
		}
		cur += 1
		outp += string(itoa64[(value>>18)&0x3f])
	}
	return outp
}

func PortableHashCheck(pw, storedHash string) bool {
	// PW: password to check (non-hashed), storedHash: Really? do you need an explanation?
	hx := CryptPrivate(pw, storedHash)
	return hx == storedHash
}