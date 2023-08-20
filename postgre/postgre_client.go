package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "usr_mehmet"
	password = "12345"
	dbname   = "productsdb"
)

type Product struct {
	ID                 int
	Title, Description string
	Price              float32
}

type client struct {
	db *sql.DB
}

type Client interface {
	InsertProduct(data Product)
	UpdateProduct(data Product)
	GetProducts() (error, []*Product)
	GetProductByID(id int) (error, *Product)
	Close()
}

func NewPostgresConnection() (Client, error) {
	var err error
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}

	return client{db: db}, err
}

func (c client) InsertProduct(data Product) {
	res, err := c.db.Exec("INSERT INTO deneme.products(title, description,price) values ($1,$2,$3)", data.Title, data.Description, data.Price)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("Eklenen kayıt sayısı(%d)", rowsAffected)
}

func (c client) UpdateProduct(data Product) {
	res, err := c.db.Exec("UPDATE deneme.products set title=$2, description=$3, price=$4 where id=$1", data.ID, data.Title, data.Description, data.Price)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("Eklenen kayıt sayısı(%d)", rowsAffected)
}

func (c client) GetProducts() (error, []*Product) {
	rows, err := c.db.Query("select * from deneme.products")
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Records Found!")
			return err, nil
		}
		log.Fatal(err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		prd := &Product{}
		err := rows.Scan(&prd.ID, &prd.Title, &prd.Description, &prd.Price)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, prd)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return err, nil
	}

	return nil, products
}

func (c client) GetProductByID(id int) (error, *Product) {
	product := &Product{}
	err := c.db.QueryRow("select * from deneme.products where id=$1", id).Scan(&product.ID, &product.Title, &product.Description, &product.Price)
	switch {
	case err == sql.ErrNoRows:
		return errors.New("no product with that ID"), nil
	case err != nil:
		return err, nil
	default:
		return nil, product
	}
}

func (c client) Close() {
	_ = c.db.Close()
}
