package provider

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
)

func TestAuthanticationOK(t *testing.T) {

	defer gock.Off()

	var btb = NewBitbucket()

	gock.New(btb.Auth.EndPoint).
		Post("/").
		Reply(200).
		JSON(map[string]interface{}{"access_token": "rKO3vh7g23qJrUM5BZuHmQ9Gw3c4hELARui5cBb4EBqIOTiI0-ueM6Qt-kXTW4YbignB__gJLCPL4D_-jBs=",
			"scopes":        "pullrequest account:write webhook repository:write project:write",
			"expires_in":    7200,
			"refresh_token": "bh8vuASNWHjMVmkArN",
			"token_type":    "bearer"})

	err := btb.Authenticate()
	if err != nil {
		t.Errorf(err.Error())
	}

	if btb.Token.AccessToken == "" {
		t.Errorf("Access Token couldn't be fetched")
	}
}

func TestAuthanticationNOTOK(t *testing.T) {

	defer gock.Off()

	var btb = NewBitbucket()

	gock.New(btb.Auth.EndPoint).
		Post("/").
		Reply(401)

	err := btb.Authenticate()
	if err == nil {
		t.Errorf("Authantication didn't work correctly")
	}
}

func TestCheckPermissions(t *testing.T) {
	var btb = NewBitbucket()
	btb.Token.Scopes = "pullrequest account:write webhook repository:write project:write"
	err := btb.CheckPermissions()
	if err != nil {
		t.Error(err.Error())
	}

}

func TestCheckPermissionsNOTOK(t *testing.T) {
	var btb = NewBitbucket()
	btb.Token.Scopes = "pullrequest"
	err := btb.CheckPermissions()
	if err == nil {
		t.Error(fmt.Errorf("CheckPermission didn't work correctly"))
	}

}

func TestGetWebHookEventsOK(t *testing.T) {

	defer gock.Off()

	var btb = NewBitbucket()

	gock.New(btb.WebHooks.EventListEndPoint).
		Get("/").
		Reply(200).
		JSON(map[string]interface{}{
			"pagelen": 30,
			"values": []HookEvent{{
				Category:    "Repository",
				Label:       "Push",
				Description: "Whenever a repository push occurs",
				Event:       "repo:push",
			},
				{
					Category:    "Repository",
					Label:       "Fork",
					Description: "Whenever a repository fork occurs",
					Event:       "repo:fork",
				}},
			"page": 1,
			"size": 18,
		})

	err := btb.GetWebHookEvents()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetWebHookEventsNOTOK(t *testing.T) {

	defer gock.Off()

	var btb = NewBitbucket()

	gock.New(btb.WebHooks.EventListEndPoint).
		Get("/").
		Reply(400)

	err := btb.GetWebHookEvents()
	if err == nil {
		t.Errorf("GetWebHookEvents didn't throw wrong response error")
	}
}
