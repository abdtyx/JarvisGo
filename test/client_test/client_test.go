package client_test

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

// Initialize a server listening on 5700, maintain a message queue

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!

	// These declarations may be useful for multi-session testing.

	// var err error

	// A bunch of filenames that may be useful.
	// aliceFile := "aliceFile.txt"
	// aliceFile2 := "aliceFile2.jpg"
	// bobFile := "bobFile.txt"
	// bobFile2 := "bobFile2.png"
	// charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"
	// jimFile := "jimFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		// userlib.DatastoreClear()
		// userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			// userlib.DebugMsg("Initializing user Alice.")
			// alice, err = client.InitUser("alice", defaultPassword)
			// Expect(err).To(BeNil())

			// userlib.DebugMsg("Getting user Alice.")
			// aliceLaptop, err = client.GetUser("alice", defaultPassword)
			// Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			// userlib.DebugMsg("Initializing user Alice.")
			// alice, err = client.InitUser("alice", defaultPassword)
			// Expect(err).To(BeNil())

			// userlib.DebugMsg("Storing file data: %s", contentOne)
			// err = alice.StoreFile(aliceFile, []byte(contentOne))
			// Expect(err).To(BeNil())

			// userlib.DebugMsg("Appending file data: %s", contentTwo)
			// err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			// Expect(err).To(BeNil())

			// userlib.DebugMsg("Appending file data: %s", contentThree)
			// err = alice.AppendToFile(aliceFile, []byte(contentThree))
			// Expect(err).To(BeNil())

			// userlib.DebugMsg("Loading file...")
			// data, err := alice.LoadFile(aliceFile)
			// Expect(err).To(BeNil())
			// Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

	})
})
