package crosscutting

var globalRandStrIDGenerator StrIDGenerator

func SetGlobalRandStrIDGenerator(gen StrIDGenerator) {
	globalRandStrIDGenerator = gen
}

func IssueRandomStrID() string {
	return globalRandStrIDGenerator.Issue()
}

type StrIDGenerator interface {
	Issue() string
}
