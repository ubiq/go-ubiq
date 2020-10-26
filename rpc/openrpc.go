package rpc

import (
	"encoding"
	"encoding/json"
	"errors"
	openrpc "github.com/octanolabs/g0penrpc"
	"github.com/ubiq/go-ubiq/common/hexutil"
	"github.com/ubiq/go-ubiq/log"
	"reflect"
)

func makeOpenRpcSpecV1(server *Server) (*openrpc.DocumentSpec1, error) {

	methods := make([]*openrpc.Method, 0)

	p, _ := openrpc.NewPointer("#/components/schemas")
	reg, err := openrpc.NewSchemaRegistry(p)
	if err != nil {
		return nil, errors.New("could not create schema registry")
	}

	bnh := reflect.TypeOf((*BlockNumberOrHash)(nil)).Elem()
	reg.AddTypeException(bnh)
	hxb := reflect.TypeOf((*hexutil.Big)(nil)).Elem()
	reg.AddTypeException(hxb)

	//  subscribeSyncStatus accepts a chan and is not available publicly so how did it show up here

	log.Warn("openrpc: discovering services")

	for serviceName, service := range server.services.services {
		if serviceName != "rpc" {
			for methodName, callback := range service.callbacks {
				method, err := methodObjectFromCallback(reg, service.name, methodName, callback)
				if err != nil {
					log.Error("openrpc: error making content", "service", serviceName, "method", methodName, "err", err)
					continue
				}
				methods = append(methods, method)
			}
		}
	}

	info := &openrpc.Info{
		Title:   "go-ubiq JSON-RPC",
		Version: "1.0",
	}

	doc := openrpc.NewDocument(methods, info)

	doc.Components.Schemas = reg

	return doc, nil
}

func methodObjectFromCallback(reg *openrpc.SchemaRegistry, service, methodName string, cb *callback) (*openrpc.Method, error) {

	argTypes := cb.argTypes

	method := &openrpc.Method{
		Name:   service + serviceMethodSeparator + methodName,
		Params: nil,
		Result: nil,
	}
	params := make([]*openrpc.ContentDescriptor, len(cb.argTypes))

	//For argtypes we need to check if custom types implement un-marshaling (text or json)

	for idx, t := range argTypes {

		var (
			paramPtr  openrpc.Pointer
			paramName string
			err       error
		)

		marshalable := checkMarshalable(t)

		paramPtr, paramName, err = reg.RegisterType(t, marshalable)
		if err != nil {
			return nil, errors.New("error handling type: " + t.String() + " : " + err.Error())
		}

		// If param is pointer it can be omitted
		isReq := t.Kind() != reflect.Ptr

		params[idx] = &openrpc.ContentDescriptor{
			Name:     paramName,
			Schema:   paramPtr,
			Required: isReq,
		}

	}
	method.Params = params

	var (
		returnPtr  openrpc.Pointer
		returnName string
		err        error
	)

	returnType := cb.fn.Type().Out(0)

	unmarshalable := checkUnmarshalable(returnType)

	returnPtr, returnName, err = reg.RegisterType(returnType, unmarshalable)
	if err != nil {
		return nil, errors.New("error handling type: " + returnType.String() + " : " + err.Error())
	}

	method.Result = &openrpc.ContentDescriptor{
		Name:     returnName,
		Required: returnType.Kind() != reflect.Ptr,
		Schema:   returnPtr,
	}

	return method, err

}

func checkUnmarshalable(t reflect.Type) (ok bool) {

	if ok = t.Kind() == reflect.String; ok {
		return
	}

	if ok = t.Kind() == reflect.Map; ok {
		return checkUnmarshalable(t.Key()) && checkUnmarshalable(t.Elem())
	}

	txt := reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

	if ok = t.Implements(txt); ok {
		return
	}

	jsn := reflect.TypeOf(new(json.Unmarshaler)).Elem()

	if ok = t.Implements(jsn); ok {
		return
	}

	if t.Kind() != reflect.Ptr {
		p := reflect.PtrTo(t)
		return checkUnmarshalable(p)
	}

	return
}

func checkMarshalable(t reflect.Type) (ok bool) {

	if ok = t.Kind() == reflect.String; ok {
		return
	}

	if ok = t.Kind() == reflect.Map; ok {
		return checkMarshalable(t.Key()) && checkMarshalable(t.Elem())
	}

	txt := reflect.TypeOf(new(encoding.TextMarshaler)).Elem()

	if ok = t.Implements(txt); ok {
		return
	}

	jsn := reflect.TypeOf(new(json.Marshaler)).Elem()

	if ok = t.Implements(jsn); ok {
		return
	}

	if t.Kind() != reflect.Ptr {
		p := reflect.PtrTo(t)
		return checkMarshalable(p)
	}

	return
}
