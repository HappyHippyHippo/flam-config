package config

const (
	providerId = "flam.config.provider"

	ParserCreatorGroup         = "flam.config.parsers.creator"
	ParserDriverYaml           = "flam.config.parsers.driver.yaml"
	ParserDriverJson           = "flam.config.parsers.driver.json"
	SourceCreatorGroup         = "flam.config.sources.creator"
	SourceDriverEnv            = "flam.config.sources.driver.env"
	SourceDriverFile           = "flam.config.sources.driver.file"
	SourceDriverObservableFile = "flam.config.sources.driver.observable-file"
	SourceDriverDir            = "flam.config.sources.driver.dir"
	SourceDriverRest           = "flam.config.sources.driver.rest"
	SourceDriverObservableRest = "flam.config.sources.driver.observable-rest"

	PathDefaultFileParser = "flam.config.defaults.file.parser"
	PathDefaultFileDisk   = "flam.config.defaults.file.disk"
	PathDefaultRestParser = "flam.config.defaults.rest.parser"
	PathBoot              = "flam.config.boot"
	PathObserverFrequency = "flam.config.observer.frequency"
	PathParsers           = "flam.config.parsers"
	PathSources           = "flam.config.sources"
)
