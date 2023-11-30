package product

import (
	"database/sql"
	"nojoke/auth"
	"nojoke/lib"
)

type ProductQuery struct {
	database *sql.DB
	logger   *lib.Logger
	admin    *auth.Admin
}

func (g *ProductQuery) HandleGuestGet(pagination *lib.Pagination) ([]Product, error) {
	g.logger.Info("Getting products ...")
	query := GetProductsWithoutCollectionQuery
	rows, err := g.database.Query(query)
	if err != nil {
		g.logger.Error("Error getting products: " + err.Error())
		return nil, err
	}
	productList := []Product{}
	for rows.Next() {
		product := Product{}
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Discount,
			&product.Rating,
			&product.Stock,
			&product.Brand,
			&product.Category_id,
			&product.Thumbnail,
			&product.Image,
			&product.Collection_id,
		)
		productList = append(productList, product)
	}
	return productList, nil
}

func (g *ProductQuery) HandleGet(collectionId int, pagination *lib.Pagination) ([]Product, error) {
	g.logger.Info("Getting products ...")
	query := GetProductsByCollectionQuery
	rows, err := g.database.Query(query, collectionId)
	if err != nil {
		g.logger.Error("Error getting products: " + err.Error())
		return nil, err
	}
	productList := []Product{}
	for rows.Next() {
		product := Product{}
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Discount,
			&product.Rating,
			&product.Stock,
			&product.Brand,
			&product.Category_id,
			&product.Thumbnail,
			&product.Image,
			&product.Collection_id,
		)
		productList = append(productList, product)
	}
	return productList, nil
}
