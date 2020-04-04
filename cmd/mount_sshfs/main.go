package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/malikbenkirane/mount_sshfs/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type sshfs struct {
	UID         int    `yaml:"uid"`
	GID         int    `yaml:"gid"`
	IsRoot      bool   `yaml:"root"`
	IsForDocker bool   `yaml:"docker"`
	MountDir    string `yaml:"mount"`
	Remote      string `yaml:"remote"`
}

var (
	// ErrRootWallBroken is the error value returned when
	// the sshfs.IsRoot is false and either sshfs.UID is 0 or sshf.GID is 0
	ErrRootWallBroken = errors.New("uid is 0 or gid is 0 and -root flag is not set")
	// ErrInvalidMountDirPath is the error value returned when
	// the mount dir path passed is invalid
	ErrInvalidMountDirPath = errors.New("invalid mount dir path")
)

// reference: https://stackoverflow.com/a/35240286
func isValidPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	var d []byte
	if err := ioutil.WriteFile(path, d, 0644); err == nil {
		os.Remove(path)
		return true
	}
	return false
}

func readConfig(filename string) (*sshfs, error) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var options *sshfs
	if err := yaml.Unmarshal(config, &options); err != nil {
		return nil, err
	}
	if err := validOptions(options); err != nil {
		return nil, err
	}
	return options, nil
}

func validOptions(options *sshfs) error {
	if !options.IsRoot && (options.UID == 0 || options.GID == 0) {
		return ErrRootWallBroken
	}
	if !isValidPath(options.MountDir) {
		return ErrInvalidMountDirPath
	}
	if err := isValidRemote(options.Remote); err != nil {
		return fmt.Errorf("Unable to verify remote: %v", err)
	}
	return nil
}

func isValidRemote(remote string) error {
	var err error
	connection, err := ssh.NewConnection(remote)
	if err != nil {
		return err
	}
	remoteDir := connection.Directory
	defer connection.Client.Close()
	session, err := connection.Client.NewSession()
	if err != nil {
		return err
	}
	session.Stdout = os.Stderr
	session.Stderr = os.Stderr
	cmd := fmt.Sprintf("stat %q", remoteDir)
	err = session.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}

// ConfigError return the error value passed prefixed with "Configureation error: "
func ConfigError(err error) error {
	return fmt.Errorf("Configuration error: %v", err)
}

func flags() (*sshfs, error) {
	isForDocker := flag.Bool("docker", false, "sshfs with docker")
	mountDir := flag.String("mount-to", "", "mount directory /mnt/{directory}")
	remoteDir := flag.String("remote-dir", "", "path to the directory on the remote")
	host := flag.String("remote-host", "", "remote host [username@]remote")
	uid := flag.Int("uid", 0, "uid for sshfs -o idmap")
	gid := flag.Int("gid", 0, "gid for sshfs -o idmap")
	isRoot := flag.Bool("root", false, "enable root idmap")
	config := flag.String("config", "", "configuration file")
	flag.Parse()
	if *config != "" {
		log.Printf("reading configuration from %s\n", *config)
		options, err := readConfig(*config)
		if err != nil {
			return nil, ConfigError(err)
		}
		return options, nil
	}
	remote := fmt.Sprintf("%s:%s", *host, filepath.Clean(*remoteDir))
	options := &sshfs{
		UID:         *uid,
		GID:         *gid,
		IsForDocker: *isForDocker,
		MountDir:    *mountDir,
		Remote:      remote,
		IsRoot:      *isRoot,
	}
	if err := validOptions(options); err != nil {
		return nil, err
	}
	return options, nil
}

func main() {
	log.SetFlags(0)
	flags, err := flags()
	if err != nil {
		log.Fatalf("Invalid flags: %v", err)
	}
	cmd := `sudo sshfs \
	-o idmap=user,uid=%d,gid=%d`
	if flags.IsForDocker {
		cmd += ` \
	-o allow_other`
	}
	cmd += ` \
	%s \
	%s
`
	fmt.Printf(cmd,
		flags.UID, flags.GID,
		flags.Remote,
		flags.MountDir)
}
