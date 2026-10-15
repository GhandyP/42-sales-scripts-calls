package playbook

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/GhandyP/42-sales-scripts-calls/internal/sales"
)

const (
	maxPassageRunes   = 1000
	defaultSearchSize = 3
	maxSearchSize     = 5
)

var paragraphSeparator = regexp.MustCompile(`\r?\n[ \t]*\r?\n(?:[ \t]*\r?\n)*`)

// Index is an immutable, request-scoped search index over supplied playbook text.
type Index struct {
	passages []indexedPassage
}

type indexedPassage struct {
	evidence sales.EvidencePassage
	terms    map[string]struct{}
}

type searchHit struct {
	passage indexedPassage
	score   int
}

// NewIndex creates deterministic, bounded excerpts from one approved playbook.
// The returned index contains no shared or persistent state.
func NewIndex(approvedPlaybook string) (*Index, error) {
	if strings.TrimSpace(approvedPlaybook) == "" {
		return nil, sales.ValidationError{Field: "approved playbook text", Problem: "is required"}
	}
	if runeCount := utf8.RuneCountInString(approvedPlaybook); runeCount > sales.MaxApprovedPlaybookRunes {
		return nil, sales.ValidationError{
			Field:   "approved playbook text",
			Problem: fmt.Sprintf("must be at most %d characters", sales.MaxApprovedPlaybookRunes),
		}
	}

	index := &Index{passages: make([]indexedPassage, 0)}
	for _, paragraph := range paragraphSeparator.Split(approvedPlaybook, -1) {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		for _, text := range splitParagraph(paragraph) {
			id := fmt.Sprintf("playbook-%03d", len(index.passages)+1)
			evidence := sales.EvidencePassage{ID: id, Text: text}
			index.passages = append(index.passages, indexedPassage{
				evidence: evidence,
				terms:    tokenSet(text),
			})
		}
	}

	return index, nil
}

// Search returns the highest-ranked excerpts. Query terms are matched as
// case-insensitive Unicode words; equal scores retain source order. A zero
// limit selects the default result count.
func (index *Index) Search(query string, limit int) ([]sales.EvidencePassage, error) {
	if index == nil {
		return nil, errors.New("playbook index is required")
	}
	queryTerms := tokenSet(query)
	if len(queryTerms) == 0 {
		return nil, errors.New("search query must contain at least one term")
	}
	if limit == 0 {
		limit = defaultSearchSize
	}
	if limit < 1 || limit > maxSearchSize {
		return nil, fmt.Errorf("search limit must be between 1 and %d", maxSearchSize)
	}

	hits := make([]searchHit, 0, len(index.passages))
	for _, passage := range index.passages {
		score := 0
		for term := range queryTerms {
			if _, found := passage.terms[term]; found {
				score++
			}
		}
		if score > 0 {
			hits = append(hits, searchHit{passage: passage, score: score})
		}
	}

	sort.SliceStable(hits, func(i, j int) bool {
		return hits[i].score > hits[j].score
	})
	if len(hits) > limit {
		hits = hits[:limit]
	}

	results := make([]sales.EvidencePassage, len(hits))
	for i, hit := range hits {
		results[i] = hit.passage.evidence
	}
	return results, nil
}

func splitParagraph(paragraph string) []string {
	runes := []rune(paragraph)
	chunks := make([]string, 0, (len(runes)+maxPassageRunes-1)/maxPassageRunes)
	for start := 0; start < len(runes); {
		end := start + maxPassageRunes
		if end >= len(runes) {
			end = len(runes)
		} else {
			boundary := end
			for boundary > start && !unicode.IsSpace(runes[boundary-1]) {
				boundary--
			}
			if boundary > start {
				end = boundary
			}
		}

		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		start = end
	}
	return chunks
}

func tokenSet(text string) map[string]struct{} {
	terms := make(map[string]struct{})
	var word strings.Builder
	flush := func() {
		if word.Len() == 0 {
			return
		}
		terms[word.String()] = struct{}{}
		word.Reset()
	}

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			word.WriteRune(unicode.ToLower(r))
			continue
		}
		flush()
	}
	flush()
	return terms
}
