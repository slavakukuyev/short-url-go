package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"sync"
)

const shortURLLength = 6 // Length of the generated short URL

type URLShortener struct {
	shortURLMap   map[string]string
	reverseURLMap map[string]string
	mutex         sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		shortURLMap:   make(map[string]string),
		reverseURLMap: make(map[string]string),
	}
}

func (us *URLShortener) generateShortURL() string {
	// Generate a random byte slice
	randomBytes := make([]byte, shortURLLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("Error generating random bytes")
	}

	// Encode the random bytes to base64 to create a URL-safe string
	shortURL := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim any special characters from the base64 encoded string
	shortURL = url.PathEscape(shortURL)

	return shortURL
}

func (us *URLShortener) ShortenURL(longURL string) string {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	// Check if the long URL is already shortened
	shortURL, exists := us.reverseURLMap[longURL]
	if exists {
		return shortURL
	}

	// Generate a new short URL
	shortURL = us.generateShortURL()
	us.shortURLMap[shortURL] = longURL
	us.reverseURLMap[longURL] = shortURL

	return shortURL
}

func (us *URLShortener) GetOriginalURL(shortURL string) (string, bool) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	// Lookup the original URL from the cache
	originalURL, found := us.shortURLMap[shortURL]
	return originalURL, found
}

func main() {
	urlShortener := NewURLShortener()

	// Example usage of the URL shortening algorithm and caching
	longURL := "https://www.example.com/some/long/path/to/a/resource"

	// Shorten the URL and store it in the cache
	shortURL := urlShortener.ShortenURL(longURL)

	fmt.Printf("Long URL: %s\n", longURL)
	fmt.Printf("Short URL: https://shorturl.com/%s\n", shortURL)

	// Retrieve the original URL from the cache using the short URL
	originalURL, found := urlShortener.GetOriginalURL(shortURL)
	if found {
		fmt.Printf("Original URL for short URL %s: %s\n", shortURL, originalURL)
	} else {
		fmt.Printf("Original URL for short URL %s not found in cache.\n", shortURL)
	}
}
