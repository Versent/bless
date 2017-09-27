package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/versent/bless"
)

var (
	app   = kingpin.New("versent-bless", "A command line client for netflix bless.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()
	login = app.Command("login", "Login and retrieve a key.")

	region              = login.Arg("region", "AWS Region.").Required().String()
	lambdaFunctionName  = login.Arg("lambda_function_name", "Lambda function name.").Required().String()
	bastionUser         = login.Arg("bastion_user", "Bastion user.").Required().String()
	bastionUserIP       = login.Arg("bastion_user_ip", "Bastion user IP.").Required().String()
	remoteUsernames     = login.Arg("remote_usernames", "Remote user names.").Required().String()
	bastionIps          = login.Arg("bastion_ips", "Bastion IPs.").Required().String()
	bastionCommand      = login.Arg("bastion_command", "Bastion command.").Required().String()
	publicKeyToSign     = login.Arg("public_key_to_sign", "Public key to sign.").Required().String()
	certificateFilename = login.Arg("certificate_filename", "Certificate filename.").Required().String()
	kmsAuthToken        = login.Arg("kmsauth_token", "KMS Auth Token.").String()
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case login.FullCommand():

		bless.ConfigureAws(&aws.Config{Region: region})

		bless.SetDebug(*debug)

		publicKeyData, err := bless.LoadPublicKey(*publicKeyToSign)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		bless.Debug("publicKeyData: %s", string(publicKeyData))

		publicKey, err := bless.ValidatePublicKey(publicKeyData)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		payload := &bless.Payload{
			BastionUser:     *bastionUser,
			BastionUserIP:   *bastionUserIP,
			RemoteUsernames: *remoteUsernames,
			BastionIps:      *bastionIps,
			BastionCommand:  *bastionCommand,
			PublicKeyToSign: publicKey,
		}

		if kmsAuthToken != nil {
			payload.KmsAuthToken = *kmsAuthToken
		}

		bless.Debug("payload: %v", payload)

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		log.Printf("payload_json is: %s", string(payloadJSON))

		resultPayload, err := bless.InvokeBlessLambda(region, lambdaFunctionName, payloadJSON)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		bless.Debug("resultPayload: %v", resultPayload)

		err = bless.WriteCertificate(*certificateFilename, resultPayload.Certificate)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		log.Printf("Wrote Certificate to: %s", *certificateFilename)

	}
}
