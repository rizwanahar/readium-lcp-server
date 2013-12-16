package index

import (
  "database/sql"
  "errors"
)

var NotFound = errors.New("Package not found")

type Index interface {
  Get(storageKey string) (Package, error)
  Add(p Package) error
}

type Package struct {
  StorageKey string
  EncryptionKey []byte
  Filename string
}

type dbIndex struct {
  db *sql.DB
  get *sql.Stmt
  add *sql.Stmt
}

func (i dbIndex) Get(storageKey string) (Package, error) {
  records, err := i.get.Query(storageKey)
  if records.Next() {
    var p Package
    err = records.Scan(&p.StorageKey, &p.EncryptionKey, &p.Filename)
    return p, err
  }

  return Package{}, NotFound
}

func (i dbIndex) Add(p Package) error {
  _, err := i.add.Exec(p.StorageKey, p.EncryptionKey, p.Filename)
  return err
}

func Open(where string) (i Index, err error) {
  db, err := sql.Open("sqlite3", where)
  if err != nil {
    return
  }
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS packages (storage_key varchar(255) PRIMARY KEY, encryption_key blob, filename varchar(255))")
    if err != nil {
      return
    }
  get, err := db.Prepare("SELECT * FROM packages WHERE storage_key = ? LIMIT 1")
    if err != nil {
      return
    }
  add, err := db.Prepare("INSERT INTO packages VALUES (?, ?, ?)")
    if err != nil {
      return
    }
  i = dbIndex{db, get, add}
  return
}
