package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func init() {
	psql_con := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sqlx.Open("postgres", psql_con)
	if err != nil {
		panic(err)
	}
	DB = db
	if err := enableUUIDExtension(db); err != nil {
		log.Fatal("Error enabling UUID extension:", err)
	}
}

func enableUUIDExtension(db *sqlx.DB) error {
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return err
}

func GetMotorData(searchKey string, keyList []string) ([]structs.MotorPolicy, error) {
	var policies []structs.MotorPolicy

	placeholders := make([]string, len(keyList))
	args := make([]interface{}, len(keyList))
	for i, key := range keyList {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = key
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s IN (%s)", "motorinsurance", searchKey, strings.Join(placeholders, ", "))

	err := DB.Select(&policies, query, args...)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	return policies, nil

}

func MotorBulkCreate(policies []structs.MotorPolicy) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO motorinsurance (
        transaction_type, rm_code, rm_name, child_id, booking_date, 
        insurer_name, insured_name, major_category, product, product_type, 
        policy_number, plan_type, premium, net_premium, od, tp, 
        commissionable_premium, registration_no, rto_code, state, rto_cluster, 
        city, insurer_biff, fuel_type, cpa, cc, gvw, ncb_type, 
        seating_capacity, vehicle_registration_year, discount_in_percentage, make, 
        model, ctg, idv, unique_id, sum_insured_val, vehicle_registration_date, utr, 
        utr_date, utr_amount, slot_payment_batch, paid_on_in, tentative_in_percentage, 
        tentative_in_amount, paid_on_out, out_percentage, out_amount, total_out_amount, 
        co_type, remarks, bu_head, manager, enricher_status, approver_status, 
        enricher_remark, approver_remark
    ) VALUES (
        :transaction_type, :rm_code, :rm_name, :child_id, :booking_date, 
        :insurer_name, :insured_name, :major_category, :product, :product_type, 
        :policy_number, :plan_type, :premium, :net_premium, :od, :tp, 
        :commissionable_premium, :registration_no, :rto_code, :state, :rto_cluster, 
        :city, :insurer_biff, :fuel_type, :cpa, :cc, :gvw, :ncb_type, 
        :seating_capacity, :vehicle_registration_year, :discount_in_percentage, :make, 
        :model, :ctg, :idv, :unique_id, :sum_insured_val, :vehicle_registration_date, :utr, 
        :utr_date, :utr_amount, :slot_payment_batch, :paid_on_in, :tentative_in_percentage, 
        :tentative_in_amount, :paid_on_out, :out_percentage, :out_amount, :total_out_amount, 
        :co_type, :remarks, :bu_head, :manager, :enricher_status, :approver_status, 
        :enricher_remark, :approver_remark
    )`

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, policy := range policies {
		_, err := stmt.Exec(policy)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func CreateMotorTable() {
	query := `
	CREATE TABLE IF NOT EXISTS motorinsurance (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		transaction_type VARCHAR(225) DEFAULT 'PRIMARY' CHECK (transaction_type IN ('PRIMARY', 'ADJUSTMENT', 'ENDORSEMENT')),
		rm_code VARCHAR(225),
		rm_name VARCHAR(225),
		child_id VARCHAR(225),
		booking_date DATE,
		insurer_name VARCHAR(225),
		insured_name VARCHAR(225),
		major_category VARCHAR(225) DEFAULT 'Motor',
		product VARCHAR(225),
		product_type VARCHAR(225),
		policy_number VARCHAR(225),
		plan_type VARCHAR(225),
		premium FLOAT,
		net_premium FLOAT,
		od FLOAT,
		tp FLOAT,
		commissionable_premium FLOAT,
		registration_no VARCHAR(225),
		rto_code VARCHAR(225),
		state VARCHAR(225),
		rto_cluster VARCHAR(225),
		city VARCHAR(225),
		insurer_biff VARCHAR(225),
		fuel_type VARCHAR(225),
		cpa VARCHAR(225),
		cc VARCHAR(225),
		gvw VARCHAR(225),
		ncb_type VARCHAR(225),
		seating_capacity INT,
		vehicle_registration_year INT,
		discount_in_percentage FLOAT,
		make VARCHAR(225),
		model VARCHAR(225),
		ctg VARCHAR(225),
		idv VARCHAR(225),
		unique_id VARCHAR(225),
		sum_insured_val VARCHAR(225),
		vehicle_registration_date DATE,
		utr VARCHAR(225),
		utr_date DATE,
		utr_amount INT,
		slot_payment_batch VARCHAR(225),
		paid_on_in VARCHAR(225),
		tentative_in_percentage DECIMAL(5, 2),
		tentative_in_amount FLOAT,
		paid_on_out VARCHAR(225),
		out_percentage DECIMAL(5, 2),
		out_amount FLOAT,
		total_out_amount FLOAT,
		co_type VARCHAR(225) CHECK (co_type IN ('ZIBPL', 'I CARE', 'Individual')),
		remarks VARCHAR(225),
		bu_head VARCHAR(225),
		manager VARCHAR(225),
		enricher_status VARCHAR(225) DEFAULT 'PENDING' CHECK (enricher_status IN ('PENDING', 'SENT_TO_NEXT_STAGE', 'REJECTED', 'PUSHBACK')),
		approver_status VARCHAR(225) DEFAULT 'PENDING' CHECK (approver_status IN ('PENDING', 'APPROVED', 'REJECTED', 'PUSHBACK')),
		enricher_remark VARCHAR(225),
		approver_remark VARCHAR(225)
	);
	`
	res, err := DB.Exec(query)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Create table result", res)
}
