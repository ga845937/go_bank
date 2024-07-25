package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = os.Chdir("./internal/db/postgres")
	if err != nil {
		log.Fatalf("folder not exists：%v\n", err)
	}

	err = os.RemoveAll("./model")
	if err != nil {
		log.Fatalf("delete model error：%v\n", err)
	}

	// err = os.RemoveAll("./schema")
	// if err != nil {
	// 	log.Fatalf("delete schema error：%v\n", err)
	// }

	// err = os.Mkdir("./schema", 0755)
	// if err != nil {
	// 	log.Fatalf("create schema folder error：%v\n", err)
	// }

	// _, err = os.Create("./schema/.keep")
	// if err != nil {
	// 	log.Fatalf("create schema/.keep error：%v\n", err)
	// }

	// postgresURI := os.Getenv("POSTGRES_URI")
	// entImportCmd := exec.Command("go", "run", "ariga.io/entimport/cmd/entimport",
	// 	"-dsn", postgresURI,
	// 	"-schema-path", "./schema")

	// entImportCmd.Stdout = os.Stdout
	// entImportCmd.Stderr = os.Stderr

	// err = entImportCmd.Run()
	// if err != nil {
	// 	log.Fatalf("exec entimport error：%v\n", err)
	// }

	entGenerateCmd := exec.Command("go", "run", "-mod=mod", "entgo.io/ent/cmd/ent", "generate", "./schema", "--target", "./model")

	entGenerateCmd.Stdout = os.Stdout
	entGenerateCmd.Stderr = os.Stderr

	err = entGenerateCmd.Run()
	if err != nil {
		log.Fatalf("exec ent generateerror：%v\n", err)
	}
}
