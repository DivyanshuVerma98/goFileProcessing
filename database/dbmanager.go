package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"github.com/lib/pq"
)

var DB *sql.DB

func init() {
	psql_con := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", psql_con)
	if err != nil {
		panic(err)
	}
	DB = db
}

func GetPolicyNo(policies []string) {
	// Generate placeholders for the policy numbers
	placeholders := []string{}
	for policy := range policies {
		// Use single quotes around each policy number
		placeholders = append(placeholders, fmt.Sprintf("'%d'", policy))
	}
	query := fmt.Sprintf("SELECT policy_no FROM motor_policies WHERE policy_no IN (%s)", strings.Join(placeholders, ","))
	fmt.Println("QUERY LINE --->", query)
	row, err := DB.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer row.Close()
	//	for row.Next() {
	//		fmt.Println(row.Scan("policy_no"))
	//	}
	//
	// fmt.Println("GET POLICY NO _____________")
	// fmt.Println(row.Scan())
}

func BulkInsert(policies []structs.MotorPolicy) {
	transaction, err := DB.Begin()
	if err != nil {
		log.Panic(err)
	}
	// Prepare the COPY command with the target table and column names
	statement, err := transaction.Prepare(pq.CopyIn("motor_policies", "transaction_type", "rm_code", "rm_name", "child_id",
		"booking_date", "insurer_name", "insured_name", "major_categorisation", "product", "product_type", "policy_no",
		"plan_type", "premium", "net_premium", "od", "tp", "commissionable_premium", "registration_no", "rto_code", "state",
		"rto_cluster", "city", "insurer_biff", "fuel_type", "cpa", "cc", "gvw", "ncb_type", "seating_capacity", "vehicle_registration_year",
		"discount_percentage", "make", "model", "utr", "utr_date", "utr_amount", "slot", "paid_on_in", "tentative_in_percentage",
		"tentative_in_amount", "paid_on_out", "out_percentage", "out_amount", "co_type", "remarks"))

	if err != nil {
		log.Panic(err)
	}
	for _, policy := range policies {
		_, err = statement.Exec(
			// Execute the COPY command for each record

			policy.TransactionType,
			policy.RmCode,
			policy.RmName,
			policy.ChildId,
			policy.BookingDate,
			policy.InsurerName,
			policy.InsuredName,
			policy.MajorCategory,
			policy.Product,
			policy.ProductType,
			policy.PolicyNo,
			policy.PlanType,
			policy.Premium,
			policy.NetPremium,
			policy.OD,
			policy.TP,
			policy.CommissionablePremium,
			policy.RegistrationNo,
			policy.RtoCode,
			policy.State,
			policy.RtoCluster,
			policy.City,
			policy.InsurerBiff,
			policy.FuelType,
			policy.CPA,
			policy.CC,
			policy.GVW,
			policy.NcbType,
			policy.SeatingCapacity,
			policy.VehicleRegistrationYear,
			policy.DiscountPercentage,
			policy.Make,
			policy.Model,
			policy.UTR,
			policy.UTRDate,
			policy.UTRAmount,
			policy.SlotPaymentBatch,
			policy.PaidOnIn,
			policy.TentativeInPercentage,
			policy.TentativeInAmount,
			policy.PaidOnOut,
			policy.OutPercentage,
			policy.OutAmount,
			policy.CoType,
			policy.Remarks,
		)
		if err != nil {
			log.Panic(err)
			transaction.Rollback()
		}
	}
	// Finalize the COPY command and execute it. This step is necessary to complete the bulk insert operation.
	_, err = statement.Exec()
	if err != nil {
		log.Panic(err)
		transaction.Rollback()
	}

	// Close the prepared statement
	// The prepared statement is closed to release any associated resources.
	err = statement.Close()
	if err != nil {
		log.Panic(err)
	}

	// Commit the transaction to persist the changes
	err = transaction.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func CreateTable() {
	query := `
	CREATE TABLE IF NOT EXISTS motor_policies (
		transaction_type TEXT,
		rm_code TEXT,
		rm_name TEXT,
		child_id TEXT,
		booking_date TEXT,
		insurer_name TEXT,
		insured_name TEXT,
		major_categorisation TEXT,
		product TEXT,
		product_type TEXT,
		policy_no TEXT,
		plan_type TEXT,
		premium TEXT,
		net_premium TEXT,
		od TEXT,
		tp TEXT,
		commissionable_premium TEXT,
		registration_no TEXT,
		rto_code TEXT,
		state TEXT,
		rto_cluster TEXT,
		city TEXT,
		insurer_biff TEXT,
		fuel_type TEXT,
		cpa TEXT,
		cc TEXT,
		gvw TEXT,
		ncb_type TEXT,
		seating_capacity TEXT,
		vehicle_registration_year TEXT,
		discount_percentage TEXT,
		make TEXT,
		model TEXT,
		utr TEXT,
		utr_date TEXT,
		utr_amount TEXT,
		slot TEXT,
		paid_on_in TEXT,
		tentative_in_percentage TEXT,
		tentative_in_amount TEXT,
		paid_on_out TEXT,
		out_percentage TEXT,
		out_amount TEXT,
		co_type TEXT,
		remarks TEXT
	);
	`
	res, err := DB.Exec(query)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Create table result", res)
}
