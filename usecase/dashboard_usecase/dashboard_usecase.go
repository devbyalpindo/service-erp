package dashboard_usecase

import "erp-service/model/dto"

type DashboardUsecase interface {
	GetDashboard(dateFrom string, dateTo string) dto.Response
}
