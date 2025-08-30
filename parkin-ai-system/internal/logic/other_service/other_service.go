package other_service

import (
	"context"
	"parkin-ai-system/api/other_service"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/os/gtime"
)

type sOtherService struct{}

func (s *sOtherService) OtherServiceAdd(ctx context.Context, req *other_service.OtherServiceAddReq) (res *other_service.OtherServiceAddRes, err error) {
	serviceData := do.OthersService{
		LotId:           req.LotId,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
		CreatedAt:       gtime.Now(),
	}
	result, err := dao.OthersService.Ctx(ctx).Data(serviceData).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	res = &other_service.OtherServiceAddRes{ServiceId: result}
	return
}

func (s *sOtherService) OtherServiceUpdate(ctx context.Context, req *other_service.OtherServiceUpdateReq) (res *other_service.OtherServiceUpdateRes, err error) {
	data := do.OthersService{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
	}
	_, err = dao.OthersService.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, err
	}
	res = &other_service.OtherServiceUpdateRes{Success: true}
	return
}

func (s *sOtherService) OtherServiceDelete(ctx context.Context, req *other_service.OtherServiceDeleteReq) (res *other_service.OtherServiceDeleteRes, err error) {
	_, err = dao.OthersService.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	res = &other_service.OtherServiceDeleteRes{Success: true}
	return
}

func (s *sOtherService) OtherServiceDetail(ctx context.Context, req *other_service.OtherServiceDetailReq) (res *other_service.OtherServiceDetailRes, err error) {
	serviceRow, err := dao.OthersService.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if serviceRow.IsEmpty() {
		return nil, nil
	}
	info := &other_service.OtherServiceInfo{
		Id:              serviceRow["id"].Int64(),
		LotId:           serviceRow["lot_id"].Int64(),
		Name:            serviceRow["name"].String(),
		Description:     serviceRow["description"].String(),
		Price:           serviceRow["price"].Float64(),
		DurationMinutes: serviceRow["duration_minutes"].Int(),
		IsActive:        serviceRow["is_active"].Bool(),
		CreatedAt:       serviceRow["created_at"].GTime().Format("Y-m-d H:i:s"),
	}
	res = &other_service.OtherServiceDetailRes{Service: info}
	return
}

func (s *sOtherService) OtherServiceList(ctx context.Context, req *other_service.OtherServiceListReq) (res *other_service.OtherServiceListRes, err error) {
	m := dao.OthersService.Ctx(ctx)
	if req.LotId != 0 {
		m = m.Where("lot_id", req.LotId)
	}
	all, err := m.All()
	if err != nil {
		return nil, err
	}
	services := make([]other_service.OtherServiceInfo, 0, len(all))
	for _, r := range all {
		services = append(services, other_service.OtherServiceInfo{
			Id:              r["id"].Int64(),
			LotId:           r["lot_id"].Int64(),
			Name:            r["name"].String(),
			Description:     r["description"].String(),
			Price:           r["price"].Float64(),
			DurationMinutes: r["duration_minutes"].Int(),
			IsActive:        r["is_active"].Bool(),
			CreatedAt:       r["created_at"].GTime().Format("Y-m-d H:i:s"),
		})
	}
	res = &other_service.OtherServiceListRes{Services: services}
	return
}

func InitOtherService() {
	service.RegisterOtherService(&sOtherService{})
}

func init() {
	InitOtherService()
}
