package utils

import (
	"fmt"
	"strings"

	"trade-api/models"
)

func BuildAggregateQuery(req *models.AggregateRequest) (string, string, []interface{}) {
	args := []interface{}{}
	argCount := 0

	// Determine which fact table to use
	needsProduct := contains(req.GroupBy, "product") || len(req.Filters.ProductIDs) > 0
	needsCountry := contains(req.GroupBy, "country") || len(req.Filters.CountryIDs) > 0

	if needsProduct && needsCountry {
		panic("Cannot query both products and countries in the same request due to data structure")
	}

	var factTable string
	var dimensionJoins []string
	var selectFields []string
	var groupByFields []string

	// Build based on fact table
	if needsProduct {
		factTable = "fact_trade_by_product_port"

		if contains(req.GroupBy, "product") {
			dimensionJoins = append(dimensionJoins, "JOIN dim_product p ON f.product_id = p.product_id")
			selectFields = append(selectFields, "p.product_id", "p.product_desc_en", "p.product_desc_ar")
			groupByFields = append(groupByFields, "p.product_id", "p.product_desc_en", "p.product_desc_ar")
		}
	} else if needsCountry {
		factTable = "fact_trade_by_country_port"

		if contains(req.GroupBy, "country") {
			dimensionJoins = append(dimensionJoins, "JOIN dim_country c ON f.country_id = c.country_id")
			selectFields = append(selectFields, "c.country_id", "c.country_name_en", "c.country_name_ar")
			groupByFields = append(groupByFields, "c.country_id", "c.country_name_en", "c.country_name_ar")
		}
	} else {
		// Default to product table if no specific dimension requested
		factTable = "fact_trade_by_product_port"
	}

	// Add port joins if needed
	if contains(req.GroupBy, "port") || len(req.Filters.PortTypes) > 0 || len(req.Filters.PortIDs) > 0 {
		dimensionJoins = append(dimensionJoins, "JOIN dim_port dp ON f.port_id = dp.port_id")
		if contains(req.GroupBy, "port") {
			selectFields = append(selectFields, "dp.port_id", "dp.port_name_en", "dp.port_name_ar")
			groupByFields = append(groupByFields, "dp.port_id", "dp.port_name_en", "dp.port_name_ar")
		}
	}

	// Add year grouping
	if contains(req.GroupBy, "year") {
		selectFields = append(selectFields, "f.year")
		groupByFields = append(groupByFields, "f.year")
	}

	// Add trade_type grouping
	if contains(req.GroupBy, "trade_type") {
		selectFields = append(selectFields, "f.trade_type")
		groupByFields = append(groupByFields, "f.trade_type")
	}

	// Always add total_value
	selectFields = append(selectFields, "SUM(f.value) as total_value")

	// Build WHERE clause
	whereClauses := []string{}

	// Date range (requires 2 parameters)
	argCount++
	whereClauses = append(whereClauses, fmt.Sprintf("f.year BETWEEN $%d AND $%d", argCount, argCount+1))
	args = append(args, req.DateRange.StartYear, req.DateRange.EndYear)
	argCount++ // Increment for the second parameter

	// Trade types
	if len(req.TradeTypes) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("f.trade_type = ANY($%d)", argCount))
		args = append(args, req.TradeTypes)
	}

	// Product filter
	if len(req.Filters.ProductIDs) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("f.product_id = ANY($%d)", argCount))
		args = append(args, req.Filters.ProductIDs)
	}

	// Country filter
	if len(req.Filters.CountryIDs) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("f.country_id = ANY($%d)", argCount))
		args = append(args, req.Filters.CountryIDs)
	}

	// Port filter
	if len(req.Filters.PortIDs) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("f.port_id = ANY($%d)", argCount))
		args = append(args, req.Filters.PortIDs)
	}

	// Port type filter
	if len(req.Filters.PortTypes) > 0 {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("dp.port_type_en = ANY($%d)", argCount))
		args = append(args, req.Filters.PortTypes)
	}

	// Build final query
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s f
		%s
		WHERE %s
		GROUP BY %s
		ORDER BY %s %s
	`,
		strings.Join(selectFields, ", "),
		factTable,
		strings.Join(dimensionJoins, " "),
		strings.Join(whereClauses, " AND "),
		strings.Join(groupByFields, ", "),
		req.Sorting.SortBy,
		strings.ToUpper(req.Sorting.SortOrder),
	)

	// Build count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT %s
			FROM %s f
			%s
			WHERE %s
			GROUP BY %s
		) as subquery
	`,
		strings.Join(groupByFields, ", "),
		factTable,
		strings.Join(dimensionJoins, " "),
		strings.Join(whereClauses, " AND "),
		strings.Join(groupByFields, ", "),
	)

	return query, countQuery, args
}

func BuildScanTargets(req *models.AggregateRequest, result *models.AggregateResult) []interface{} {
	targets := []interface{}{}

	if contains(req.GroupBy, "product") {
		targets = append(targets, &result.ProductID, &result.ProductDescEN, &result.ProductDescAR)
	}
	if contains(req.GroupBy, "country") {
		targets = append(targets, &result.CountryID, &result.CountryNameEN, &result.CountryNameAR)
	}
	if contains(req.GroupBy, "port") {
		targets = append(targets, &result.PortID, &result.PortNameEN, &result.PortNameAR)
	}
	if contains(req.GroupBy, "year") {
		targets = append(targets, &result.Year)
	}
	if contains(req.GroupBy, "trade_type") {
		targets = append(targets, &result.TradeType)
	}

	targets = append(targets, &result.TotalValue)
	return targets
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
