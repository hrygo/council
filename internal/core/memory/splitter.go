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
	currentDocs := []string{}
	currentLen := 0

	for _, split := range splits {
		sepLen := 0
		if len(currentDocs) > 0 {
			sepLen = len(separator)
		}

		if currentLen+sepLen+len(split) > s.ChunkSize && currentLen > 0 {
			// Flush current docs
			doc := strings.Join(currentDocs, separator)
			goodSplits = append(goodSplits, doc)

			// Handle overlap: keep last N characters or last M segments?
			// LangChain style: keep segments until overlapping length is reached.
			for currentLen > s.ChunkOverlap && len(currentDocs) > 0 {
				removed := currentDocs[0]
				currentDocs = currentDocs[1:]
				currentLen -= len(removed)
				if len(currentDocs) > 0 {
					currentLen -= len(separator)
				}
			}
		}

		if len(split) > s.ChunkSize {
			if len(newSeparators) > 0 {
				subSplits := s.split(split, newSeparators)
				goodSplits = append(goodSplits, subSplits...)
			} else {
				goodSplits = append(goodSplits, split)
			}
			// Reset current since we just flushed a giant chunk
			currentDocs = []string{}
			currentLen = 0
		} else {
			currentDocs = append(currentDocs, split)
			if currentLen == 0 {
				currentLen = len(split)
			} else {
				currentLen += len(separator) + len(split)
			}
		}
	}

	if len(currentDocs) > 0 {
		goodSplits = append(goodSplits, strings.Join(currentDocs, separator))
	}

	return goodSplits
}
