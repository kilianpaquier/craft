package initialize

import (
	"bufio"
	"context"
	"io"
	"os"
	"slices"
	"strconv"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

var reader io.Reader = os.Stdin

// Run inits craft configuration with user inputs from terminal.
func Run(ctx context.Context) (craft models.CraftConfig) {
	log := logrus.WithContext(ctx)
	log.Info("starting craft repository configuration")

	scanner := bufio.NewScanner(reader)

	// read main maintainer information
	craft.Maintainers = append(craft.Maintainers, readMaintainer(ctx, scanner))

	// read api configuration
	craft.API = readAPI(ctx, scanner)

	// Helm chart generation
	craft.NoChart = !readChart(ctx, scanner)

	return craft
}

// readMaintainer reads from input scanner the main maintainer with Q&A method.
func readMaintainer(ctx context.Context, scanner *bufio.Scanner) (maintainer models.Maintainer) {
	// main maintainer name
	for {
		name := ask(ctx, scanner, "Who's the maintainer name (it can be a group name, anything) ?")
		if name == nil {
			logrus.Warn("maintainer name is mandatory")
			continue
		}
		maintainer.Name = *name
		break
	}

	// main maintainer email
	maintainer.Email = ask(ctx, scanner, "Who's the maintainer email (optional) ?")

	// main maintainer url
	maintainer.URL = ask(ctx, scanner, "Who's the maintainer url (optional) ?")

	return maintainer
}

// readAPI reads user inputs for API generation with version (Q&A method).
func readAPI(ctx context.Context, scanner *bufio.Scanner) *models.API {
	// API generation
	for {
		api := ask(ctx, scanner, "Would you like to generate a golang based API (optional, default is falsy) 0/1 ?")
		if api == nil {
			break // response not provided, going through next questions
		}

		value, err := strconv.ParseBool(*api)
		if err != nil {
			logrus.WithError(err).Warn("invalid api value, must be a boolean")
			continue
		}
		if !value {
			return nil // no api is wanted
		}
		break
	}

	// API version
	for {
		version := ask(ctx, scanner, "Would you like to specify an OpenAPI version for your API (optional, default is 'v2') v2/v3 ?")
		if version == nil {
			return &models.API{OpenAPIVersion: lo.ToPtr("v2")} // response not provided, going with default openapi version
		}

		if !slices.Contains([]string{"", "v2", "v3"}, *version) {
			logrus.Warn("openapi version must be either 'v2' or 'v3'")
			continue
		}
		return &models.API{OpenAPIVersion: version} // api is wanted with valid version
	}
}

// readChart reads from input scanner the answers related to chart generation.
func readChart(ctx context.Context, scanner *bufio.Scanner) bool {
	// Helm chart generation
	for {
		chart := ask(ctx, scanner, "Would you like to generate an Helm chart (optional, default is truthy) 0/1 ?")
		if chart == nil {
			return true // no response provided, going through with chart activated
		}

		value, err := strconv.ParseBool(*chart)
		if err != nil {
			logrus.WithError(err).Warn("invalid chart answer, must be a boolean")
			continue
		}
		return value
	}
}

// ask asks a question and retrieves answer from input scanner.
func ask(ctx context.Context, scanner *bufio.Scanner, question string, args ...any) *string {
	logrus.WithContext(ctx).Infof(question, args...)
	scanner.Scan()
	if answer := scanner.Text(); answer != "" {
		return &answer
	}
	return nil
}
