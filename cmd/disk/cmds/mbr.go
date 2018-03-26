package cmds

import (
	"flag"
	"fmt"

	"github.com/barsoomia/disk/mbr"
)

var (
	flags                       *flag.FlagSet
	help, create, update        *bool
	addPart, delPart, startSect *int
	bootcode, lastSect          *string
)

func init() {
	flags = flag.NewFlagSet("mbr", flag.ContinueOnError)
	help = flags.Bool("help", false, "Show this help")
	create = flags.Bool("create", false, "Create new MBR")
	update = flags.Bool("update", false, "Update MBR")
	addPart = flag.Int("add-part", 0, "Add partition")
	delPart = flag.Int("del-part", 0, "Delete partition")
	startSect = flag.Int("start-sect", 0, "start sector")
	lastSect = flag.String("last-sect", "", "last sector (modififers +K, +M, +G works)")
	bootcode = flags.String("bootcode", "", "Bootsector binary code")
}

func addPartition(disk string, partn int, startsect int, lastsect string) error {
	mbrdata, err := mbr.FromFile(disk)
	if err != nil {
		return err
	}

	p := mbr.NewEmptyPartition()
	mbrdata.SetPart(partn, p)
	return nil
}

func MBR(args []string) error {
	flags.Parse(args[1:])

	if *help {
		flags.PrintDefaults()
		return nil
	}

	disks := flags.Args()
	if len(disks) != 1 {
		return fmt.Errorf("Require one device file")
	}

	if *create {
		if *update {
			return fmt.Errorf("-create conflicts with -update")
		}

		return mbr.Create(disks[0], *bootcode)
	}

	if *update {
		if *addPart <= 0 || *delPart == 0 {
			return fmt.Errorf("-update requires flag -add-part or --del-part")
		}

		if *addPart != 0 {
			partn := *addPart
			if *startSect == -1 {
				return fmt.Errorf("-add-part requires -start-sect")
			}

			if *lastSect == "" {
				return fmt.Errorf("-add-part requires -last-sect")
			}
			return addPartition(disks[0], partn, *startSect, *lastSect)
		}

		return fmt.Errorf("-del-part not implemented")
	}

	for _, disk := range disks {
		err := mbr.Info(disk)
		if err != nil {
			return err
		}
	}

	return nil
}