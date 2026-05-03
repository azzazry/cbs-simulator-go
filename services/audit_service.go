package services

import (
	"cbs-simulator/database"
	"cbs-simulator/models"
	"fmt"
)

func LogAudit(log models.AuditLog) error {
	query := `INSERT INTO audit_logs (cif, action, resource, resource_id, ip_address, user_agent,
	          request_method, request_path, request_body, response_status, details)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := database.DB.Exec(query, log.CIF, log.Action, log.Resource, log.ResourceID,
		log.IPAddress, log.UserAgent, log.RequestMethod, log.RequestPath,
		log.RequestBody, log.ResponseStatus, log.Details)
	return err
}

func GetAuditLogs(cif, action string, page, pageSize int) ([]models.AuditLog, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	args := []interface{}{}
	argN := 1
	where := " WHERE 1=1"

	if cif != "" {
		where += fmt.Sprintf(` AND cif = $%d`, argN)
		args = append(args, cif)
		argN++
	}
	if action != "" {
		where += fmt.Sprintf(` AND action = $%d`, argN)
		args = append(args, action)
		argN++
	}

	var total int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM audit_logs`+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %v", err)
	}

	dataQuery := `SELECT id, COALESCE(cif,''), action, COALESCE(resource,''), COALESCE(resource_id,''),
	              COALESCE(ip_address,''), COALESCE(user_agent,''), COALESCE(request_method,''),
	              COALESCE(request_path,''), COALESCE(request_body,''), COALESCE(response_status,0),
	              COALESCE(details,''), created_at FROM audit_logs` + where +
		fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argN, argN+1)
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
