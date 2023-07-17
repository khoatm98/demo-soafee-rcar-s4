package renesasrefappota_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/aoscloud/aos_common/aoserrors"
	log "github.com/sirupsen/logrus"
	"github.com/syucream/posix_mq"

	"github.com/khoatm98/demo-soafee-rcar-s4/updatemodules/renesasrefappota"
)

/***********************************************************************************************************************
 * Consts
 **********************************************************************************************************************/

const (
	commandQueue = "/ota_master_queue"
	statusQueue  = "/ota_mater_result"

	queueTimeout = 1 * time.Second
)

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type testOtaMaster struct {
	sync.Mutex

	sendMQ       *posix_mq.MessageQueue
	recvMQ       *posix_mq.MessageQueue
	recvCommands []int64
	statusMap    map[int64]int64
}

type testStateStorage struct {
	state []byte
}

/***********************************************************************************************************************
 * Vars
 **********************************************************************************************************************/

var tmpDir string

/***********************************************************************************************************************
 * Init
 **********************************************************************************************************************/

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05.000",
		FullTimestamp:    true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

/*******************************************************************************
 * Main
 ******************************************************************************/

func TestMain(m *testing.M) {
	var err error

	tmpDir, err = ioutil.TempDir("", "um_")
	if err != nil {
		log.Fatalf("Error create tmp dir: %v", err)
	}

	ret := m.Run()

	if err = os.RemoveAll(tmpDir); err != nil {
		log.Fatalf("Error deleting tmp dir: %v", err)
	}

	os.Exit(ret)
}

/***********************************************************************************************************************
 * Tests
 **********************************************************************************************************************/

func TestUpdate(t *testing.T) {
	master, err := newTestOtaMaster(statusQueue, commandQueue,
		map[int64]int64{0: 0, 1: 0, 2: 0, 3: 0, 4: 0})
	if err != nil {
		t.Fatalf("Can't create test OTA master: %v", err)
	}
	defer master.close()

	targetFile := filepath.Join(tmpDir, "target.dat")

	module, err := renesasrefappota.New("test", moduleConfig(targetFile), &testStateStorage{})
	if err != nil {
		t.Fatalf("Can't create test module: %v", err)
	}
	defer module.Close()

	if id := module.GetID(); id != "test" {
		t.Errorf("Wrong module id: %s", id)
	}

	const updateVersion = "2.1.0"

	// Init

	if err = module.Init(); err != nil {
		t.Fatalf("Error init module: %v", err)
	}

	// Prepare

	const imageContent = "this is image content"

	imageFile := filepath.Join(tmpDir, "image.dat")

	if err = createImage(imageFile, imageContent); err != nil {
		t.Fatalf("Error create image: %v", err)
	}

	if err = module.Prepare(imageFile, updateVersion, nil); err != nil {
		t.Fatalf("Error prepare module: %v", err)
	}

	readContent, err := ioutil.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("Error read content: %s", err)
	}

	if imageContent != string(readContent) {
		t.Error("Wrong image content")
	}

	if !reflect.DeepEqual(master.getRecvCommands(), []int64{0, 1}) {
		t.Error("Wrong commands received")
	}

	// Update

	if _, err = module.Update(); err != nil {
		t.Errorf("Error update module: %v", err)
	}

	if !reflect.DeepEqual(master.getRecvCommands(), []int64{2, 3}) {
		t.Error("Wrong commands received")
	}

	version, err := module.GetVendorVersion()
	if err != nil {
		t.Errorf("Can't get vendor version: %v", err)
	}

	if version != updateVersion {
		t.Errorf("Wrong vendor version: %s", version)
	}

	// Apply

	if _, err = module.Apply(); err != nil {
		t.Errorf("Error apply module: %v", err)
	}

	if master.getRecvCommands() != nil {
		t.Error("Wrong commands received")
	}

	if version, err = module.GetVendorVersion(); err != nil {
		t.Errorf("Can't get vendor version: %v", err)
	}

	if version != updateVersion {
		t.Errorf("Wrong vendor version: %s", version)
	}
}

func TestRevert(t *testing.T) {
	master, err := newTestOtaMaster(statusQueue, commandQueue,
		map[int64]int64{0: 0, 1: 0, 2: 1, 3: 0, 4: 0})
	if err != nil {
		t.Fatalf("Can't create test OTA master: %v", err)
	}
	defer master.close()

	module, err := renesasrefappota.New(
		"test", moduleConfig(filepath.Join(tmpDir, "target.dat")), &testStateStorage{})
	if err != nil {
		t.Fatalf("Can't create test module: %v", err)
	}
	defer module.Close()

	if id := module.GetID(); id != "test" {
		t.Errorf("Wrong module id: %s", id)
	}

	const updateVersion = "2.1.0"

	// Init

	if err = module.Init(); err != nil {
		t.Fatalf("Error init module: %v", err)
	}

	// Prepare

	imageFile := filepath.Join(tmpDir, "image.dat")

	if err = createImage(imageFile, "Some image content"); err != nil {
		t.Fatalf("Error create image: %v", err)
	}

	if err = module.Prepare(imageFile, updateVersion, nil); err != nil {
		t.Fatalf("Error prepare module: %v", err)
	}

	if !reflect.DeepEqual(master.getRecvCommands(), []int64{0, 1}) {
		t.Error("Wrong commands received")
	}

	// Update

	if _, err = module.Update(); err == nil {
		t.Error("Update should fail")
	}

	if !reflect.DeepEqual(master.getRecvCommands(), []int64{2}) {
		t.Error("Wrong commands received")
	}

	version, err := module.GetVendorVersion()
	if err != nil {
		t.Errorf("Can't get vendor version: %v", err)
	}

	if version != "" {
		t.Errorf("Wrong vendor version: %s", version)
	}

	// Revert

	if _, err = module.Revert(); err != nil {
		t.Errorf("Error revert module: %v", err)
	}

	if !reflect.DeepEqual(master.getRecvCommands(), []int64{4}) {
		t.Error("Wrong commands received")
	}

	if version, err = module.GetVendorVersion(); err != nil {
		t.Errorf("Can't get vendor version: %v", err)
	}

	if version != "" {
		t.Errorf("Wrong vendor version: %s", version)
	}
}

/***********************************************************************************************************************
 * testOtaMaster
 **********************************************************************************************************************/

func newTestOtaMaster(sendQueue, receiveQueue string, statusMap map[int64]int64) (master *testOtaMaster, err error) {
	localMaster := &testOtaMaster{
		statusMap: statusMap,
	}

	defer func() {
		if err != nil {
			localMaster.close()
		}
	}()

	if localMaster.sendMQ, err = posix_mq.NewMessageQueue(
		sendQueue, posix_mq.O_CREAT|posix_mq.O_WRONLY, 0o600, nil); err != nil {
		return nil, aoserrors.Wrap(err)
	}

	if localMaster.recvMQ, err = posix_mq.NewMessageQueue(
		receiveQueue, posix_mq.O_CREAT|posix_mq.O_RDONLY, 0o600, nil); err != nil {
		return nil, aoserrors.Wrap(err)
	}

	go func() {
		for {
			data, _, err := localMaster.recvMQ.Receive()
			if err != nil {
				log.Errorf("Receive message error: %v", err)

				return
			}

			buffer := bytes.NewBuffer(data)

			var command int64

			if err = binary.Read(buffer, binary.LittleEndian, &command); err != nil {
				log.Errorf("Read message error: %v", err)
			}

			localMaster.Lock()
			localMaster.recvCommands = append(localMaster.recvCommands, command)
			localMaster.Unlock()

			status, ok := localMaster.statusMap[command]
			if !ok {
				continue
			}

			buffer = bytes.NewBuffer(nil)

			if err = binary.Write(buffer, binary.LittleEndian, status); err != nil {
				log.Errorf("Write message error: %v", err)
			}

			if err = localMaster.sendMQ.Send(buffer.Bytes(), 0); err != nil {
				log.Errorf("Send message error: %v", err)
			}
		}
	}()

	return localMaster, nil
}

func (master *testOtaMaster) getRecvCommands() []int64 {
	master.Lock()
	defer master.Unlock()

	recvCommands := master.recvCommands

	master.recvCommands = nil

	return recvCommands
}

func (master *testOtaMaster) close() {
	if master.recvMQ != nil {
		_ = master.recvMQ.Unlink()
	}

	if master.sendMQ != nil {
		_ = master.sendMQ.Unlink()
	}
}

/***********************************************************************************************************************
 * testStateStorage
 **********************************************************************************************************************/

func (storage *testStateStorage) GetModuleState(id string) (state []byte, err error) {
	return storage.state, nil
}

func (storage *testStateStorage) SetModuleState(id string, state []byte) (err error) {
	storage.state = state

	return nil
}

/***********************************************************************************************************************
 * Private
 **********************************************************************************************************************/

func moduleConfig(targetFile string) json.RawMessage {
	return json.RawMessage(
		fmt.Sprintf(`{"sendQueueName":"%s","receiveQueueName":"%s","timeout":"%s","targetFile":"%s"}`,
			commandQueue, statusQueue, queueTimeout.String(), targetFile))
}

func createImage(imageFile, content string) error {
	if err := ioutil.WriteFile(imageFile, []byte(content), 0o600); err != nil {
		return aoserrors.Wrap(err)
	}

	if output, err := exec.Command("gzip", "-f", imageFile).CombinedOutput(); err != nil {
		return aoserrors.Errorf("%s (%s)", err, (string(output)))
	}

	if output, err := exec.Command("mv", imageFile+".gz", imageFile).CombinedOutput(); err != nil {
		return aoserrors.Errorf("%s (%s)", err, (string(output)))
	}

	return nil
}
