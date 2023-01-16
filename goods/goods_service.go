package goods

import (
	"context"
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
func ParseXlsx(file []byte, sellerID string) ([]Goods, error, int) {
	err_count := 0
	xl, err := xlsxreader.NewReader(file)
	if err != nil {
		err_count++
		return nil, err, err_count
	}

	var goods []Goods
	for row := range xl.ReadRows(xl.Sheets[0]) {
		var g Goods
		var err error

		OfferID, err := strconv.ParseFloat(row.Cells[0].Value, 64)
		if err != nil {
			err_count++
			fmt.Printf("Error with convertation id to int: %v", err)
			continue
		}

		g.Name = row.Cells[1].Value
		Price, err := strconv.ParseFloat(row.Cells[2].Value, 64)
		if err != nil {
			err_count++
			fmt.Printf("Error with convertation price to float: %v", err)
			continue
		}

		Quantity, err := strconv.ParseFloat(row.Cells[3].Value, 64)
		if err != nil {
			err_count++
			fmt.Printf("Error with convertation sellerid to int: %v", err)
			continue
		}

		Available, err := strconv.ParseBool(row.Cells[4].Value)
		if err != nil {
			err_count++
			fmt.Printf("Error with convertation available to bool: %v", err)
			continue
		}

		SellerID, err := strconv.Atoi(sellerID)
		if err != nil {
			err_count++
			fmt.Printf("Error with convertation sellerid to int: %v", err)
			continue
		}
		conds := []bool{OfferID <= 0, g.Name == "", Price <= 0, int(Quantity) <= 0, SellerID <= 0}
		skip := false
		for id, flag := range conds {
			if flag {
				switch id {
				case 0:
					err = fmt.Errorf("Incorrect format of OfferID")
				case 1:
					err = fmt.Errorf("Incorrect format of Name")
				case 2:
					err = fmt.Errorf("Incorrect format of Price")
				case 3:
					err = fmt.Errorf("Incorrect format of Quantity")
				case 4:
					err = fmt.Errorf("Incorrect format of SellerID")
				case 5:
					err = fmt.Errorf("Incorrect format of Available")
				}
				err_count++
				skip = true
				fmt.Printf("%v", err)

			}
		}
		if skip {
			continue
		}
		g.OfferID = int(OfferID)
		g.Name = row.Cells[1].Value
		g.Price = Price
		g.Quantity = int(Quantity)
		g.Available = Available
		g.SellerID = SellerID

		goods = append(goods, g)
	}
	return goods, nil, err_count
}

func SaveGoods(goods []Goods, err_count int) (Stats, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return Stats{
			Errors: err_count,
		}, err
	}
	defer db.Close()

	stats := Stats{
		Errors: err_count,
	}
	for _, good := range goods {
		match, err := models.Products(Where("seller_id = ? AND offer_id = ?",
			good.SellerID, good.OfferID)).Exists(context.Background(), db)
		if err != nil {
			stats.Errors++
			continue
		}

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
		}
	}
	return stats, nil
}

func GetGoods(offerId int, sellerId int, substring string) ([]Goods, error) {
	rows, err := doQuery(offerId, sellerId, substring)
	if err != nil {
		return nil, err
	}

	var goods []Goods
	for _, row := range rows {
		goods = append(goods, Goods{
			OfferID:   row.OfferID,
			Name:      row.Name,
			Price:     row.Price,
			Quantity:  row.Quantity,
			Available: row.Available,
			SellerID:  row.SellerID,
		})
	}
	return goods, nil
}

func doQuery(offerId int, sellerId int, substring string) (models.ProductSlice, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, fmt.Errorf("Error with connect to database: %v", err)
	}
	defer db.Close()
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
	rows, err := models.Products(getQuery...).All(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
