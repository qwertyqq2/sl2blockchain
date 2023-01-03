package network

import (
	"log"
	"net"
	"strings"
	"time"
)

type Package struct {
	Option int
	Data   string
}

const (
	ENDBYTES = "\0005\0000\0001"
	WAITTAIM = 10
	BUFSIZE  = 4 << 10
	MAXSIZE  = 2 << 20
)

type Listener net.Listener
type Conn net.Conn

func Listen(address string, handle func(Conn, *Package)) (Listener, error) {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return nil, ErrIncorrectAddress
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+splited[1])
	if err != nil {
		return nil, err
	}
	go serve(listener, handle)
	return Listener(listener), nil
}

func serve(listener net.Listener, handle func(Conn, *Package)) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handleConn(conn, handle)
	}
}

func handleConn(conn net.Conn, handle func(Conn, *Package)) {
	defer conn.Close()
	pack := readPack(conn)
	if pack == nil {
		log.Println(ErrNotPack)
		return
	}
	handle(Conn(conn), pack)
}

func readPack(conn net.Conn) *Package {
	var (
		size = uint64(0)
		buf  = make([]byte, BUFSIZE)
		data string
	)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return nil
		}
		size += uint64(length)
		if size > MAXSIZE {
			log.Println(err)
			return nil
		}
		data += string(buf[:length])
		if strings.Contains(data, ENDBYTES) {
			data = strings.Split(data, ENDBYTES)[0]
			break
		}
	}
	deserializePack, err := DeserializePack(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return deserializePack
}

func Send(address string, pack *Package) *Package {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer conn.Close()
	serilizePack, err := SerializePack(pack)
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = conn.Write([]byte(serilizePack + ENDBYTES))
	if err != nil {
		log.Println(err)
		return nil
	}
	var (
		ch  = make(chan bool)
		res = new(Package)
	)
	go func() {
		res = readPack(conn)
		ch <- true
	}()
	select {
	case <-ch:

	case <-time.After(WAITTAIM * time.Second):
		log.Println(ErrTimeWait)
		return nil
	}
	return res
}

func Handle(option int, conn Conn, pack *Package, handle func(*Package) string) error {
	if pack.Option != option {
		return ErrNotOpt
	}
	serializePack, err := SerializePack(&Package{
		Option: option,
		Data:   handle(pack),
	})
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(serializePack + ENDBYTES))
	return err
}
