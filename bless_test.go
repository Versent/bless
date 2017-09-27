package bless

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/versent/bless/mocks"
)

func TestLoadPublicKey(t *testing.T) {

	tfile := mustWriteTempfile("data")

	data, err := LoadPublicKey(tfile.Name())
	require.Nil(t, err)

	require.Len(t, data, 4)

	_, err = LoadPublicKey("badfile")
	require.Error(t, err)

}

func TestValidatePublicKey(t *testing.T) {
	data, err := ValidatePublicKey([]byte{0xa, 0xb, 0xc, 0xd})
	require.Nil(t, err)
	require.Len(t, data, 4)

	_, err = ValidatePublicKey([]byte{})
	require.Error(t, err)
}

func TestWriteCertificate(t *testing.T) {
	err := WriteCertificate(fmt.Sprintf("%s/%s-%d", os.TempDir(), "abc", time.Now().Unix()), "data")
	require.Nil(t, err)
	err = WriteCertificate(fmt.Sprintf("nothere/%s-%d", "abc", time.Now().Unix()), "data")
	require.Error(t, err)
}

func TestInvokeBlessLambda(t *testing.T) {

	lambdaMock := &mocks.LambdaAPI{}

	lambdaSvc = lambdaMock

	result := &lambda.InvokeOutput{
		StatusCode: aws.Int64(200),
		Payload:    []byte(`{"certificate":"data"}`),
	}

	lambdaMock.On("Invoke", mock.AnythingOfType("*lambda.InvokeInput")).Return(result, nil)

	res, err := InvokeBlessLambda(aws.String("us-west-2"), aws.String("whatever"), []byte("whatever"))
	require.Nil(t, err)
	require.Len(t, res.Certificate, 4)
}

func mustWriteTempfile(data string) *os.File {
	tfile, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}

	_, err = tfile.WriteString("data")
	if err != nil {
		panic(err)
	}

	return tfile
}
