package goods

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/thedatashed/xlsxreader"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"merchant-experience/database"
	"merchant-experience/models"
	"strconv"
)

// Stats Определение статистики для возврата клиенту
type Stats struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Deleted int `json:"deleted"`
	Errors  int `json:"errors"`
}

// Goods Определение структуры товара
type Goods struct {
	OfferID   int
	Name      string
	Price     float64
	Quantity  int
	Available bool
	SellerID  int
}

// ParseXlsx Парсинг XLSX файла по заданной структуре Goods
func ParseXlsx(file []byte, sellerID string) ([]Goods, error) {
	xl, err := xlsxreader.NewReader(file)
	if err != nil {
		return nil, err
	}

	var goods []Goods
	for row := range xl.ReadRows(xl.Sheets[0]) {
		id, err1 := strconv.ParseFloat(row.Cells[0].Value, 64)
		price, err2 := strconv.ParseFloat(row.Cells[2].Value, 64)
		quantity, err3 := strconv.ParseFloat(row.Cells[3].Value, 64)
		available, err4 := strconv.ParseBool(row.Cells[4].Value)
		sellerId, err5 := strconv.Atoi(sellerID)
		//Пропустим товары с некорректными данными
		if err1 != nil {
			fmt.Printf("Error with convertation id to int: %v", err1)
			continue
		}

		if err2 != nil {
			fmt.Printf("Error with convertation price to float: %v", err2)
			continue
		}

		if err3 != nil {
			fmt.Printf("Error with convertation sellerid to int: %v", err3)
			continue
		}
		if err4 != nil {
			fmt.Printf("Error with convertation available to bool: %v", err4)
			continue
		}
		if err5 != nil {
			fmt.Printf("Error with convertation sellerid to int: %v", err5)
			continue
		}

		g := Goods{
			OfferID:   int(id),
			Name:      row.Cells[1].Value,
			Price:     price,
			Quantity:  int(quantity),
			Available: available,
			SellerID:  sellerId,
		}
		goods = append(goods, g)
	}
	return goods, nil
}

func SaveGoods(goods []Goods) (Stats, error) {
	// Connect to the database
	db, err := database.ConnectToDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	if err != nil {
		return Stats{}, err
	}

	// Save or update the goods in the database
	stats := Stats{}
	for _, good := range goods {
		match, _ := models.Products(Where("seller_id = ? AND offer_id = ?", good.SellerID, good.OfferID)).Exists(context.Background(), db)
		product := models.Product{
			OfferID:   good.OfferID,
			Name:      good.Name,
			Price:     good.Price,
			Quantity:  good.Quantity,
			Available: good.Available,
			SellerID:  good.SellerID,
		}

		if !match && good.Available {
			err := product.Insert(context.Background(), db, boil.Infer())
			if err != nil {
				stats.Errors++
				continue
			}
			stats.Created++
		} else if good.Available == false && match {
			product := models.Product{
				OfferID:  good.OfferID,
				SellerID: good.SellerID}
			_, err := product.Delete(context.Background(), db)
			if err != nil {
				stats.Errors++
				continue
			}
			stats.Deleted++
		} else if match && good.Available {
			_, err := product.Update(context.Background(), db, boil.Infer())
			if err != nil {
				stats.Errors++
				continue
			}
			stats.Updated++
		} else {
			continue
		}
	}
	return stats, nil
}

func GetGoods(offerId int, sellerId int, substring string) []Goods {
	rows := doQuery(offerId, sellerId, substring)
	var goods []Goods
	for _, row := range rows {
		g := Goods{
			OfferID:   row.OfferID,
			Name:      row.Name,
			Price:     row.Price,
			Quantity:  row.Quantity,
			Available: row.Available,
			SellerID:  row.SellerID,
		}
		goods = append(goods, g)
	}
	return goods
}

func doQuery(offerId int, sellerId int, substring string) models.ProductSlice {
	db, err := database.ConnectToDB()
	if err != nil {
		fmt.Printf("Error with connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	boil.SetDB(db)

	var getQuery []QueryMod
	if offerId != 0 {
		getQuery = append(getQuery, Where("offer_id = ?", offerId))
	}
	if sellerId != 0 {
		getQuery = append(getQuery, Where("seller_id = ?", sellerId))
	}
	if substring != "" {
		getQuery = append(getQuery, Where("name LIKE ?", "%"+substring+"%"))
	}
	rows, _ := models.Products(getQuery...).All(context.Background(), db)
	return rows
}
