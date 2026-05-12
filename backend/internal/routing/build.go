package routing

import (
	"encoding/json"
	"math/rand/v2"
	"sort"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/repository"
)

// UpstreamTarget is one hop for OpenAI-compatible relay (DB channel or env fallback).
type UpstreamTarget struct {
	ChannelID *int64
	BaseURL   string
	APIKey    string
	RateLimit int
	ModelMap  map[string]string
}

// BuildUpstreamTargets orders active DB channels for the model, then appends env fallback (if configured).
func BuildUpstreamTargets(repos *repository.Repos, upstreamBase, upstreamKey, modelName string) ([]UpstreamTarget, error) {
	chs, err := repos.ListActiveChannelsForModel(modelName)
	if err != nil {
		return nil, err
	}
	ordered := orderChannelsForFailover(chs)
	out := make([]UpstreamTarget, 0, len(ordered)+1)
	for _, ch := range ordered {
		mm, _ := ParseModelMapping(ch.ModelMapping)
		id := ch.ID
		out = append(out, UpstreamTarget{
			ChannelID: &id,
			BaseURL:   ch.BaseURL,
			APIKey:    ch.APIKey,
			RateLimit: ch.RateLimit,
			ModelMap:  mm,
		})
	}
	if upstreamBase != "" {
		out = append(out, UpstreamTarget{
			ChannelID: nil,
			BaseURL:   upstreamBase,
			APIKey:    upstreamKey,
			RateLimit: 0,
			ModelMap:  nil,
		})
	}
	return out, nil
}

func orderChannelsForFailover(chs []model.Channel) []model.Channel {
	if len(chs) == 0 {
		return chs
	}
	sort.SliceStable(chs, func(i, j int) bool {
		if chs[i].Priority != chs[j].Priority {
			return chs[i].Priority > chs[j].Priority
		}
		return chs[i].ID < chs[j].ID
	})
	out := make([]model.Channel, 0, len(chs))
	for i := 0; i < len(chs); {
		j := i
		p := chs[i].Priority
		for j < len(chs) && chs[j].Priority == p {
			j++
		}
		block := append([]model.Channel(nil), chs[i:j]...)
		weightedShuffle(block)
		out = append(out, block...)
		i = j
	}
	return out
}

func weightInt(w int) int {
	if w <= 0 {
		return 1
	}
	return w
}

func weightedShuffle(block []model.Channel) {
	if len(block) <= 1 {
		return
	}
	list := append([]model.Channel(nil), block...)
	res := make([]model.Channel, 0, len(list))
	for len(list) > 0 {
		sum := 0
		for _, c := range list {
			sum += weightInt(c.Weight)
		}
		pick := rand.IntN(sum)
		acc := 0
		for idx, c := range list {
			acc += weightInt(c.Weight)
			if pick < acc {
				res = append(res, c)
				list = append(list[:idx], list[idx+1:]...)
				break
			}
		}
	}
	copy(block, res)
}

// ParseModelMapping decodes JSONB model_mapping (object of string -> string).
func ParseModelMapping(raw []byte) (map[string]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// RemapChatBody updates the JSON "model" field when mapping provides an upstream name.
func RemapChatBody(body []byte, clientModel string, mapping map[string]string) ([]byte, error) {
	upstream := clientModel
	if mapping != nil {
		if u, ok := mapping[clientModel]; ok && u != "" {
			upstream = u
		}
	}
	if upstream == clientModel {
		return body, nil
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	payload["model"] = upstream
	return json.Marshal(payload)
}
