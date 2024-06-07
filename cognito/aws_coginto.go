package cognito

import (
    "fmt"
    "net/http"
    "strings"
    "io/ioutil"
    "encoding/json"
)

type CognitoResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
}

func GetJWTToken(userPoolID, clientID, clientSecret, username, password string) (string, error) {
    url := fmt.Sprintf("https://%s.auth.us-east-1.amazoncognito.com/oauth2/token", userPoolID)

    payload := strings.NewReader(fmt.Sprintf("grant_type=password&client_id=%s&username=%s&password=%s&scope=aws.cognito.signin.user.admin", clientID, username, password))
    req, err := http.NewRequest("POST", url, payload)
    if err != nil {
        return "", err
    }

    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.SetBasicAuth(clientID, clientSecret)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return "", err
    }

    var cognitoResponse CognitoResponse
    if err := json.Unmarshal(body, &cognitoResponse); err != nil {
        return "", err
    }

    return cognitoResponse.AccessToken, nil
}
