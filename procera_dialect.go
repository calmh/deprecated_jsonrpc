package jsonrpc

/*
ProceraDialect implements the Procera modified JSON-RPC 2.0 dialect. It differs
from the standard JSON-RPC 2.0 by omitting the "jsonrpc":"2.0" field and adding
an optional "tags" field. For all requests and notifications, it is required
that the first parameter is a []string containing the tags. It may be nil.
*/
var ProceraDialect = Dialect{proceraRequest, proceraNotification}

func proceraRequest(id int, method string, params []interface{}) map[string]interface{} {
	tags := params[0]
	params = params[1:]
	return map[string]interface{}{
		"id":     id,
		"method": method,
		"params": params,
		"tags":   tags,
	}
}

func proceraNotification(method string, params []interface{}) map[string]interface{} {
	tags := params[0]
	params = params[1:]
	return map[string]interface{}{
		"method": method,
		"params": params,
		"tags":   tags,
	}
}
