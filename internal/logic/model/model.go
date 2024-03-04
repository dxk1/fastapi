package model

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/iimeta/fastapi/internal/consts"
	"github.com/iimeta/fastapi/internal/dao"
	"github.com/iimeta/fastapi/internal/errors"
	"github.com/iimeta/fastapi/internal/model"
	"github.com/iimeta/fastapi/internal/model/entity"
	"github.com/iimeta/fastapi/internal/service"
	"github.com/iimeta/fastapi/utility/logger"
	"github.com/iimeta/fastapi/utility/redis"
	"go.mongodb.org/mongo-driver/bson"
)

type sModel struct {
	modelCacheMap *gmap.StrAnyMap
}

func init() {
	service.RegisterModel(New())
}

func New() service.IModel {
	return &sModel{
		modelCacheMap: gmap.NewStrAnyMap(true),
	}
}

// 根据model获取模型信息
func (s *sModel) GetModel(ctx context.Context, m string) (*model.Model, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "GetModel time: %d", gtime.TimestampMilli()-now)
	}()

	result, err := dao.Model.FindOne(ctx, bson.M{"model": m, "status": 1})
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	return &model.Model{
		Id:                 result.Id,
		Corp:               result.Corp,
		Name:               result.Name,
		Model:              result.Model,
		Type:               result.Type,
		PromptRatio:        result.PromptRatio,
		CompletionRatio:    result.CompletionRatio,
		DataFormat:         result.DataFormat,
		IsEnableModelAgent: result.IsEnableModelAgent,
		ModelAgents:        result.ModelAgents,
		IsPublic:           result.IsPublic,
		Remark:             result.Remark,
		Status:             result.Status,
	}, nil
}

// 根据model和secretKey获取模型信息
func (s *sModel) GetModelBySecretKey(ctx context.Context, m, secretKey string) (md *model.Model, err error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "GetModelBySecretKey time: %d", gtime.TimestampMilli()-now)
	}()

	app, err := service.Common().GetCacheApp(ctx, service.Session().GetAppId(ctx))
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	if len(app.Models) == 0 {
		err = errors.ERR_MODEL_NOT_FOUND
		logger.Error(ctx, err)
		return nil, err
	}

	key, err := service.Common().GetCacheAppKey(g.RequestFromCtx(ctx).GetCtx(), secretKey)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	keyModelList := make([]*model.Model, 0)
	if len(key.Models) > 0 {

		models, err := s.GetCacheList(ctx, key.Models...)
		if err != nil || len(models) != len(key.Models) {

			if models, err = s.List(ctx, key.Models); err != nil {
				logger.Error(ctx, err)
				return nil, err
			}

			if err = s.SaveCacheList(ctx, models); err != nil {
				logger.Error(ctx, err)
				return nil, err
			}
		}

		for _, v := range models {
			if v.Name == m {
				keyModelList = append(keyModelList, &model.Model{
					Id:                 v.Id,
					Corp:               v.Corp,
					Name:               v.Name,
					Model:              v.Model,
					Type:               v.Type,
					PromptRatio:        v.PromptRatio,
					CompletionRatio:    v.CompletionRatio,
					DataFormat:         v.DataFormat,
					IsEnableModelAgent: v.IsEnableModelAgent,
					ModelAgents:        v.ModelAgents,
					IsPublic:           v.IsPublic,
					Remark:             v.Remark,
					Status:             v.Status,
				})
			}
		}

		for _, v := range models {
			if v.Model == m {
				keyModelList = append(keyModelList, &model.Model{
					Id:                 v.Id,
					Corp:               v.Corp,
					Name:               v.Name,
					Model:              v.Model,
					Type:               v.Type,
					PromptRatio:        v.PromptRatio,
					CompletionRatio:    v.CompletionRatio,
					DataFormat:         v.DataFormat,
					IsEnableModelAgent: v.IsEnableModelAgent,
					ModelAgents:        v.ModelAgents,
					IsPublic:           v.IsPublic,
					Remark:             v.Remark,
					Status:             v.Status,
				})
			}
		}

		if len(keyModelList) == 0 {
			err = errors.ERR_MODEL_NOT_FOUND
			logger.Error(ctx, err)
			return nil, err
		}
	}

	models, err := s.GetCacheList(ctx, app.Models...)
	if err != nil || len(models) != len(app.Models) {

		if models, err = s.List(ctx, app.Models); err != nil {
			logger.Error(ctx, err)
			return nil, err
		}

		if err = s.SaveCacheList(ctx, models); err != nil {
			logger.Error(ctx, err)
			return nil, err
		}
	}

	if len(models) == 0 {
		err = errors.ERR_MODEL_NOT_FOUND
		logger.Error(ctx, err)
		return nil, err
	}

	appModelList := make([]*model.Model, 0)
	for _, v := range models {
		if v.Name == m {
			appModelList = append(appModelList, &model.Model{
				Id:                 v.Id,
				Corp:               v.Corp,
				Name:               v.Name,
				Model:              v.Model,
				Type:               v.Type,
				PromptRatio:        v.PromptRatio,
				CompletionRatio:    v.CompletionRatio,
				DataFormat:         v.DataFormat,
				IsEnableModelAgent: v.IsEnableModelAgent,
				ModelAgents:        v.ModelAgents,
				IsPublic:           v.IsPublic,
				Remark:             v.Remark,
				Status:             v.Status,
			})
		}
	}

	for _, v := range models {
		if v.Model == m {
			appModelList = append(appModelList, &model.Model{
				Id:                 v.Id,
				Corp:               v.Corp,
				Name:               v.Name,
				Model:              v.Model,
				Type:               v.Type,
				PromptRatio:        v.PromptRatio,
				CompletionRatio:    v.CompletionRatio,
				DataFormat:         v.DataFormat,
				IsEnableModelAgent: v.IsEnableModelAgent,
				ModelAgents:        v.ModelAgents,
				IsPublic:           v.IsPublic,
				Remark:             v.Remark,
				Status:             v.Status,
			})
		}
	}

	if len(appModelList) == 0 {
		err = errors.ERR_MODEL_NOT_FOUND
		logger.Error(ctx, err)
		return nil, err
	}

	isModelDisabled := false
	for _, keyModel := range keyModelList {
		if keyModel.Name == m {

			if keyModel.Status == 2 {
				isModelDisabled = true
				continue
			}

			for _, appModel := range appModelList {
				if keyModel.Id == appModel.Id {
					return keyModel, nil
				}
			}
		}
	}

	for _, keyModel := range keyModelList {
		if keyModel.Model == m {

			if keyModel.Status == 2 {
				isModelDisabled = true
				continue
			}

			for _, appModel := range appModelList {
				if keyModel.Id == appModel.Id {
					return keyModel, nil
				}
			}
		}
	}

	if isModelDisabled {
		err = errors.ERR_MODEL_DISABLED
		logger.Error(ctx, err)
		return nil, err
	}

	return nil, errors.ERR_MODEL_NOT_FOUND
}

// 模型列表
func (s *sModel) List(ctx context.Context, ids []string) ([]*model.Model, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "sModel List time: %d", gtime.TimestampMilli()-now)
	}()

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"status": 1,
	}

	results, err := dao.Model.Find(ctx, filter, "-updated_at")
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	items := make([]*model.Model, 0)
	for _, result := range results {
		items = append(items, &model.Model{
			Id:                 result.Id,
			Corp:               result.Corp,
			Name:               result.Name,
			Model:              result.Model,
			Type:               result.Type,
			PromptRatio:        result.PromptRatio,
			CompletionRatio:    result.CompletionRatio,
			DataFormat:         result.DataFormat,
			IsEnableModelAgent: result.IsEnableModelAgent,
			ModelAgents:        result.ModelAgents,
			Remark:             result.Remark,
			Status:             result.Status,
		})
	}

	return items, nil
}

// 保存模型列表到缓存
func (s *sModel) SaveCacheList(ctx context.Context, models []*model.Model) error {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "sModel SaveCacheList time: %d", gtime.TimestampMilli()-now)
	}()

	fields := g.Map{}
	for _, model := range models {
		fields[model.Id] = model
		s.modelCacheMap.Set(model.Id, model)
	}

	if len(fields) > 0 {
		if _, err := redis.HSet(ctx, consts.API_MODELS_KEY, fields); err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	return nil
}

// 获取缓存中的模型列表
func (s *sModel) GetCacheList(ctx context.Context, ids ...string) ([]*model.Model, error) {

	now := gtime.TimestampMilli()
	defer func() {
		logger.Debugf(ctx, "sModel GetCacheList time: %d", gtime.TimestampMilli()-now)
	}()

	items := make([]*model.Model, 0)

	for _, id := range ids {
		if modelCacheValue := s.modelCacheMap.Get(id); modelCacheValue != nil {
			items = append(items, modelCacheValue.(*model.Model))
		}
	}

	if len(items) == len(ids) {
		return items, nil
	}

	reply, err := redis.HMGet(ctx, consts.API_MODELS_KEY, ids...)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	if reply == nil || len(reply) == 0 {
		if len(items) != 0 {
			return items, nil
		}
		return nil, errors.New("models is nil")
	}

	for _, str := range reply.Strings() {

		if str == "" {
			continue
		}

		result := new(model.Model)
		if err = gjson.Unmarshal([]byte(str), &result); err != nil {
			logger.Error(ctx, err)
			return nil, err
		}

		if s.modelCacheMap.Get(result.Id) != nil {
			continue
		}

		if result.Status == 1 {
			items = append(items, result)
			s.modelCacheMap.Set(result.Id, result)
		}
	}

	if len(items) == 0 {
		return nil, errors.New("models is nil")
	}

	return items, nil
}

// 更新缓存中的模型列表
func (s *sModel) UpdateCacheModel(ctx context.Context, m *entity.Model) {
	if err := s.SaveCacheList(ctx, []*model.Model{{
		Id:                 m.Id,
		Corp:               m.Corp,
		Name:               m.Name,
		Model:              m.Model,
		Type:               m.Type,
		PromptRatio:        m.PromptRatio,
		CompletionRatio:    m.CompletionRatio,
		DataFormat:         m.DataFormat,
		IsPublic:           m.IsPublic,
		IsEnableModelAgent: m.IsEnableModelAgent,
		ModelAgents:        m.ModelAgents,
		Remark:             m.Remark,
		Status:             m.Status,
		Creator:            m.Creator,
		Updater:            m.Updater,
	}}); err != nil {
		logger.Error(ctx, err)
	}
}

// 移除缓存中的模型列表
func (s *sModel) RemoveCacheModel(ctx context.Context, id string) {

	s.modelCacheMap.Remove(id)

	if _, err := redis.HDel(ctx, consts.API_MODELS_KEY, id); err != nil {
		logger.Error(ctx, err)
	}
}

// 变更订阅
func (s *sModel) Subscribe(ctx context.Context, msg string) error {

	message := new(model.SubMessage)
	if err := gjson.Unmarshal([]byte(msg), &message); err != nil {
		logger.Error(ctx, err)
		return err
	}
	logger.Infof(ctx, "sModel Subscribe: %s", gjson.MustEncodeString(message))

	var model *entity.Model
	switch message.Action {
	case consts.ACTION_UPDATE, consts.ACTION_STATUS:
		if err := gjson.Unmarshal(gjson.MustEncode(message.NewData), &model); err != nil {
			logger.Error(ctx, err)
			return err
		}
		s.UpdateCacheModel(ctx, model)
	case consts.ACTION_DELETE:
		if err := gjson.Unmarshal(gjson.MustEncode(message.OldData), &model); err != nil {
			logger.Error(ctx, err)
			return err
		}
		s.RemoveCacheModel(ctx, model.Id)
	}

	return nil
}
