package database

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"worker-nicepay/infrastructure/configuration"
	"worker-nicepay/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var YugabyteDBClient *gorm.DB

func InitializeYugabyteDB() {
	conn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		configuration.AppConfig.YugabyteHost,
		configuration.AppConfig.YugabyteUsername,
		configuration.AppConfig.YugabytePassword,
		configuration.AppConfig.YugabyteDatabase,
		configuration.AppConfig.YugabytePort,
	)

	fmt.Println(conn)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
			NameReplacer:  strings.NewReplacer("DataModel", ""),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	YugabyteDBClient = db

	registerUUIDv7BeforeCreate(YugabyteDBClient)

	if err := YugabyteDBClient.AutoMigrate(
		&models.PaymentXenditQrisesDataModel{},
		&models.PaymentXenditVasDataModel{},
		&models.PaymentXenditEWalletsDataModel{},
		&models.EWalletProvidersDataModel{},
		&models.VAProvidersDataModel{},
		&models.PaymentsDataModel{},
		&models.MerchantsDataModel{},
		&models.PaymentMethodsDataModel{},
		&models.CurrenciesDataModel{},
		&models.CountriesDataModel{},
	); err != nil {
		log.Fatal(err)
	}
}

func registerUUIDv7BeforeCreate(db *gorm.DB) {
	if db == nil {
		return
	}

	db.Callback().Create().Before("gorm:create").Register("uuidv7_before_create", func(tx *gorm.DB) {
		if tx == nil || tx.Statement == nil {
			return
		}
		rv := tx.Statement.ReflectValue
		if !rv.IsValid() {
			return
		}
		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return
			}
			rv = rv.Elem()
		}

		setIDIfZero := func(v reflect.Value) {
			for v.Kind() == reflect.Ptr {
				if v.IsNil() {
					return
				}
				v = v.Elem()
			}
			if !v.IsValid() || v.Kind() != reflect.Struct {
				return
			}
			f := v.FieldByName("ID")
			if !f.IsValid() || !f.CanSet() || f.Type() != reflect.TypeOf(uuid.UUID{}) {
				return
			}
			if id, ok := f.Interface().(uuid.UUID); ok && id == uuid.Nil {
				newID, err := uuid.NewV7()
				if err != nil {
					return
				}
				f.Set(reflect.ValueOf(newID))
			}
		}

		switch rv.Kind() {
		case reflect.Struct:
			setIDIfZero(rv)
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				setIDIfZero(rv.Index(i))
			}
		}
	})
}
