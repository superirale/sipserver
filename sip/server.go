package sip

import (
	"bufio"
	"net"
	"github.com/superirale/sipserver/utils"
)

// UseUDP udp server implementation
func UseUDP(port string)  {
	udpAddr, err := net.ResolveUDPAddr("udp4", port)
	utils.CheckError(err)

	conn, err := net.ListenUDP("udp4", udpAddr)
	utils.CheckError(err)

	for {
		handleUDPClient(conn)
	}
}


func UseTCP(port string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	utils.CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	utils.CheckError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// run as a goroutine
		go handleTCPClient(conn)
	}

}


func handleUDPClient(conn *net.UDPConn) {
	defer conn.Close()

	buf := make([]byte, (4096))

	for {

		n, remoteAddr, err := conn.ReadFromUDP(buf)
		localAddr := conn.LocalAddr().String()
		addresses := map[string]string{"local": localAddr, "remote": remoteAddr.String()}

		if err != nil {
			continue
		}

		resp := RequestHandler(buf[0:n], n, addresses)

		_, err2 := conn.WriteToUDP([]byte(resp), remoteAddr)
		if err2 != nil {
			return
		}

	}
}


func handleTCPClient(conn net.Conn)  {
	defer conn.Close()

	buf := make([]byte, (4096))
	for {

		n, err := conn.Read(buf)

		if err != nil {
			continue
		}

		localAddr := conn.LocalAddr().String()
		remoteAddr := conn.RemoteAddr().String()
		addresses := map[string]string{"local": localAddr, "remote": remoteAddr}

		resp := RequestHandler(buf[0:n], n, addresses)

		_, err2 := conn.Write([]byte(resp+"\n"))
		if err2 != nil {
			return
		}
	}
}


// ForwardConn function to send requests to clients
func ForwardConn(msg []byte, remoteAddr, connType string) string  {

	var response string

	if connType == "udp" {
		rAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
		utils.CheckError(err)

		conn, conErr := net.DialUDP("udp", nil, rAddr)
		utils.CheckError(conErr)
		defer conn.Close()
		conn.Write(msg)

		_, writeErr := conn.Write(msg)
		utils.CheckError(writeErr)

		var buf [512]byte
		n, err := conn.Read(buf[0:])
		utils.CheckError(err)

		response = string(buf[:n])

	}

	if connType == "tcp" {
		rAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
		utils.CheckError(err)

		conn, conErr := net.DialTCP("tcp", nil, rAddr)
		utils.CheckError(conErr)
		defer conn.Close()
		conn.Write(msg)

		_, writeErr := conn.Write(msg)
		utils.CheckError(writeErr)

		connbuf := bufio.NewReader(conn)

		for{

			//  read from server
			str, err := connbuf.ReadString('\n')
			if len(str)>0 {
				response = str
			}
			utils.CheckError(err)

		}
	}
	return response
}