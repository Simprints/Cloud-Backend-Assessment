package users

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

const creds = `{
  "type": "service_account",
  "project_id": "simprints-backend-hiring",
  "private_key_id": "e9fd63af042e239ad4f8a41aba81e8652ec334e1",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDW9t6vQqewSvdu\nMTvepmckl+pdFHSLCeJuIrM+ee9o1sEUC7insWSNEyRJk36eNZaBUWTTf+lYqVgD\nACM9eV/42r/uSbxnLLcnKw02rLWiO6gTfhZKqlNI/fgC+E88V05JP3FB2Q5Azpxd\njIX1GvwOFw0Al+wzkOVgJpz1ekml1Gdfki1B4yIWYTmoVSfRSyguq+RwBM/k34xE\nLn5ZAM1/c4qf9/fkGjwD/3Vj5Cdd/hb0UNuZ1L0YSXT/jrSzdGgzT+gkCP346nRJ\nF63C1zH/bWHA3H6Ca9tLFA7yh4a3I8TP/DACRWzV0iGSYAhRU5DLboAYAtBn+4PE\nXbdC6GZxAgMBAAECggEACFFaXcH4otzNvO7tBAEgl9UEgahFV61w3GXdlW7qrX30\nz+GhuiRmMYLCelHZmgXVDmkvu3LVPMPlts0CMBJinSqat13VS1E6v4peCyXs68t+\n5g0wj+BOE3KXTocc6sbuja1FpMCBOY79FC4YL1slUmbALywRvM+QHpOz6lg6Yg8X\nIxGu+n6hJxvr2cy+wPQk3Bl5uBNrAIUNoKqx7Q1lch+z13NWd9SABrff7Th1PT/G\nR+Jf3+DlXzJAnXiuMaSPvGt220sYxZS5G8WtDSQzcgliTv9kqy0xCwpkg7q5KGl6\nNtU7E9ToYaHyLjzoQLDNUZyzWpAS/3T3nNbpVzsl2QKBgQD765exTas1mZrLjki4\nFFgbcRXz+yoC2V2KytS/LnixYbmdgTgfes/g3CQtCLw7UUeQbDKtpDXQ0N82LrGV\n9OkEC9m5OKY/xJwWORycwpQdOiVKoK7SjRGZn092Gf8nOPOu2G6uxkkzWh3NDmBs\nK+LzX9PT6hV3uzJ7bmWS4q2w6QKBgQDachDdh9yMRPEJHY9AW/VG46LA0JsK1V35\n+z4NLTMYq0uqZqv3Vf/UpCZ6DWjp3vSVE/HvzGFOBArtsQWIz8GTk/OWYZeCO+DD\n/5+Fd4K1poVyTzBRSNDKAe+5EylTeoo4/bQnKrLmVUWhb0G4rixIU7bxQQbeSyTD\n3B5b9eDUSQKBgQDcrEL0zVRsX2F5benFVgzX/Pd+AUWLuVx3d7VkwxB2UWSG0+qV\nqL7v+ea2jDBWxZwqppy9/lol0NG2ZLCq6x4yrS7LURRQR6lyzhSCPPABqi2AccCy\nL2B7cVHp4lvfv8O2JWDPOGJm2UnBlhZgqxDin86ukx67Av/1n37abDY6AQKBgQCK\nuDkBlU93PCidE0pvInaGR/SI4XAz1v9Qyj3DfFqgZdctJPo7nT9TN9K/W1iue8ly\nCjJvh6ibNHIEM5BCKzzQjPn5G4xtRb0cem5BAX3eARtpVeRnGgiM3+Ht878gpga0\n3lfTL4hgQPJw7AgeUW0JmS/p0NOdwrZcMqKM332hEQKBgQDEhVIFJGfTHpn3QXD6\ndUKDJgG44Rb73HFdd4XLxruKziKYvTRbSNtCkq6+2fSZGYV6gO2d04P+UApno0Uj\nDe6QsZlL72ZRK0UI7WYP5aR5KOVRzXXbeHQz0h1EJYEcPXlCW4/aFg352H5gh6tg\nDfXKzKzDe8KrVex3ztKfTxcxMw==\n-----END PRIVATE KEY-----\n",
  "client_email": "senior-backend-swe@simprints-backend-hiring.iam.gserviceaccount.com",
  "client_id": "101501086024594090018",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/senior-backend-swe%40simprints-backend-hiring.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, userId string) (User, error)
	DeleteUser(ctx context.Context, userId string) error
}

type repository struct {
	firestoreClient *firestore.Client
}

func NewRepository(ctx context.Context) (Repository, error) {
	creds, err := google.CredentialsFromJSON(ctx, []byte(creds), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}
	project := os.Getenv("GCLOUD_PROJECT")
	if project == "" {
		return nil, fmt.Errorf("missing required environment variable GCLOUD_PROJECT")
	}
	firestoreClient, err := firestore.NewClient(ctx, project, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	repo := &repository{
		firestoreClient: firestoreClient,
	}
	return repo, nil
}

func (r *repository) CreateUser(ctx context.Context, user User) error {
	_, err := r.firestoreClient.Collection("users").Doc(user.Id).Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUser(ctx context.Context, userId string) (User, error) {
	user := User{}
	snapshot, err := r.firestoreClient.Collection("users").Doc(userId).Get(ctx)
	if err != nil {
		return user, err
	}
	err = snapshot.DataTo(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) DeleteUser(ctx context.Context, userId string) error {
	_, err := r.firestoreClient.Collection("users").Doc(userId).Delete(ctx, firestore.Exists)
	if err != nil {
		return err
	}
	return nil
}
