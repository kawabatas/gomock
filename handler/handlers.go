package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rakyll/statik/fs"

	"github.com/kawabatas/gomock/model"
	_ "github.com/kawabatas/gomock/statik"
)

const routesFile string = "routes.json"

func HandleStub(w http.ResponseWriter, r *http.Request) {
	statikFS, err := fs.New()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "statik/fs New error.")
		return
	}

	jsonBytes, err := readFile(statikFS, routesFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	routes := &model.Routes{}
	if err := json.Unmarshal(jsonBytes, &routes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "JSON Unmarshal error: %v", routesFile)
		return
	}

	fmt.Printf("\nreceive request: [%v] %v\n", r.Method, r.URL.Path)
	contents, err := findContens(routes.Routes, r.URL.Path, r.Method)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
		return
	}

	reqParam := requestParameter(r)
	fmt.Printf("request parameter: %v\n", reqParam)
	for _, handler := range contents.Handlers {
		if isMatchRequest(r, reqParam, handler) {
			fmt.Printf("handle pattern: %+v\n", handler.Content.Body)
			setContent(w, handler.Content, statikFS)
			return
		}
	}

	fmt.Printf("default pattern: %+v\n", contents.Default.Body)
	setContent(w, contents.Default, statikFS)
}

func readFile(fs http.FileSystem, fileName string) ([]byte, error) {
	file, err := fs.Open("/" + fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("statik/fs Open error: %v", fileName))
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	return bytes, nil
}

func findContens(routes []model.Route, path string, method string) (model.Contents, error) {
	var contents model.Contents
	for _, route := range routes {
		if route.Path == path && route.Method == strings.ToLower(method) {
			contents = model.Contents{
				Default:  route.Contents.Default,
				Handlers: route.Contents.Handlers,
			}
			break
		}
	}
	// empty check
	if contents.Default == (model.Content{}) {
		return model.Contents{}, errors.New(fmt.Sprintf("Not Found: [%v] %v", method, path))
	}
	return contents, nil
}

func requestParameter(r *http.Request) map[string]string {
	if r.Method == http.MethodPost {
		return postParameter(r.Body)
	} else if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodDelete {
		params := getParameter(r.URL.Query())
		return params
	}
	return map[string]string{}
}

func getParameter(query url.Values) map[string]string {
	var queryParams = map[string]string{}
	for k := range query {
		queryParams[k] = query.Get(k)
	}
	return queryParams
}

func postParameter(b io.ReadCloser) map[string]string {
	var postParamsBox map[string]interface{}
	if err := json.NewDecoder(b).Decode(&postParamsBox); err != nil {
		return map[string]string{}
	}
	var postParams = map[string]string{}
	for k, v := range postParamsBox {
		vs := fmt.Sprint(v)
		postParams[k] = vs
	}
	return postParams
}

func isMatchRequest(request *http.Request, reqParams map[string]string, handler model.Handler) bool {
	if len(handler.Header)+len(handler.Param) == 0 {
		return false
	}
	for k, v := range handler.Header {
		if !isMatchRegex(fmt.Sprintf("%v", v), request.Header.Get(k)) {
			return false
		}
	}
	for k, v := range handler.Param {
		if !isMatchRegex(fmt.Sprintf("%v", v), reqParams[k]) {
			return false
		}
	}
	return true
}

func isMatchRegex(regexPattern string, target string) bool {
	regex := regexp.MustCompile(regexPattern)
	return regex.MatchString(target)
}

func setContent(w http.ResponseWriter, content model.Content, fs http.FileSystem) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(content.Status)
	response, err := readFile(fs, content.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, string(response))
}
