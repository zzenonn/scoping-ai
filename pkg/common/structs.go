package common

// If a question is multiple choice, is it checkbox or radio button?
type Options struct {
	MultiAnswer     bool     `json:"multi_answer,omitempty" firestore:"multi_answer,omitempty"`
	PossibleOptions []string `json:"possible_options,omitempty" firestore:"possible_options,omitempty"`
}

// Question representation
type Question struct {
	Category *string  `json:"category,omitempty" firestore:"category,omitempty"`
	Text     *string  `json:"text,omitempty" firestore:"text,omitempty"`
	Options  *Options `json:"options,omitempty" firestore:"options,omitempty"`
}
