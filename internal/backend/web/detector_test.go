package web

import (
	"testing"
)

func TestIsJobPage(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		body     string
		expected bool
	}{
		{
			name:     "URL with careers keyword",
			url:      "https://example.com/careers",
			body:     "<html><body>Welcome to our site</body></html>",
			expected: true,
		},
		{
			name:     "URL with jobs keyword",
			url:      "https://example.com/jobs",
			body:     "<html><body>Welcome to our site</body></html>",
			expected: true,
		},
		{
			name:     "URL with hiring keyword",
			url:      "https://example.com/hiring",
			body:     "<html><body>Welcome to our site</body></html>",
			expected: true,
		},
		{
			name:     "Body with careers keyword",
			url:      "https://example.com/about",
			body:     "<html><body><h1>Careers</h1><p>Join our team</p></body></html>",
			expected: true,
		},
		{
			name:     "Body with jobs keyword",
			url:      "https://example.com/about",
			body:     "<html><body><h1>Jobs</h1><p>We are hiring</p></body></html>",
			expected: true,
		},
		{
			name:     "Body with employment keyword",
			url:      "https://example.com/about",
			body:     "<html><body><h1>Employment</h1><p>Work with us</p></body></html>",
			expected: true,
		},
		{
			name:     "No job-related content",
			url:      "https://example.com/about",
			body:     "<html><body><h1>About Us</h1><p>We are a company</p></body></html>",
			expected: false,
		},
		{
			name:     "Script and style tags ignored",
			url:      "https://example.com/about",
			body:     "<html><body><script>var careers = 'jobs';</script><style>.careers { color: red; }</style><h1>About</h1></body></html>",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsJobPage(tt.url, tt.body)
			if result != tt.expected {
				t.Errorf("IsJobPage(%q, %q) = %v, want %v", tt.url, tt.body, result, tt.expected)
			}
		})
	}
}

func TestMatchesJobTitle(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		titles   []string
		expected bool
	}{
		{
			name:     "Exact title match with context",
			body:     "<html><body><h1>Software Engineer</h1><p>Apply now for this position</p></body></html>",
			titles:   []string{"software"},
			expected: true,
		},
		{
			name:     "Partial title match with context",
			body:     "<html><body><h1>Senior Software Engineer</h1><p>Join our team</p></body></html>",
			titles:   []string{"software"},
			expected: true,
		},
		{
			name:     "Multiple titles, one matches",
			body:     "<html><body><h1>Software Engineer</h1><p>Apply now</p></body></html>",
			titles:   []string{"developer", "software", "designer"},
			expected: true,
		},
		{
			name:     "Title without context words",
			body:     "<html><body><h1>Software Engineer</h1><p>This is just a description</p></body></html>",
			titles:   []string{"software"},
			expected: false,
		},
		{
			name:     "No titles provided",
			body:     "<html><body><h1>Software Engineer</h1><p>Apply now</p></body></html>",
			titles:   []string{},
			expected: true,
		},
		{
			name:     "Case insensitive matching",
			body:     "<html><body><h1>SOFTWARE ENGINEER</h1><p>Apply now</p></body></html>",
			titles:   []string{"software"},
			expected: true,
		},
		{
			name:     "Context words in nearby text",
			body:     "<html><body><h1>Software Engineer</h1><p>We are hiring for this full-time position</p></body></html>",
			titles:   []string{"software"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchesJobTitle(tt.body, tt.titles)
			if result != tt.expected {
				t.Errorf("MatchesJobTitle(%q, %v) = %v, want %v", tt.body, tt.titles, result, tt.expected)
			}
		})
	}
}

func TestExtractVisibleText(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple HTML",
			html:     "<html><body><h1>Title</h1><p>Content</p></body></html>",
			expected: "title content ",
		},
		{
			name:     "HTML with script and style",
			html:     "<html><head><script>var x = 1;</script><style>.class { color: red; }</style></head><body><h1>Title</h1></body></html>",
			expected: "title ",
		},
		{
			name:     "Nested elements",
			html:     "<html><body><div><h1>Title</h1><p>Content with <strong>bold</strong> text</p></div></body></html>",
			expected: "title content with bold text ",
		},
		{
			name:     "Empty HTML",
			html:     "",
			expected: "",
		},
		{
			name:     "Whitespace handling",
			html:     "<html><body>  <h1>  Title  </h1>  <p>  Content  </p>  </body></html>",
			expected: "title content ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVisibleText(tt.html)
			if result != tt.expected {
				t.Errorf("extractVisibleText(%q) = %q, want %q", tt.html, result, tt.expected)
			}
		})
	}
}

func TestHasNearbyContext(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		token    string
		expected bool
	}{
		{
			name:     "Token with nearby context",
			words:    []string{"we", "are", "hiring", "software", "engineer", "for", "this", "position"},
			token:    "software",
			expected: true,
		},
		{
			name:     "Token without nearby context",
			words:    []string{"software", "engineer", "is", "a", "good", "profession", "choice"},
			token:    "engineer",
			expected: false,
		},
		{
			name:     "Token at beginning with context",
			words:    []string{"software", "engineer", "apply", "now", "for", "this", "job"},
			token:    "software",
			expected: true,
		},
		{
			name:     "Token at end with context",
			words:    []string{"we", "are", "hiring", "a", "software", "engineer"},
			token:    "engineer",
			expected: true,
		},
		{
			name:     "Token not found",
			words:    []string{"we", "are", "hiring", "developers"},
			token:    "software",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNearbyContext(tt.words, tt.token)
			if result != tt.expected {
				t.Errorf("hasNearbyContext(%v, %q) = %v, want %v", tt.words, tt.token, result, tt.expected)
			}
		})
	}
}