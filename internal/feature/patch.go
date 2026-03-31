package feature

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func appendConfigFromTemplate(fsys fs.FS, configPath, templatePath string) error {
	yamlSnippet, err := fs.ReadFile(fsys, templatePath)
	if err != nil {
		return err
	}
	existingContent, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	newContent := string(existingContent) + "\n" + string(yamlSnippet)
	return os.WriteFile(configPath, []byte(newContent), 0644)
}

func updateConfigGoForMySQL(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.Contains(contentStr, "DBConfig") && strings.Contains(contentStr, "DB DBConfig") {
		return nil
	}
	contentStr = strings.Replace(contentStr,
		"	// DB DBConfig `mapstructure:\"db\"` // Example for when DB is added",
		"	DB DBConfig `mapstructure:\"db\"`",
		1)
	contentStr = strings.Replace(contentStr,
		"// DBConfig might be added by `ginGen add mysql`",
		"// DBConfig holds database configuration",
		1)
	contentStr = strings.Replace(contentStr,
		"// type DBConfig struct {",
		"type DBConfig struct {",
		1)
	contentStr = strings.Replace(contentStr,
		"//   MySQL MySQLConfig `mapstructure:\"mysql\"`",
		"	MySQL MySQLConfig `mapstructure:\"mysql\"`",
		1)
	contentStr = strings.Replace(contentStr,
		"// }",
		"}",
		1)
	contentStr = strings.Replace(contentStr,
		"// type MySQLConfig struct {",
		"type MySQLConfig struct {",
		1)
	contentStr = strings.Replace(contentStr,
		"// 	DSN             string `mapstructure:\"dsn\"`",
		"	DSN             string `mapstructure:\"dsn\"`",
		1)
	contentStr = strings.Replace(contentStr,
		"// 	MaxIdleConns    int    `mapstructure:\"max_idle_conns\"`",
		"	MaxIdleConns    int    `mapstructure:\"max_idle_conns\"`",
		1)
	contentStr = strings.Replace(contentStr,
		"// 	MaxOpenConns    int    `mapstructure:\"max_open_conns\"`",
		"	MaxOpenConns    int    `mapstructure:\"max_open_conns\"`",
		1)
	contentStr = strings.Replace(contentStr,
		"// 	ConnMaxLifetime int    `mapstructure:\"conn_max_lifetime\"` // in seconds",
		"	ConnMaxLifetime int    `mapstructure:\"conn_max_lifetime\"` // in seconds",
		1)
	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func patchMySQL(fsys fs.FS, projectPath string) error {
	configYaml := filepath.Join(projectPath, "configs", "config.yaml")
	if err := appendConfigFromTemplate(fsys, configYaml, "templates/addfeature/mysql/config_entry.yaml.tmpl"); err != nil {
		return err
	}
	return updateConfigGoForMySQL(filepath.Join(projectPath, "internal", "config", "config.go"))
}

func updateConfigGoForRedis(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.Contains(contentStr, "Redis RedisConfig") {
		return nil
	}
	if strings.Contains(contentStr, "DB DBConfig") {
		contentStr = strings.Replace(contentStr,
			"	DB DBConfig `mapstructure:\"db\"`",
			"	DB DBConfig `mapstructure:\"db\"`\n	Redis RedisConfig `mapstructure:\"redis\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	Redis RedisConfig `mapstructure:\"redis\"`",
			1)
	}
	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func patchRedis(fsys fs.FS, projectPath string) error {
	configYaml := filepath.Join(projectPath, "configs", "config.yaml")
	if err := appendConfigFromTemplate(fsys, configYaml, "templates/addfeature/redis/config_entry.yaml.tmpl"); err != nil {
		return err
	}
	return updateConfigGoForRedis(filepath.Join(projectPath, "internal", "config", "config.go"))
}

func updateConfigGoForKafka(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.Contains(contentStr, "Kafka KafkaConfig") {
		return nil
	}
	if strings.Contains(contentStr, "Redis RedisConfig") {
		contentStr = strings.Replace(contentStr,
			"	Redis RedisConfig `mapstructure:\"redis\"`",
			"	Redis RedisConfig `mapstructure:\"redis\"`\n	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			1)
	} else if strings.Contains(contentStr, "DB DBConfig") {
		contentStr = strings.Replace(contentStr,
			"	DB DBConfig `mapstructure:\"db\"`",
			"	DB DBConfig `mapstructure:\"db\"`\n	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			1)
	}
	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func patchKafka(fsys fs.FS, projectPath string) error {
	configYaml := filepath.Join(projectPath, "configs", "config.yaml")
	if err := appendConfigFromTemplate(fsys, configYaml, "templates/addfeature/kafka/config_entry.yaml.tmpl"); err != nil {
		return err
	}
	return updateConfigGoForKafka(filepath.Join(projectPath, "internal", "config", "config.go"))
}

func updateConfigGoForPostgres(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.Contains(contentStr, "PostgreSQL PostgreSQLConfig") {
		return nil
	}
	if strings.Contains(contentStr, "DB DBConfig") {
		return nil
	}
	if strings.Contains(contentStr, "Redis RedisConfig") {
		contentStr = strings.Replace(contentStr,
			"	Redis RedisConfig `mapstructure:\"redis\"`",
			"	Redis RedisConfig `mapstructure:\"redis\"`\n	PostgreSQL PostgreSQLConfig `mapstructure:\"postgres\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	PostgreSQL PostgreSQLConfig `mapstructure:\"postgres\"`",
			1)
	}
	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func patchPostgres(fsys fs.FS, projectPath string) error {
	configYaml := filepath.Join(projectPath, "configs", "config.yaml")
	if err := appendConfigFromTemplate(fsys, configYaml, "templates/addfeature/postgres/config_entry.yaml.tmpl"); err != nil {
		return err
	}
	return updateConfigGoForPostgres(filepath.Join(projectPath, "internal", "config", "config.go"))
}

func updateConfigGoForJWT(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.Contains(contentStr, "JWT JWTConfig") {
		return nil
	}
	if strings.Contains(contentStr, "Kafka KafkaConfig") {
		contentStr = strings.Replace(contentStr,
			"	Kafka KafkaConfig `mapstructure:\"kafka\"`",
			"	Kafka KafkaConfig `mapstructure:\"kafka\"`\n	JWT JWTConfig `mapstructure:\"jwt\"`",
			1)
	} else if strings.Contains(contentStr, "Redis RedisConfig") {
		contentStr = strings.Replace(contentStr,
			"	Redis RedisConfig `mapstructure:\"redis\"`",
			"	Redis RedisConfig `mapstructure:\"redis\"`\n	JWT JWTConfig `mapstructure:\"jwt\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	JWT JWTConfig `mapstructure:\"jwt\"`",
			1)
	}
	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func patchJWT(fsys fs.FS, projectPath string) error {
	configYaml := filepath.Join(projectPath, "configs", "config.yaml")
	if err := appendConfigFromTemplate(fsys, configYaml, "templates/addfeature/jwt/config_entry.yaml.tmpl"); err != nil {
		return err
	}
	return updateConfigGoForJWT(filepath.Join(projectPath, "internal", "config", "config.go"))
}

func updateConfigGoForLogger(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.Contains(contentStr, "Logger LoggerConfig") {
		return nil
	}
	if strings.Contains(contentStr, "JWT JWTConfig") {
		contentStr = strings.Replace(contentStr,
			"	JWT JWTConfig `mapstructure:\"jwt\"`",
			"	JWT JWTConfig `mapstructure:\"jwt\"`\n	Logger LoggerConfig `mapstructure:\"logger\"`",
			1)
	} else {
		contentStr = strings.Replace(contentStr,
			"	Server ServerConfig `mapstructure:\"server\"`",
			"	Server ServerConfig `mapstructure:\"server\"`\n	Logger LoggerConfig `mapstructure:\"logger\"`",
			1)
	}
	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

func patchLogger(fsys fs.FS, projectPath string) error {
	configYaml := filepath.Join(projectPath, "configs", "config.yaml")
	if err := appendConfigFromTemplate(fsys, configYaml, "templates/addfeature/logger/config_entry.yaml.tmpl"); err != nil {
		return err
	}
	return updateConfigGoForLogger(filepath.Join(projectPath, "internal", "config", "config.go"))
}
