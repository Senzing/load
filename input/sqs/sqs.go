package sqs

import (
	"context"
	"fmt"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-queueing/queues/sqs"
	"github.com/senzing/go-sdk-abstract-factory/factory"
)

// ----------------------------------------------------------------------------

// read and process records from the given queue until a system interrupt
func Read(ctx context.Context, urlString, engineConfigJson string, engineLogLevel, numberOfWorkers, visibilityPeriodInSeconds int, logLevel string, jsonOutput bool) {

	jsonOutput = jsonOutput
	logger = getLogger()
	setLogLevel(ctx, logLevel)

	// Work with G2engine.
	g2engine := createG2Engine(ctx, engineConfigJson, engineLogLevel)
	defer g2engine.Destroy(ctx)

	startErr := sqs.StartManagedConsumer(ctx, urlString, numberOfWorkers, g2engine, false, int32(visibilityPeriodInSeconds), logLevel, jsonOutput)
	if startErr != nil {
		log(5000, startErr.Error())
	}
	log(2999)
}

// ----------------------------------------------------------------------------

// create a G2Engine object, on error this function panics.
// see failOnError
func createG2Engine(ctx context.Context, engineConfigJson string, engineLogLevel int) g2api.G2engine {
	senzingFactory := &factory.SdkAbstractFactoryImpl{}
	g2Config, err := senzingFactory.GetG2config(ctx)
	if err != nil {
		log(2004, err.Error())
	}
	g2engine, err := senzingFactory.GetG2engine(ctx)

	if err != nil {
		log(2005, err.Error())
	}
	if g2Config.GetSdkId(ctx) == "base" {
		err = g2engine.Init(ctx, "load", engineConfigJson, engineLogLevel)
		if err != nil {
			log(2006, err.Error())
		}
	}
	return g2engine
}

var logger logging.LoggingInterface = nil
var jsonOutput bool = false

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func getLogger() logging.LoggingInterface {
	var err error = nil
	if logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		logger, err = logging.NewSenzingToolsLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return logger
}

// Log message.
func log(messageNumber int, details ...interface{}) {
	if jsonOutput {
		getLogger().Log(messageNumber, details...)
	} else {
		fmt.Println(fmt.Sprintf(IDMessages[messageNumber], details...))
	}
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func setLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = getLogger().SetLogLevel(logLevelName)
	return err
}
