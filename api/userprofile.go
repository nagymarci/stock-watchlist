package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	userprofileModel "github.com/nagymarci/stock-user-profile/model"
)

type UserprofileClient struct {
	host string
}

func NewUserprofileClient(h string) *UserprofileClient {
	return &UserprofileClient{
		host: h,
	}
}

func (uc *UserprofileClient) GetUserprofile(userId string) (userprofileModel.Userprofile, error) {
	resp, err := http.Get(uc.host + userId)

	if err != nil {
		return userprofileModel.Userprofile{}, fmt.Errorf("Failed to get userprofile [%s] with error [%v]", userId, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 299 {
		var response string
		fmt.Fscan(resp.Body, &response)
		return userprofileModel.Userprofile{}, fmt.Errorf("Failed to get [%s], status code [%d], response [%v]", userId, resp.StatusCode, response)
	}

	userprofile := userprofileModel.Userprofile{}

	err = json.NewDecoder(resp.Body).Decode(&userprofile)

	if err != nil {
		return userprofileModel.Userprofile{}, fmt.Errorf("Failed to deserialize data for [%s], error: [%v]", userId, err)
	}

	return userprofile, nil
}
