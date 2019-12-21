/**
PUG
Copyright (C) 2019  Yesterday17

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package bilibili

import (
	"github.com/Yesterday17/pug/utils/net"
	"os"
	"strings"
)

func (m *Module) PreUpload(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	json, err := net.GetJSON(net.BuildUrl("member.bilibili.com", true, "preupload", map[string]string{
		"name":    file.Name(),
		"size":    bigInt(stat.Size()).String(),
		"r":       m.Route.os,
		"profile": m.Route.profile,
		"ssl":     "0",
		"version": "2.7.0",
		"build":   "2070000",
	})+m.Route.query, nil)
	if err != nil {
		return err
	}

	m.UposUri = json.Get("upos_uri").String()          // Upos Uri
	m.Auth = json.Get("auth").String()                 // Auth
	m.BizID = bigInt(json.Get("biz_id").Int())         // Biz ID
	m.ChunkSize = bigInt(json.Get("chunk_size").Int()) // Chunk Size
	m.Threads = bigInt(json.Get("threads").Int())      // Threads

	// EndPoint
	if json.Get("endpoints.#").Int() > 0 {
		m.EndPoint = json.Get("endpoints.0").String()
	} else {
		m.EndPoint = json.Get("endpoint").String()
	}
	return nil
}

func (m *Module) UploadsPost() error {
	url := strings.ReplaceAll(m.UposUri, "upos:\\/\\/", m.EndPoint)[2:] // 2: here to remove \/\/
	json, err := net.PostJSON(net.BuildUrl(url, true, "", map[string]string{
		"uploads": "true",
		"output":  "json",
	}), net.Headers{
		"X-Upos-Auth": m.Auth,
	}, nil)
	if err != nil {
		return err
	}

	m.UploadID = json.Get("upload_id").String()
	m.Key = json.Get("key").String()
	return nil
}
