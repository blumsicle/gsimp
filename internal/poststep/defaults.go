package poststep

func DefaultPostSteps() []PostStep {
	return []PostStep{
		GoGetUpdatePostStep{},
		GoModTidyPostStep{},
		GitInitPostStep{},
	}
}
