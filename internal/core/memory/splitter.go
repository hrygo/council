package memory

import (
	"strings"
)

// RecursiveCharacterSplitter splits text recursively by separators
type RecursiveCharacterSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string
}

func NewRecursiveCharacterSplitter(chunkSize, chunkOverlap int) *RecursiveCharacterSplitter {
	return &RecursiveCharacterSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		Separators:   []string{"\n\n", "\n", " ", ""}, // Order matters: priority top to bottom
	}
}

func (s *RecursiveCharacterSplitter) SplitText(text string) []string {
	finalChunks := []string{}
	goodSplits := s.split(text, s.Separators)
	finalChunks = append(finalChunks, goodSplits...)
	return finalChunks
}

func (s *RecursiveCharacterSplitter) split(text string, separators []string) []string {
	separator := separators[len(separators)-1]
	newSeparators := []string{}

	for i, sep := range separators {
		if sep == "" {
			separator = ""
			break
		}
		if strings.Contains(text, sep) {
			separator = sep
			newSeparators = separators[i+1:]
			break
		}
	}

	splits := strings.Split(text, separator)

	// Now merge splits that occur to fit chunk size
	goodSplits := []string{}
	currentDoc := ""

	for _, split := range splits {
		sepLen := len(separator)
		if currentDoc != "" {
			// Check if we can add
			if len(currentDoc)+sepLen+len(split) > s.ChunkSize {
				// Cannot add, flush currentDoc
				if currentDoc != "" {
					goodSplits = append(goodSplits, currentDoc)
				}
				// If this single split is too big, recurse on it
				if len(split) > s.ChunkSize && len(newSeparators) > 0 {
					subSplits := s.split(split, newSeparators)
					goodSplits = append(goodSplits, subSplits...)
					currentDoc = ""
				} else {
					// Otherwise start new doc
					currentDoc = split
				}
			} else {
				// Can add
				currentDoc += separator + split
			}
		} else {
			// First entry
			if len(split) > s.ChunkSize && len(newSeparators) > 0 {
				subSplits := s.split(split, newSeparators)
				goodSplits = append(goodSplits, subSplits...)
			} else {
				currentDoc = split
			}
		}
	}

	if currentDoc != "" {
		goodSplits = append(goodSplits, currentDoc)
	}

	return goodSplits
}
