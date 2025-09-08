/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 19:24:23
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-03 11:43:24
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package httpclient_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/Mrzhouyl/go-aisdk/httpclient"
	"github.com/Mrzhouyl/go-aisdk/httpclient/test/checks"
)

type failingWriter struct {
}

var errMockFailingWriterError = errors.New("mock writer failed")

func (*failingWriter) Write([]byte) (int, error) {
	return 0, errMockFailingWriterError
}

func TestFormBuilderWithFailingWriter(t *testing.T) {
	var (
		file *os.File
		err  error
	)
	if file, err = os.CreateTemp(t.TempDir(), ""); err != nil {
		t.Fatalf("error creating tmp file: %v", err)
	}
	defer file.Close()

	builder := httpclient.NewFormBuilder(&failingWriter{})
	err = builder.CreateFormFile("file", file)
	checks.ErrorIs(t, err, errMockFailingWriterError, "formbuilder should return error if writer fails")
}

func TestFormBuilderWithClosedFile(t *testing.T) {
	var (
		file *os.File
		err  error
	)
	if file, err = os.CreateTemp(t.TempDir(), ""); err != nil {
		t.Fatalf("error creating tmp file: %v", err)
	}
	file.Close()

	builder := httpclient.NewFormBuilder(&bytes.Buffer{})
	err = builder.CreateFormFile("file", file)
	checks.HasError(t, err, "formbuilder should return error if file is closed")
	checks.ErrorIs(t, err, os.ErrClosed, "formbuilder should return error if file is closed")
}
