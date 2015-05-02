/*
	(c) 2015 Vlad Didenko. All Rights Reserved.

	This file is part of Port ALLocator Daemon a.k.a. pald.

	pald is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	pald is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with pald.  If not, see <http://www.gnu.org/licenses/>.
*/

package server

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const (
	testPort uint16 = 8001
)

var (
	testUrl string = "http://localhost:" + strconv.Itoa(int(testPort))
)

func TestPaldHttp(t *testing.T) {

	testCases := []struct {
		request  string
		httpCode int
		respFore string
	}{
		{request: "/set?service=a0", httpCode: http.StatusOK, respFore: "49200"},
		{request: "/set?service=a1", httpCode: http.StatusOK, respFore: "49201"},
		{request: "/get?service=a2", httpCode: http.StatusNotFound, respFore: "Name \"a2\" not found in the port registry"},
		{request: "/set?service=a2", httpCode: http.StatusOK, respFore: "49202"},
		{request: "/get?service=a1", httpCode: http.StatusOK, respFore: "49201"},
		{request: "/set?service=a2", httpCode: http.StatusPreconditionFailed, respFore: "Name \"a2\" is already taken"},
		{request: "/set?service=a2&port=49000", httpCode: http.StatusPreconditionFailed, respFore: "Name \"a2\" is already taken"},
		{request: "/set?service=f0&port=49000", httpCode: http.StatusOK, respFore: "OK"},
		{request: "/get?service=f0", httpCode: http.StatusOK, respFore: "49000"},
		{request: "/del?port=49201", httpCode: http.StatusOK, respFore: "OK"},
		{request: "/set?service=a3", httpCode: http.StatusOK, respFore: "49201"},
		{request: "/set?service=a4", httpCode: http.StatusPreconditionFailed, respFore: "No ports available"},
		{request: "/set?svc=er", httpCode: http.StatusBadRequest, respFore: "Service name is missing"},
		{request: "/get?svc=er", httpCode: http.StatusBadRequest, respFore: "Service name is missing"},
		{request: "/set?service=f1&port=49001", httpCode: http.StatusOK, respFore: "OK"},
		{request: "/del?port=49O01", httpCode: http.StatusBadRequest, respFore: "strconv.ParseUint:"},
		{request: "/set?sevice=er", httpCode: http.StatusBadRequest, respFore: "Service name is missing"},
		{request: "/del?svc=er", httpCode: http.StatusBadRequest, respFore: "Port number is missing"},
	}

	go Run(testPort, 49200, 49202)

	for _, tc := range testCases {

		resp, err := http.Get(testUrl + tc.request)
		defer resp.Body.Close()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tc.httpCode {
			t.Errorf("Received code %d instead of %d from %q request", resp.StatusCode, tc.httpCode, tc.request)
		}

		line, err := bufio.NewReader(resp.Body).ReadString('\n')
		if !strings.HasPrefix(line, tc.respFore) {
			t.Errorf("Wrong response body. Expected to start with %q, but it is %q", tc.respFore, line)
		}
	}
}
