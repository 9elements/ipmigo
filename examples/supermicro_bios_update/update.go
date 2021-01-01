package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/k-sone/ipmigo"
)

func supermicroBiosUpdate(c *ipmigo.Client, rom []byte) (err error) {
	var id uint16
	var chunksize uint32

	startcmd := &ipmigo.SupermicroOEMStartBIOSUpgradeCommand{
		ImageSize: uint32(len(rom)),
	}
	if err = c.Execute(startcmd); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Transfer ID: 0x%x\n", startcmd.ID)
	fmt.Printf("MaxChunkSize: 0x%x\n", startcmd.MaxChunkSize)
	id = startcmd.ID
	chunksize = startcmd.MaxChunkSize

	defer func() {
		if err != nil {
			cancelCmd := &ipmigo.SupermicroOEMCancelBIOSCommand{
				ID: id,
			}
			c.Execute(cancelCmd)
		}
	}()

	fmt.Printf("\n============\nUploading BIOS\n============\n")

	reader := bytes.NewBuffer(rom)

	uploadcounter := uint8(0)
	progresscounter := uint8(0)
	for i := uint32(0); i < uint32(len(rom)); {
		uploadcmd := &ipmigo.SupermicroOEMUploadBIOSCommand{
			ID:     id,
			Offset: i,
			Data:   reader.Next(int(chunksize)),
		}
		if err = c.Execute(uploadcmd); err != nil {
			fmt.Println(err)
			return
		}
		i += uint32(chunksize)
		progresscounter = uint8(i * uint32(100) / uint32(len(rom)))
		if progresscounter > uploadcounter+1 {
			for i := uint8(0); i < progresscounter-uploadcounter; i += 2 {
				fmt.Printf(">")
			}
			uploadcounter = progresscounter & 0xfe
		}
	}

	fmt.Printf("\n============\nUpdating BIOS\n============\n")
	flashcmd := &ipmigo.SupermicroOEMFlashBIOSCommand{
		ID: id,
	}
	if err = c.Execute(flashcmd); err != nil {
		fmt.Println(err)
		return
	}

	count := uint8(0)
	for progress := uint8(0); progress < 100; {
		progressCmd := &ipmigo.SupermicroOEMBIOSUpdateProgressCommand{
			ID: id,
		}
		if err = c.Execute(progressCmd); err != nil {
			fmt.Println(err)
			return
		}
		progress = progressCmd.CompletionCode
		if progress > count+1 {
			for i := uint8(0); i < progress-count; i += 2 {
				fmt.Printf(">")
			}
			count = progress & 0xfe
		}
		time.Sleep(time.Second)
	}

	finalCmd := &ipmigo.SupermicroOEMFinalizeBIOSCommand{
		ID: id,
	}
	if err = c.Execute(finalCmd); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("| 100%%\n")
	fmt.Printf("Done\n")
	return
}

// Print sensor data repository entries and readings.
func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: <IP> <username> <password> <BIOS file>\n")
		os.Exit(1)
	}
	addr := os.Args[1]
	if !strings.HasPrefix(addr, ":623") {
		addr = addr + ":623"
	}
	c, err := ipmigo.NewClient(ipmigo.Arguments{
		Version:       ipmigo.V2_0,
		Address:       addr,
		Timeout:       30 * time.Second,
		Retries:       1,
		Username:      os.Args[2],
		Password:      os.Args[3],
		CipherSuiteID: 3,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := c.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	cmdProductID := &ipmigo.SupermicroOEMProductIDCommand{}
	if err := c.Execute(cmdProductID); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Remote Product ID: 0x%x\n", cmdProductID.BoardModelID)
	cmdVer := &ipmigo.SupermicroOEMBIOSVersionCommand{}
	if err := c.Execute(cmdVer); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("BIOS Version: %s\n", cmdVer.BIOSVersion)
	cmdDate := &ipmigo.SupermicroOEMBIOSDate2Command{}
	if err := c.Execute(cmdDate); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("BIOS Date: %s\n", cmdDate.BIOSDate)

	rom, err := ioutil.ReadFile(os.Args[4])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Writing new ROM with size: %d\n", len(rom))

	err = supermicroBiosUpdate(c, rom)
	if err != nil {
		fmt.Println(err)
		return
	}
}
