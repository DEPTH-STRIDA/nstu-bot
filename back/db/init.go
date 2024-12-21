package db

var App *DataBase

func Init() error {
	var err error

	App, err = NewDatabase()
	if err != nil {
		return err
	}

	App.Migrate()

	return nil
}
