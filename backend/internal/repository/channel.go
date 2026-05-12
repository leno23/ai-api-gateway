package repository

import (
	"github.com/leno23/ai-api-gateway/internal/model"
)

// ListActiveChannelsForModel returns channels that serve the given model (empty Models = wildcard).
func (r *Repos) ListActiveChannelsForModel(modelName string) ([]model.Channel, error) {
	var rows []model.Channel
	q := r.DB.Where("status = ?", model.ChannelStatusActive).Order("priority DESC, id ASC")
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]model.Channel, 0, len(rows))
	for _, ch := range rows {
		if channelMatchesModel(ch, modelName) {
			out = append(out, ch)
		}
	}
	return out, nil
}

func channelMatchesModel(ch model.Channel, modelName string) bool {
	if len(ch.Models) == 0 {
		return true
	}
	for _, m := range ch.Models {
		if m == modelName {
			return true
		}
	}
	return false
}

// ListAllChannels for admin listing.
func (r *Repos) ListAllChannels() ([]model.Channel, error) {
	var rows []model.Channel
	if err := r.DB.Order("priority DESC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repos) GetChannel(id int64) (*model.Channel, error) {
	var ch model.Channel
	if err := r.DB.First(&ch, id).Error; err != nil {
		return nil, err
	}
	return &ch, nil
}

func (r *Repos) CreateChannel(ch *model.Channel) error {
	return r.DB.Create(ch).Error
}

func (r *Repos) SaveChannel(ch *model.Channel) error {
	return r.DB.Save(ch).Error
}

func (r *Repos) DeleteChannel(id int64) error {
	return r.DB.Delete(&model.Channel{}, id).Error
}
