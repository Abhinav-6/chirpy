package database

import (
	"encoding/json"
	"errors"
	"strconv"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps     map[string]Chirp
	TotalChirp int
}

type Chirp struct {
	ID   int    `json:"Id"`
	Body string `json:"body"`
}

func NewDb() (*DB, error) {
	db := DB{"database/data.json", &sync.RWMutex{}}
	err := db.ensureDB()
	if err != nil {
		x, er := json.Marshal(DBStructure{make(map[string]Chirp), 0})
		if er != nil {
			panic(er)
		}
		err := os.WriteFile("database/data.json", x, 0644)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

	}
	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	d, err := db.loadDB()
	assert(err)
	d.TotalChirp++
	d.Chirps[strconv.Itoa(d.TotalChirp)] = Chirp{Body: body, ID: d.TotalChirp}
	err = db.writeDB(d)
	assert(err)
	return d.Chirps[strconv.Itoa(d.TotalChirp)], nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dat, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	var ch []Chirp
	for _, d := range dat.Chirps {
		ch = append(ch, d)
	}
	return ch, err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile("database/data.json")
	assert(err)
	var dat DBStructure
	err = json.Unmarshal(data, &dat)
	assert(err)
	return dat, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	d, err := json.Marshal(dbStructure)
	assert(err)
	err = os.WriteFile("database/data.json", d, os.ModeAppend)
	assert(err)
	return nil

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return os.ErrNotExist
	}
	return nil
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
