package bless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

var (
	lambdaSvc lambdaiface.LambdaAPI = lambda.New(session.New())

	debug = false
)

// Payload used to encode request payload for bless lambda
type Payload struct {
	BastionUser     string `json:"bastion_user"`
	BastionUserIP   string `json:"bastion_user_ip"`
	RemoteUsernames string `json:"remote_usernames"`
	BastionIps      string `json:"bastion_ips"`
	BastionCommand  string `json:"command"`
	PublicKeyToSign string `json:"public_key_to_sign"`
	KmsAuthToken    string `json:"kms_auth_token,omitempty"`
}

// Result used to decode response from bless lambda
type Result struct {
	Certificate string `json:"certificate,omitempty"`
}

// ConfigureAws enable override of the default aws sdk configuration
func ConfigureAws(config *aws.Config) {
	lambdaSvc = lambda.New(session.New(config))
}

// SetDebug enable or disable debugging
func SetDebug(enable bool) {
	debug = enable
}

// LoadPublicKey load the public key from the supplied path
func LoadPublicKey(publicKey string) ([]byte, error) {
	data, err := ioutil.ReadFile(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load public key")
	}

	return data, nil
}

// ValidatePublicKey validate the public key
func ValidatePublicKey(publicKeyData []byte) (string, error) {
	if len(publicKeyData) == 0 {
		return "", errors.New("Empty public key supplied")
	}

	return string(publicKeyData), nil
}

// InvokeBlessLambda invoke the bless lambda function
func InvokeBlessLambda(region, lambdaFunctionName *string, payloadJSON []byte) (*Result, error) {

	input := &lambda.InvokeInput{
		FunctionName:   lambdaFunctionName,
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("None"),
		Payload:        payloadJSON,
	}

	res, err := lambdaSvc.Invoke(input)
	if err != nil {
		return nil, errors.Wrap(err, "Lambda invoke failed")
	}

	if aws.Int64Value(res.StatusCode) != 200 {
		return nil, errors.Wrapf(err, "Lambda Invoke failed: %v", res)
	}

	resultPayload := &Result{}

	err = json.Unmarshal(res.Payload, resultPayload)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to deserialise JSON result")
	}

	return resultPayload, nil
}

// WriteCertificate write the generated certificate out to a file
func WriteCertificate(certificateFilename string, certificateContent string) error {
	err := ioutil.WriteFile(certificateFilename, bytes.NewBufferString(certificateContent).Bytes(), 0600)
	if err != nil {
		return errors.Wrap(err, "Failed to write certificate file")
	}

	return nil
}

// Debug write debug messages if it is enabled
func Debug(message string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
