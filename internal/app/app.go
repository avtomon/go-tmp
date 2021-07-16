package app

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"os/signal"
	"parser/internal/domain"
	"parser/internal/repository"
	"parser/internal/service"
	"sync"
	"syscall"
	"time"
)

const (
	LogWriteErrorAllowableAmount = 3
)

func getDbConnect() *sqlx.DB {
	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		panic("database user environment variable must be setted")
	}

	dbPassword, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		panic("database password environment variable must be setted")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		panic("database password environment variable must be setted")
	}

	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		panic("database host with port environment variable must be setted")
	}

	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		panic(err)
	}

	return db
}

func Run() {

	loadEnvs()

	wg := sync.WaitGroup{}
	db := getDbConnect()

	destructor := func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}

		err = service.CloseLogFile()
		if err != nil {
			panic(err)
		}
	}

	defer destructor()

	siteConfigs, err := repository.GetSiteConfigs(db)
	if err != nil {
		panic("config receive error: " + err.Error())
	}

	for _, site := range *siteConfigs {
		var urls []string

		for _, catalogUrl := range site.CatalogUrls {
			maxPage := service.Search(catalogUrl)
			urls = append(urls, service.GenerateUrls(catalogUrl, maxPage)...)
		}

		shuffleUrls(&urls)

		wg.Add(1)
		site := site
		go func() {
			err := parseAndSave(db, &urls, &site)
			if err != nil {
				panic("parse data saving error: " + err.Error())
			}

			wg.Done()
		}()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		select {
		case sig := <-c:
			fmt.Printf("\n%s\nExiting...\n", sig)
			destructor()
			time.Sleep(time.Duration(1) * time.Second)
			os.Exit(1)
		}
	}()

	wg.Wait()
}

func shuffleUrls(urls *[]string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*urls), func(i, j int) { (*urls)[i], (*urls)[j] = (*urls)[j], (*urls)[i] })
}

func loadEnvs() {
	if err := godotenv.Load(); err != nil {
		panic("no .env file found")
	}
}

func parseAndSave(db *sqlx.DB, urls *[]string, siteConfig *domain.SiteConfig) error {
	result, errs := service.ParseSite(urls, siteConfig)
	logWriteErrorCount := 0
	for _, err := range *errs {
		errorMessage := fmt.Sprintf("parse error: %v", err.Error())
		fmt.Println(errorMessage)
		err := service.ToLog(errorMessage)
		if err != nil {
			logWriteErrorCount++
		}

		if logWriteErrorCount == LogWriteErrorAllowableAmount {
			return err
		}
	}

	return service.Save(db, result)
}
