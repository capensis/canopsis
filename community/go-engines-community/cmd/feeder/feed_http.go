package main

// TODO :
//    * allow the user to use environment variables to set :
//        - the authkey (CPS_AUTH_KEY)
//        - the host address (CPS_WEB_ADDR)
//        - the host port (CPS_WEB_PORT)
//   * failed at the first encountered error

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/rs/zerolog"
)

var ErrParseArg int = 1
var ErrFileNotFound int = 2
var ErrLogin int = 3
var ErrHTTP int = 4

var client http.Client

type HTTPSender struct {
	URL     string
	AuthKey string
	Files   string
	client  *http.Client
	Logger  zerolog.Logger
}

// Login try to authenticate to canopsis designated by the given url
// with the given authkey. The session cookie will be store in the
// cookiejar inside client. In case of error while trying to connect,
// this function, will exit wit ERR_LOGIN_ERROR
func (h *HTTPSender) Login() {
	cookies, _ := cookiejar.New(nil)
	client = http.Client{Jar: cookies}
	resp, err := client.Get(h.URL + "/autologin?authkey=" + h.AuthKey) //nolint:noctx
	if err != nil {
		h.Logger.Error().Err(err).Msg("Can not login")
		os.Exit(ErrLogin)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.Logger.Warn().Msgf("Can not login: Login return %d HTTP status code", resp.StatusCode)
		os.Exit(ErrLogin)
	}
	h.client = &client
}

// validJSON check if the given file is a valid JSON, return true if it is.
// False otherwise.
func validJSON(filename string) bool {
	var data []byte
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Can not load %s.\n", filename)
		os.Exit(ErrFileNotFound)
	}
	return json.Valid(data)
}

// sendEvent try to send an event to canopsis using /event route.
// eventFiles will contains every JSON files to be send, the authkey,
// the authentication key used to log into canopsis. An error, will
// be return if an error occurred during the upload of an event.
func (h *HTTPSender) sendEvent(eventFiles []string) {
	for _, filename := range eventFiles {
		fd, errFile := os.Open(filename)
		if errFile != nil {
			h.Logger.Error().Err(errFile).Msgf("Cannot load %s", filename)
			os.Exit(ErrFileNotFound)
		}

		if !validJSON(filename) {
			h.Logger.Warn().Msgf("%s is not a valid json, skipping.", filename)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, h.URL+"/api/v2/event", fd) //nolint:noctx
		if err != nil {
			h.Logger.Error().Err(err).Msg("new request error")
			os.Exit(ErrHTTP)
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := h.client.Do(req)
		if err != nil {
			h.Logger.Error().Err(err).Msg("")
		}

		data, err := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			h.Logger.Error().Err(err).Str("data", string(data)).Msg("")
			os.Exit(ErrHTTP)
		}

		defer resp.Body.Close()
		if err != nil {
			h.Logger.Error().Err(err).Str("data", string(data)).Msg("")
			os.Exit(ErrHTTP)
		}

		h.Logger.Info().Msg(string(data))
	}
}

func getJSONFromDir(filename string) []string {
	return []string{}
}

// getJSONFileNames return the list of files name to send to canopsis.
// If filename designate a directory, get JsonFilenames will return
// every json files inside the givent directory.
func getJSONFileNames(filename string) []string {
	var fileStat os.FileInfo
	var err error
	fileStat, err = os.Lstat(filename)
	if err != nil {
		log.Printf("Can not load %s\n", filename)
		os.Exit(ErrFileNotFound)
	}

	if fileStat.IsDir() {
		return getJSONFromDir(filename)
	}
	return []string{filename}
}

func (f *Feeder) modeSendEventHTTP() {
	if f.flags.AuthKey == "" {
		f.logger.Warn().Msg("An authentication key must be provided through CPS_AUTH_KEY environment variable or -authkey argument.")
		flag.PrintDefaults()
		os.Exit(ErrParseArg)
	}

	if f.flags.File == "" {
		f.logger.Warn().Msg("You must provide a JSON file or a directory (-file)")
		flag.PrintDefaults()
		os.Exit(ErrParseArg)
	}

	s := HTTPSender{
		AuthKey: f.flags.AuthKey,
		URL:     f.flags.PubHTTPURL,
		Files:   f.flags.File,
		Logger:  f.logger,
	}

	fileList := getJSONFileNames(f.flags.File)

	s.Login()
	s.sendEvent(fileList)
}
