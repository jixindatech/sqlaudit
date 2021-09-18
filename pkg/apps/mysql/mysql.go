package mysql

import (
	"errors"
)

/*
 * Mysql part
 */

func readDataWithLength(info *MysqlInfo, data []byte) ([]byte, error) {
	if len(info.Data) > 0 {
		data = append(info.Data, data...)
	}

	if len(data) < 4 {
		return nil, errors.New("not enough data")
	}

	length := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)
	if length < 1 {
		return nil, errors.New("invalid data length")
	}

	info.PacketNo = data[3]

	dataLen := len(data)
	if length == dataLen-4 {
		return data[4:], nil
	} else if length < dataLen {
		return nil, errors.New("invalid data length")
	} else {
		info.Data = append(info.Data, data...)
	}

	return nil, errors.New("not enough data")
}

func ProcessServer(session interface{}, data []byte) error {
	info := session.(*MysqlInfo)
	if info.Status == UNKOWNSTATUS {
		data, err := readDataWithLength(info, data)
		if err != nil {
			return err
		}

		if info.PacketNo != 0 {
			return errors.New("invalid server packet no")
		}

		info.Status = SERVERHELLO
		_ = processServerHello(info, data)

	} else if info.Status == CLIENTHEELO {
		data, err := readDataWithLength(info, data)
		if err != nil {
			info.Status = UNKOWNSTATUS
			return err
		}

		//TODO: mysql 8.0 ?
		if info.PacketNo != 2 {
			return errors.New("Invalid server packet no")
		}

		if data[0] == OK_HEADER {
			info.Status = STATUSOK
		} else if data[0] == 0xfe {
			info.Status = SERVERSWITCHREQUEST
			return nil
		}
	} else if info.Status == CLIENTSWITCHRESPONSE {
		data, err := readDataWithLength(info, data)
		if err != nil {
			info.Status = UNKOWNSTATUS
			return err
		}

		//TODO: mysql 8.0 ?
		if info.PacketNo != 4 {
			return errors.New("Invalid server packet no")
		}

		if data[0] == OK_HEADER {
			info.Status = STATUSOK
		}
	} else {
		//TODO: skip now
	}

	return nil
}

func ProcessClient(session interface{}, data []byte) error {
	info := session.(*MysqlInfo)
	data, err := readDataWithLength(info, data)
	if err != nil {
		return err
	}

	// TODO: https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse
	// Depending on the servers support for the CLIENT_PROTOCOL_41 capability and the clients
	// understanding of that flag the client has to send either a Protocol::HandshakeResponse41 or Protocol::HandshakeResponse320.
	if info.Status == SERVERHELLO {
		if info.PacketNo != 1 {
			return errors.New("Invalid client packet no")
		}

		info.Status = CLIENTHEELO
		err = processClientHello(info, data)
	} else if info.Status == SERVERSWITCHREQUEST {
		if info.PacketNo != 3 {
			return errors.New("Invalid server packet no")
		}

		info.Status = CLIENTSWITCHRESPONSE
		// process client data
	} else if info.Status == STATUSOK {
		err = processRequest(info, data)
	}

	return err
}
