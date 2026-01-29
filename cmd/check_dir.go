package cmd

func CheckDir() (bool, error) {
	envFiles, err := Scan("", "")
	if err != nil {
		return false, err
	}
	if len(envFiles) == 0 {
		return false, nil
	}
	validEnvs, err := Validate(envFiles)
	if err != nil {
		return false, err
	}
	return len(validEnvs) != 0, nil
}
