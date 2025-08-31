package parking_order

type IParkingOrder interface {
	ParkingOrderAdd(req *ParkingOrderAddReq) (*ParkingOrderAddRes, error)
	ParkingOrderList(req *ParkingOrderListReq) (*ParkingOrderListRes, error)
	ParkingOrderUpdate(req *ParkingOrderUpdateReq) (*ParkingOrderUpdateRes, error)
	ParkingOrderDelete(req *ParkingOrderDeleteReq) (*ParkingOrderDeleteRes, error)
}
