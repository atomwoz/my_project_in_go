package database

import "atomwoz.com/remitly_task/internal/models"

// FetchSwiftRecord retrieves a SWIFT record from the database
func FetchSwiftRecord(code string) (*models.SwiftModel, error) {
	var record models.SwiftModel
	err := DB.Table(DEFAULT_TABLE_NAME).
		Where("swift_code = ?", code).First(&record).Error

	if err != nil {
		return nil, err
	}
	return &record, nil
}

// DeleteSwiftRecord deletes a SWIFT record from the database
func DeleteSwiftRecord(record *models.SwiftModel) error {
	return DB.Table(DEFAULT_TABLE_NAME).Delete(record).Error
}

// InsertSwiftRecord inserts a SWIFT record into the database
func InsertSwiftRecord(record *models.SwiftModel) error {
	return DB.Table(DEFAULT_TABLE_NAME).Create(record).Error
}
