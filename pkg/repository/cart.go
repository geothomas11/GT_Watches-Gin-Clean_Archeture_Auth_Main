package repository

import (
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
	return &CartRepository{
		db: DB,
	}

}
func (cr *CartRepository) AddToCart(userID int, productId int, Quantity int, productprice float64) error {

	query := "INSERT INTO carts (user_id,quantity,total_price) VALUES(?,?,?,?)"

	if err := cr.db.Exec(query, userID, productId, Quantity, productprice).Error; err != nil {
		return err
	}
	return nil

}

func (cr *CartRepository) CheckProductAvailable(product_id int) (bool, error) {
	var count int
	querry := "SELECT COUNT(*) FROM products where id = ?"

	err := cr.db.Raw(querry, product_id).Scan(&count).Error
	if err != nil {
		return false, errors.New("product does not exist")
	}
	if count < 1 {
		return false, errors.New("product does not exist")
	}
	return true, nil
}

func (cr *CartRepository) CheckStock(product_id int) (int, error) {
	qurry := "SELECT stock FROM products where id = ?"
	var stock int
	err := cr.db.Raw(qurry, product_id).Scan(&stock).Error
	if err != nil {
		return 0, errors.New("error in getting stock")
	}
	return stock, nil
}

func (cr *CartRepository) QuantityOfProductInCart(userId int, productId int) (int, error) {
	var productQty int
	querry := "SELECT quantity FROM carts WHERE user_id = ? AND product_id = ?"
	err := cr.db.Raw(querry, userId, productId).Scan(&productQty).Error
	if err != nil {
		return 0, errors.New("error in getting quantity")
	}
	return productQty, nil
}

func (cr *CartRepository) GetTotalPriceFromCart(userID int) (float64, error) {
	var totalPrice float64
	err := cr.db.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&totalPrice).Error
	if err != nil {
		return 0.0, err
	}
	return totalPrice, nil

}
func (cr *CartRepository) UpdateCart(quantity int, price float64, userID int, product_id int) error {
	if err := cr.db.Exec(`UPDATE carts
	SET quantity = ?, total_price = ? 
	WHERE user_id = ? and product_id = ?`, quantity, price, product_id, userID).Error; err != nil {
		return err
	}

	return nil

}

func (cr *CartRepository) TotalPriceForProductInCart(userID int, productID int) (float64, error) {

	var totalPrice float64
	if err := cr.db.Raw("SELECT SUM(total_price) as total_price FROM carts  WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}

func (cr *CartRepository) DisplayCart(userID int) ([]models.Cart, error) {

	var count int
	if err := cr.db.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}

	var cartResponse []models.Cart
	qurry := `SELECT carts.user_id,users.name as name,carts.product_id,
				products.product_name as product_name,carts.quantity,carts.total_price 
				from carts INNER JOIN users ON carts.user_id = users.id 
				INNER JOIN products ON carts.product_id = products.id WHERE user_id = ?`

	if err := cr.db.Raw(qurry, userID).First(&cartResponse).Error; err != nil {
		return []models.Cart{}, err
	}
	return cartResponse, nil

}

func (cr *CartRepository) GetTotalPrice(userID int) (models.CartTotal, error) {

	var cartTotal models.CartTotal
	err := cr.db.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.TotalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = cr.db.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.FinalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = cr.db.Raw("SELECT name as user_name FROM users WHERE id = ?", userID).Scan(&cartTotal.UserName).Error
	if err != nil {
		return models.CartTotal{}, err
	}

	return cartTotal, nil

}
