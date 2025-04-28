package tests

import (
	"fullsteak/internal/article"
	"testing"
)

func TestArticleCreateOne(t *testing.T) {
	var err error

	a := &article.Article{
		Title:       "tytuł artykułu",
		Description: "streszczenie/opis",
		Content:     "Writing a custom session middleware in Go Echo involves intercepting requests to handle session creation, validation, and management. Below is a complete guide to creating a basic custom session middleware.",
		Public:      true,
	}

	repo := article.NewRepository(TestDB)
	err = repo.CreateOne(a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArticleGetAll(t *testing.T) {
	repo := article.NewRepository(TestDB)
	articles, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if articles[0].Title != "tytuł artykułu" {
		t.Fail()
	}
	article2 := &article.Article{
		Title:       "tytuł artykułu po aktualizacji",
		Description: "streszczenie/opis",
		Content:     "Writing a custom session middleware in Go Echo involves intercepting requests to handle session creation, validation, and management. Below is a complete guide to creating a basic custom session middleware.",
		Public:      true,
	}
	err = repo.UpdateOneById(article2, articles[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	articleFinal, err := repo.GetOneById(articles[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	if articleFinal.Title != "tytuł artykułu po aktualizacji" {
		t.Fail()
	}
}
