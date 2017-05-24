package publishedcollection

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetListReturnsInternalError(t *testing.T) {
	Convey("Get list returns internal error code when a statement is called", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT schedule_id").ExpectQuery().
			WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/report", nil)
		So(err, ShouldBeNil)
		api.GetList(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestGetListReturnsJsonMessage(t *testing.T) {
	Convey("Get list returns an array of published collection messages", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT schedule_id").ExpectQuery().WillReturnRows(makeCollectionRow())
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/report", nil)
		So(err, ShouldBeNil)
		api.GetList(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestGetCollectionReturnsInternalError(t *testing.T) {
	Convey("Get a collection returns a internal error from a statement", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT schedule_id").ExpectQuery().
			WillReturnRows(makeCollectionRow())
		mock.ExpectPrepare("SELECT uri").ExpectQuery().
			WillReturnError(fmt.Errorf("Testing internal server error"))
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "http://localhost:9090/report/test1", nil)
		So(err, ShouldBeNil)
		api.GetCollection(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestGetCollectionReturnsJsonMessage(t *testing.T) {
	Convey("Get collection returns a single published collection message", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT schedule_id").ExpectQuery().WillReturnRows(makeCollectionRow())
		mock.ExpectPrepare("SELECT uri").ExpectQuery().WithArgs(sqlmock.AnyArg()).WillReturnRows(makeResultRow())
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		r, err := http.NewRequest("GET", "http://localhost:9090/report/test1", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		So(err, ShouldBeNil)
		api.GetCollection(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestHealthIsOk(t *testing.T) {
	Convey("Health check returns OK", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		mock.ExpectPrepare("SELECT 1").ExpectQuery().WillReturnRows(makeCollectionRow())
		db.Begin()
		defer db.Close()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		r, err := http.NewRequest("GET", "http://localhost:9090/health", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		So(err, ShouldBeNil)
		api.Health(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestHealthIsBad(t *testing.T) {
	Convey("Health check returns a error code", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		prepareMockStmts(mock)
		db.Begin()
		api, err := NewAPI(db)
		So(err, ShouldBeNil)
		defer api.Close()
		r, err := http.NewRequest("GET", "http://localhost:9090/health", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		So(err, ShouldBeNil)
		db.Close()
		api.Health(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func makeCollectionRow() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"scheduleID", "collectionPath", "startTime", "completeTime"}).
		AddRow(1, "collection-name", 100, 1493891808831316672)
	return rows
}

func makeResultRow() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"uri", "fileCompleteTime"}).
		AddRow("/aboutus", 1493891808831316672)
	return rows
}

func prepareMockStmts(m sqlmock.Sqlmock) {
	m.ExpectBegin()
	m.MatchExpectationsInOrder(false)
	m.ExpectPrepare("SELECT schedule_id")
	m.ExpectPrepare("SELECT schedule_id")
	m.ExpectPrepare("SELECT uri")
	m.ExpectPrepare("SELECT 1")
}
