package api

import (
	"testing"

	"github.com/StageAutoControl/controller/pkg/cntl"
	internalTesting "github.com/StageAutoControl/controller/pkg/internal/testing"
	"github.com/jinzhu/copier"
)

func TestDMXDeviceController_Create_WithID(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"
	entity := ds.DMXDevices[key]

	createReply := &cntl.DMXDevice{}
	if err := controller.Create(req, entity, createReply); err != nil {
		t.Errorf("failed to call controller: %v", err)
	}

	if createReply.ID != key {
		t.Errorf("Expected createReply to have id %s, but has %s", key, createReply.ID)
	}
}

func TestDMXDeviceController_Create_WithoutID(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"
	entity := ds.DMXDevices[key]

	createEntity := &cntl.DMXDevice{}
	if err := copier.Copy(createEntity, entity); err != nil {
		t.Fatal(err)
	}

	createEntity.ID = ""

	createReply := &cntl.DMXDevice{}
	if err := controller.Create(req, entity, createReply); err != nil {
		t.Errorf("failed to call controller: %v", err)
	}

	if createReply.ID != key {
		t.Errorf("Expected createReply to have id %s, but has %s", key, createReply.ID)
	}
}

func TestDMXDeviceController_Get_NotExisting(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"

	reply := &cntl.DMXDevice{}

	idReq := &IDRequest{ID: key}
	if err := controller.Get(req, idReq, reply); err != errNotExists {
		t.Errorf("expected to get errNotExists, but got %v", err)
	}
}

func TestDMXDeviceController_Get_Existing(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"
	entity := ds.DMXDevices[key]

	createReply := &cntl.DMXDevice{}
	if err := controller.Create(req, entity, createReply); err != nil {
		t.Errorf("failed to call controller: %v", err)
	}

	if createReply.ID != key {
		t.Errorf("Expected createReply to have id %s, but has %s", key, createReply.ID)
	}

	reply := &cntl.DMXDevice{}
	idReq := &IDRequest{ID: key}
	t.Log("idReq has ID:", idReq.ID)
	if err := controller.Get(req, idReq, reply); err != nil {
		t.Errorf("failed to call controller: %v", err)
	}

	if reply.ID != key {
		t.Errorf("Expected reply to have id %s, but has %s", key, reply.ID)
	}
}

func TestDMXDeviceController_Update_NotExisting(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"
	entity := ds.DMXDevices[key]

	reply := &cntl.DMXDevice{}

	if err := controller.Update(req, entity, reply); err != errNotExists {
		t.Errorf("expected to get errNotExists, but got %v", err)
	}
}

func TestDMXDeviceController_Update_Existing(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"
	entity := ds.DMXDevices[key]

	createReply := &cntl.DMXDevice{}
	if err := controller.Create(req, entity, createReply); err != nil {
		t.Errorf("failed to call controller: %v", err)
	}

	if createReply.ID != key {
		t.Errorf("Expected createReply to have id %s, but has %s", key, createReply.ID)
	}

	reply := &cntl.DMXDevice{}
	if err := controller.Update(req, entity, reply); err != nil {
		t.Errorf("expected to get no error, but got %v", err)
	}

	if reply.ID != key {
		t.Errorf("Expected reply to have id %s, but has %s", key, reply.ID)
	}
}
func TestDMXDeviceController_Delete_NotExisting(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"

	reply := &SuccessResponse{}
	idReq := &IDRequest{ID: key}
	if err := controller.Delete(req, idReq, reply); err != errNotExists {
		t.Errorf("expected to get errNotExists, but got %v", err)
	}
}

func TestDMXDeviceController_Delete_Existing(t *testing.T) {
	defer internalTesting.Cleanup(t, path)
	controller := newDMXDeviceController(logger, store)
	key := "35cae00a-0b17-11e7-8bca-bbf30c56f20e"
	entity := ds.DMXDevices[key]

	createReply := &cntl.DMXDevice{}
	if err := controller.Create(req, entity, createReply); err != nil {
		t.Errorf("failed to call controller: %v", err)
	}

	if createReply.ID != key {
		t.Errorf("Expected createReply to have id %s, but has %s", key, createReply.ID)
	}

	reply := &SuccessResponse{}
	idReq := &IDRequest{ID: key}
	if err := controller.Delete(req, idReq, reply); err != nil {
		t.Errorf("expected to get no error, but got %v", err)
	}

	if !reply.Success {
		t.Error("Expected to get result true, but got false")
	}
}
