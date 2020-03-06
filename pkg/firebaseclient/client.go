package firebaseclient

import (
	"context"
	"encoding/json"

	"go-tutorial-2020/internal/config"
	"go-tutorial-2020/pkg/errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	sharedClient = &firestore.Client{}
	credentials  = map[string]string{
		"type":                        "service_account",
		"project_id":                  "neogenesis-2a947",
		"private_key_id":              "e39364a249fa1b2841302761ae3f782895e0716f",
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCiZGRPGf7J412H\n5JsVFFwDIClHnntkr1HrrSbom9t9MMwrvBoqSIiw0Vdaad+2VdYWhAsEskzmEc89\nO1+e5kLk287KM7x4YPHYgVa8bo0/GDrkRYCSJf+4o2PGF73Y0XSY5QsIELX1wFuD\nZ16oHs1HbQXJqhvtCuBoxdpJvE4/srdvsrxj/drmN+s3eaCzQXSnRbUGSlqr8bAv\nLx1YTNdBjScbGewpESJjDiJCdJLzOAX8Onek6uNwLs7RE8eFNdGtHaEO0mWRJ80j\ncee9hvaamB+7OKm6r3pMCzakQaogR2hlCHwKfcRFl5JNUKWLdkigKeCfgTY5FOCL\nuww2aVXfAgMBAAECggEAMPvYWBXuyYYoSdv2vZqJELZMTVqsPNg3fUAbKvIMeIxW\nIeEZAWUkujVDRLYH8a+dpugIjM+dq452tTEqLDuntPHqxApsEOGpQdXtzGQKfhw4\nc2/VkwAcTV0XtQKnVPYFVjNMqw+jW3A9CnuNTWgRIrjrlIOX7d4oh+IacMB41/R9\nAq+Cwz29uFnzqnw1lmYY6cPX8fRFIdr9ZUbgLxRVZTJ43sZVTOqKYnvH/VSIzvY7\nu4OUuMF5jrEfaHloZ0v7EJ5u7iSVnLR6Yip5dLb4tRx4uPmY9T1YGmRNeggInEVq\nS3QP5znMsiqYd1dH7fgatqxVOrCvgOtgQ+PF5nnJEQKBgQDWEloa9oDJ+j0ZKTU9\nRQR1ShwHxinuSISwSRXCum5OYk7jQaaKIOtmO77oMsGBAB2qnPXRERWd6i7sepQI\nTU2w2I95FJHiJ3wla29JNd70iiZm9g1mouy9Juf783hfrLNicog4grWbDjOHkMd4\nRfS1JFwpwsMICn2h6LtqE6cXcQKBgQDCMs43Xwr6X/Tw6OA9TFo7/ZopziHmCBat\nx7O+PaUVMLMkuJ0WTDAugDVXHpQf4j3c7sFjM4dVxf03FKussN+LVDV7Kqvlo/Tz\nTjLPivWeJ5v5EtacB0CMksxthWQciJ8YM6xY2urmbX+GwIB1Q2HGeGCyZpbSun11\ng9I6LWG6TwKBgDnp6RqWSa1a16Cw90hGHbilfgPUZo+iatNOUmbGDQWDrxfoMOf5\nk/WqgDFNWfBOWbhIknAnERQRwPQVDWtZjoUjcV0uZXErgXiWIhtKSbEalt0P//we\nY7Ggju4opg4sKLOfjJ7NLdhu9R3d2zj8pAFFfvGFUUIhpG9jCSDfhDrhAoGAVzBY\nRaP2WdMbNc5YXy0YljaWMI7LyWt4Qy2WFaO3qnvi8mmwnYI3X6lQSX0BJA9/luHb\nEZ7g9DGgLkqpiS8gLn3wRQwzgTbLkzFYvrW08Pz3mixLDmJzKKn//mwVNnpgc40/\no+Ul8a7XwrhK9Fr8Ww9Q2sgUjygLi+dTS6t0ZSECgYA9oMzp7vD4A8wM6KvF8GyL\ntIa5uIjbSMsMwAe5mawCgFk4gvBn7BBLXOCXoF3KY6EZu6hARHqAgZ4LlsoUfd2V\nJiQonaZBj3y4/RmNBBAgdLlohBBKQLhKuoQsPSLdE43cdS2QuXSDcvGj6bVvmwmQ\ngoNPjyAs9/qgQQ3UClrQEg==\n-----END PRIVATE KEY-----\n",
		"client_email":                "firebase-adminsdk-7zlk2@neogenesis-2a947.iam.gserviceaccount.com",
		"client_id":                   "110912029148578276583",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-7zlk2%40neogenesis-2a947.iam.gserviceaccount.com",
	}
)

// Client ...
type Client struct {
	Client *firestore.Client
}

// NewClient ...
func NewClient(cfg *config.Config) (*Client, error) {
	var c Client
	cb, err := json.Marshal(credentials)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to marshal credentials!")
	}
	option := option.WithCredentialsJSON(cb)
	c.Client, err = firestore.NewClient(context.Background(), cfg.Firebase.ProjectID, option)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to initiate firebase client!")
	}
	return &c, err
}
