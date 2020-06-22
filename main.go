package main

import (
	"github.com/raju/go-qemu-ex/qemu"
	"github.com/raju/go-qemu-ex/qmp"
		"log"
		"fmt"
		// "net"
		"io"
		//"time"
		)
		const (
			GiB = 1073741824 // 10 GiB 
		)

func createImage(imageName string){
		img := qemu.NewImage(imageName, qemu.ImageFormatQCOW2, 10*GiB)
		img.SetBackingFile(imageName)
	
		err := img.Create()
		if err != nil {
			log.Fatal(err)
		}
}
func createMachine(socketPath string, imageName string, isoImagePath string){
	img, err := qemu.OpenImage(imageName)
	if err != nil {
		//log.Fatal(err)
		createImage(imageName)
		newimg, err := qemu.OpenImage(imageName)
		if err != nil {
			log.Fatal(err)			
		}
		img = newimg
	}

	fmt.Println("image", img.Path, "format", img.Format, "size", img.Size)

	err = img.CreateSnapshot("backup")
	if err != nil {
		log.Fatal(err)
	}

	snaps, err := img.Snapshots()
	if err != nil {
		log.Fatal(err)
	}

	for _, snapshot := range snaps {
		fmt.Println(snapshot.Name, snapshot.Date)
	}
	// img, err := qemu.OpenImage("debian.qcow2")
	// if err != nil {
	// 	log.Fatal(err)
	// }


	// d := qemu.Drive{
	// 	Path:"debian.qcow2",
	// 	Format:"qcow2",/tmp/qmp-socket
	// }

	m := qemu.NewMachine(1, 2024) // 1 CPU, 512MiB RAM
	// m.AddDrive(d)
	m.AddDriveImage(img)
	m.AddCDRom(isoImagePath)
	m.AddMonitorUnix(socketPath)
	m.SetDisplay("none")
	//m.SetDisplay("vga")

	pid, err := m.Start("x86_64", false) // x86_64 arch (using qemu-system-x86_64), with kvm
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("QEMU started on PID", pid)
}

func reader(r io.Reader) {
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf[:])
        if err != nil {
            return
        }
        println("Client got:", string(buf[0:n]))
    }
}

func connectMachine(socketPath string){
		// Connection to QMP
		c, err := qmp.Open("unix", socketPath)
		if err != nil {
			log.Fatal(err)
		}

		defer c.Close()

		// Execute simple QMP command
		result, err := c.Command("query-status", nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("query status")
		fmt.Println(result)

		// Execute QMP command with arguments

		args := map[string]interface{} {"device": "ide1-cd0"}
		fmt.Println(args)

		result, err = c.Command("eject", args)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("eject status")
		fmt.Println(result)
		// Execute HMP command
		result, err = c.HumanMonitorCommand("savevm checkpoint")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("monistor status")
		fmt.Println(result)
		fmt.Println("connection done")
}
func  main()  {
	unixSocket := "qmp1.sock"
	// isoImagePath := "/home/swamym/Downloads/ubuntu-18.04.4-desktop-amd64.iso"
	// imageName := "ubuntu-debian181.qcow2"
 	// createMachine(unixSocket,imageName,isoImagePath)
 	fmt.Println("machine created")
 	 connectMachine(unixSocket)

}