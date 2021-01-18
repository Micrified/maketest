package maketest

import (

	// Standard packages
	"errors"
	"encoding/json"

	// Custom packages
	"temporal"
	"gen"
	"types"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


// Encapsulates a test, and used to produce a test run script
type Test struct {
	Test_name           string    // The name of the test (name of results file)
	App_name            string    // Name of ROS program being built (from rules)
	App_rules           string    // Rules to build program with
	Is_custom_timing    bool      // Whether custom timing should be used
	App_timing          string    // Custom time parameter
	Is_duration         bool      // Whether a duration is set
	Duration_s          int       // Timeout to use (in seconds)
	Generate_directory  string    // Directory containing rosgraph
	Workspace_directory string    // Directory containing ROS workspace
	ROS_directory       string    // Directory containing ROS installation
	Analysis_directory  string    // Directory containing analysis program
	Results_directory   string    // Directory in which results should be put
	Logfile_directory   string    // Name of directory in which logs are placed
	Logfile_name        string    // Name of the logfile to use for analysis
}

// Describes all necessary directories for running the test
type Environment struct {
	Generate_directory  string
	Workspace_directory string
	ROS_directory       string
	Analysis_directory  string
	Results_directory   string
	Logfile_directory   string
	Logfile_name        string
}

/*
 *******************************************************************************
 *                        Testing Function Definitions                         *
 *******************************************************************************
*/

// Performs a test with the given name, using supplied rules and environment
// configuration
func Maketest (name, path string, rules types.Rules, is_custom_timing bool,
	timing []temporal.Temporal, environment Environment) error {
	var rules_data []byte
	var timing_data []byte
	var err error

	// Convert the rules to JSON
	rules_data, err = json.Marshal(rules)
	if nil != err {
		return errors.New("Unable to marshal rules: " + err.Error())
	}

	// Convert the timing to JSON
	if is_custom_timing {
		timing_data, err = json.Marshal(timing)
		if nil != err {
			return errors.New("Unable to marshal timing: " + err.Error())
		}
	}

	// Convert duration to seconds for use with timeout
	is_duration := (rules.Max_duration_us != -1)
	duration_s  := 1 + (rules.Max_duration_us / 1000000)

	// Build the test data structure
	test := Test{
		Test_name:           name,
		App_name:            rules.Name,
		App_rules:           string(rules_data),
		Is_custom_timing:    is_custom_timing,
		App_timing:          string(timing_data),
		Is_duration:         is_duration,
		Duration_s:          duration_s,
		Generate_directory:  environment.Generate_directory,
		Workspace_directory: environment.Workspace_directory,
		ROS_directory:       environment.ROS_directory,
		Analysis_directory:  environment.Analysis_directory,
		Results_directory:   environment.Results_directory,
		Logfile_directory:   environment.Logfile_directory,
		Logfile_name:        environment.Logfile_name,
	}

	return gen.GenerateTemplate(test, "/home/micrified/Go/src/maketest/templates/autotest.tmpl", path + "/" + name + ".sh")
}