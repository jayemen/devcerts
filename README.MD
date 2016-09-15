# devcerts
Generates development certificates, signed by a configured development root certificate. The UI is provided via a local web-server.

The web-server port is specified via the 'PORT' environment variable. The root CA certificate and private key are read from the files 'ca.crt' and 'ca.key', respectively, from  the directory from which devcerts is launched. These files must be in PEM format.