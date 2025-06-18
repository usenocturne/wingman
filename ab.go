package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
)

const (
	ABMetadataMiscPartitionOffset = 2048
	MiscbufSize                   = 2080
	ABMagic                       = "\000AB0"
	ABMagicLen                    = 4
	ABMajorVersion                = 1
	ABMinorVersion                = 0
	ABDataSize                    = 32
	ABMaxPriority                 = 15
	ABMaxTriesRemaining           = 7
	MiscDevicePath                = "/dev/misc"
)

type ABSlotData struct {
	Priority       uint8    `json:"priority"`
	TriesRemaining uint8    `json:"tries_remaining"`
	SuccessfulBoot uint8    `json:"successful_boot"`
	Reserved       [1]uint8 `json:"-"`
}

type ABData struct {
	Magic        [ABMagicLen]uint8 `json:"-"`
	VersionMajor uint8             `json:"version_major"`
	VersionMinor uint8             `json:"version_minor"`
	Reserved1    [2]uint8          `json:"-"`
	Slots        [2]ABSlotData     `json:"slots"`
	Reserved2    [12]uint8         `json:"-"`
	CRC32        uint32            `json:"crc32"`
}

func (info *ABData) Validate() bool {
	if !bytes.Equal(info.Magic[:], []byte(ABMagic)) {
		fmt.Printf("Magic %s is incorrect.\n", string(info.Magic[:]))
		return false
	}
	if info.VersionMajor > ABMajorVersion {
		fmt.Printf("No support for given major version.\n")
		return false
	}
	return true
}

func (info *ABData) Reset() {
	*info = ABData{}
	copy(info.Magic[:], ABMagic)
	info.VersionMajor = ABMajorVersion
	info.VersionMinor = ABMinorVersion
	info.Slots[0].Priority = ABMaxPriority
	info.Slots[0].TriesRemaining = ABMaxTriesRemaining
	info.Slots[1].Priority = ABMaxPriority - 1
	info.Slots[1].TriesRemaining = ABMaxTriesRemaining
}

func (info *ABData) DumpInfo() {
	activeSlot := info.GetActiveSlot()
	activeSlotLetter := "A"
	if activeSlot == 1 {
		activeSlotLetter = "B"
	}
	fmt.Printf("active slot number: %d\n", activeSlot)
	fmt.Printf("active slot letter: %s\n\n", activeSlotLetter)

	fmt.Printf("slot a priority = %d\n", info.Slots[0].Priority)
	fmt.Printf("slot a tries_remaining = %d\n", info.Slots[0].TriesRemaining)
	fmt.Printf("slot a successful_boot = %d\n\n", info.Slots[0].SuccessfulBoot)

	fmt.Printf("slot b priority = %d\n", info.Slots[1].Priority)
	fmt.Printf("slot b tries_remaining = %d\n", info.Slots[1].TriesRemaining)
	fmt.Printf("slot b successful_boot = %d\n", info.Slots[1].SuccessfulBoot)
}

func (info *ABData) GetActiveSlot() int {
	if info.Slots[0].Priority > info.Slots[1].Priority {
		return 0
	}
	return 1
}

func (info *ABData) SetActiveSlot(slot int) {
	otherSlot := 1 - slot

	info.Slots[slot].Priority = ABMaxPriority
	info.Slots[slot].TriesRemaining = ABMaxTriesRemaining
	info.Slots[slot].SuccessfulBoot = 0

	if info.Slots[otherSlot].Priority == ABMaxPriority {
		info.Slots[otherSlot].Priority = ABMaxPriority - 1
	}
}

func (info *ABData) Failover() {
	newSlot := 1 - info.GetActiveSlot()
	fmt.Printf("Failing over to slot %d...\n", newSlot)
	info.SetActiveSlot(newSlot)
}

func (info *ABData) SetSuccessfulBoot(slot int) {
	info.Slots[slot].TriesRemaining = ABMaxTriesRemaining
	info.Slots[slot].SuccessfulBoot = 1
}

func (info *ABData) calculateCRC32() uint32 {
	data := make([]byte, ABDataSize-4)
	binary.BigEndian.PutUint32(data[len(data)-4:], 0)
	copy(data, info.Magic[:])
	return crc32.ChecksumIEEE(data)
}

func OpenAndLoadABData() (*ABData, error) {
	file, err := os.OpenFile(MiscDevicePath, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open misc partition: %v", err)
	}
	defer file.Close()

	miscBuf := make([]byte, MiscbufSize)
	if _, err := file.ReadAt(miscBuf, 0); err != nil {
		return nil, fmt.Errorf("failed to read misc partition: %v", err)
	}

	info := &ABData{}
	abData := miscBuf[ABMetadataMiscPartitionOffset : ABMetadataMiscPartitionOffset+ABDataSize]
	if err := binary.Read(bytes.NewReader(abData), binary.BigEndian, info); err != nil {
		return nil, fmt.Errorf("failed to parse AB data: %v", err)
	}

	if !info.Validate() {
		return nil, fmt.Errorf("invalid AB data")
	}

	return info, nil
}

func (info *ABData) Save() error {
	info.CRC32 = info.calculateCRC32()

	file, err := os.OpenFile(MiscDevicePath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open misc partition: %v", err)
	}
	defer file.Close()

	miscBuf := make([]byte, MiscbufSize)
	if _, err := file.ReadAt(miscBuf, 0); err != nil {
		return fmt.Errorf("failed to read misc partition: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, info); err != nil {
		return fmt.Errorf("failed to serialize AB data: %v", err)
	}

	copy(miscBuf[ABMetadataMiscPartitionOffset:], buf.Bytes())

	if _, err := file.WriteAt(miscBuf, 0); err != nil {
		return fmt.Errorf("failed to write misc partition: %v", err)
	}

	return nil
}

func (info *ABData) DumpJSON() error {
	type JSONOutput struct {
		ActiveSlot       int           `json:"active_slot"`
		ActiveSlotLetter string        `json:"active_slot_letter"`
		VersionMajor     uint8         `json:"version_major"`
		VersionMinor     uint8         `json:"version_minor"`
		Slots            [2]ABSlotData `json:"slots"`
		CRC32            uint32        `json:"crc32"`
	}

	activeSlot := info.GetActiveSlot()
	activeSlotLetter := "A"
	if activeSlot == 1 {
		activeSlotLetter = "B"
	}

	output := JSONOutput{
		ActiveSlot:       activeSlot,
		ActiveSlotLetter: activeSlotLetter,
		VersionMajor:     info.VersionMajor,
		VersionMinor:     info.VersionMinor,
		Slots:            info.Slots,
		CRC32:            info.CRC32,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %v", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
