package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinghzhu/kservice/pkg/api/v1/handler"
	"github.com/jinghzhu/kservice/pkg/logger"

	apitypes "github.com/jinghzhu/kservice/pkg/api/types"
)

const (
	// router
	routerV1 = "/api/v1"
	// api endpoint
	epPostPod      = "/pods"
	epGetPodStatus = "/pods/{key}/status"
	epGetPodLogs   = "/pods/{key}/logs"
	epGetPodInfo   = "/pods/{key}/info"
)

type wrappedHandlerFunc func(context.Context, *http.Request) ([]byte, int, error)

func GetJson(b []byte) (v interface{}, err error) {
	err = json.Unmarshal(b, &v)
	return v, err
}

type httpRequest struct {
	Headers []interface{}
	Body    interface{}
}

func parseRequest(ctx context.Context, r *http.Request) (*httpRequest, error) {
	req := &httpRequest{}
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		logger.ErrorFields("Fail to dump request", logger.Fields{
			apitypes.LogCtxID: ctx.Value(apitypes.LogCtxID),
			logger.ERROR:      err,
		})

		return req, err
	}
	str := string(b)
	z := strings.Split(str, "\n")
	dec := json.NewDecoder(strings.NewReader(z[len(z)-1]))
	if err := dec.Decode(&req.Body); err != nil {
		logger.ErrorFields("Fail to decode bytes for request ", logger.Fields{
			apitypes.LogCtxID: ctx.Value(apitypes.LogCtxID),
			logger.ERROR:      err,
		})

		return req, err
	}

	return req, nil
}

func handlerWrapper(fn wrappedHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := SetRequestContext(apitypes.ContextRoot)
		defer cancel()
		rd, err := parseRequest(ctx, r)
		body := rd.Body
		result, status, err := fn(ctx, r)
		reqID := fmt.Sprintf("%v", ctx.Value(apitypes.LogCtxID))
		if err != nil {
			logger.ErrorFields("Fail to warp handler", logger.Fields{
				apitypes.LogCtxID: reqID,
				"Request":         body,
				"Status":          status,
				"Error":           err,
			})
			http.Error(w, err.Error()+" "+reqID, status)

			return
		}
		data, err := GetJson(result)
		if err != nil {
			errMsg := fmt.Sprintf("%s:\n%b", err, result)
			status = 500
			logger.ErrorFields("Failed to Unmarshal the result",
				logger.Fields{
					apitypes.LogCtxID: reqID,
					"Request":         body,
					"Error":           err,
					"Result":          result,
					"Status":          status,
				})
			http.Error(w, errMsg, status)
		}
		json.NewEncoder(w).Encode(data)
		logger.InfoFields("Successfully handle with request", logger.Fields{
			apitypes.LogCtxID: reqID,
			"Request":         body,
			"Status":          status,
		})
	}
}

func DefaultRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	SetRouterV1(router)

	return router
}

// v1 api router
func SetRouterV1(r *mux.Router) {
	routerV1 := r.PathPrefix(routerV1).Subrouter()
	routerV1.HandleFunc(epPostPod, handlerWrapper(handler.CreatePod)).Methods(http.MethodPost)
	routerV1.HandleFunc(epGetPodStatus, handlerWrapper(handler.GetPodStatus)).Methods(http.MethodGet)
	routerV1.HandleFunc(epGetPodLogs, handlerWrapper(handler.GetPodLog)).Methods(http.MethodGet)
	routerV1.HandleFunc(epGetPodInfo, handlerWrapper(handler.GetPodInfo)).Methods(http.MethodGet)
}

func SetRequestContext(ctx context.Context) (context.Context, context.CancelFunc) {
	id, _ := uuid.NewRandom()
	ctx1, cancel := context.WithCancel(context.Background())
	return context.WithValue(ctx1, apitypes.LogCtxID, id), cancel
}
