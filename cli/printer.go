package cli

import (
	"fmt"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/model"
)

func print(query *model.Query, results *model.QueryResults) {
	formattedResults := prepareResults(query.Fields, results)

	if len(formattedResults) == 0 {
		fmt.Println("no results!")
	} else {
		fmt.Println()
		fmt.Println("---------------------------")
	}

	for _, result := range formattedResults {
		fmt.Printf(result)
	}

	if len(formattedResults) > 0 {
		fmt.Println("---------------------------")
		fmt.Println()
	}

}

func prepareResults(fields *common.StringSlice, results *model.QueryResults) []string {
	var formattedResults = make([]string, 0)

	sources := results.Sources
	meta := results.Page.Meta
	sections := results.Page.Sections

	for _, source := range sources {
		formattedResults = append(formattedResults, fmt.Sprintf("Name:  %s\n", source.Name))
		formattedResults = append(formattedResults, fmt.Sprintf("URL:  %s\n", source.Url))
		formattedResults = append(formattedResults, fmt.Sprintf("Last update: %s\n", source.LastUpdate))
		formattedResults = append(formattedResults, fmt.Sprintf("Title: %s\n", source.Page.Meta.Title))
		formattedResults = append(formattedResults, fmt.Sprintf("Description: %s\n", source.Page.Meta.Description))
		formattedResults = append(formattedResults, fmt.Sprintf("Twitter: %s\n", source.Page.Meta.Twitter))
		formattedResults = append(formattedResults, fmt.Sprintf("Locale: %s\n", source.Page.Meta.Locale))
		formattedResults = append(formattedResults, fmt.Sprintf("Total sections: %d\n", len(source.Page.Sections)))
	}

	if meta != nil {
		if fields.Contains(model.SourceMetaTitle) || fields.Contains(model.Wildcard) {
			formattedResults = append(formattedResults, fmt.Sprintf("Title:  %s\n", meta.Title))
		}
		if fields.Contains(model.SourceMetaTwitter) || fields.Contains(model.Wildcard) {
			formattedResults = append(formattedResults, fmt.Sprintf("Twitter:  %s\n", meta.Twitter))
		}
		if fields.Contains(model.SourceMetaUrl) || fields.Contains(model.Wildcard) {
			formattedResults = append(formattedResults, fmt.Sprintf("Url:  %s\n", meta.Url))
		}
		if fields.Contains(model.SourceMetaDescription) || fields.Contains(model.Wildcard) {
			formattedResults = append(formattedResults, fmt.Sprintf("Description:  %s\n", meta.Description))
		}
		if fields.Contains(model.SourceMetaLocale) || fields.Contains(model.Wildcard) {
			formattedResults = append(formattedResults, fmt.Sprintf("Locale:  %s\n", meta.Locale))
		}
		if meta.SiteMap != nil && (fields.AnyPrefix(model.SourcePageSitemap) || fields.Contains(model.Wildcard)) {
			for _, site := range meta.SiteMap.Sites {
				if fields.Contains(model.SourcePageSitemapUrl) || fields.Contains(model.Wildcard) {
					formattedResults = append(formattedResults, fmt.Sprintf("Sitemap URL:  %s\n", site.Loc))
				}
				if fields.Contains(model.SourcePageSitemapLastMod) || fields.Contains(model.Wildcard) {
					formattedResults = append(formattedResults, fmt.Sprintf("Sitemap Last Modified:  %s\n", site.Lastmod))
				}
				if fields.Contains(model.SourcePageSitemapChangeFreq) || fields.Contains(model.Wildcard) {
					formattedResults = append(formattedResults, fmt.Sprintf("Sitemap Change frequency:  %s\n", site.ChangeFreq))
				}
			}
		}
	}

	if len(sections) > 0 {
		for _, section := range sections {
			if fields.Contains(model.SourcePageText) || fields.Contains(model.Wildcard) {
				formattedResults = append(formattedResults, fmt.Sprintf("Text:  %s\n", section.Text))
			}
			if fields.Contains(model.SourcePageLinks) || fields.Contains(model.Wildcard) {
				formattedResults = append(formattedResults, fmt.Sprintf("Links:  %s\n", section.Links))
			}
		}
	}

	return formattedResults
}
