package Publishing

func formatGetOrdenResponse(res GetOrdersResponseModel) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"ordenes": res.Orders,
		},
	}
}
