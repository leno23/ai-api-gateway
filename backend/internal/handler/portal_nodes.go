package handler

import "encoding/json"

func DefaultPortalNodes() []PortalNode {
	return []PortalNode{
		{Name: "主站", URL: "http://localhost:8080", Region: "local"},
		{Name: "香港", URL: "https://hk.example.com", Region: "hk"},
		{Name: "美区", URL: "https://us.example.com", Region: "us"},
	}
}

func ParsePortalNodes(raw string) []PortalNode {
	if raw == "" {
		return DefaultPortalNodes()
	}
	var nodes []PortalNode
	if err := json.Unmarshal([]byte(raw), &nodes); err != nil || len(nodes) == 0 {
		return DefaultPortalNodes()
	}
	return nodes
}
