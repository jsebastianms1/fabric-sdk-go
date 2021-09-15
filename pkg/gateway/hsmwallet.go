/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"io/ioutil"
	"os"
	"path/filepath"
)


// hsmWalletStore stores identity information used to connect to a Hyperledger Fabric network.
// Instances are created using NewHSMWallet()
type hsmWalletStore struct {
	path string
	hsmConf string
}

// NewHSMWallet creates an instance of a wallet, one part is held on memory
// This implementation is not backed by a persistent store.
//  Parameters:
//  path specifies where on the filesystem to store the wallet.
//	cert is the public certificate for the wallet
//	hsmConf is the path where yaml with PCS11 config is stored
//
//  Returns:
//  A Wallet object.
func NewHSMWallet(path string, hsmConf string) (*Wallet, error) {
	cleanPath := filepath.Clean(path)
	err := os.MkdirAll(cleanPath, os.ModePerm)

	if err != nil {
		return nil, err
	}

	store := &hsmWalletStore{cleanPath, hsmConf }
	return &Wallet{store}, nil

}

// Put an identity into the wallet.
func (hsmw *hsmWalletStore) Put(label string, content []byte) error {
	pathname := filepath.Join(hsmw.path, label) + dataFileExtension
	print(pathname)
	f, err := os.OpenFile(filepath.Clean(pathname), os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	if _, err := f.Write(content); err != nil {
		_ = f.Close() // ignore error; Write error takes precedence
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

// Get an identity from the wallet.
func (hsmw *hsmWalletStore) Get(label string) ([]byte, error) {
	pathname := filepath.Join(hsmw.path, label) + dataFileExtension

	return ioutil.ReadFile(filepath.Clean(pathname))
}

// Remove an identity from the wallet. If the identity does not exist, this method does nothing.
func (hsmw *hsmWalletStore) Remove(label string) error {
	pathname := filepath.Join(hsmw.path, label) + dataFileExtension
	_ = os.Remove(filepath.Clean(pathname))
	return nil
}

// Exists tests the existence of an identity in the wallet.
func (hsmw *hsmWalletStore) Exists(label string) bool {
	pathname := filepath.Join(hsmw.path, label) + dataFileExtension

	_, err := os.Stat(filepath.Clean(pathname))
	return err == nil
}

// List all of the labels in the wallet.
func (hsmw *hsmWalletStore) List() ([]string, error) {
	files, err := ioutil.ReadDir(hsmw.path)

	if err != nil {
		return nil, err
	}

	var labels []string
	for _, file := range files {
		name := file.Name()
		if filepath.Ext(name) == dataFileExtension {
			labels = append(labels, name[:len(name)-extensionLength])
		}
	}

	return labels, nil
}
