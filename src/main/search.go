package main

import (
	"strings"
	"fmt"
)

type (
	// SearchMapping describes a search mapping for a list of entities
	SearchMapping struct {
		mapping map[string]FieldMapping
	}

	// FieldMapping describes how a field is mapped
	FieldMapping struct {
		DBField string
		SearchType int
	}
)

const (
	SearchTypeEqual = iota
)

// NewMapping creates a new search mapping
func NewMapping() SearchMapping {
	return SearchMapping { mapping: make(map[string]FieldMapping) }
}

// DefineFieldMapping defines a new field mapping on a search mapping
func (sm SearchMapping) DefineFieldMapping(field string, mapping FieldMapping) {
	sm.mapping[field] = mapping
}

// CreateQuery changes a search query based on the search mapping, and returns the changed query and its parameters
func (sm SearchMapping) CreateQuery(query string, search map[string]string, params... interface{}) (string, []interface{}) {
	conditions := []string{}

	for param, fieldMapping := range sm.mapping {
		if searchValue, ok := search[param]; ok && len(searchValue) > 0 {
			clauses := []string{}

			for _, val := range strings.Split(searchValue, "|") {
				switch(fieldMapping.SearchType) {
				case SearchTypeEqual:
					clauses = append(clauses, fmt.Sprintf("%s = $%d", fieldMapping.DBField, len(params)+1))
				}
				params = append(params, val)
			}

			conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(clauses, " OR ")))
		}
	}

	if len(conditions) > 0 {
		query = strings.Replace(query, "%MAPPING_CONDITIONS%", strings.Join(conditions, " AND "), -1)
	} else {
		query = strings.Replace(query, "%MAPPING_CONDITIONS%", "TRUE", -1)
	}

	return query, params
}