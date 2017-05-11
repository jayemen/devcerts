package certutil

import "testing"
import "crypto/x509"

const certPem = `-----BEGIN CERTIFICATE-----
MIIDnTCCAoWgAwIBAgIBATANBgkqhkiG9w0BAQsFADBoMQswCQYDVQQGEwJDQTEQ
MA4GA1UECAwHT250YXJpbzERMA8GA1UEBwwIS2luZ3N0b24xETAPBgNVBAoMCGpt
bi5saW5rMQswCQYDVQQLDAJJVDEUMBIGA1UEAwwLam1uLmxpbmsgQ0EwIBcNMTYw
OTEyMDIwNzU0WhgPMjExNjA4MTkwMjA3NTRaMGgxCzAJBgNVBAYTAkNBMRAwDgYD
VQQIDAdPbnRhcmlvMREwDwYDVQQHDAhLaW5nc3RvbjERMA8GA1UECgwIam1uLmxp
bmsxCzAJBgNVBAsMAklUMRQwEgYDVQQDDAtqbW4ubGluayBDQTCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAORHiweYhJDILPx1QSezSF92rPdYDk0M8zSu
OnApABUFLi5v0CvKuRJPq0DQ5C1WaZAUNEW+U9WsbhhqTnIkKDEQiup0vJElhRXh
N9OS37gsTbczkFSeFKmUcnQhbpn/5Zek7BTGbEW4GzaAwFyl827jml5mngo2ufIN
p8l3FzAdXCqG2cZo8Cy/gWFQR9hn+CfoWqO+JeX3Zx4cC7RBC3QM27SDLD/JV78Z
hhGNY+oekDapEzD4z/VXvyUzQR0vyEYO4bhmTsfgqyqE+Kn6C+Gtb6hKG7OwOyrl
9Ds4lF0BebaOgZU/7SUX4dt8oRRX2zaGeLKKB3WSPSl/lhAtGb0CAwEAAaNQME4w
HQYDVR0OBBYEFAMhKU6xuxGzodunQSsh+jr3016DMB8GA1UdIwQYMBaAFAMhKU6x
uxGzodunQSsh+jr3016DMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
AJB+BnaJrVpVO/OU9CO+IrvtlaO5f8TO5uB3RbenmlddM4b+CRikWPPE4zDLALA/
rM4ioPzWXQ3fon6dx2u/00g1cJCK73bhK3OUfDd+Dtfj1yM5LLoWo5uJq+xDQa3O
L5XO/XxZ4xJ6t4t3QhQz6cbAVFkzlmmMvBNk41KBvYkAFANAk/qbaRcR/IDiaKq6
TlyeSg9H7wrc1mKOwjmRrgVpcUMfqx88//qSExweeY/DQ6RKGoqFHc5bjeWlw2EB
3r1ZbqoJ/WVX84SdxsUdsSu3YF2inF9CD3Zcv4GMqUMY5h8QWdPkJpZfDX0cIqsY
6P24jUaeXaZiVEChVz0npTw=
-----END CERTIFICATE-----`

const keyPem = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDkR4sHmISQyCz8
dUEns0hfdqz3WA5NDPM0rjpwKQAVBS4ub9AryrkST6tA0OQtVmmQFDRFvlPVrG4Y
ak5yJCgxEIrqdLyRJYUV4TfTkt+4LE23M5BUnhSplHJ0IW6Z/+WXpOwUxmxFuBs2
gMBcpfNu45peZp4KNrnyDafJdxcwHVwqhtnGaPAsv4FhUEfYZ/gn6FqjviXl92ce
HAu0QQt0DNu0gyw/yVe/GYYRjWPqHpA2qRMw+M/1V78lM0EdL8hGDuG4Zk7H4Ksq
hPip+gvhrW+oShuzsDsq5fQ7OJRdAXm2joGVP+0lF+HbfKEUV9s2hniyigd1kj0p
f5YQLRm9AgMBAAECggEAWmHpLgy5EAnxpdNXBLz7PrDiMtxubRtff9Ar2xSgr7Hp
YwFqTqxpMlLQ30zVyw2XpjAZsjN1RfiLbqdIf/DI6QQ/vCyULHPKiasuS1qvsV/5
NTv5PUeJrsrTth82h0rGQJBP2LnnTINkYuP/Ra9+/rym9hFzKWAZpzi44g+A4s89
fIHotrIEOO+VjUa1WMFSOPI+nmbJ+8PmyIIFqjuXvWmk9M0IvizubjxWLXiya0Wr
gRiDuJruXznsvqI6k6BmZ1V8fSTyyneKKwSikCF5B7CpZIdwb4nqMnB5fAxh7JzN
S86N2re91N4uXsBxbyB+hkpz7XjKYH6KdXV9JlVKgQKBgQD8uYicPo8JFRp6UdEg
WJShLlgwFceOKdV+koL6KdbjBTo/LJrf3ZzHIKc/v84R5T/HwYbE2MOGd6V67VcD
RmBZNmnzYqVFzQgM6ZOoVasMaJGRqJUxm1YHPgGD8ezVOURuN4gT80weH7EkPKJ6
IohSGAcBxNOptn1J4yqCKoVXRQKBgQDnPOg/QteLqmlVi4ct0ZtuE//F/rx1iMWT
+Ix/Z5ht4Poha4N0K6iNqtmgKPQBWxjirGwYRLn5PYaUo69R8nv85U0MJLnkj4+u
6+kW+zMV9LElpyRty4h5+DtCWHiOyXVcO4s/A32SCz3nSOmKg+zXOr6BYxv8ft2X
4I5ulgiEGQKBgQDD9xHpJdE/169aXgrtLALEIO2dC6ZbpDC6Ht2VIdBQ5QLPbcUC
BhPFjJpjolUmJz+Xo4bfKL4kjK4ybctk7LNVOg5Z/YnuYBf3+z7V7ufdjVAjRDe7
6ZmBsCD4sSVWCTv4wvKvlZ7WVPjFAodycUiHb74vLvJ5zNnF63JQ0KvoiQKBgF3H
ut0lK6uuAig6fSlc7+911u6iwCXewVqgm8Jz7kLp0ifJpbdwmVxTJQ2qbkM6gd40
VWaGQPJPPIx90fWnJRfMmzHIl0eV3YzwikjSucY2xb1iiwioWgI1ZTskDEjEdX9h
eriknsGjI4jwbh7KIDyty2NIIaqGfTJCVSGOYYfJAoGBAJZ3Lgev9Bh7Dp22A9Y+
wicoMWpkCryIZfciiHbPlif1AVgvCEXAopgG3Na/VKNNWBiJjHC5lWpdy3gQ7RVS
1eIvNsshk8BhMPRujtXbNgUNW8m2dvJilveQsaPFKgnUNhTq0FWahnkVjjt4YTD9
Q62AzR0UM8vcpf6475tiAvma
-----END PRIVATE KEY-----`

type wrapper struct {
	*testing.T
}

func wrap(t *testing.T) wrapper {
	return wrapper{t}
}

func (t *wrapper) test(condition bool, msg string) bool {
	if !condition {
		t.Error(msg)
	}

	return condition
}

func (t *wrapper) eq(want interface{}, got interface{}) bool {
	if want != got {
		t.Errorf("expected '%+v' got '%+v'", want, got)
		return false
	}

	return true
}

func (t *wrapper) nilErr(err error) bool {
	if err != nil {
		t.Error(err)
	}

	return err == nil
}

func (t *wrapper) sliceEq(want []string, got []string) bool {
	if !t.eq(len(want), len(got)) {
		return false
	}

	eq := true
	for i := range want {
		eq = !t.eq(want[i], got[i]) && eq
	}

	return eq
}

func TestLoad(testing *testing.T) {
	t := wrap(testing)

	store, err := Load([]byte(certPem), []byte(keyPem))

	if !t.nilErr(err) {
		return
	}

	t.test(len(store.Chain) == 0, "the chain should be empty")
	t.eq("jmn.link CA", store.Cert.Subject.CommonName)
	t.eq("CA", store.Cert.Subject.Country[0])
	t.eq("Ontario", store.Cert.Subject.Province[0])
	t.eq("Kingston", store.Cert.Subject.Locality[0])
	t.eq("IT", store.Cert.Subject.OrganizationalUnit[0])
	t.eq("jmn.link", store.Cert.Subject.Organization[0])
}

func TestNew(testing *testing.T) {
	t := wrap(testing)

	store, err := Load([]byte(certPem), []byte(keyPem))

	if !t.nilErr(err) {
		return
	}

	child, err := store.New("test.com", []string{"a.b.c.com"}, []string{"1.2.3.4"})

	if !t.nilErr(err) {
		return
	}

	t.eq("test.com", child.Cert.Subject.CommonName)
	t.sliceEq([]string{"CA"}, child.Cert.Subject.Country)
	t.sliceEq([]string{"Ontario"}, child.Cert.Subject.Province)
	t.sliceEq([]string{"Kingston"}, child.Cert.Subject.Locality)
	t.sliceEq([]string{"IT"}, child.Cert.Subject.OrganizationalUnit)
	t.sliceEq([]string{"jmn.link"}, child.Cert.Subject.Organization)
	t.sliceEq([]string{"a.b.c.com"}, child.Cert.DNSNames)
	t.eq("1.2.3.4", child.Cert.IPAddresses[0].String())
	t.eq(1, len(child.Cert.IPAddresses))

	intermediates := x509.NewCertPool()
	roots := x509.NewCertPool()

	for i, c := range child.Chain {
		if i != 0 {
			intermediates.AddCert(c)
		} else {
			roots.AddCert(c)
		}
	}

	_, err = child.Cert.Verify(x509.VerifyOptions{
		Intermediates: intermediates,
		Roots:         roots,
	})

	if !t.nilErr(err) {
		return
	}
}
