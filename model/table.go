package model

func (table *ClassTableModel) New() error {
	d := DB.Self.Create(table)
	return d.Error
}

func (table *ClassTableModel) Delete() error {
	d := DB.Self.Delete(table)
	return d.Error
}

func (table *ClassTableModel) Rename() error {
	return nil
}
