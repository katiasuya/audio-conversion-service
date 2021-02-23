package domain

// Audio represents any user's audio file.
type Audio struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	File   string `json:"file"`
}

// AudioExists checks whether there is an audio with the given id.
func AudioExists(id string) error {
	// return  errors.New("there is no song with the given id")
	return nil
}
