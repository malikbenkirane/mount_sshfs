package ssh

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

var debug = false

var (
	HostKeyNotFound = errors.New("Host key not found")
	InvalidRemote   = errors.New("Invalid remote")
)

type Connection struct {
	Client    *ssh.Client
	Directory string
}

// NewConnection establish a connection to the remote.
// The remote format is "user@host:dir".
// You will need to close the client.
func NewConnection(remote string) (*Connection, error) {
	remoteFields := strings.Split(remote, "@")
	if len(remoteFields) != 2 {
		return nil, InvalidRemote
	}
	remoteFields = append(strings.Split(remoteFields[1], ":"), remoteFields...)
	if debug {
		log.Printf("[rhost rdir ruser rhost:rdir] %v", remoteFields)
	}
	if len(remoteFields) != 4 {
		return nil, InvalidRemote
	}
	connection := &Connection{}
	host := remoteFields[0]
	connection.Directory = remoteFields[1]
	remoteUser := remoteFields[2]
	var err error
	hkCallback, err := GetHostKey(host)
	if err != nil {
		return nil, err
	}
	localUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	identity := filepath.Join(localUser.HomeDir, ".ssh", "id_rsa")
	log.Printf("Reading identity from %q", identity)
	key, err := ioutil.ReadFile(identity)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hkCallback,
	}
	if debug {
		log.Printf("%+v", config)
	}
	host = host + ":22"
	connection.Client, err = ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %q: %v", host, err)
	}
	return connection, nil
}

// GetHostKey return the host key parsed with ssh.ParseAuthorizedKey
// See also
// - golang.org/x/crypto/ssh documentation
// - https://networkbit.ch/golang-ssh-client/
func GetHostKey(host string) (ssh.HostKeyCallback, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	file := filepath.Join(user.HomeDir, ".ssh", "known_hosts")
	return knownhosts.New(file)
}
