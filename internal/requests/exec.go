package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/client"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/utils"
	"gorm.io/datatypes"
	"moul.io/http2curl"
)

type Result struct {
	Err      error
	Response any
	Loading  bool
}

func (m Model) Exec() tea.Cmd {
	return func() tea.Msg {
		// workspace := common.CurrWorkspace
		request := common.CurrRequest

		url := utils.ReplaceByEnv(request.Endpoint)
		res := client.Request(url, strings.ToUpper(request.Method))

		rawbody, _ := request.Body.MarshalJSON()
		bodystring := utils.ReplaceByEnv(string(rawbody))

		var body any
		if err := json.Unmarshal([]byte(bodystring), &body); err != nil {
			panic(err)
		}

		if _, ok := body.(map[string]any); ok {
			res.Body([]byte(bodystring))
		} else {
			res.Body(nil)
		}

		rawheaders, _ := request.Headers.MarshalJSON()
		headersstring := utils.ReplaceByEnv(string(rawheaders))

		headers := []map[string]string{}
		json.Unmarshal([]byte(headersstring), &headers)
		for _, item := range headers {
			for k, v := range item {
				res.Header(k, v)
			}
		}

		rawparams, _ := request.QueryParams.MarshalJSON()
		paramsstring := utils.ReplaceByEnv(string(rawparams))

		params := []map[string]string{}
		json.Unmarshal([]byte(paramsstring), &params)
		for _, item := range params {
			for k, v := range item {
				res.Params(k, v)
			}
		}

		data, err := res.Exec()
		if err != nil {
			return Result{
				Err:     err,
				Loading: false,
			}
		}

		var response interface{}
		readbody, _ := ioutil.ReadAll(data.Body)
		json.Unmarshal(readbody, &response)

		result := repositories.Response{
			RequestId: request.ID,
			Url:       url,
			Method:    request.Method,
			Status:    data.Status,
			Params:    datatypes.NewJSONType(params),
			Headers:   datatypes.NewJSONType(headers),
			Response:  datatypes.NewJSONType(response),
			Body:      datatypes.NewJSONType(body),
			Curl:      Curl(data.Request),
		}

		repositories.NewResponse().Create(&result)

		return Result{
			Response: result,
		}
	}
}

func Curl(req *http.Request) string {
	command, _ := http2curl.GetCurlCommand(req)
	return command.String()
}
