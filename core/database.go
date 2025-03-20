package core

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase() *Database {
	return &Database{}
}


func (d *Database) InitializeDB() error {
	// Open a SQLite database (it creates the file if it doesn't exist)
	db, err := sql.Open("sqlite3", "./go_motion.db")
	if err != nil {
		return err
	}

	// Create the 'definitions' table to store serialized definitions
	definitionsTableCreation := `
	CREATE TABLE IF NOT EXISTS definitions (
		id TEXT PRIMARY KEY,
		xml TEXT
	);
	`
	_, err = db.Exec(definitionsTableCreation)
	if err != nil {
		return err
	}

	// Create the 'process_instances' table to store process instance data
	processInstancesTableCreation := `
	CREATE TABLE IF NOT EXISTS process_instances (
		id TEXT PRIMARY KEY,
		process_model_name TEXT,
		current_element TEXT,
		started_at DATETIME,
		finished_at DATETIME,
		state TEXT
	);
	`
	_, err = db.Exec(processInstancesTableCreation)
	if err != nil {
		return err
	}

	d.Db = db
	return nil
}

func (d *Database) SaveDefinitionToDB(definition *Definitions) error {
	// Serialize Definitions struct to XML
	data, err := xml.Marshal(definition)
	if err != nil {
		return err
	}

	fmt.Println()

	// Insert the serialized data into the definitions table
	_, err = d.Db.Exec("INSERT INTO definitions (id, xml) VALUES (?,?)", definition.ID, string(data))
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) LoadAllXMLs() ([]string, error) {
	var xmls []string

	// Query all XML data from the definitions table
	rows, err := d.Db.Query("SELECT xml FROM definitions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and append each XML string to the slice
	for rows.Next() {
		var xmlData string
		if err := rows.Scan(&xmlData); err != nil {
			return nil, err
		}
		xmls = append(xmls, xmlData)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return xmls, nil
}

// Save ProcessInstance to 'process_instances' table
func (d *Database) SaveProcessInstanceToDB(processInstance *ProcessInstance) error {
	fmt.Println("Saving ProcessInstance to DB...")
	// Insert the ProcessInstance data into the table
	fmt.Println(processInstance)
	_, err := d.Db.Exec(
		"INSERT INTO process_instances (id, process_model_name, current_element, started_at, state) VALUES (?,?,?,?,?)",
		processInstance.Id,
		processInstance.ProcessModel.Name,
		processInstance.CurrentElement,
		processInstance.StartedAt.Format(time.RFC3339),
		processInstance.State,
	)

	if err != nil {
		return err
	}

	return nil
}

// Update the state and current_element of an existing process instance
func (d *Database) PersistProcessInstance(processInstance *ProcessInstance) error {
	fmt.Printf("Updating ProcessInstance: %s", processInstance.Id)
	_, err := d.Db.Exec(
		"UPDATE process_instances SET state = ?, current_element = ? WHERE id = ?",
		processInstance.State,
		processInstance.CurrentElement,
		processInstance.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) FinishProcessInstance(processInstance *ProcessInstance) error {
	fmt.Printf("Finish ProcessInstance: %s", processInstance.Id)
	_, err := d.Db.Exec(
		"UPDATE process_instances SET state = ?, current_element = ?, finished_at = ? WHERE id = ?",
		processInstance.State,
		processInstance.CurrentElement,
		processInstance.FinishedAt.Format(time.RFC3339),
		processInstance.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

