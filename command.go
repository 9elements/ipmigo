package ipmigo

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Completion Code (Section 5.2)
type CompletionCode uint8

const (
	CompletionOK               CompletionCode = 0x00
	CompletionUnspecifiedError CompletionCode = 0xff

	CompletionNodeBusy                 CompletionCode = 0xc0
	CompletionInvalidCommand           CompletionCode = 0xc1
	CompletionInvalidCommandForLUN     CompletionCode = 0xc2
	CompletionTimeout                  CompletionCode = 0xc3
	CompletionOutOfSpace               CompletionCode = 0xc4
	CompletionReservationCancelled     CompletionCode = 0xc5
	CompletionRequestDataTruncated     CompletionCode = 0xc6
	CompletionRequestDataInvalidLength CompletionCode = 0xc7
	CompletionRequestDataFieldExceedEd CompletionCode = 0xc8
	CompletionParameterOutOfRange      CompletionCode = 0xc9
	CompletionCantReturnDataBytes      CompletionCode = 0xca
	CompletionRequestDataNotPresent    CompletionCode = 0xcb
	CompletionInvalidDataField         CompletionCode = 0xcc
	CompletionIllegalSendorOrRecord    CompletionCode = 0xcd
	CompletionCantBeProvided           CompletionCode = 0xce
	CompletionDuplicatedRequest        CompletionCode = 0xcf
	CompletionSDRInUpdateMode          CompletionCode = 0xd0
	CompletionFirmwareUpdateMode       CompletionCode = 0xd1
	CompletionBMCInitialization        CompletionCode = 0xd2
	CompletionDestinationUnavailable   CompletionCode = 0xd3
	CompletionInsufficientPrivilege    CompletionCode = 0xd4
	CompletionNotSupportedPresentState CompletionCode = 0xd5
	CompletionIllegalCommandDisabled   CompletionCode = 0xd6
)

func (c CompletionCode) String() string {
	switch c {
	case CompletionOK:
		return "Command Completed Normally"
	case CompletionUnspecifiedError:
		return "Unspecified error"
	case CompletionNodeBusy:
		return "Node Busy"
	case CompletionInvalidCommand:
		return "Invalid Command"
	case CompletionInvalidCommandForLUN:
		return "Command invalid for given LUN"
	case CompletionTimeout:
		return "Timeout"
	case CompletionOutOfSpace:
		return "Out of space"
	case CompletionReservationCancelled:
		return "Reservation Canceled or Invalid Reservation ID"
	case CompletionRequestDataTruncated:
		return "Request data truncated"
	case CompletionRequestDataInvalidLength:
		return "Request data length invalid"
	case CompletionRequestDataFieldExceedEd:
		return "Request data field length limit exceeded"
	case CompletionParameterOutOfRange:
		return "Parameter out of range"
	case CompletionCantReturnDataBytes:
		return "Cannot return number of requested data bytes"
	case CompletionRequestDataNotPresent:
		return "Requested Sensor, data, or record not present"
	case CompletionInvalidDataField:
		return "Invalid data field in Request"
	case CompletionIllegalSendorOrRecord:
		return "Command illegal for specified sensor or record type"
	case CompletionCantBeProvided:
		return "Command response could not be provided"
	case CompletionDuplicatedRequest:
		return "Cannot execute duplicated request"
	case CompletionSDRInUpdateMode:
		return "SDR Repository in update mode"
	case CompletionFirmwareUpdateMode:
		return "Device in firmware update mode"
	case CompletionBMCInitialization:
		return "BMC initialization or initialization agent in progress"
	case CompletionDestinationUnavailable:
		return "Destination unavailable"
	case CompletionInsufficientPrivilege:
		return "Cannot execute command due to insufficient privilege level"
	case CompletionNotSupportedPresentState:
		return "Command not supported in present state"
	case CompletionIllegalCommandDisabled:
		return "Command sub-function has been disabled or is unavailable"
	default:
		return fmt.Sprintf("0x%02x", int(c))
	}
}

type Command interface {
	Name() string
	Code() uint8
	NetFnRsLUN() NetFnRsLUN
	Marshal() (buf []byte, err error)
	Unmarshal(buf []byte) (rest []byte, err error)
	String() string
}

type RawCommand struct {
	name       string
	code       uint8
	netFnRsLUN NetFnRsLUN
	input      []byte
	output     []byte
}

func (c *RawCommand) Name() string             { return c.name }
func (c *RawCommand) Code() uint8              { return c.code }
func (c *RawCommand) NetFnRsLUN() NetFnRsLUN   { return c.netFnRsLUN }
func (c *RawCommand) Input() []byte            { return c.input }
func (c *RawCommand) Output() []byte           { return c.output }
func (c *RawCommand) Marshal() ([]byte, error) { return c.input, nil }

func (c *RawCommand) Unmarshal(buf []byte) ([]byte, error) {
	c.output = make([]byte, len(buf))
	copy(c.output, buf)
	return nil, nil
}

func (c *RawCommand) String() string {
	return fmt.Sprintf(`{"Name":"%s","Code":%d,"NetFnRsRUN":%d,"Input":"%s","Output":"%s"}`,
		c.name, c.code, c.netFnRsLUN, hex.EncodeToString(c.input), hex.EncodeToString(c.output))
}

func NewRawCommand(name string, code uint8, fn NetFnRsLUN, input []byte) *RawCommand {
	return &RawCommand{
		name:       name,
		code:       code,
		netFnRsLUN: fn,
		input:      input,
	}
}

func cmdToJSON(c Command) string {
	s := fmt.Sprintf(`{"Name":"%s","Code":%d,"NetFnRsLUN":%d,`, c.Name(), c.Code(), c.NetFnRsLUN())
	return strings.Replace(toJSON(c), `{`, s, 1)
}

func cmdValidateLength(c Command, msg []byte, min int) error {
	if l := len(msg); l < min {
		return &MessageError{
			Message: fmt.Sprintf("Invalid %s Response size : %d/%d", c.Name(), l, min),
			Detail:  hex.EncodeToString(msg),
		}
	}
	return nil
}
