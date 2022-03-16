package candidate

type Candidate struct {
	ID         int          `json:"candidate-id,omitempty"`
	Contact    Contact      `json:"contact,omitempty"`
	Experience []Experience `json:"experience,omitempty"`
	Projects   []Project    `json:"projects,omitempty"`
	DevEnvs    []DevEnv     `json:"dev-env,omitempty"`
	Hobbies    []Hobby      `json:"hobbies,omitempty"`
}

type Contact struct {
	Name    string   `json:"name,omitempty"`
	Address string   `json:"address,omitempty"`
	Email   string   `json:"email,omitempty"`
	URLs    []string `json:"urls,omitempty"`
}

type Experience struct {
	Date        string   `json:"job-date,omitempty"`
	Company     string   `json:"job-company,omitempty"`
	Title       string   `json:"job-title,omitempty"`
	Location    string   `json:"job-location,omitempty"`
	Description string   `json:"job-description,omitempty"`
	Highlights  []string `json:"job-highlights,omitempty"`
}

type Project struct {
	Date        string   `json:"project-date,omitempty"`
	Name        string   `json:"project-name,omitempty"`
	URL         string   `json:"project-url,omitempty"`
	Description string   `json:"project-description,omitempty"`
	Highlights  []string `json:"project-highlights,omitempty"`
}
type DevEnv struct {
	Name     string `json:"dev-name,omitempty"`
	Category string `json:"dev-category,omitempty"`
	URL      string `json:"dev-url,omitempty"`
}

type Hobby struct {
	Name string `json:"hobby-name,omitempty"`
}
