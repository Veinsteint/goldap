package common

import (
	"fmt"
	"log"
	"os"
	"time"

	"goldap-server/config"
	"goldap-server/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	switch config.Conf.Database.Driver {
	case "mysql":
		DB = ConnMysql()
	case "sqlite3":
		DB = ConnSqlite()
	}
	dbAutoMigrate()
}

func dbAutoMigrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Group{},
		&model.Menu{},
		&model.Api{},
		&model.OperationLog{},
		&model.FieldRelation{},
		&model.SSHKey{},
		&model.IPGroup{},
		&model.IPGroupUserPermission{},
		&model.SudoRule{},
		&model.GroupUserPermission{},
		&model.PendingUser{},
		&model.UserPreConfig{},
		&model.SystemConfig{},
	)
	if err != nil {
		Log.Errorf("Auto migration failed: %v", err)
	}

	if config.Conf.Database.Driver == "mysql" {
		applyMysqlMigrations()
	}
}

func applyMysqlMigrations() {
	migrations := []string{
		"ALTER TABLE `users` MODIFY COLUMN password TEXT NOT NULL",
	}

	indexMigrations := []struct {
		check  string
		action string
	}{
		{
			"SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'users' AND index_name = 'idx_mobile'",
			"ALTER TABLE `users` DROP INDEX `idx_mobile`",
		},
	}

	columnMigrations := []struct {
		column string
		sql    string
	}{
		{"ip_ranges", "ALTER TABLE `groups` ADD COLUMN `ip_ranges` TEXT"},
		{"uid_number", "ALTER TABLE `users` ADD COLUMN `uid_number` INT DEFAULT NULL"},
		{"gid_number", "ALTER TABLE `users` ADD COLUMN `gid_number` INT DEFAULT 0"},
		{"home_directory", "ALTER TABLE `users` ADD COLUMN `home_directory` VARCHAR(255) DEFAULT NULL"},
		{"login_shell", "ALTER TABLE `users` ADD COLUMN `login_shell` VARCHAR(255) DEFAULT '/bin/bash'"},
		{"gecos", "ALTER TABLE `users` ADD COLUMN `gecos` VARCHAR(255) DEFAULT NULL"},
	}

	for _, sql := range migrations {
		_ = DB.Exec(sql).Error
	}

	for _, m := range indexMigrations {
		var count int64
		if DB.Raw(m.check).Scan(&count).Error == nil && count > 0 {
			_ = DB.Exec(m.action).Error
		}
	}

	for _, m := range columnMigrations {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = '%s'", m.column)
		if m.column == "ip_ranges" {
			query = "SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'groups' AND column_name = 'ip_ranges'"
		}
		if DB.Raw(query).Scan(&count).Error == nil && count == 0 {
			_ = DB.Exec(m.sql).Error
		}
	}

	// Add unique index for uid_number
	var count int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'users' AND index_name = 'idx_uid_number'").Scan(&count)
	if count == 0 {
		_ = DB.Exec("ALTER TABLE `users` ADD UNIQUE INDEX `idx_uid_number` (`uid_number`)").Error
	}
}

// silentLogger returns a GORM logger that only logs errors
func silentLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Error, // Only log errors
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

func ConnSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.Conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   silentLogger(),
	})
	if err != nil {
		Log.Panicf("SQLite connection failed: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db
}

func ConnMysql() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
		config.Conf.Mysql.Charset,
		config.Conf.Mysql.Collation,
		config.Conf.Mysql.Query,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   silentLogger(),
	})
	if err != nil {
		Log.Panicf("MySQL connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		Log.Panicf("Failed to get DB object: %v", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
