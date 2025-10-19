package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"trade-api/models"
	"trade-api/utils"
)

func GetTradeSummary(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startYear := c.QueryInt("start_year", 0)
		endYear := c.QueryInt("end_year", 0)

		if startYear == 0 || endYear == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "start_year and end_year are required")
		}

		if startYear > endYear {
			return fiber.NewError(fiber.StatusBadRequest, "start_year must be less than or equal to end_year")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		query := `
			SELECT 
				year,
				COALESCE(import_value, 0) as import_value,
				COALESCE(export_value, 0) as export_value,
				COALESCE(reexport_value, 0) as reexport_value,
				COALESCE(tradebalance_value, 0) as tradebalance_value,
				COALESCE(import_value, 0) + COALESCE(export_value, 0) + COALESCE(reexport_value, 0) as total_trade_value
			FROM fact_yearly_summary
			WHERE year BETWEEN $1 AND $2
			ORDER BY year
		`

		rows, err := db.Query(ctx, query, startYear, endYear)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to query trade summary")
		}
		defer rows.Close()

		summaries := []models.TradeSummary{}
		for rows.Next() {
			var s models.TradeSummary
			if err := rows.Scan(&s.Year, &s.ImportValue, &s.ExportValue, &s.ReExportValue,
				&s.TradeBalanceValue, &s.TotalTradeValue); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan summary")
			}
			summaries = append(summaries, s)
		}

		return c.JSON(summaries)
	}
}

func GetTradeBalance(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startYear := c.QueryInt("start_year", 0)
		endYear := c.QueryInt("end_year", 0)

		if startYear == 0 || endYear == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "start_year and end_year are required")
		}

		if startYear > endYear {
			return fiber.NewError(fiber.StatusBadRequest, "start_year must be less than or equal to end_year")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		query := `
			SELECT 
				COALESCE(SUM(import_value), 0) as total_import,
				COALESCE(SUM(export_value), 0) as total_export,
				COALESCE(SUM(reexport_value), 0) as total_reexport,
				COALESCE(SUM(export_value), 0) + COALESCE(SUM(reexport_value), 0) - COALESCE(SUM(import_value), 0) as trade_balance
			FROM fact_yearly_summary
			WHERE year BETWEEN $1 AND $2
		`

		var balance models.TradeBalance
		balance.StartYear = startYear
		balance.EndYear = endYear

		err := db.QueryRow(ctx, query, startYear, endYear).Scan(
			&balance.TotalImport,
			&balance.TotalExport,
			&balance.TotalReExport,
			&balance.TradeBalance,
		)
		if err != nil {
			log.Printf("Trade balance query error: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to query trade balance: %v", err))
		}

		return c.JSON(balance)
	}
}

func AggregateTradeData(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.AggregateRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate request
		if err := validateAggregateRequest(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// Set defaults
		if req.Pagination.Page == 0 {
			req.Pagination.Page = 1
		}
		if req.Pagination.Limit == 0 {
			req.Pagination.Limit = 25
		}
		if req.Pagination.Limit > 1000 {
			req.Pagination.Limit = 1000
		}
		if req.Sorting.SortBy == "" {
			req.Sorting.SortBy = "total_value"
		}
		if req.Sorting.SortOrder == "" {
			req.Sorting.SortOrder = "desc"
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Build query
		query, countQuery, args := utils.BuildAggregateQuery(&req)

		log.Printf("Count Query: %s", countQuery)
		log.Printf("Args: %+v", args)

		// Get total count
		var totalCount int64
		if err := db.QueryRow(ctx, countQuery, args...).Scan(&totalCount); err != nil {
			log.Printf("Count query error: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get total count: %v", err))
		}

		// Get data
		offset := (req.Pagination.Page - 1) * req.Pagination.Limit
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, req.Pagination.Limit, offset)

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to execute query: "+err.Error())
		}
		defer rows.Close()

		results := []models.AggregateResult{}
		for rows.Next() {
			result := models.AggregateResult{}
			scanTargets := utils.BuildScanTargets(&req, &result)

			if err := rows.Scan(scanTargets...); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan result: "+err.Error())
			}
			results = append(results, result)
		}

		totalPages := int(totalCount) / req.Pagination.Limit
		if int(totalCount)%req.Pagination.Limit > 0 {
			totalPages++
		}

		response := models.PaginatedResponse{
			Data: results,
			Pagination: models.PaginationMeta{
				CurrentPage: req.Pagination.Page,
				PageSize:    req.Pagination.Limit,
				TotalCount:  totalCount,
				TotalPages:  totalPages,
			},
		}

		return c.JSON(response)
	}
}

func validateAggregateRequest(req *models.AggregateRequest) error {
	if req.DateRange.StartYear == 0 || req.DateRange.EndYear == 0 {
		return fmt.Errorf("date_range.start_year and date_range.end_year are required")
	}
	if req.DateRange.StartYear > req.DateRange.EndYear {
		return fmt.Errorf("start_year must be less than or equal to end_year")
	}
	if len(req.GroupBy) == 0 {
		return fmt.Errorf("group_by is required and must contain at least one field")
	}

	validGroupBy := map[string]bool{
		"year": true, "product": true, "country": true, "port": true, "trade_type": true,
	}
	for _, g := range req.GroupBy {
		if !validGroupBy[g] {
			return fmt.Errorf("invalid group_by field: %s. Valid options: year, product, country, port, trade_type", g)
		}
	}

	validTradeTypes := map[string]bool{"Import": true, "Export": true, "Re-Export": true}
	for _, tt := range req.TradeTypes {
		if !validTradeTypes[tt] {
			return fmt.Errorf("invalid trade_type: %s. Valid options: Import, Export, Re-Export", tt)
		}
	}

	validSortBy := map[string]bool{
		"total_value": true, "year": true, "product_desc_en": true,
		"country_name_en": true, "port_name_en": true, "trade_type": true,
	}
	if !validSortBy[req.Sorting.SortBy] {
		return fmt.Errorf("invalid sort_by field: %s", req.Sorting.SortBy)
	}

	if req.Sorting.SortOrder != "asc" && req.Sorting.SortOrder != "desc" {
		return fmt.Errorf("invalid sort_order: %s. Valid options: asc, desc", req.Sorting.SortOrder)
	}

	return nil
}
