// SPDX-License-Identifier: MIT-0
// LICENSE: https://spdx.org/licenses/MIT-0.html

package dht

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"
)

// MakeHTTPSServer returns an *http.Server with a
// random self-signed tls certificate and funky logger
func MakeHTTPSServer(mux *http.ServeMux) *http.Server {

	makeSSLCert := func() tls.Certificate {

        priv, err := rsa.GenerateKey(rand.Reader, 4096)
        if err != nil {
            panic(err)
        }
    
        certTemplate := x509.Certificate{
            SerialNumber: big.NewInt(1658),
            Subject: pkix.Name{
                Organization:  []string{""},
                Country:       []string{""},
                Province:      []string{""},
                Locality:      []string{""},
                StreetAddress: []string{""},
                PostalCode:    []string{""},
            },
            IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
            NotBefore:    time.Now(),
            NotAfter:     time.Now().AddDate(10, 0, 0),
            SubjectKeyId: []byte{1, 2, 3, 4, 6},
            ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
            KeyUsage:     x509.KeyUsageDigitalSignature,
        }
    
        derBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, priv.Public(), priv)
        if err != nil {
            panic(err)
        }
    
        certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
        keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
        tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
        if err != nil {
            panic(err)
        }
    
        return tlsCert
    }
    
	TLScfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{makeSSLCert()},
    }
    
    logger := func(handler http.Handler) http.Handler {
        log := func(w http.ResponseWriter, r *http.Request) {
            log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.Host, r.URL)
            handler.ServeHTTP(w, r)
        }
        logHandler := http.HandlerFunc(log)
        return logHandler
    }
    
	return &http.Server{
		Handler:      logger(mux),
		TLSConfig:    TLScfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}
