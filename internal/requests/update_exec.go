package requests

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/pkg/client"
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
		request := m.Requests.Current

		url := utils.ReplaceByOperator(request.Endpoint, m.Workspace.ID)
		res := client.Request(url, strings.ToUpper(request.Method))

		rawbody, _ := request.Body.MarshalJSON()
		bodystring := utils.ReplaceByOperator(string(rawbody), m.Workspace.ID)

		var body any
		if err := json.Unmarshal([]byte(bodystring), &body); err != nil {
			panic(err)
		}

		if _, ok := body.(map[string]any); ok {
			res.Body([]byte(bodystring))
		} else {
			res.Body(nil)
		}

		headers := utils.GetAllParentsHeaders(request.ParentID, request.Headers.Data())
		headers = utils.ProcessParamsOperators(headers, m.Workspace.ID)

		for _, item := range headers {
			for k, v := range item {
				res.Header(k, v)
			}
		}

		params := utils.GetAllParentsParams(request.ParentID, request.QueryParams.Data())
		params = utils.ProcessParamsOperators(params, m.Workspace.ID)

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
		readbody, _ := io.ReadAll(data.Body)
		json.Unmarshal(readbody, &response)

		result := offline.Response{
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

		offline.NewResponse().Create(&result)

		return Result{
			Response: result,
		}
	}
}

func Curl(req *http.Request) string {
	command, _ := http2curl.GetCurlCommand(req)
	return command.String()
}
