// Code generated by generator, DO NOT EDIT.
package models

func (t *TodoItem) GetId() int {
	return t.Id
}

func (t *TodoItem) SetId(newId int) {
	t.Id = newId
}

func (t *TodoItem) GetTitle() string {
	return t.Title
}

func (t *TodoItem) SetTitle(newTitle string) {
	t.Title = newTitle
}

func (t *TodoItem) GetDescription() string {
	return t.Description
}

func (t *TodoItem) SetDescription(newDescription string) {
	t.Description = newDescription
}

func (t *TodoItem) GetStatus() string {
	return t.Status
}

func (t *TodoItem) SetStatus(newStatus string) {
	t.Status = newStatus
}

func (t *TodoItem) GetDone() bool {
	return t.Done
}

func (t *TodoItem) SetDone(newDone bool) {
	t.Done = newDone
}
