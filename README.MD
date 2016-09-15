# devcerts
Starts a local web-server, which can be used to generate development certificates, signed by a single development root certificate.

The web-server port is specified via the 'PORT' environment variable. The root CA certificate and private key are read from the files 'ca.crt' and 'ca.key', respectively, from  the directory from which devcerts is launched. These files must be in PEM format.