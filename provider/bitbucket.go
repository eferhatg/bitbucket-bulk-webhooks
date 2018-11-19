package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	funk "github.com/thoas/go-funk"
)

//Bitbucket container struct
type Bitbucket struct {
	Token         Token
	WebHookEvents WebHookEvents
	Username      string
	Repositories  []Repository
	Auth          Auth
	WebHooks      WebHooks
	RepoEndPoint  string
}

//Auth keeps auth arguments
type Auth struct {
	EndPoint string
	Args     map[string][]string
	Key      string
	Secret   string
}

//WebHooks keeps WebHooks arguments
type WebHooks struct {
	EventListEndPoint string
	WebHookEndPoint   string
}

type ExistingWebHooks struct {
	urls []string `json:"values"`
}

//Repositories keeps repositories arguments
type Repositories struct {
	Previous     string       `json:"previous"`
	Pagelen      int64        `json:"pagelen"`
	Repositories []Repository `json:"values"`
	Size         int64        `json:"size"`
	Page         int64        `json:"page"`
	EndPoint     string
}

//Repository keeps repository arguments
type Repository struct {
	Scm         string `json:"scm"`
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	Fullname    string `json:"full_name"`
	IsPrivate   string `json:"true"`
	Name        string `json:"name"`
}

//Token is definition of oauth2 argument
type Token struct {
	AccessToken  string `json:"access_token"`
	Scopes       string `json:"scopes"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

//WebHookEvents keeps events
type WebHookEvents struct {
	Events []HookEvent
}

//HookEvent is definition of hookevents
type HookEvent struct {
	Scope string
	Event string
}

//WebhookPayload is payload of webhook
type WebhookPayload struct {
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"`
}

//NewBitbucket is constructor of Bitbucket
func NewBitbucket() *Bitbucket {
	var btb = &Bitbucket{}
	btb.Auth.Key = os.Getenv("BITBUCKET_KEY")
	btb.Auth.Secret = os.Getenv("BITBUCKET_SECRET")
	btb.Username = os.Getenv("BITBUCKET_USERNAME")

	eventList := strings.Split(os.Getenv("WEBHOOK_EVENTS"), ",")
	btb.WebHookEvents.Events = []HookEvent{}
	// Display all elements.
	for i := range eventList {
		hookEvent := new(HookEvent)
		hookEvent.Scope = strings.Split(eventList[i], ":")[0]
		hookEvent.Event = eventList[i]
		btb.WebHookEvents.Events = append(btb.WebHookEvents.Events, *hookEvent)
	}

	btb.WebHooks.WebHookEndPoint = "https://api.bitbucket.org/2.0/repositories/"
	btb.Auth.EndPoint = fmt.Sprintf("https://%s:%s@bitbucket.org/site/oauth2/access_token", os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"))
	btb.Auth.Args = url.Values{"grant_type": {"client_credentials"}}
	btb.WebHooks.EventListEndPoint = "https://api.bitbucket.org/2.0/hook_events/repository"
	btb.RepoEndPoint = "https://api.bitbucket.org/2.0/repositories/" + btb.Username + "?pagelen=100&page=%s&fields=-values.links,-values.project,-values.mainbranch,-values.owner,-values.website,-values.has_wiki,-values.language,-values.fork_policy,-values.created_on,-values.has_issues,-values.updated_on,-values.size,-values.slug,-values.type"
	return btb
}

//Authenticate auths user to bitbucket api
func (bt *Bitbucket) Authenticate() error {

	res, err := http.PostForm(bt.Auth.EndPoint, bt.Auth.Args)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Authenticate request failed. Status code: %d", res.StatusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, &bt.Token)
	if err != nil {
		return err
	}

	return nil

}

//CheckPermissions checking auth token permissions
func (bt *Bitbucket) CheckPermissions() error {
	scopes := strings.Fields(bt.Token.Scopes)
	requiredPermissions := []string{"account:write", "webhook", "repository:write", "project:write"}
	missingPermissions := []string{}

	for _, val := range requiredPermissions {
		if !funk.Contains(scopes, val) {
			missingPermissions = append(missingPermissions, val)
		}
	}

	if len(missingPermissions) > 0 {
		return fmt.Errorf(strings.Join(missingPermissions, " ") + " couldn't be found")
	}

	return nil
}

//CheckPermissions checking auth token permissions
func (bt *Bitbucket) CheckScopes() error {

	scopeArray := strings.Split(bt.Token.Scopes, " ")
	missingScopes := []string{}
	for _, val := range bt.WebHookEvents.Events {
		if val.Scope == "repo" && (!funk.Contains(scopeArray, "repository:write") && !funk.Contains(scopeArray, "repository:admin")) {
			if !funk.Contains(missingScopes, "repository:write or repository:admin") {
				missingScopes = append(missingScopes, "repository:write or repository:admin")
			}

		} else if val.Scope == "issue" && !funk.Contains(scopeArray, "issues:write") {
			if !funk.Contains(missingScopes, "issues:write") {
				missingScopes = append(missingScopes, "issues:write")
			}
		} else if val.Scope == "pullrequest" && !funk.Contains(scopeArray, "pullrequest:write") {
			if !funk.Contains(missingScopes, "pullrequest:write") {
				missingScopes = append(missingScopes, "pullrequest:write")
			}
		}
	}
	if len(missingScopes) > 0 {
		return fmt.Errorf("You need these scope permissions to add webhook spesific events: " + strings.Join(missingScopes, ", "))
	}
	return nil
}

// //GetWebHookEvents fetches usable web hook
// func (bt *Bitbucket) GetWebHookEvents() error {
// 	res, err := http.Get(bt.WebHooks.EventListEndPoint)

// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return fmt.Errorf("GetWebHookEvents request failed. Status code: %d", res.StatusCode)
// 	}

// 	bodyBytes, _ := ioutil.ReadAll(res.Body)
// 	err = json.Unmarshal(bodyBytes, &bt.WebHookEvents)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// GetRepositories gets the repositories
func (bt *Bitbucket) GetRepositories() {

	repos := getRepositoryPage(fmt.Sprintf(bt.RepoEndPoint, strconv.Itoa(1)), bt.Token.AccessToken)

	n := int(math.Ceil(float64((repos.Size / 100) + 1)))
	bt.Repositories = append(bt.Repositories, repos.Repositories...)

	for i := 2; i <= n; i++ {
		repos = getRepositoryPage(fmt.Sprintf(bt.RepoEndPoint, strconv.Itoa(1)), bt.Token.AccessToken)
		bt.Repositories = append(bt.Repositories, repos.Repositories...)
	}
}

func getRepositoryPage(url string, accessToken string) Repositories {

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+accessToken)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		repos := Repositories{}
		err = json.Unmarshal(body, &repos)
		if err != nil {
			panic(err)
		}
		return repos
	}
	return Repositories{}

}

func (bt *Bitbucket) CheckWebHookExists(repo Repository) (bool, error) {

	url := bt.GetWebHookUrl(repo) + "?fields=values.url"
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+bt.Token.AccessToken)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Couldn't fetch repository webhooks! statuscode: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	wh := ExistingWebHooks{}
	err = json.Unmarshal(body, &wh)
	if err != nil {
		return false, err
	}

	retVal := false
	for _, val := range wh.urls {
		if val == os.Getenv("WEBHOOK_URL") {
			retVal = true
			break
		}
	}

	return retVal, nil

}

func (bt *Bitbucket) AddWebHook(r Repository) {

	var jsonStr = bt.GetEventPayload(r)

	req, err := http.NewRequest("POST", bt.GetWebHookUrl(r), bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bt.Token.AccessToken)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		_, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

	}

}

func (bt *Bitbucket) GetEventPayload(r Repository) []byte {

	desc := strings.Replace(os.Getenv("WEBHOOK_DESCRIPTION"), "{REPO_NAME}", r.Name, -1)
	active, _ := strconv.ParseBool(os.Getenv("WEBHOOK_ACTIVE"))
	eventList := strings.Split(os.Getenv("WEBHOOK_EVENTS"), ",")

	whPl := WebhookPayload{
		Description: desc,
		URL:         os.Getenv("WEBHOOK_URL"),
		Active:      active,
		Events:      eventList,
	}
	b, _ := json.Marshal(whPl)

	return b

}

func (bt *Bitbucket) GetWebHookUrl(r Repository) string {

	return bt.WebHooks.WebHookEndPoint + r.Fullname + "/" + "hooks"

}
