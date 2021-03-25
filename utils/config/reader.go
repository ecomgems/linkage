package config

// Function GetConfiguration reads the configuration from
// the file and converts it to the Configuration object.
func GetConfiguration(fileName string) (Configuration, error) {
	var (
		config  Configuration
		content []byte
		err     error
	)

	content, err = getContentFromFileName(fileName)
	if err != nil {
		return Configuration{}, err
	}

	config, err = getConfigurationFromContent(content)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}

// Function getContentFromFileName open file by it's name and
// reads content into slice of bytes.
func getContentFromFileName(fileName string) ([]byte, error) {
	return []byte{}, nil
}

// Function getConfigurationFromContent converts slice of bytes
// into Configuration object.
func getConfigurationFromContent(content []byte) (Configuration, error) {
	return Configuration{}, nil
}
