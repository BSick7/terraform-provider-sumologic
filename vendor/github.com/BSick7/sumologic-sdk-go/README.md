# sumologic-sdk-go

This sdk provides a go client for interacting with the SumoLogic API.

## Setup

In order to interact, a session and client must be created.

```
session := api.DefaultSession()
client := api.NewClient(session)
```

By default, a session is configured via env vars `SUMO_ACCESS_ID` and `SUMO_ACCESS_KEY`.
The address and credentials can be altered by configuring the session.

```
session.SetCredentials(accessID, accessKey)
session.SetAddress(sumoAddress)
```


The default API endpoint is `https://api.sumologic.com/api/v1`.
The following will alter the associated session by requesting the correct endpoint from SumoLogic.

```
client.Discover()
```
