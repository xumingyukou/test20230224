package sdk

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
)

// 支持plain,file,awskms,env四种schema
// plain:/abcdee,file:/etc/secret.conf
// awskms:/region/namespace/key/kk,env:/API_KEY
type SecretConfig struct {
	url string
	u   *url.URL

	V []byte
}

func (s *SecretConfig) read() ([]byte, error) {
	if s.u == nil {
		return nil, errors.New("Invalid secret url: " + s.url)
	}

	switch s.u.Scheme {
	case "plain":
		return s.readPlain()
	case "file":
		return s.readFile()
	case "awskms":
		return s.readAwsKms()
	case "env":
		return s.readEnv()
	default:
		return nil, fmt.Errorf("invalid secret schema: %s", s.url)
	}
}

func (s *SecretConfig) readPlain() ([]byte, error) {
	v := strings.TrimPrefix(s.u.Path, "/")
	return []byte(v), nil
}

func (s *SecretConfig) readFile() ([]byte, error) {
	path := s.u.Path
	if runtime.GOOS == "windows" || strings.HasPrefix(path, "/.") {
		// windows系统或者相对路径强制去掉开头的/
		path = strings.TrimPrefix(path, "/")
	}

	return os.ReadFile(path)
}

func (s *SecretConfig) readAwsKms() ([]byte, error) {
	l := strings.TrimLeft(s.u.Path, "/")

	if i := strings.Index(l, "/"); i < 0 {
		return nil, errors.New("No aws kms region specified: " + s.u.Path)
	} else {
		if ss, d, err := getSecret(l[:i], l[i+1:]); err != nil {
			return nil, err
		} else if ss == "" && d == nil {
			return nil, errors.New("Get aws kms key failed: " + s.u.Path)
		} else {
			if ss != "" {
				return []byte(ss), nil
			} else {
				return d, nil
			}
		}
	}
}

func (s *SecretConfig) readEnv() ([]byte, error) {
	v := os.Getenv(strings.TrimLeft(s.u.Path, "/"))
	if len(v) == 0 {
		return nil, errors.New("No env var found: " + s.u.Path)
	} else {
		return []byte(v), nil
	}
}

func (s *SecretConfig) UnmarshalText(text []byte) error {
	var err error
	s.url = string(text)
	u, err := url.Parse(s.url)
	if err != nil {
		return err
	}
	s.u = u
	s.V, err = s.read()

	return err
}

// TODO: 增加对access key的支持
func getSecret(region string, secretName string) (secretString string, decodedBinarySecret []byte, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		err = fmt.Errorf("failed to load configuration, %v", err)
		return
	}

	cli := secretsmanager.NewFromConfig(cfg)
	ctx := context.Background()

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err2 := cli.GetSecretValue(ctx, input)
	if err2 != nil {
		if aerr, ok := err2.(*smithy.OperationError); ok {
			err = aerr
		}
		err = err2
		return
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err2 := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err2 != nil {
			err = err2
			return
		}
		decodedBinarySecret = decodedBinarySecretBytes[:len]
	}

	return
}
