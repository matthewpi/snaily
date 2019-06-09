package utils

import (
	"github.com/matthewpi/snaily/bot"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

// IsProfane checks if a string contains profane language.
func IsProfane(message string) (bool, error) {
	snaily := bot.GetBot()

	var err error
	message, err = sanitize(message)
	if err != nil {
		return false, err
	}

	for _, word := range snaily.Config.Filter.Words {
		if match := strings.Contains(message, word); match {
			return true, nil
		}
	}

	return false, nil
}

// sanitize sanitizes a message so it can be checked for profanity.
func sanitize(message string) (string, error) {
	// Remove accents from the string.
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	message, _, err := transform.String(t, message)
	if err != nil {
		return "", err
	}

	// Remove zero-width spaces.
	message = strings.ReplaceAll(message, "\u200b", "")

	// Remove whitespace.
	message = removeWhitespace(message)

	// Convert the string to lowercase.
	message = strings.ToLower(message)

	// Replace repeating characters.
	message = removeRepeating(message)

	// Convert numbers to letters.
	message = strings.ReplaceAll(message, "0", "o")
	message = strings.ReplaceAll(message, "1", "i")
	message = strings.ReplaceAll(message, "3", "e")
	message = strings.ReplaceAll(message, "4", "a")
	message = strings.ReplaceAll(message, "5", "s")
	message = strings.ReplaceAll(message, "6", "b")
	message = strings.ReplaceAll(message, "7", "l")
	message = strings.ReplaceAll(message, "8", "b")

	// Replace symbols with letters.
	message = strings.ReplaceAll(message, "@", "a")
	message = strings.ReplaceAll(message, "!", "a")
	message = strings.ReplaceAll(message, "$", "s")
	message = strings.ReplaceAll(message, "_", " ")
	message = strings.ReplaceAll(message, "-", " ")
	message = strings.ReplaceAll(message, "*", " ")
	message = strings.ReplaceAll(message, "()", "0")

	return message, nil
}

// loop through characters
// store an array of these characters
// check if the current character matches the previous one in the array
// update the twoConsecutive boolean
// if a third character matches keep checking for more and more of the same character,
// if a third character doesn't match then reset the values and start checking the next character.
func removeRepeating(message string) string {
	msg := message
	previous := make([]int32, 0)
	twoConsecutive := false
	consecutiveIndex := 0

	// Loop through all characters in the string.
	for in, char := range message {
		// Skip the first index, but add it to the array.
		if in != 0 {
			// Check if the current character matches the previous one.
			if char == previous[len(previous)-1] {
				// Update the twoConsecutive boolean.
				if !twoConsecutive {
					twoConsecutive = true
					consecutiveIndex = len(previous) - 1
				}
			} else {
				// Check if there are two consecutive characters.
				if twoConsecutive {
					// Replace the consecutive characters with one.
					msg = strings.Replace(msg, message[consecutiveIndex+1:in], "", 1)
				}

				twoConsecutive = false
			}
		}

		// Add the character to the array.
		previous = append(previous, char)
	}

	// Make sure the last part of the string can actually be replaced do to how our for-loop is structured.
	if twoConsecutive {
		// Replace the consecutive characters with one.
		msg = strings.Replace(msg, message[consecutiveIndex+1:], "", 1)
	}

	return msg
}

func removeWhitespace(message string) string {
	var output []rune
	for _, r := range message {
		if !unicode.IsSpace(r) {
			output = append(output, r)
		}
	}

	return string(output)
}
