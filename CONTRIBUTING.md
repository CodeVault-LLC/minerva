# Contribution Guidelines

## General Guidelines

- **Use our Logger for Logging**
  Always use our logger for logging. This ensures that logs are consistent and can be easily filtered.

  ```go
  import "github.com/codevault-llc/minerva/pkg/logger"

  logger.Log.Info("This is an info log")
  logger.Log.Error("This is an error log")
  ```

## Database Interactions

- **Use Transactions for Data Integrity**
  When adding an entity to the database, always use a transaction to ensure atomicity. This ensures that if any part of the process fails, the database remains unaffected.

  ```go
  tx, err := repository.db.Beginx()
  if err != nil {
  	return 0, err
  }

  query, values, err := database.StructToQuery(scan, "scans")
  if err != nil {
  	logger.Log.Error("Failed to generate query", zap.Error(err))
  	return 0, err
  }

  returnId, err := database.InsertStruct(tx, query, values)
  if err != nil {
  	logger.Log.Error("Failed to insert scan", zap.Error(err))
  	tx.Rollback()
  	return 0, err
  }

  err = tx.Commit()
  if err != nil {
  	logger.Log.Error("Failed to commit transaction", zap.Error(err))
  	return 0, err
  }

  return returnId, nil
  ```

- **Use Prepared Statements for Performance**
  When performing multiple queries with the same structure, use prepared statements to improve performance.

  ```go
  query := `SELECT * FROM scans WHERE id = ?`
  stmt, err := repository.db.Preparex(query)
  if err != nil {
  	logger.Log.Error("Failed to prepare statement", zap.Error(err))
  	return nil, err
  }

  var scan Scan
  err = stmt.Get(&scan, id)
  if err != nil {
  	logger.Log.Error("Failed to get scan", zap.Error(err))
  	return nil, err
  }

  return &scan, nil
  ```
