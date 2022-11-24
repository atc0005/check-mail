// Package sasl is based on code from the github.com/emersion/go-sasl package
// at commit 4132e15e133dd337ee91a3b320fa6c0596caa819. Specifically, code was
// copied from these files:
//
// https://github.com/emersion/go-sasl/blob/4132e15e133dd337ee91a3b320fa6c0596caa819/xoauth2.go
// https://github.com/emersion/go-sasl/blob/4132e15e133dd337ee91a3b320fa6c0596caa819/sasl.go
//
// This package restores support for the XOAUTH2 authentication mechanism that
// was removed from the upstream project per
// https://github.com/emersion/go-sasl/issues/18.
package sasl
