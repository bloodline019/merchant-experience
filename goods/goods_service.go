package goods

import (
	"database/sql"
	"fmt"
	"github.com/thedatashed/xlsxreader"
	"merchant-experience/database"
	"strconv"
)

// Определение статистики для возврата клиенту
type Stats struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Deleted int `json:"deleted"`
	Errors  int `json:"errors"`
}

// Определение структуры товара
type Goods struct {
	OfferID   int
	Name      string
	Price     float64
	Quantity  int
	Available bool
	SellerID  int
}

// Парсинг XLSX файла по заданной структуре Goods
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
		avaliable, err4 := strconv.ParseBool(row.Cells[4].Value)
		seller_id, err5 := strconv.Atoi(sellerID)
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
			fmt.Printf("Error with convertation avaliable to bool: %v", err4)
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
			Available: avaliable,
			SellerID:  seller_id,
		}
		goods = append(goods, g)
	}
	return goods, nil
}

func SaveGoods(goods []Goods) (Stats, error) {
	// Connect to the database
	db, err := database.ConnectToDB()
	defer db.Close()
	if err != nil {
		return Stats{}, err
	}

	// Save or update the goods in the database
	stats := Stats{}
	entryCheckQuery := "SELECT * FROM products WHERE seller_id = $1 AND offer_id = $2"
	insertQuery := "INSERT INTO products (offer_id, name, price, quantity, available, seller_id) VALUES ($1, $2, $3, $4, $5, $6)"
	updateQuery := "UPDATE products SET name = $1, price = $2, quantity = $3, available = $4 WHERE seller_id = $5 AND offer_id = $6"
	deleteQuery := "DELETE FROM products WHERE seller_id = $1 AND offer_id = $2"
	for _, good := range goods {
		rows, _ := db.Query(entryCheckQuery, good.SellerID, good.OfferID)
		// If the entry not exists, insert it
		if rows.Next() == false {
			_, err := db.Exec(insertQuery, good.OfferID, good.Name, good.Price, good.Quantity, good.Available, good.SellerID)
			if err != nil {
				stats.Errors++
				continue
			}
			stats.Created++
		} else {
			if good.Available == false {
				_, err := db.Exec(deleteQuery, good.SellerID, good.OfferID)
				if err != nil {
					stats.Errors++
					continue
				}
				stats.Deleted++
			} else {
				_, err := db.Exec(updateQuery, good.Name, good.Price, good.Quantity, good.Available, good.SellerID, good.OfferID)
				if err != nil {
					stats.Errors++
					continue
				}
				stats.Updated++
			}
		}

	}
	return stats, nil
}

func GetGoods(offer_id int, seller_id int, substring string) []Goods {
	rows := doQuery(offer_id, seller_id, substring)
	defer rows.Close()
	var goods []Goods
	for rows.Next() {
		var good Goods
		err := rows.Scan(&good.OfferID, &good.Name, &good.Price, &good.Quantity, &good.Available, &good.SellerID)
		if err != nil {
			fmt.Printf("Error with scan row: %v", err)
		}
		goods = append(goods, good)
	}
	return goods
}

func doQuery(offer_id int, seller_id int, substring string) *sql.Rows {
	db, err := database.ConnectToDB()
	if err != nil {
		fmt.Printf("Error with connect to database: %v", err)
	}
	defer db.Close()

	getQuery := "SELECT * FROM products WHERE "
	var rows *sql.Rows
	var args []interface{}
	if offer_id != 0 {
		getQuery += "offer_id = $1"
		args = append(args, offer_id)
	}
	if seller_id != 0 {
		if len(args) == 0 {
			getQuery += "seller_id = $1"
		} else {
			getQuery += "AND seller_id = $2"
		}
		args = append(args, seller_id)
	}

	if substring != "" {
		if len(args) == 0 {
			getQuery += "name LIKE $1"
		} else {
			getQuery += "AND name LIKE $" + strconv.Itoa(len(args)+1)
		}
		args = append(args, substring+"%")
	}
	rows, _ = db.Query(getQuery, args...)
	return rows
}
