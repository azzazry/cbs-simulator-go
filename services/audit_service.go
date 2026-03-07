package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
)

// LogAudit inserts an audit log entry
func LogAudit(log models.AuditLog) error {
	query := `INSERT INTO audit_logs (cif, action, resource, resource_id, ip_address, user_agent, 
	          request_method, request_path, request_body, response_status, details) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(query, log.CIF, log.Action, log.Resource, log.ResourceID,
		log.IPAddress, log.UserAgent, log.RequestMethod, log.RequestPath,
		log.RequestBody, log.ResponseStatus, log.Details)
	return err
}

// GetAuditLogs returns paginated audit logs with optional filters
func GetAuditLogs(cif, action string, page, pageSize int) ([]models.AuditLog, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Count total
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`
	args := []interface{}{}

	if cif != "" {
		countQuery += ` AND cif = ?`
		args = append(args, cif)
	}
	if action != "" {
		countQuery += ` AND action = ?`
		args = append(args, action)
	}

	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %v", err)
	}

	// Fetch logs
	dataQuery := `SELECT id, COALESCE(cif,''), action, COALESCE(resource,''), COALESCE(resource_id,''), 
	              COALESCE(ip_address,''), COALESCE(user_agent,''), COALESCE(request_method,''), 
	              COALESCE(request_path,''), COALESCE(request_body,''), COALESCE(response_status,0), 
	              COALESCE(details,''), created_at FROM audit_logs WHERE 1=1`

	if cif != "" {
		dataQuery += ` AND cif = ?`
	}
	if action != "" {
		dataQuery += ` AND action = ?`
	}
	dataQuery += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := database.DB.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %v", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.CIF, &log.Action, &log.Resource, &log.ResourceID,
			&log.IPAddress, &log.UserAgent, &log.RequestMethod, &log.RequestPath,
			&log.RequestBody, &log.ResponseStatus, &log.Details, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}
