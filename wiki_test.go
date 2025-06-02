package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type WikiTestSuite struct {
	suite.Suite
}

// テストのエントリーポイント
func TestWikiSuite(t *testing.T) {
	suite.Run(t, new(WikiTestSuite))
}

func (s *WikiTestSuite) TestPlay() {

	s.Run("ページ内容をファイルに保存test", func() {
		title := "TestPage"
		body := []byte("これはテストです。")
		p := &Page{Title: title, Body: body}

		// 保存
		err := p.save()
		s.NoError(err, "保存時にエラーが発生しないこと")

		// 読み込み
		loaded, err := loadPage(title)
		s.NoError(err, "読み込み時にエラーが発生しないこと")
		s.Equal(p.Title, loaded.Title, "タイトルが一致すること")
		s.Equal(p.Body, loaded.Body, "本文が一致すること")
	})

	// s.Run("viewが表示されるかの確認", func() {

	// })
}
