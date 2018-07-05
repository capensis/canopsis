# Auth APIs

#### Url

  `GET` /auth

#### GET example

Parameters:
```javascript
{
    'login':         // username
    'password':      // password
    'json_response': // if True, the response is a JSON object
    'crypted':       // if True, password is encrypted using 'CRYPT' method
    'shadow':        // if True, password is encrypted using 'SHA1' method
    // if both crypted and shadow are False, password isn't encrypted
}
```

#### Response example

If `json_response` is True, the response has the following structure:

```javascript
{
    "authkey": "...",
    "mail": "...",
    "contact": {
        "name": "Administrator",
        "address": ""
    },
    "role": "admin",
    "crecord_name": "root"
}
```
