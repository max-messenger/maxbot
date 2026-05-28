package middleware

import (
	"errors"

	"github.com/max-messenger/maxbot"
)

type RecoverFunc = func(error, maxbot.Context)

func Recover(onError ...RecoverFunc) maxbot.MiddlewareFunc {
	return func(next maxbot.HandlerFunc) maxbot.HandlerFunc {
		return func(c maxbot.Context) error {
			var f RecoverFunc
			if len(onError) > 0 {
				f = onError[0]
			}
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						f(err, c)
					} else if s, ok := r.(string); ok {
						f(errors.New(s), c)
					}
				}
			}()

			return next(c)
		}
	}
}
