package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

type CodeReq struct {
	Code string `json:"code"`
}

type AppTemplate struct {
	AppID         string
	AppURL        string
	AppName       string
	PEM           string
	WebhookSecret string
	Response      string
}

type AppResult struct {
	ID            int    `json:"id"`
	PEM           string `json:"pem"`
	URL           string `json:"html_url"`
	Name          string `json:"name"`
	WebhookSecret string `json:"webhook_secret"`
}

func MakeHandler(inputMap map[string]string, resCh chan AppResult) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		if r.URL.Path == "/" || r.URL.Path == "" {

			var outBuffer bytes.Buffer
			tmpl, err := template.ParseFiles(path.Join("./pkg/github", "index.html"))
			err = tmpl.Execute(&outBuffer, &inputMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write(outBuffer.Bytes())
			return
		}

		if r.URL.Path == "/callback" {
			code := r.URL.Query().Get("code")

			req, _ := http.NewRequest(http.MethodPost,
				fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code), nil)

			req.Header.Add("Accept", "application/vnd.github.fury-preview+json")
			res, err := http.DefaultClient.Do(req)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if res.Body != nil {
				defer res.Body.Close()
				result, _ := ioutil.ReadAll(res.Body)

				appRes := AppResult{}

				err := json.Unmarshal(result, &appRes)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("<html>Thank you for creating your GitHub App: %s</html>", appRes.Name)))

				resCh <- appRes
				close(resCh)

			}
		}
	}
}
