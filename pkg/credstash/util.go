package credstash

import (
	"errors"
	"fmt"
	"strconv"
)

func formatCredstashVersion(inputVersion string) (string, error){
	_, err := strconv.Atoi(inputVersion)
	if err != nil {
		log.Error(err, "Could not parse credstash version into number",
			"Secret.Version", inputVersion)
		return "", err
	}

	// we already have a padded version so nothing to do
	if len(inputVersion) == credstashVersionLength {
		return inputVersion, nil
	}

	// version is too longßß
	if len(inputVersion) > credstashVersionLength {
		return "", errors.New(
			fmt.Sprintf("Version string is longer than supported. Maximum length is %d characters",
				credstashVersionLength))
	}

	// pad version with leading zeros until we reach credstashVersionLength
	// format becomes something like %019s which means pad the string until there's 19 0s
	format := fmt.Sprintf("%s%ds","%0", credstashVersionLength)
	newVersion := fmt.Sprintf(format, inputVersion)
	return newVersion, nil
}