//addVehicle.go
//This file contains an example of adding a new vehicle to Icarus Server
package main

import (
	"fmt"

	icarus "git.ironzone.ace/icarus/icarusClient"
)

const (
	defaultCert string = `-----BEGIN CERTIFICATE-----
MIIC9TCCAd2gAwIBAgIRAKRQOLvrvmORBJxJKkhCpWgwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xOTAxMDkxNzQ1NTdaFw0yMDAxMDkxNzQ1
NTdaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDGzdsktG2DGiQEt7ce1sWlcSc1QNbpLcRemcrGxJKw2JeWYY42R5Le
+umtLV/xy0+ZIA47iHETj0IFFYjWdixmz5/yHnnnJbz8uKinbk3eTmaR6y+EwSAp
gFjXFYjRgif4wPk0qnkgHaI+TJXn2dbBnpv0cX34aUKwaCa/qh0XEZ2nqmjjeowx
mqpD4etICnaMKdJg2Z+da/YG8ExFnwYpzNS9QdfujAxHJ7DoMhPZnyc/sCmaBq+X
eAZbMHtWFuv/24lA/KyJBmCEQGp2x9tn+HmM89SQOj1yOwOqZB87+rjENhp87rgh
MD4vB93/Mzk48tC1LYlr7cLoBh22tskRAgMBAAGjRjBEMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaH
BH8AAAEwDQYJKoZIhvcNAQELBQADggEBAKYVTX6fzjOU63QxdtSs9Ot6GdAqaknQ
baTiRnAsXJuRCDZpRIEQKFECC7tdEBbyCh5FiyjpVqxn+U2/OcdNHPYxdHRevRWM
vmNqxeOjha62Bp/JKoN0WR/NfZ7oSHFQYm5kGxTSt6n/BKovWOHI4je0PHD/YITt
Xw5IgIEPDS2eecE+jrHBU596X1jeHSuk+XdQ9Hmo1WFYE9UK7985Oy77+zaNSKdu
ZbBWN877w5AMdmswrC6HcDCXY0Sb3Noadfl9VcDqD69hq0vWPqAomrbjyyTUQCY3
ru4kYaW0u8509KHlpN6ixQWmGAmVhMWtf0g1kPpzJ9HihRtAYUp5Jj4=
-----END CERTIFICATE-----
`
	defaultKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAxs3bJLRtgxokBLe3HtbFpXEnNUDW6S3EXpnKxsSSsNiXlmGO
NkeS3vrprS1f8ctPmSAOO4hxE49CBRWI1nYsZs+f8h555yW8/Liop25N3k5mkesv
hMEgKYBY1xWI0YIn+MD5NKp5IB2iPkyV59nWwZ6b9HF9+GlCsGgmv6odFxGdp6po
43qMMZqqQ+HrSAp2jCnSYNmfnWv2BvBMRZ8GKczUvUHX7owMRyew6DIT2Z8nP7Ap
mgavl3gGWzB7Vhbr/9uJQPysiQZghEBqdsfbZ/h5jPPUkDo9cjsDqmQfO/q4xDYa
fO64ITA+Lwfd/zM5OPLQtS2Ja+3C6AYdtrbJEQIDAQABAoIBAQCMQh4TFkyRCzdQ
MME8O7Bz2ZIc6yL0njqFt6EtfPA1XooMKcWom/SN5p5IdNPVBmihEtGXxNpqP08H
wTqqe/M1kdQ5gLDmmGRuNGWgwpyjc9K/rhr3YT2sqgWDsYi2r0o+IP9w3bjZJK8b
nvLAAZuXPKyw2AVU5gaL6N81p/IgHCoQ6GAkE19Bk2EK/vRoSKtmZsmsUY3RrezC
HTSX4R4it2eMXBAxEF6stPK0tkv2BvOItoe5BGfWnd4VG2okVcHFkRPiaIiz60iy
02m6hei2dTKRe/GxaNNDNW8jNqTyJ4ZbUK8MeFVyCmFe+Zu8/LLF6bffj46HE00u
yWovzFMxAoGBANq5gT27iGukc2MgXDc/yLKcwmXTU7hGH7TRBCu11Ms9fgInCsrz
x9SRhI5KvysOVe80uPdtyoMx9+I/hzDDNBj0IrIXlO3tpgs26dqUN20CqRbepGwA
aDXxy2D2lD4EHm1pDHU5SXngUECO1fyb9ErG4qeIdE4cfQ8h6ZhRblQNAoGBAOiv
QvVNvfEm7Qz1BnCuwBTJfkScTJqLuel525c5PIj68w/B6D/cbYaKXr5kSFKroht9
8hr4lSHawLCAiBHc4mzXmyI1enxTH1+Sck2nzLxOWCqhExzmpNoKLI0T+Xhy0XH1
dBQuD7QBdWsP4Y+shBdw/ehVPRqmwqY94Lgr5HQVAoGBAKWNhaJ5QK/hIKlmBAaZ
k8qF1qqGAzdWdIdDMbn3/mH7YFY2wPeO/7EIl+Gv9/SZ/Dd7m4lEo+UbvDmWxjgF
eHhuyZgtOz/AAk84uFcGmtE7E0tJKADLahVyt/LjkJ9ENNexjIlp3BCQ1Y2Xz6ZN
UOIMmeAe65F4BLygeZQeBrk9AoGBAOPyxmrgBUMY+kOmSu/bEkuK9Ysrf5QrbC8A
9RHZvacICVQXh3oAbL/QEG7+eSecAsxh/utTOW4YCose765oMN2l/tFtiJgBKowL
QLU4vMaBDbh9Yeb/QOJl8y0mM1A/U1YLuvMGCNY0U55VyYhh3mnEhMm1r43LbodD
uUFTppPdAoGBAMHfXeOBXNGaCJf9wAF+qkfTG17z+rplO1ny9vX8p31i3F3wMJG0
rDkznuYCa6PqBeav2UjbsfeqkynBggKM0wqMUcrSSCzKJEpaOjdaq1fw92Wf4tgq
aJWMolGAClXedrb1jV2zZgl9xwi8kG1Y9EL4uVuCz2k03dFGuznMTW2v
-----END RSA PRIVATE KEY-----
`
)

func main() {
	//Create a new query pointed at the IcarusServer instance
	query := icarus.NewQuery("10.59.144.202", "179")

	//Add new vehicle to IcarusServer
	addSeq := query.AddNewVehicle("10.58.147.217", "22263", "KIRINKI-104", "25594", 0, make([]string, 1), []byte(defaultCert), []byte(defaultKey), 0, icarus.DefaultC3poTime, icarus.DefaultDaedalusTime)

	//Authenticate to the server
	resp, ok := query.Authenticate("valinar", "gitlabsux;(bloodybreach")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	//Uncomment this line to show the JSON query being sent to the server
	//query.ShowQuery()

	//Execute the query and read the responses
	responseChan, _ := query.Execute()

	fmt.Println("Waiting for responses:")
	response := <-responseChan
	fmt.Println(response.Get(addSeq))

	//Uncomment this line to show the JSON response returned by the server
	//response.ShowResponse()
}
