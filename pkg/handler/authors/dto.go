package authors

type ApiAuthor struct {
	ID   int64
	Name string `json:"name,omitempty" binding:"required,max=32"`
	Bio  string `json:"bio,omitempty" binding:"required"`
}

type ApiAuthorPartialUpdate struct {
	Name *string `json:"name,omitempty" binding:"omitempty,max=32"`
	Bio  *string `json:"bio,omitempty" binding:"omitempty"`
}

type PathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}
