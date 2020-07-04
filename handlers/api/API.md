
### Endpoints

TODO: config api is weird
GET /api/config/tileserver : get TileServerURL  
GET /api/config/nolabeltileserver : get NoLabelTileServerURL  

POST /api/maps : new Map from JSON  
GET  /api/maps/{id} : get Map by MapID  

POST /api/challenges : new Challenge from JSON (also inserts ChallengePlaces)  
GET /api/challenges/{id} : get Challenge by ChallengeID (also retrieves ChallengePlaces)  

POST /api/results : new ChallengeResult from JSON (Guesses will be empty)  
GET /api/results/{id} : get ChallengeResult by ChallengeResultID (also retrieves Guesses)  

POST /api/guesses : appends Guess from JSON to ChallengeResult.Guesses (if valid)  

### Responses

All request and response bodies contain either nothing, a JSON object containing only error: message, or a JSON object encoded directly from the corresponding type in `domain`.  
**The API does not guarantee that arrays arrive in order!**

```
Successful request responses:  
    GET:  
        200 OK  
        Body: Requested object from store as JSON  
    POST:  
        201 Created  
        Body: JSON after store insertion (including any generated IDs)  

Failed request responses:  
    GET:  
        400 Bad Request, if lacking ID  
        404 Not Found, if endpoint doesn't exist or ID not in store  
        401 and 403 may be used in the future  
        Body: {error: __description of error__}  
    POST:  
        404 Not Found, if endpoint doesn't exist  
        500 ISE, otherwise  
        401 and 403 may be used in the future  
        Body: {error: __description of error__}  
```