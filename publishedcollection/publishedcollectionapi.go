package publishedcollection

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONSDigital/go-ns/log"
	"github.com/gorilla/mux"
)

type API struct {
	CollectionStmt       *sql.Stmt
	SingleCollectionStmt *sql.Stmt
	CompleteFilesStmt    *sql.Stmt
	DB                   *sql.DB
}

func NewAPI(db *sql.DB) (*API, error) {
	collectionSQL := `SELECT schedule_id, collection_id, collection_path, start_time, complete_time
  from schedule WHERE complete_time IS NOT NULL ORDER BY complete_time DESC LIMIT 100`
	singleCollectionSQL := `SELECT schedule_id, collection_path, start_time, complete_time
  from schedule WHERE complete_time IS NOT NULL AND collection_id = $1`
	completeFilesSQL := "SELECT uri, complete_time FROM schedule_file WHERE schedule_id = $1"
	collectionStmt, err := createStmt(collectionSQL, db)
	if err != nil {
		return nil, err
	}
	singleCollectionStmt, err := createStmt(singleCollectionSQL, db)
	if err != nil {
		return nil, err
	}
	completeFilesStmt, err := createStmt(completeFilesSQL, db)
	if err != nil {
		return nil, err
	}

	return &API{CollectionStmt: collectionStmt,
		SingleCollectionStmt: singleCollectionStmt,
		CompleteFilesStmt:    completeFilesStmt, DB: db}, err
}

func (api *API) GetList(w http.ResponseWriter, r *http.Request) {
	var publishedCollections []PublishedCollection
	endPoint := "/report"
	collectionRows, err := api.CollectionStmt.Query()
	if err != nil {
		log.ErrorC("Failed run query", err, log.Data{"endpoint": endPoint})
		http.Error(w, "Failed to query collections", http.StatusInternalServerError)
		return
	}
	for collectionRows.Next() {
		var (
			collectionPath, collectionID        sql.NullString
			scheduleID, startTime, completeTime sql.NullInt64
		)
		collectionRows.Scan(&scheduleID, &collectionID, &collectionPath, &startTime, &completeTime)
		publishedCollection := PublishedCollection{CollectionID: collectionID.String,
			CollectionName: collectionPath.String, PublishDate: convertTime(startTime.Int64),
			PublishStartDate: convertTime(startTime.Int64), PublishEndDate: convertTime(completeTime.Int64)}
		publishedCollections = append(publishedCollections, publishedCollection)
	}
	data, _ := json.Marshal(publishedCollections)
	w.Write(data)
}

func (api *API) GetCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionID := vars["colllectionId"]
	endPoint := "/report/" + collectionID
	collectionRow := api.SingleCollectionStmt.QueryRow(collectionID)
	var (
		collectionPath                      sql.NullString
		scheduleID, startTime, completeTime sql.NullInt64
	)
	collectionRow.Scan(&scheduleID, &collectionPath, &startTime, &completeTime)
	results, err := api.getResults(scheduleID.Int64, startTime.Int64)
	if err != nil {
		log.ErrorC("Failed run query", err, log.Data{"endpoint": endPoint})
		http.Error(w, "Failed to query results", http.StatusInternalServerError)
		return
	}
	publishedCollection := PublishedCollection{CollectionName: collectionPath.String, PublishDate: convertTime(startTime.Int64),
		PublishStartDate: convertTime(startTime.Int64), PublishEndDate: convertTime(completeTime.Int64),
		Results: make([]PublishedItem, len(results))}
	copy(publishedCollection.Results, results)
	data, _ := json.Marshal(publishedCollection)
	w.Write(data)
}

func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	err := api.DB.Ping()
	if err != nil {
		http.Error(w, "Unable to access database", http.StatusInternalServerError)
	}
}

func (api *API) getResults(scheduleID int64, startTime int64) ([]PublishedItem, error) {
	var results []PublishedItem
	fileRows, err := api.CompleteFilesStmt.Query(scheduleID)
	if err != nil {
		return nil, err
	}
	for fileRows.Next() {
		var uri sql.NullString
		var fileCompleteTime sql.NullInt64
		fileRows.Scan(&uri, &fileCompleteTime)
		duration := (fileCompleteTime.Int64 - startTime) / 1000 / 1000
		results = append(results, PublishedItem{Duration: duration,
			Uri: uri.String, Size: 0})
	}
	return results, nil
}

func (api *API) Close() {
	api.CollectionStmt.Close()
	api.SingleCollectionStmt.Close()
	api.CompleteFilesStmt.Close()
}

func createStmt(sqlStatement string, db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare(sqlStatement)
}

func convertTime(epoch int64) string {
	return time.Unix(0, epoch).UTC().String()
}
