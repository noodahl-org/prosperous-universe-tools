package fio

import (
	"log"
	"prosperous-universe-tools/clients/fio/models"
)

func (f *FIOClient) Login() error {
	auth := models.AuthResponse{}
	resp, err := f.client.R().
		SetResult(&auth).
		SetBody(models.AuthRequest{
			Username: f.conf.Username,
			Password: f.conf.Password,
		}).
		Post("/auth/login")
	if err != nil {
		return err
	}
	log.Printf("login: %v", resp.StatusCode())
	f.client.SetHeader("Authorization", auth.AuthToken)
	return nil
}
