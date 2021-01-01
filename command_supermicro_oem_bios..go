package ipmigo

import (
	"bytes"
	"encoding/binary"
)

// SupermicroOEMStartBIOSUpgradeCommand is the Supermicro OEM command to start BIOS upgrade over IPMI
type SupermicroOEMStartBIOSUpgradeCommand struct {
	// Request Data
	ImageSize uint32

	// Response Data
	ID           uint16
	MaxChunkSize uint32
}

func (c *SupermicroOEMStartBIOSUpgradeCommand) Name() string { return "Start BIOS Upgrade" }
func (c *SupermicroOEMStartBIOSUpgradeCommand) Code() uint8  { return 0x61 }
func (c *SupermicroOEMStartBIOSUpgradeCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnOEM, 0)
}
func (c *SupermicroOEMStartBIOSUpgradeCommand) String() string { return cmdToJSON(c) }
func (c *SupermicroOEMStartBIOSUpgradeCommand) Marshal() ([]byte, error) {
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, c.ImageSize)

	return a, nil
}

func (c *SupermicroOEMStartBIOSUpgradeCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 6); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf)

	if err := binary.Read(reader, binary.LittleEndian, &c.ID); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.MaxChunkSize); err != nil {
		return nil, err
	}

	return nil, nil
}

// SupermicroOEMUploadBIOSCommand is the Supermicro OEM command to upload a BIOS chunk over IPMI
type SupermicroOEMUploadBIOSCommand struct {
	// Request Data
	ID     uint16
	Offset uint32
	Data   []byte
}

func (c *SupermicroOEMUploadBIOSCommand) Name() string           { return "Upload BIOS chunk" }
func (c *SupermicroOEMUploadBIOSCommand) Code() uint8            { return 0x62 }
func (c *SupermicroOEMUploadBIOSCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMUploadBIOSCommand) String() string         { return cmdToJSON(c) }
func (c *SupermicroOEMUploadBIOSCommand) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, c.ID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Offset); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *SupermicroOEMUploadBIOSCommand) Unmarshal(buf []byte) ([]byte, error) { return nil, nil }

// SupermicroOEMFlashBIOSCommand is the Supermicro OEM command to flash all uploaded BIOS chunk over IPMI
type SupermicroOEMFlashBIOSCommand struct {
	// Request Data
	ID   uint16
	Flag uint32
}

func (c *SupermicroOEMFlashBIOSCommand) Name() string           { return "Flash BIOS" }
func (c *SupermicroOEMFlashBIOSCommand) Code() uint8            { return 0x63 }
func (c *SupermicroOEMFlashBIOSCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMFlashBIOSCommand) String() string         { return cmdToJSON(c) }
func (c *SupermicroOEMFlashBIOSCommand) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, c.ID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Flag); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *SupermicroOEMFlashBIOSCommand) Unmarshal(buf []byte) ([]byte, error) { return nil, nil }

// SupermicroOEMCancelBIOSCommand is the Supermicro OEM command to cancel a BIOS flash over IPMI
type SupermicroOEMCancelBIOSCommand struct {
	// Request Data
	ID uint16
}

func (c *SupermicroOEMCancelBIOSCommand) Name() string           { return "Cancel BIOS Update" }
func (c *SupermicroOEMCancelBIOSCommand) Code() uint8            { return 0x64 }
func (c *SupermicroOEMCancelBIOSCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMCancelBIOSCommand) String() string         { return cmdToJSON(c) }
func (c *SupermicroOEMCancelBIOSCommand) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, c.ID); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func (c *SupermicroOEMCancelBIOSCommand) Unmarshal(buf []byte) ([]byte, error) { return nil, nil }

// SupermicroOEMFinalizeBIOSCommand is the Supermicro OEM command to finalize a BIOS flash over IPMI
type SupermicroOEMFinalizeBIOSCommand struct {
	// Request Data
	ID uint16
}

func (c *SupermicroOEMFinalizeBIOSCommand) Name() string           { return "Finalize BIOS Update" }
func (c *SupermicroOEMFinalizeBIOSCommand) Code() uint8            { return 0x65 }
func (c *SupermicroOEMFinalizeBIOSCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMFinalizeBIOSCommand) String() string         { return cmdToJSON(c) }
func (c *SupermicroOEMFinalizeBIOSCommand) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, c.ID); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *SupermicroOEMFinalizeBIOSCommand) Unmarshal(buf []byte) ([]byte, error) { return nil, nil }

// SupermicroOEMBIOSUpdateProgressCommand is the Supermicro OEM command to get the current BIOS flash state over IPMI
type SupermicroOEMBIOSUpdateProgressCommand struct {
	// Request Data
	ID uint16
	// Response Data
	CompletionCode uint8
}

func (c *SupermicroOEMBIOSUpdateProgressCommand) Name() string { return "BIOS Update Progress" }
func (c *SupermicroOEMBIOSUpdateProgressCommand) Code() uint8  { return 0x66 }
func (c *SupermicroOEMBIOSUpdateProgressCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnOEM, 0)
}
func (c *SupermicroOEMBIOSUpdateProgressCommand) String() string { return cmdToJSON(c) }
func (c *SupermicroOEMBIOSUpdateProgressCommand) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, c.ID); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *SupermicroOEMBIOSUpdateProgressCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 1); err != nil {
		return nil, err
	}

	c.CompletionCode = buf[0]

	return nil, nil
}

// SupermicroOEMProductIDCommand is the Supermicro OEM command to get the board's product ID
type SupermicroOEMProductIDCommand struct {
	// Response Data
	BoardModelID uint16
}

func (c *SupermicroOEMProductIDCommand) Name() string             { return "Get Product ID" }
func (c *SupermicroOEMProductIDCommand) Code() uint8              { return 0x21 }
func (c *SupermicroOEMProductIDCommand) NetFnRsLUN() NetFnRsLUN   { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMProductIDCommand) String() string           { return cmdToJSON(c) }
func (c *SupermicroOEMProductIDCommand) Marshal() ([]byte, error) { return []byte{}, nil }

func (c *SupermicroOEMProductIDCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 2); err != nil {
		return nil, err
	}

	c.BoardModelID = binary.LittleEndian.Uint16(buf)

	return nil, nil
}

// SupermicroOEMBIOSVersionCommand is the Supermicro OEM command to get the board's BIOS version
type SupermicroOEMBIOSVersionCommand struct {
	// Response Data
	BIOSVersion string
}

func (c *SupermicroOEMBIOSVersionCommand) Name() string           { return "Get BIOS Version" }
func (c *SupermicroOEMBIOSVersionCommand) Code() uint8            { return 0xac }
func (c *SupermicroOEMBIOSVersionCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMBIOSVersionCommand) String() string         { return cmdToJSON(c) }
func (c *SupermicroOEMBIOSVersionCommand) Marshal() ([]byte, error) {
	a := make([]byte, 2)
	binary.LittleEndian.PutUint16(a, 0)

	return a, nil
}

func (c *SupermicroOEMBIOSVersionCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 1); err != nil {
		return nil, err
	}

	c.BIOSVersion = string(buf)

	return nil, nil
}

// SupermicroOEMBIOSDate2Command is the Supermicro OEM command to get the board's BIOS date
type SupermicroOEMBIOSDate2Command struct {
	// Response Data
	BIOSDate string
}

func (c *SupermicroOEMBIOSDate2Command) Name() string           { return "Get BIOS Date" }
func (c *SupermicroOEMBIOSDate2Command) Code() uint8            { return 0xac }
func (c *SupermicroOEMBIOSDate2Command) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnOEM, 0) }
func (c *SupermicroOEMBIOSDate2Command) String() string         { return cmdToJSON(c) }
func (c *SupermicroOEMBIOSDate2Command) Marshal() ([]byte, error) {
	a := make([]byte, 2)
	binary.LittleEndian.PutUint16(a, 1)

	return a, nil
}

func (c *SupermicroOEMBIOSDate2Command) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 1); err != nil {
		return nil, err
	}

	c.BIOSDate = string(buf)

	return nil, nil
}
