package main

import "fmt"

func makeup(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}

	return fmt.Sprintf(format, args...)
}

func (b *App) infof(format string, args ...interface{}) {
	b.logger.Info(makeup(format, args...))
}

func (b *App) errorf(format string, args ...interface{}) {
	b.logger.Info(makeup(format, args...))
}

func (b *App) debugf(format string, args ...interface{}) {
	b.logger.Info(makeup(format, args...))
}
