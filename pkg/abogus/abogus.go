package abogus

import (
	"crypto/rc4"
	"encoding/base64"
	"math/rand"
	"net/url"
	"time"

	"github.com/tjfoc/gmsm/sm3"
)

var uaCode = []byte{
	76, 98, 15, 131, 97, 245, 224, 133, 122, 199, 241, 166, 79, 32, 90, 191,
	128, 126, 122, 98, 66, 11, 14, 40, 49, 110, 110, 173, 67, 96, 138, 252,
}

const browserStr = "1536|742|1536|864|0|0|0|0|1536|864|1536|864|1536|742|24|24|MacIntel"

var customBase64Encoder = base64.NewEncoding("Dkdpgh2ZmsQB80/MfvV36XI1R45-WUAlEixNLwoqYTOPuzKFjJnry79HbGcaStCe").WithPadding(base64.StdPadding)

func sm3Hash(data []byte) []byte {
	h := sm3.Sm3Sum(data)
	return h
}

func randomList(r float64, b, c, d, e, f, g int) []byte {
	v1 := int(r) & 255
	v2 := int(r) >> 8
	return []byte{
		byte(v1&b | d),
		byte(v1&c | e),
		byte(v2&b | f),
		byte(v2&c | g),
	}
}

func generateString1(rn1, rn2, rn3 float64) string {
	l1 := randomList(rn1, 170, 85, 1, 2, 5, 45&170)
	l2 := randomList(rn2, 170, 85, 1, 0, 0, 0)
	l3 := randomList(rn3, 170, 85, 1, 0, 5, 0)
	return string(l1) + string(l2) + string(l3)
}

func list4(a, b, c, d, e, f, g, h, i, j, k, m, n, o, p, q, r int) []byte {
	return []byte{
		44, byte(a), 0, 0, 0, 0, 24, byte(b), byte(n), 0, byte(c), byte(d), 0, 0, 0, 1, 0, 239,
		byte(e), byte(o), byte(f), byte(g), 0, 0, 0, 0, byte(h), 0, 0, 14, byte(i), byte(j), 0,
		byte(k), byte(m), 3, byte(p), 1, byte(q), 1, byte(r), 0, 0, 0,
	}
}

func generateString2(params string, method string, startTime, endTime int64) string {
	paramsArray := sm3Hash(sm3Hash([]byte(params + "cus")))
	methodArray := sm3Hash(sm3Hash([]byte(method + "cus")))

	a := list4(
		int((endTime>>24)&255),
		int(paramsArray[21]),
		int(uaCode[23]),
		int((endTime>>16)&255),
		int(paramsArray[22]),
		int(uaCode[24]),
		int((endTime>>8)&255),
		int((endTime>>0)&255),
		int((startTime>>24)&255),
		int((startTime>>16)&255),
		int((startTime>>8)&255),
		int((startTime>>0)&255),
		int(methodArray[21]),
		int(methodArray[22]),
		int((endTime>>32)&255),
		int((startTime>>32)&255),
		len(browserStr),
	)

	e := byte(0)
	for _, b := range a {
		e ^= b
	}

	a = append(a, []byte(browserStr)...)
	a = append(a, e)

	cipher, err := rc4.NewCipher([]byte("y"))
	if err != nil {
		panic(err)
	}
	dst := make([]byte, len(a))
	cipher.XORKeyStream(dst, a)
	return string(dst)
}

// GenerateABogusWithOptions generates a_bogus signature with custom options for testing
func GenerateABogusWithOptions(params string, method string, startTime, endTime int64, rn1, rn2, rn3 float64) string {
	str1 := generateString1(rn1, rn2, rn3)
	str2 := generateString2(params, method, startTime, endTime)
	finalStr := str1 + str2
	return customBase64Encoder.EncodeToString([]byte(finalStr))
}

// GenerateABogus generates the a_bogus signature for a query string and User-Agent
func GenerateABogus(params string, userAgent string) string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	// end_time = start_time + randint(4, 8)
	// We'll mimic this behavior
	rand.Seed(time.Now().UnixNano())
	endTime := now + int64(rand.Intn(5)+4)

	rn1 := rand.Float64() * 10000.0
	rn2 := rand.Float64() * 10000.0
	rn3 := rand.Float64() * 10000.0

	return GenerateABogusWithOptions(params, "GET", now, endTime, rn1, rn2, rn3)
}

// GenerateABogusFromMap generates the a_bogus signature from URL query parameters represented as a map
func GenerateABogusFromMap(params map[string]string, userAgent string) string {
	// Construct the query string in sorted order (standard urlencode behavior)
	u := url.Values{}
	for k, v := range params {
		u.Set(k, v)
	}
	encoded := u.Encode()
	return GenerateABogus(encoded, userAgent)
}
