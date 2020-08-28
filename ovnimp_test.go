package goovn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadTransact(t *testing.T) {
	ovndbapi := getOVNClient(DBSB)
	t.Logf("Adding Chassis to OVN SB DB")
	ocmd, err := ovndbapi.ChassisAdd(CHASSIS_NAME, CHASSIS_HOSTNAME, ENCAP_TYPES, IP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Chassis to OVN failed with err %v", err)
	}
	t.Logf("Adding Chassis to OVN Done")

	t.Logf("Adding second Chassis to OVN SB DB but with same ENCAP_TYPES and IP")
	ocmd, err = ovndbapi.ChassisAdd(CHASSIS2_NAME, CHASSIS2_HOSTNAME, ENCAP_TYPES, IP, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// expecting constraint violation error with following details -- "Transaction causes multiple
	// rows in \"Encap\" table to have identical values (stt and \"10.0.0.11\") for index on columns
	// \"type\" and \"ip\".  First row, with UUID 9860cf40-bd82-4c24-9514-05b225434934, existed in
	// the database before this transaction and was not modified by the transaction.  Second row,
	// with UUID 10d7d018-7444-48de-89fc-cb062f88e520, was inserted by this transaction."
	err = ovndbapi.Execute(ocmd)
	assert.Error(t, err)

	t.Logf("Deleting Chassis:%v", CHASSIS_NAME)
	ocmd, err = ovndbapi.ChassisDel(CHASSIS_NAME)
	if err != nil && err != ErrorNotFound {
		t.Fatal(err)
	}

	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("err executing command:%v", err)
	}
}

func TestConvertGoSetToStringArray(t *testing.T) {
	// 1. create a logical switch and add a port to it.
	// 2. get the newly added port's uuid
	// 3. make sure that portUUID is in logical_switch's ports field.
	ovndbapi := getOVNClient(DBNB)
	t.Logf("Adding LogicalSwitch to OVN NB DB")
	ocmd, err := ovndbapi.LSAdd(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Logical Switch to OVN failed with err %v", err)
	}
	t.Logf("Adding Logical Switch to OVN Done")

	ocmd, err = ovndbapi.LSPAdd(LSW, LSP)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Adding Logical Switch Port to OVN failed with err %v", err)
	}
	t.Logf("Adding Logical Switch Port to OVN Done")

	lspInfo, err := ovndbapi.LSPGet(LSP)
	if err != nil {
		t.Fatal(err)
	}
	lsInfo, err := ovndbapi.LSGet(LSW)
	if err != nil {
		t.Fatal(err)
	}
	uuidFound := false
	for _, port := range lsInfo[0].Ports {
		if port == lspInfo.UUID {
			uuidFound = true
			break
		}
	}
	if !uuidFound {
		t.Fatalf("couldn't find port uuid %s in %s", lspInfo.UUID, LSW)
	}
	t.Logf("Found Logical Switch Port's UUID in Logical Switch")

	t.Logf("Deleting the logical switch " + LSW)
	ocmd, err = ovndbapi.LSDel(LSW)
	if err != nil {
		t.Fatal(err)
	}
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		t.Fatalf("Deleting Logical Switch from OVN failed with err %v", err)
	}
	t.Logf("Deleted the logical switch " + LSW)
}
