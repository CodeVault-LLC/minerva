# Contribution Guidelines

## Database Interactions

- **Use Transactions for Data Integrity**
  When adding an entity to the database, always use a transaction to ensure atomicity. This ensures that if any part of the process fails, the database remains unaffected.

  ```go
  tx := repository.db.Begin()
	if err := tx.Create(&scan).Error; err != nil {
		tx.Rollback()
		return entities.ScanModel{}, err
	}

	tx.Commit()
	return scan, nil
  ```
