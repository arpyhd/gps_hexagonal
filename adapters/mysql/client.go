package client

import (
	"database/sql"
	"fmt"
	"gps_hexagonal/helpers/logger"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ConnectionPool struct {
	// the underlying connection pool
	Pool *sql.DB

	// the maximum number of connections in the pool
	maxConnections int

	// the current number of connections in the pool
	numConnections int

	// the mutex to synchronize access to the connection pool
	mutex *sync.Mutex
}

// NewConnectionPool creates a new ConnectionPool instance
func NewConnectionPool(dsn string, maxConnections int) (*ConnectionPool, error) {
	// create a new connection pool
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// set the maximum number of connections in the pool
	pool.SetMaxOpenConns(maxConnections)

	// create a new ConnectionPool instance
	p := &ConnectionPool{
		Pool:           pool,
		maxConnections: maxConnections,
		numConnections: 0,
		mutex:          &sync.Mutex{},
	}

	return p, nil
}

// GetConnection acquires a connection from the pool
func (p *ConnectionPool) GetConnection() (*sql.DB, error) {
	// acquire the mutex lock
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// check if the pool is full
	if p.numConnections == p.maxConnections {
		return nil, fmt.Errorf("connection pool is full")
	}

	// increment the number of connections in the pool
	p.numConnections++

	// return a connection from the underlying pool
	return p.Pool, nil
}

// ReleaseConnection releases a connection back to the pool
func (p *ConnectionPool) ReleaseConnection(conn *sql.DB) {
	// acquire the mutex lock
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// decrement the number of connections in the pool
	p.numConnections--
}

type GpsAdapter struct {
	Db *sql.DB
}

func NewGpsAdapter(Db *sql.DB) *GpsAdapter {
	g := &GpsAdapter{Db: Db}
	return g
}

func (gp *GpsAdapter) GetId(device_id string) (int, string, error) {
	var gps_id int
	var imei string
	row := gp.Db.QueryRow("select gps_id,imei FROM gps WHERE device_id = ?", device_id)
	switch err := row.Scan(&gps_id, &imei); err {
	case sql.ErrNoRows:
		stmt, err := gp.Db.Prepare("insert into gps (device_id) values(?)")
		if err != nil {
			return 0, "", err
		}
		res, err := stmt.Exec(device_id)
		if err != nil {
			return 0, "", err
		}
		gps_id, err := res.LastInsertId()
		if err != nil {
			return 0, "", err
		}
		return int(gps_id), "", nil
	case nil:
		return gps_id, imei, nil
	default:
		return 0, "", err
	}
}

func (gp *GpsAdapter) SetImei(gps_id int, imei string) error {
	stmt, err := gp.Db.Prepare("update  gps set imei=(?) where gps_id=(?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(imei, gps_id)
	if err != nil {
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		return err
	}
	logger.Log.Println("Sql Row affected gps_id: ", gps_id, " nr: ", a)
	return nil
}

func (gp *GpsAdapter) InsertEvent(gps_id int, message map[string]string) error {
	var statement string
	now := time.Now()
	statement = "Insert into ievents_" + fmt.Sprintf("%d%02d",
		now.Year(), now.Month()) + "  set "
	for key, value := range message {
		statement = statement + " " + key + " = '" + value + "' , "
	}
	statement = statement + " gps_id =" + strconv.Itoa(gps_id) + " , "
	statement = statement + " insert_date = NOW() "
	logger.Log.Println(statement)
	stmt, err := gp.Db.Prepare(statement)
	if err != nil {
		return err
	}
	res, err := stmt.Exec()
	if err != nil {
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		return err
	}
	logger.Log.Println("Sql Row affected gps_id: ", gps_id, " nr: ", a)
	return nil
}

func (gp *GpsAdapter) ReadCommand(gps_id int) (string, error) {
	var command string
	var params string
	var device_id string
	now := time.Now()
	row := gp.Db.QueryRow("select id,command,params FROM commands_"+fmt.Sprintf("%d%02d",
		now.Year(), now.Month())+" WHERE gps_id = ? AND send_date is NULL limit 1", gps_id)
	switch err := row.Scan(&device_id, &command, &params); err {
	case sql.ErrNoRows:
		logger.Log.Println("No command to send for gps_id: ", gps_id)
		return "", nil
	case nil:
		stmt, err := gp.Db.Prepare("update  gps_commands_" + fmt.Sprintf("%d%02d",
			now.Year(), now.Month()) + " set send_date=NOW()  where id=(?)")
		if err != nil {
			return "", err
		}
		res, err := stmt.Exec(gps_id)
		if err != nil {
			return "", err
		}
		a, err := res.RowsAffected()
		if err != nil {
			return "", err
		}
		logger.Log.Println("Sql Row affected:", a)
		response := command + params
		logger.Log.Println("Response in write for gps_id: ", gps_id, " response: ", response)
		return response, nil
	default:
		return "", err
	}
}
