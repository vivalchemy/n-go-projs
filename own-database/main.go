package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
	log     Logger
}

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: map[string]*sync.Mutex{},
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using %s (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Error: Missing collection - collection cannot be empty")
	}
	if resource == "" {
		return fmt.Errorf("Error: Missing resource - resource cannot be empty")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v any) error {
	if collection == "" {
		return fmt.Errorf("Error: Missing collection - collection cannot be empty")
	}
	if resource == "" {
		return fmt.Errorf("Error: Missing resource - resource cannot be empty")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := os.ReadFile(record)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Error: Missing collection - collection cannot be empty")
	}

	dir := filepath.Join(d.dir, collection)
	files, _ := os.ReadDir(dir)

	var records []string

	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		records = append(records, string(b))
	}
	return records, nil
}

// keep the resource empty to delete the collection itself
func (d *Driver) Delete(collection, resource string) error {
	if collection == "" {
		return fmt.Errorf("Error: Missing collection - collection cannot be empty")
	}

	path := filepath.Join(d.dir, collection, resource)

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	switch fi, err := stat(path); {
	case fi == nil, err != nil:
		return err
	case fi.Mode().IsDir():
		return os.RemoveAll(path)
	case fi.Mode().IsRegular():
		return os.RemoveAll(path + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		return os.Stat(path + ".json")
	}
	return
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

func main() {
	dir := "./data"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println(err)
	}

	employees := []User{
		{"John", "18", "123456789", "Google", Address{"New York", "NY", "USA", "123456"}},
		{"Jane", "25", "987654321", "Apple", Address{"San Francisco", "CA", "USA", "987654"}},
		{"Bob", "30", "876543210", "Microsoft", Address{"Los Angeles", "CA", "USA", "876543"}},
		{"Alice", "22", "234567891", "Amazon", Address{"Seattle", "WA", "USA", "234567"}},
		{"Charlie", "27", "345678912", "Facebook", Address{"Menlo Park", "CA", "USA", "345678"}},
		{"David", "35", "456789123", "Tesla", Address{"Palo Alto", "CA", "USA", "456789"}},
		{"Emma", "28", "567891234", "Netflix", Address{"Los Gatos", "CA", "USA", "567891"}},
		{"Frank", "31", "678912345", "Twitter", Address{"San Francisco", "CA", "USA", "678912"}},
		{"Grace", "24", "789123456", "Adobe", Address{"San Jose", "CA", "USA", "789123"}},
		{"Hank", "29", "891234567", "Intel", Address{"Santa Clara", "CA", "USA", "891234"}},
		{"Ivy", "26", "912345678", "Spotify", Address{"Stockholm", "Stockholm", "Sweden", "912345"}},
		{"Jack", "33", "102345678", "IBM", Address{"Armonk", "NY", "USA", "102345"}},
		{"Kelly", "21", "112345678", "Snapchat", Address{"Santa Monica", "CA", "USA", "112345"}},
		{"Leo", "40", "122345678", "Samsung", Address{"Seoul", "Seoul", "South Korea", "122345"}},
		{"Mia", "23", "132345678", "LG", Address{"Seoul", "Seoul", "South Korea", "132345"}},
		{"Noah", "36", "142345678", "Nvidia", Address{"Santa Clara", "CA", "USA", "142345"}},
		{"Olivia", "32", "152345678", "Oracle", Address{"Redwood City", "CA", "USA", "152345"}},
		{"Peter", "37", "162345678", "Cisco", Address{"San Jose", "CA", "USA", "162345"}},
		{"Quinn", "34", "172345678", "AMD", Address{"Austin", "TX", "USA", "172345"}},
		{"Rachel", "27", "182345678", "Dropbox", Address{"San Francisco", "CA", "USA", "182345"}},
	}

	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(records)

	allUsers := []User{}
	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println(err)
		}
		allUsers = append(allUsers, employeeFound)
	}
	fmt.Println(allUsers)

	// delete a record
	// if err := db.Delete("users", "John"); err != nil {
	// 	fmt.Println(err)
	// }
	//
	// delete a collection
	// if err := db.Delete("users", ""); err != nil {
	// 	fmt.Println(err)
	// }

}
