package jsonrpc

/*
StandardDialect implements the JSON-RPC 2.0 standard dialect as described in
http://www.jsonrpc.org/specification.
*/
var StandardDialect = Dialect{standardRequest, standardNotification}

func standardRequest(id int, method string, params []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":      id,
		"method":  method,
		"params":  params,
		"jsonrpc": "2.0",
	}
}

func standardNotification(method string, params []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"method":  method,
		"params":  params,
		"jsonrpc": "2.0",
	}
}
