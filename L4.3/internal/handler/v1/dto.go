package v1

// CreateRequestV1 represents the request body for creating a new event.
type CreateRequestV1 struct {
	UserID    int    `json:"user_id" example:"1"`                  // UserID is the ID of the user who owns the event.
	EventDate string `json:"date" example:"2028-12-04"`            // EventDate is the date of the event in YYYY-MM-DD format.
	Text      string `json:"text,omitempty" example:"Touch grass"` // Text is the optional description of the event.
}

// CreateResponseV1 represents the response returned after creating an event.
type CreateResponseV1 struct {
	EventID string `json:"event_id" example:"3383503d-fb71-4b8c-85bd-a914c84252a9"` // EventID is the unique identifier of the newly created event.
}

// UpdateRequestV1 represents the request body for updating an existing event.
type UpdateRequestV1 struct {
	UserID  int    `json:"user_id" example:"1"`                                     // UserID is the ID of the user who owns the event.
	EventID string `json:"event_id" example:"3383503d-fb71-4b8c-85bd-a914c84252a9"` // EventID is the unique identifier of the event to update.
	Text    string `json:"text,omitempty" example:"Grind leetcode"`                 // Text is the new optional description for the event.
	NewDate string `json:"new_date,omitempty" example:"2028-12-05"`                 // NewDate is the new optional date for the event in YYYY-MM-DD format.
}

// UpdateResponseV1 represents the response returned after updating an event.
type UpdateResponseV1 struct {
	Updated bool `json:"event_updated" example:"true"` // Updated indicates whether the event was successfully updated.
}

// DeleteRequestV1 represents the request body for deleting an existing event.
type DeleteRequestV1 struct {
	UserID  int    `json:"user_id" binding:"required" example:"1"`                                     // UserID is the ID of the user who owns the event.
	EventID string `json:"event_id" binding:"required" example:"3383503d-fb71-4b8c-85bd-a914c84252a9"` // EventID is the unique identifier of the event to delete.
}

// DeleteResponseV1 represents the response returned after deleting an event.
type DeleteResponseV1 struct {
	Deleted bool `json:"event_deleted" example:"true"` // Deleted indicates whether the event was successfully deleted.
}

// EventDtoV1 represents an event in responses containing event info.
type EventDtoV1 struct {
	Text      string `json:"text" example:"Touch grass"`                              // Text is the description of the event.
	EventDate string `json:"date" example:"2028-12-04"`                               // EventDate is the date of the event in YYYY-MM-DD format.
	EventID   string `json:"event_id" example:"3383503d-fb71-4b8c-85bd-a914c84252a9"` // EventID is the unique identifier of the event.
}

// ListOfEventsResponseV1 represents a response containing a list of events.
type ListOfEventsResponseV1 struct {
	Events []EventDtoV1 `json:"events"` // Events is the list of events returned by the API.
}

// ErrorResponse represents a standard bad request response.
type ErrorResponse400 struct {
	Code    int    `json:"code" example:"400"`                                         // Code is the HTTP status code.
	Message string `json:"message" example:"invalid date format, expected YYYY-MM-DD"` // Message is a human-readable description of the error.
}

// ErrorResponse represents a standard internal error response.
type ErrorResponse500 struct {
	Code    int    `json:"code" example:"500"`                      // Code is the HTTP status code.
	Message string `json:"message" example:"internal server error"` // Message is a human-readable description of the error.
}
