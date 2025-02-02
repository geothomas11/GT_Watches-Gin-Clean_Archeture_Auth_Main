package repository

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}
func (ad *adminRepository) LoginHandler(adminDetails models.AdminLogin) (domain.Users, error) {
	var adminCompareDetails domain.Users
	if err := ad.DB.Raw("select * from users where email=?", adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.Users{}, err
	}
	return adminCompareDetails, nil

}
func (ad *adminRepository) GetUserByID(id int) (domain.Users, error) {
	var users domain.Users
	if err := ad.DB.Raw("select * from users where id=?", id).Scan(&users).Error; err != nil {
		return domain.Users{}, err
	}
	return users, nil

}

//	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)

func (ad *adminRepository) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {
	// pagination purpose -
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * 2
	var Get_Users []models.UserDetailsAtAdmin
	if err := ad.DB.Raw("SELECT id,name,email,phone,blocked FROM users limit ? offset ?", 3, offset).Scan(&Get_Users).Error; err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return Get_Users, nil

}
func (ad *adminRepository) UpdateBlockUserByID(user models.UpdateBlock) error {

	err := ad.DB.Exec("update users set blocked = ? where id = ?", user.Blocked, user.ID).Error
	if err != nil {
		return err
	}
	return nil

}
func (ad *adminRepository) IsUserExist(id int) (bool, error) {

	var count int
	err := ad.DB.Raw("select count(*) from users where id = ?", id).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil

}

// admin Dashboard
func (ar *adminRepository) DashboardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := ar.DB.Raw("select count (*) from users where is_admin ='false' ").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		err = errors.New("cannot get total users from db")
		return models.DashBoardUser{}, err
	}
	return userDetails, nil

}

func (ar *adminRepository) DashboardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := ar.DB.Raw("SELECT COUNT(*) FROM users WHERE is_admin = 'false' ").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		err = errors.New("cannot get products from db")
		return models.DashBoardProduct{}, err

	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM products WHERE stock <= 0").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		err = errors.New("cannot get stock from db")
		return models.DashBoardProduct{}, err
	}
	return productDetails, nil
}

func (ar *adminRepository) DashboardAmountDetails() (models.DashBoardAmount, error) {
	var amountDetails models.DashBoardAmount
	query := `SELECT coalesce(sum(final_price),0) FROM orders WHERE payment_status ='paid' `
	err := ar.DB.Raw(query).Scan(&amountDetails.CreditedAmount).Error
	if err != nil {
		err = errors.New("cannot get total amount from  db")
		return models.DashBoardAmount{}, err
	}
	query = `SELECT coalese(sum(final_price),0) FROM orders WHERE payment_status ='not_paid'
	 and
	  shipment_status = 'pending'
	  or 
	  shipment_status = 'processing'
	  or 
	  shipment_status = 'shipped'
	    `
	err = ar.DB.Raw(query).Scan(&amountDetails.PendingAmount).Error
	if err != nil {
		err = errors.New("cannot get pending amount from db")
		return models.DashBoardAmount{}, err
	}
	return amountDetails, nil
}

func (ar *adminRepository) DashboardOrderDetails() (models.DashBoardOrder, error) {
	var orderDetails models.DashBoardOrder
	err := ar.DB.Raw("SELECT count(*) FROM orders WHERE payment_status = 'paid' ").Scan(&orderDetails.CompletedOrder).Error
	if err != nil {
		err = errors.New("cannot get total order from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("SELECT COUNT(*) from orders WHERE shipment_status = 'pending or shipment_status = 'processing'").Scan(&orderDetails.PendingOrder).Error
	if err != nil {
		err = errors.New("cannoot get pending orders from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM orders WHERE shipment_status = 'cancelled' ").Scan(&orderDetails.CancelledOrder).Error
	if err != nil {
		err = errors.New("cannot get cancelled order from db")
		return models.DashBoardOrder{}, err

	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM orders ").Scan(&orderDetails.TotalOrderItem).Error
	if err != nil {
		err = errors.New("cannot get total order items from db")
		return models.DashBoardOrder{}, err
	}
	return orderDetails, nil

}
func (ar *adminRepository) DashboardTotalRevenueDetails() (models.DashBoardRevenue, error) {
	var revenueDetails models.DashBoardRevenue
	startTime := time.Now().AddDate(0, 0, 1)
	err := ar.DB.Raw("SELECT coalesce(sum(final_price),0) FROM orders WHERE payment_status ='paid' and created_at >=?", startTime).Scan(&revenueDetails.TodayRevenue).Error
	if err != nil {
		err = errors.New("cannot get today revenue from db")
		return models.DashBoardRevenue{}, err
	}
	startTime = time.Now().AddDate(0, -1, 1)
	err = ar.DB.Raw("SELECT COALESCE(sum(final_price),0) FROM orders WHERE payment_status = 'paid' and created_at >= ?", startTime).Scan(&revenueDetails.MonthRevenue).Error
	if err != nil {
		err = errors.New("cannot get month revenue from db")
		return models.DashBoardRevenue{}, err
	}
	startTime = time.Now().AddDate(-1, 1, 1)
	err = ar.DB.Raw("SELECT COALESCE(sum(final_price),0) FROM orders WHERE payment_status = 'paid' and created_at >= ?", startTime).Scan(&revenueDetails.YearRevenue).Error
	if err != nil {
		err = errors.New("cannot get year revenue from db")
		return models.DashBoardRevenue{}, err

	}
	return models.DashBoardRevenue{}, nil
}
