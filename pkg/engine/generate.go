package engine

import (
	"context"
	"errors"
)

// ErrFailedGeneration is returned when at least one file couldn't be properly generated.
//
// Every generation error is logged during processing to avoid a big aggregated error at the end.
var ErrFailedGeneration = errors.New("some error(s) occurred during generation")

// Generate is the main function from generate package.
// It takes a configuration and various options.
//
// It executes all parsers given in options (or default ones)
// and then iterates over provided templates to apply or remove those.
func Generate[T any](ctx context.Context, destdir string, config T, parsers []Parser[T], generators []Generator[T]) (T, error) {
	// parse repository
	errs := make([]error, 0, len(parsers))
	for _, parser := range parsers {
		errs = append(errs, parser(ctx, destdir, &config))
	}
	if err := errors.Join(errs...); err != nil {
		return config, err
	}

	// execute generators
	var errcount int
	for _, generator := range generators {
		if err := generator(ctx, destdir, config); err != nil {
			if !errors.Is(err, ErrFailedGeneration) {
				GetLogger().Errorf(err.Error())
			}
			errcount++
		}
	}
	if errcount > 0 {
		return config, ErrFailedGeneration
	}
	return config, nil
}
