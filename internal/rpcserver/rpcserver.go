package rpcserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/pastelnetwork/go-pastel"
	"github.com/pastelnetwork/supernode/internal/common"
	"github.com/pastelnetwork/supernode/internal/restserver"
)

type RpcServer struct {
	pslNode *pastel.Client
	config  *restserver.Config
	funcs   map[string]interface{}
}

type RpcMethod struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      string   `json:"id"`
}

func New(psl *pastel.Client, cfg *restserver.Config) *RpcServer {
	return &RpcServer{pslNode: psl, config: cfg}
}

func (s *RpcServer) Start(ctx context.Context, app *common.Application) func() error {

	s.AddHandler("getinfo", s.Getinfo)

	rpcAddress := fmt.Sprintf("%s:%d", s.config.Hostname, s.config.Port)
	server := s.InitServer(rpcAddress)

	return app.CreateServer(ctx, "rpc_server",
		//initServer
		func(ctx context.Context) error {
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				return fmt.Errorf("error starting Rest server: %w", err)
			}
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			return server.Shutdown(ctx)
		})
}

func (rpcServer *RpcServer) AddHandler(handlerName string, handler func(method RpcMethod) ([]byte, error)) {
	if rpcServer.funcs == nil {
		rpcServer.funcs = make(map[string]interface{})
	}
	rpcServer.funcs[handlerName] = handler
}

func (rpcServer *RpcServer) InitServer(address string) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var rpcMethod RpcMethod
		err := json.NewDecoder(r.Body).Decode(&rpcMethod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if rpcServer.funcs == nil {
			err := fmt.Errorf("rpcMethod methods are not implemented  - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		m := reflect.ValueOf(rpcServer.funcs[strings.ToLower(rpcMethod.Method)])
		//m := reflect.ValueOf(&rpcMethod).MethodByName(strings.Title(strings.ToLower(rpcMethod.Method)))
		if m.Kind() == reflect.Invalid || m.IsNil() || m.IsZero() {
			err := fmt.Errorf("rpcMethod method not implemented  - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(rpcMethod)

		retValue := m.Call(in)
		if retValue == nil || len(retValue) != 2 ||
			retValue[0].Kind() != reflect.Slice ||
			retValue[1].Kind() != reflect.Interface {
			err := fmt.Errorf("invalid rpcMethod method implementation - %s", rpcMethod.Method)
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		if !retValue[1].IsNil() {
			err := retValue[1].Interface().(error)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		retBytes := retValue[0].Interface().([]byte)
		var buf bytes.Buffer
		json.Indent(&buf, retBytes, "", "  ")

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	})
	return &http.Server{Addr: address, Handler: mux}
}

func (s *RpcServer) Getinfo(method RpcMethod) ([]byte, error) {
	type info struct {
		Method    string `json:"method"`
		PSLNode   bool   `json:"psl_node"`
		RpcServer bool   `json:"rpc_server"`
	}

	i := info{method.Method, true, true}
	return json.Marshal(i)
}
