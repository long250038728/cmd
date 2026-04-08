package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/gen"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/spf13/cobra"
)

/**
*   go build -o ormgen
*  ./ormgen -s zhubaoe -t "zby_customer,zby_user"
*  ./ormgen -s zhubaoe -t "zby_customer,zby_user" -o true
 */

// Config 应用配置结构体
type Config struct {
	Sources map[string]SourceConfig `yaml:"sources"`
}

// SourceConfig 数据源配置结构体
type SourceConfig struct {
	ConfigPath string `yaml:"config_path"`
	DbName     string `yaml:"db_name"`
}

var (
	source     string
	tables     string
	output     bool
	appConfig  Config
	configFile = "./config/config.yaml"
	outputFile = "./gen/models.go"
)

// 初始化配置
func initConfig() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	configFile = filepath.Join(exeDir, "config", "config.yaml")
	configurator.NewYaml().MustLoad(configFile, &appConfig)
}

func init() {
	// 加载配置文件
	initConfig()

	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "Source name (zhubaoe or emperor)")
	rootCmd.PersistentFlags().StringVarP(&tables, "tables", "t", "", "Table names, separated by commas")
	rootCmd.PersistentFlags().BoolVarP(&output, "output", "o", false, "Output to models.go file")

	rootCmd.MarkPersistentFlagRequired("source")
	rootCmd.MarkPersistentFlagRequired("tables")
}

var rootCmd = &cobra.Command{
	Use:   "ormgen",
	Short: "ORM struct generator",
	Long:  `Generate ORM structs from database tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var db *orm.Gorm
		var b []byte
		var configPath string
		var dbName string

		// 根据source获取配置信息
		configInfo, exists := appConfig.Sources[source]
		if !exists {
			fmt.Printf("Invalid source: %s. Available sources: %v\n", source, getAvailableSources())
			os.Exit(1)
		}

		// 使用source对应的配置，路径相对于可执行文件目录
		exePath, _ := os.Executable()
		exePath, _ = filepath.EvalSymlinks(exePath)
		exeDir := filepath.Dir(exePath)
		configPath = filepath.Join(exeDir, configInfo.ConfigPath)
		dbName = configInfo.DbName

		// 加载数据库配置文件
		var ormConfig orm.Config
		configurator.NewYaml().MustLoad(configPath, &ormConfig)

		// 创建数据库连接
		db, err = orm.NewMySQLGorm(&ormConfig)
		if err != nil {
			fmt.Printf("Failed to connect to database: %v\n", err)
			os.Exit(1)
		}

		// 解析表名
		tableList := strings.Split(tables, ",")
		for i, table := range tableList {
			tableList[i] = strings.TrimSpace(table)
		}

		// 生成ORM结构体
		if b, err = gen.NewModelsGen(db).Gen(dbName, tableList); err != nil {
			fmt.Printf("Failed to generate models: %v\n", err)
			os.Exit(1)
		}

		// 输出结果
		if output {
			// 检查并创建gen目录
			genDir := "./gen"
			if _, err := os.Stat(genDir); os.IsNotExist(err) {
				if err := os.MkdirAll(genDir, os.ModePerm); err != nil {
					fmt.Printf("Failed to create directory: %v\n", err)
					os.Exit(1)
				}
			}
			// 写入文件
			if err := os.WriteFile(outputFile, b, os.ModePerm); err != nil {
				fmt.Printf("Failed to write file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Models generated successfully to %s!\n", outputFile)
		} else {
			// 输出到控制台
			fmt.Println(string(b))
			fmt.Println("\nModels generated successfully!")
		}
	},
}

// getAvailableSources 返回所有可用的source名称
func getAvailableSources() []string {
	sources := make([]string, 0, len(appConfig.Sources))
	for s := range appConfig.Sources {
		sources = append(sources, s)
	}
	return sources
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
