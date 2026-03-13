package crosscutting

import "app/pkg/crosscutting/errors"

type RollbackFunc func() error

type RoollbackManager struct {
	funcs []RollbackFunc
}

func (manager RoollbackManager) NeedsRollback() bool {
	return len(manager.funcs) > 0
}

func (manager *RoollbackManager) Rollback(err error) error {
	for _, f := range manager.funcs {
		tmpErr := f()
		if tmpErr != nil {
			if err != nil {
				err = errors.New("rollback errors")
			}
			err = errors.AddSubErr(err, tmpErr)
		}
	}
	return err
}

func (manager *RoollbackManager) AddRollbackFunc(f RollbackFunc) {
	manager.funcs = append(manager.funcs, f)
}
