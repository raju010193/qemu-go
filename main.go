package main

import (
	"github.com/raju/go-qemu-ex/qemu"
		"log"
		"fmt"
		"net"
		"io"
		//"time"
		)
		const (
			GiB = 10737418240 // 10 GiB 
		)

func createImage(imageName string){
		img := qemu.NewImage(imageName, qemu.ImageFormatQCOW2, 5*GiB)
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
	c, err := net.Dial("unix", socketPath)
    if err != nil {
        panic(err)
	}
	fmt.Println(c)
    defer c.Close()

	reader(c)
	fmt.Println("its done")
    // for {
		
	// 	a, err := c.Write([]byte("hi"))
	// 	fmt.Println("ineer data", a)
    //     if err != nil {
    //         log.Fatal("write error:", err)
    //         break
    //     }
    //     time.Sleep(10)
	// }
	fmt.Println("its last one")
}
func  main()  {
	unixSocket := "qmp.sock"
	isoImagePath := "/home/swamym/Downloads/ubuntu-20.04-desktop-amd64.iso"
	imageName := "ubuntu-debian.qcow2"
 	createMachine(unixSocket,imageName,isoImagePath)
 	fmt.Println("machine created")
 	 connectMachine(unixSocket)

}