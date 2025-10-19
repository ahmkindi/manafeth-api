package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"trade-api/models"
)

func GetProducts(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		search := c.Query("search")
		limit := c.QueryInt("limit", 50)
		if limit > 500 {
			limit = 500
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := `
			SELECT product_id, product_desc_en, product_desc_ar
			FROM dim_product
			WHERE 1=1
		`
		args := []interface{}{}
		argCount := 0

		if search != "" {
			argCount++
			query += fmt.Sprintf(" AND (product_desc_en ILIKE $%d OR product_desc_ar ILIKE $%d)", argCount, argCount)
			args = append(args, "%"+search+"%")
		}

		argCount++
		query += fmt.Sprintf(" ORDER BY product_desc_en LIMIT $%d", argCount)
		args = append(args, limit)

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to query products")
		}
		defer rows.Close()

		products := []models.Product{}
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ProductID, &p.ProductDescEN, &p.ProductDescAR); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan product")
			}
			products = append(products, p)
		}

		return c.JSON(products)
	}
}

func GetCountries(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		search := c.Query("search")
		limit := c.QueryInt("limit", 50)
		if limit > 500 {
			limit = 500
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := `
			SELECT country_id, country_name_en, country_name_ar
			FROM dim_country
			WHERE 1=1
		`
		args := []interface{}{}
		argCount := 0

		if search != "" {
			argCount++
			query += fmt.Sprintf(" AND (country_name_en ILIKE $%d OR country_name_ar ILIKE $%d)", argCount, argCount)
			args = append(args, "%"+search+"%")
		}

		argCount++
		query += fmt.Sprintf(" ORDER BY country_name_en LIMIT $%d", argCount)
		args = append(args, limit)

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to query countries")
		}
		defer rows.Close()

		countries := []models.Country{}
		for rows.Next() {
			var country models.Country
			if err := rows.Scan(&country.CountryID, &country.CountryNameEN, &country.CountryNameAR); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan country")
			}
			countries = append(countries, country)
		}

		return c.JSON(countries)
	}
}

func GetPorts(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		search := c.Query("search")
		portType := c.Query("port_type")
		limit := c.QueryInt("limit", 50)
		if limit > 500 {
			limit = 500
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := `
			SELECT port_id, port_name_en, port_name_ar, port_type_en, port_type_ar, mode_id
			FROM dim_port
			WHERE 1=1
		`
		args := []interface{}{}
		argCount := 0

		if search != "" {
			argCount++
			query += fmt.Sprintf(" AND (port_name_en ILIKE $%d OR port_name_ar ILIKE $%d)", argCount, argCount)
			args = append(args, "%"+search+"%")
		}

		if portType != "" {
			argCount++
			query += fmt.Sprintf(" AND port_type_en = $%d", argCount)
			args = append(args, portType)
		}

		argCount++
		query += fmt.Sprintf(" ORDER BY port_name_en LIMIT $%d", argCount)
		args = append(args, limit)

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to query ports")
		}
		defer rows.Close()

		ports := []models.Port{}
		for rows.Next() {
			var port models.Port
			if err := rows.Scan(&port.PortID, &port.PortNameEN, &port.PortNameAR,
				&port.PortTypeEN, &port.PortTypeAR, &port.ModeID); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan port")
			}
			ports = append(ports, port)
		}

		return c.JSON(ports)
	}
}
