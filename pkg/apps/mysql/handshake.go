package mysql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func processServerHello(info *MysqlInfo, data []byte) error {
	versionEnd := bytes.IndexByte(data, 0)
	if versionEnd == -1 {
		return errors.New("invalid server version")
	}

	info.Protocol = data[0]
	info.Version = string(data[1:versionEnd])

	return nil
}

func processClientHello(info *MysqlInfo, data []byte) error {
	pos := 0

	//skip capability
	info.Capability = binary.LittleEndian.Uint32(data[:4])
	pos += 4

	//skip max packet size
	pos += 4

	//charset, skip, if you want to use another charset, use set names
	//c.collation = CollationId(data[pos])
	pos++

	//skip reserved 23[00]
	pos += 23

	//user name
	info.User = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])

	pos += len(info.User) + 1

	//auth length and auth
	authLen := int(data[pos])
	pos++
	// auth := data[pos : pos+authLen]
	pos += authLen
	fmt.Println("Authlen:", authLen)

	if info.Capability&CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA == 1 {
		//TODO:
	} else if info.Capability&CLIENT_SECURE_CONNECTION == 1 {
		//TODO:
	} else {
		//TODO:
	}

	var db string
	if info.Capability&CLIENT_CONNECT_WITH_DB > 0 {
		if len(data[pos:]) == 0 {
			return nil
		}

		db = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
		pos += len(db) + 1
	}
	info.Db = db

	return nil
}
