package main

import (
	"log"

	"github.com/jdkato/prose/v2"
	"github.com/mcnijman/go-emailaddress"
	"golang.org/x/exp/slices"
)

func filterEmails(emails []string) (filteredEmails []string) {
	for _, email := range emails {
		parsedEmail, err := emailaddress.Parse(email)
		if err != nil {
			log.Printf("Email parsing failed for %s: %v", email, err)
			continue
		}

		if !slices.Contains(EMAIL_DOMAIN_BLACKLIST, parsedEmail.Domain) {
			filteredEmails = append(filteredEmails, parsedEmail.String())
		}
	}

	return
}

func extractAndFilterEmailsFromText(text string) (filteredEmails []string) {
	parsedEmails := emailaddress.FindWithIcannSuffix([]byte(text), false)

	for _, e := range parsedEmails {
		if !slices.Contains(EMAIL_DOMAIN_BLACKLIST, e.Domain) {
			filteredEmails = append(filteredEmails, e.String())
		}
	}

	return
}

func extractEntitiesFromText(text string) (entities []string) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Printf("Text failed for prose %s: %v", text, err)
		return
	}

	for _, ent := range doc.Entities() {
		entities = append(entities, ent.Text)
	}

	return
}
