package yarn

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

// If there is a open ticket, does nothing, else kinits a new session
func kCheck(username, keytab string) error {
	if err := exec.Command("klist", "-s").Run(); err != nil { // returns no err if there's an open ticket
		if err = kInit(username, keytab); err != nil {
			return err
		}
	}

	return nil
}

func kInit(username, keytab string) error {

	if len(username) < 1 {
		return errors.New("A username is required to authenticate with Kerberos.")
	}
	if len(keytab) < 0 {
		return errors.New("A keytab is required to authenticate with Kerberos.")
	}

	cmd := exec.Command("kinit", username, "-k", "-t", keytab)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	// Scanning for errors while cmd is running
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			log.Printf(in.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Println("Failed to authenticate.")
		return err
	}

	return nil
}

// StartCloud starts a yarn cloud by shelling out to hadoop
//
// This process needs to store the job-ID to kill the process in the future
func StartCloud(size int, kerberos bool, name, username, keytab string) (string, error) {
	if kerberos {
		if err := kCheck(username, keytab); err != nil {
			return "", err
		}
	}

	cmdArgs := []string{
		"jar",              //
		"h2odriver.jar",    //FIXME: This should be a pack method
		"-jobname",         //
		"H2O_" + name,      //
		"-n",               //
		strconv.Itoa(size), //
		"-mapperXmx",       //
		"10g",              // FIXME: This may be modifialbe down the road
		"-output",          //
		name + "_out",      //
		"-disown",          //
	}

	log.Println("Attempting to start cloud...")
	cmdOut, err := exec.Command("hadoop", cmdArgs...).CombinedOutput()

	if err != nil {
		log.Println("Failed to launch hadoop.")
		log.Println("\n" + string(cmdOut)) // This captures error from the drive.jar
		return "", err                     // This captures erros from Stderr
	}
	hpOut := (string(cmdOut))
	// Capture only the address and ID respectively
	reNode := regexp.MustCompile(`H2O node (\d+\.\d+\.\d+\.\d+:\d+)`)
	reApID := regexp.MustCompile(`application_(\d+_\d+)`)

	for _, node := range reNode.FindAllStringSubmatch(hpOut, size) {
		address := node[1]
		log.Println("Node started at:", address)
	}
	apID := reApID.FindStringSubmatch(hpOut)[1]

	fmt.Println("")
	log.Println("Started cloud with ID:", apID)

	return apID, nil
}

// StopCloud kills a hadoop cloud by shelling out a command based on the job-ID
//
// In the future this
func StopCloud(kerberos bool, name, id, username, keytab string) error {
	if kerberos {
		if err := kCheck(username, keytab); err != nil {
			return err
		}
	}

	log.Println("Attempting to stop cloud...")
	cmdStop := exec.Command("hadoop", "job", "-kill", "job_"+id)
	if out, err := cmdStop.CombinedOutput(); err != nil {
		log.Println("Failed to shutdown hadoop.")
		log.Println("\n" + string(out))
		return err
	}
	cmdClean := exec.Command("hadoop", "fs", "-rmdir", name+"_out")
	log.Println("Stopped cloud:", "job_"+id)
	if out, err := cmdClean.Output(); err != nil {
		log.Fatalln("Failed to remove outdir.")
		log.Println("\n" + string(out))
		return err
	}

	return nil
}
