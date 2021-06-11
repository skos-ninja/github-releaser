package app

import "context"

func (a *app) GetZen(ctx context.Context) (string, error) {
	zen, _, err := a.github.Zen(ctx)
	if err != nil {
		return "", err
	}

	return zen, nil
}
