package seeders

import (
	"devin/database"
)

func SeedDB() {
	db := database.NewGORMInstance()
	defer db.Close()
	db.LogMode(false)
	db.Exec(`insert into address_countries (id, name, phone_prefix, alpha2_code, alpha3_code, locale_code) values (1, 'Iran', '+98', 'IR', 'IRR', 'fa_IR');`)
	db.Exec(`insert into address_countries (id, name, phone_prefix, alpha2_code, alpha3_code, locale_code) values (2, 'Germany', '+49', 'DE', 'DEU', 'de_DE');`)

	db.Exec(`insert into address_provinces (id, name, country_id) values (1, 'Fars', 1);`)
	db.Exec(`insert into address_provinces (id, name, country_id) values (2, 'Tehran', 1);`)
	db.Exec(`insert into address_provinces (id, name, country_id) values (3, 'Berlin', 2);`)

	db.Exec(`insert into address_cities (id, name, country_id, province_id) values (1, 'Shiraz', 1, 1);`)
	db.Exec(`insert into address_cities (id, name, country_id, province_id) values (2, 'Marvdasht', 1, 1);`)

	db.Exec(`insert into address_cities (id, name, country_id, province_id) values (3, 'Tehran', 1, 2);`)
	db.Exec(`insert into address_cities (id, name, country_id, province_id) values (4, 'Karaj', 1, 2);`)

	db.Exec(`insert into calendar_systems (id, name, component_name, filter_name) values (1, 'Jalali', 'jdate', 'jdate');`)
	db.Exec(`insert into calendar_systems (id, name, component_name, filter_name) values (2, 'Gregorian', 'gregorian', 'gregorian');`)

	db.Exec(`insert into date_formats (id, name) values (1, '2006-01-02');`)
	db.Exec(`insert into date_formats (id, name) values (2, '2006/01/02');`)

	db.Exec(`insert into time_formats (id, name) values (1, '15:04:05');`)
	db.Exec(`insert into time_formats (id, name) values (2, '15:04');`)
}
