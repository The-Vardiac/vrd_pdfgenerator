package services

import "github.com/williamluisan/vrd_pdfgenerator/repository"

type GenerateFileEmailTemplate repository.GenerateFileEmailTemplate

func (generateFileEmailTemplate *GenerateFileEmailTemplate) GetGenerateFileHTMLEmailTemplate() (html string) {
	html = `
		<html>
			<body>
				<p>Hi, ` + generateFileEmailTemplate.UserName + `</p>
				<p>
					Kindly find your PDF file (` + generateFileEmailTemplate.Filename + `) on this url: <br/>
					` + generateFileEmailTemplate.FileUrl + `
				</p>
				<p>
					Thank you!<br/>
					The Vardiac Team.
				</p>
			</body>
		</html>
	`

	return html
}