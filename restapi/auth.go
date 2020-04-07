package restapi

import (
	"github.com/edwardmartinsjr/pararius-scraper/config"
	"github.com/sirupsen/logrus"
)

// NewAuth - Creates a new instance of Auth configuration
func (a *Auth) NewAuth(logger logrus.FieldLogger, authenticateConfig *config.Authenticate, projectConfig *config.Project) (*Auth, error) {
	return &Auth{
		logger:             logger,
		authenticateConfig: authenticateConfig,
		projectConfig:      projectConfig,
	}, nil
}

// Auth - Contains a collection of user authenticate configuration and properties
type Auth struct {
	authenticateConfig *config.Authenticate
	projectConfig      *config.Project
	logger             logrus.FieldLogger
}

// UserAuthenticate - Authenticate user application by using the auth endpoint
func (a *Auth) UserAuthenticate() (cr *ClientResp, err error) {
	// Instantiate new API Client
	client := New(nil)

	var jsonStr = []byte(`{
		"method": "ums",
		"email": "` + a.authenticateConfig.Email + `",
		"password": "` + a.authenticateConfig.Password + `"
	  }`)

	cp, err := client.post(a.projectConfig.ProjectID, authRequest, jsonStr, "")
	if err != nil {
		a.logger.Error("couldn't authenticate user: %s", err)
		return nil, err
	}
	return cp, nil
}

// Features - Dataset features
type Features struct {
	Bedroom string
	Area    string
	Price   string
}

// Insert - Add records to Dataset using REST API's.
func (a *Auth) Insert(features Features, refreshToken string) error {
	// Instantiate new API Client
	client := New(nil)

	var jsonStr = []byte(`[
		{
			"bedroom": "` + features.Bedroom + `",
			"area": "` + features.Area + `",
			"price": "` + features.Price + `"
		}
	   ]`)

	_, err := client.post(a.projectConfig.ProjectID, ds+a.projectConfig.TrainingDataSet, jsonStr, refreshToken)
	if err != nil {
		a.logger.Error("couldn't authenticate user: %s", err)
	}

	return err
}
