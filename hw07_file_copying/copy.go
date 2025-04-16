package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer func() {
		if closeErr := srcFile.Close(); closeErr != nil {
			err = errorJoin(err, fmt.Errorf("close source file: %w", closeErr))
		}
	}()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("get file info: %w", err)
	}

	// For regular files, the expression mode & os.ModeType will return 0
	if srcInfo.Size() < 0 || srcInfo.Mode()&os.ModeType != 0 {
		return ErrUnsupportedFile
	}
	fileSize := srcInfo.Size()

	if offset > fileSize {
		return fmt.Errorf("%w: offset %d exceeds file size %d",
			ErrOffsetExceedsFileSize, offset, fileSize)
	}

	if limit == 0 || offset+limit > fileSize {
		limit = fileSize - offset
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer func() {
		if closeErr := destFile.Close(); closeErr != nil {
			err = errorJoin(err, fmt.Errorf("close destination file: %w", closeErr))
		}
	}()

	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	sr := io.NewSectionReader(srcFile, offset, limit)
	proxyReader := bar.NewProxyReader(sr)

	if _, err := io.CopyN(destFile, proxyReader, limit); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("copy data: %w", err)
	}

	return nil
}

func errorJoin(errs ...error) error {
	return errors.Join(errs...)
}
