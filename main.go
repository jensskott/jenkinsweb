package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/tsuru/config"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Takes arguments from commandline to pass in config file for yml
	// If no argument is presented it uses config.yml
	// Takes a bool to see if a new job is going to be created
	// Takes a xml file for jenkins job config
	// If no xml file is defined it used job.xml
	var cfgfile string
	var jobName string
	var xmlfile string
	createJob := flag.Bool("new-job", false, "Create a new job?")
	flag.StringVar(&cfgfile, "config-file", "config.yml", "Your config file")
	flag.StringVar(&jobName, "job-name", "", "Your new jobs name")
	flag.StringVar(&xmlfile, "xml-file", "job.xml", "Your xml file for jobs")
	flag.Parse()
	jobxml, err := ioutil.ReadFile(xmlfile)
	check(err)

	// Check if a new job is created and has a name
	if *createJob == true && len(strings.TrimSpace(jobName)) == 0 {
		log.Fatal("No job name provided")
	}

	// Read config file from flag input
	config.ReadConfigFile(cfgfile)
	jenkinsURL, err := config.GetString("jenkins:url")
	jenkinsUser, err := config.GetString("jenkins:user")
	jenkinsPasswd, err := config.GetString("jenkins:password")
	check(err)

	//Connect to jenkins
	jenkins, err := gojenkins.CreateJenkins(jenkinsURL, jenkinsUser, jenkinsPasswd).Init()
	check(err)

	// Create a jenkins job is flag is true and job doesnt exist
	if *createJob == true {
		job, _ := jenkins.GetJob(jobName)
		if job != nil {
			log.Fatal("Job allready exists")
		}
		jenkins.CreateJob(string(jobxml), jobName)
		fmt.Println("Jenkins job", jobName, "created!")
	}
}
