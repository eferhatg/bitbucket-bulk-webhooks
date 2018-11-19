# bitbucket-bulk-webhooks

> An attempt to adding same webhook to all repositories of an account

## Prerequisites

* You need a Bitbucket Oauth consumer
    * You can create from Account Settings -> OAuth -> Add Consumer
    * You need to give webhook permission and corresponding permissions to events
        * As an example: if you want to add webhook to issue:created event you have to have issue:write permission
* You need to fill .env file

## .env file details
```
BITBUCKET_KEY=[OAUTH_KEY]  
BITBUCKET_SECRET=[OAUTH_SECRET] 
BITBUCKET_USERNAME=[USERNAME OF REPOSITORIES' OWNER]   
WEBHOOK_URL=[WEBHOOK URL]
WEBHOOK_DESCRIPTION=Something is going on {REPO_NAME} 
WEBHOOK_ACTIVE=true
WEBHOOK_EVENTS=[WEBHOOK_EVENTS]

```

* Current **WEBHOOK_EVENTS** list;  
    * repo:push,repo:fork,repo:updated,repo:commit_comment_created,repo:commit_status_created,pullrequest:created,pullrequest:updated,pullrequest:approved,pullrequest:unapproved,pullrequest:fulfilled,pullrequest:rejected,pullrequest:comment_created,pullrequest:comment_updated,pullrequest:comment_deleted,issue:created,issue:updated,issue:comment_created
* You can make passive webhooks initially with **WEBHOOK_ACTIVE**
* You can use repository name with {REPO_NAME} slug in **WEBHOOK_DESCRIPTION**



## Usage

After filling .env file simply

```Shell
go run main.go

```

## License

MIT © [eferhatg](https://github.com/eferhatg)
