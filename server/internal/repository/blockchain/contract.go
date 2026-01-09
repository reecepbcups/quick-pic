// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// QuickPicStorageMetaData contains all meta data concerning the QuickPicStorage contract.
var QuickPicStorageMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"areFriends\",\"inputs\":[{\"name\":\"userId1\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userId2\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createFriendRequest\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createFriendship\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userAId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userBId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createMessage\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedContent\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"contentType\",\"type\":\"uint8\",\"internalType\":\"enumQuickPicStorage.ContentType\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createUser\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"username\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"passwordHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publicKey\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"userNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deleteMessage\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"friendRequestByUsers\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"friendRequestIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"friendRequests\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumQuickPicStorage.FriendRequestStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"friendshipByUsers\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"friendshipIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"friendships\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userAId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userBId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFriendRequest\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumQuickPicStorage.FriendRequestStatus\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFriendsOfUser\",\"inputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"friendIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFriendship\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"friendshipId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userAId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userBId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessage\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedContent\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"contentType\",\"type\":\"uint8\",\"internalType\":\"enumQuickPicStorage.ContentType\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessagesForUser\",\"inputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessagesSentByUser\",\"inputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingRequestsForUser\",\"inputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUser\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"username\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"passwordHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserByUsername\",\"inputs\":[{\"name\":\"username\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"usernameOut\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"passwordHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserFriendships\",\"inputs\":[{\"name\":\"userId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messageIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messages\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedContent\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"contentType\",\"type\":\"uint8\",\"internalType\":\"enumQuickPicStorage.ContentType\"},{\"name\":\"signature\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messagesFromUser\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messagesToUser\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextUserNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingRequestsTo\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFriendRequestStatus\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumQuickPicStorage.FriendRequestStatus\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateUser\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"passwordHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publicKey\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userExists\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userFriendships\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usernameExists\",\"inputs\":[{\"name\":\"username\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usernameToId\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"users\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"username\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"passwordHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publicKey\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"FriendRequestCreated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FriendRequestUpdated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumQuickPicStorage.FriendRequestStatus\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FriendshipCreated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"userAId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"userBId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageCreated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"fromUserId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"toUserId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageDeleted\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UserCreated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"userNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"username\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UserUpdated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false}]",
}

// QuickPicStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use QuickPicStorageMetaData.ABI instead.
var QuickPicStorageABI = QuickPicStorageMetaData.ABI

// QuickPicStorage is an auto generated Go binding around an Ethereum contract.
type QuickPicStorage struct {
	QuickPicStorageCaller     // Read-only binding to the contract
	QuickPicStorageTransactor // Write-only binding to the contract
	QuickPicStorageFilterer   // Log filterer for contract events
}

// QuickPicStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type QuickPicStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuickPicStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QuickPicStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuickPicStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuickPicStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuickPicStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuickPicStorageSession struct {
	Contract     *QuickPicStorage  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuickPicStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuickPicStorageCallerSession struct {
	Contract *QuickPicStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// QuickPicStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuickPicStorageTransactorSession struct {
	Contract     *QuickPicStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// QuickPicStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type QuickPicStorageRaw struct {
	Contract *QuickPicStorage // Generic contract binding to access the raw methods on
}

// QuickPicStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuickPicStorageCallerRaw struct {
	Contract *QuickPicStorageCaller // Generic read-only contract binding to access the raw methods on
}

// QuickPicStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuickPicStorageTransactorRaw struct {
	Contract *QuickPicStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQuickPicStorage creates a new instance of QuickPicStorage, bound to a specific deployed contract.
func NewQuickPicStorage(address common.Address, backend bind.ContractBackend) (*QuickPicStorage, error) {
	contract, err := bindQuickPicStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorage{QuickPicStorageCaller: QuickPicStorageCaller{contract: contract}, QuickPicStorageTransactor: QuickPicStorageTransactor{contract: contract}, QuickPicStorageFilterer: QuickPicStorageFilterer{contract: contract}}, nil
}

// NewQuickPicStorageCaller creates a new read-only instance of QuickPicStorage, bound to a specific deployed contract.
func NewQuickPicStorageCaller(address common.Address, caller bind.ContractCaller) (*QuickPicStorageCaller, error) {
	contract, err := bindQuickPicStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageCaller{contract: contract}, nil
}

// NewQuickPicStorageTransactor creates a new write-only instance of QuickPicStorage, bound to a specific deployed contract.
func NewQuickPicStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*QuickPicStorageTransactor, error) {
	contract, err := bindQuickPicStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageTransactor{contract: contract}, nil
}

// NewQuickPicStorageFilterer creates a new log filterer instance of QuickPicStorage, bound to a specific deployed contract.
func NewQuickPicStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*QuickPicStorageFilterer, error) {
	contract, err := bindQuickPicStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageFilterer{contract: contract}, nil
}

// bindQuickPicStorage binds a generic wrapper to an already deployed contract.
func bindQuickPicStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := QuickPicStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuickPicStorage *QuickPicStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QuickPicStorage.Contract.QuickPicStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuickPicStorage *QuickPicStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.QuickPicStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuickPicStorage *QuickPicStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.QuickPicStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuickPicStorage *QuickPicStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QuickPicStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuickPicStorage *QuickPicStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuickPicStorage *QuickPicStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.contract.Transact(opts, method, params...)
}

// AreFriends is a free data retrieval call binding the contract method 0xe1c32b2f.
//
// Solidity: function areFriends(bytes32 userId1, bytes32 userId2) view returns(bool)
func (_QuickPicStorage *QuickPicStorageCaller) AreFriends(opts *bind.CallOpts, userId1 [32]byte, userId2 [32]byte) (bool, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "areFriends", userId1, userId2)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AreFriends is a free data retrieval call binding the contract method 0xe1c32b2f.
//
// Solidity: function areFriends(bytes32 userId1, bytes32 userId2) view returns(bool)
func (_QuickPicStorage *QuickPicStorageSession) AreFriends(userId1 [32]byte, userId2 [32]byte) (bool, error) {
	return _QuickPicStorage.Contract.AreFriends(&_QuickPicStorage.CallOpts, userId1, userId2)
}

// AreFriends is a free data retrieval call binding the contract method 0xe1c32b2f.
//
// Solidity: function areFriends(bytes32 userId1, bytes32 userId2) view returns(bool)
func (_QuickPicStorage *QuickPicStorageCallerSession) AreFriends(userId1 [32]byte, userId2 [32]byte) (bool, error) {
	return _QuickPicStorage.Contract.AreFriends(&_QuickPicStorage.CallOpts, userId1, userId2)
}

// FriendRequestByUsers is a free data retrieval call binding the contract method 0xc6b50640.
//
// Solidity: function friendRequestByUsers(bytes32 , bytes32 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) FriendRequestByUsers(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "friendRequestByUsers", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FriendRequestByUsers is a free data retrieval call binding the contract method 0xc6b50640.
//
// Solidity: function friendRequestByUsers(bytes32 , bytes32 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) FriendRequestByUsers(arg0 [32]byte, arg1 [32]byte) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendRequestByUsers(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// FriendRequestByUsers is a free data retrieval call binding the contract method 0xc6b50640.
//
// Solidity: function friendRequestByUsers(bytes32 , bytes32 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) FriendRequestByUsers(arg0 [32]byte, arg1 [32]byte) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendRequestByUsers(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// FriendRequestIds is a free data retrieval call binding the contract method 0xb4ecac47.
//
// Solidity: function friendRequestIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) FriendRequestIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "friendRequestIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FriendRequestIds is a free data retrieval call binding the contract method 0xb4ecac47.
//
// Solidity: function friendRequestIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) FriendRequestIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendRequestIds(&_QuickPicStorage.CallOpts, arg0)
}

// FriendRequestIds is a free data retrieval call binding the contract method 0xb4ecac47.
//
// Solidity: function friendRequestIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) FriendRequestIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendRequestIds(&_QuickPicStorage.CallOpts, arg0)
}

// FriendRequests is a free data retrieval call binding the contract method 0xe44834eb.
//
// Solidity: function friendRequests(bytes32 ) view returns(bytes32 id, bytes32 fromUserId, bytes32 toUserId, uint8 status, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCaller) FriendRequests(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Id         [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Status     uint8
	CreatedAt  *big.Int
	Exists     bool
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "friendRequests", arg0)

	outstruct := new(struct {
		Id         [32]byte
		FromUserId [32]byte
		ToUserId   [32]byte
		Status     uint8
		CreatedAt  *big.Int
		Exists     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.FromUserId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ToUserId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Status = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.CreatedAt = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// FriendRequests is a free data retrieval call binding the contract method 0xe44834eb.
//
// Solidity: function friendRequests(bytes32 ) view returns(bytes32 id, bytes32 fromUserId, bytes32 toUserId, uint8 status, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageSession) FriendRequests(arg0 [32]byte) (struct {
	Id         [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Status     uint8
	CreatedAt  *big.Int
	Exists     bool
}, error) {
	return _QuickPicStorage.Contract.FriendRequests(&_QuickPicStorage.CallOpts, arg0)
}

// FriendRequests is a free data retrieval call binding the contract method 0xe44834eb.
//
// Solidity: function friendRequests(bytes32 ) view returns(bytes32 id, bytes32 fromUserId, bytes32 toUserId, uint8 status, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCallerSession) FriendRequests(arg0 [32]byte) (struct {
	Id         [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Status     uint8
	CreatedAt  *big.Int
	Exists     bool
}, error) {
	return _QuickPicStorage.Contract.FriendRequests(&_QuickPicStorage.CallOpts, arg0)
}

// FriendshipByUsers is a free data retrieval call binding the contract method 0xfec729e1.
//
// Solidity: function friendshipByUsers(bytes32 , bytes32 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) FriendshipByUsers(opts *bind.CallOpts, arg0 [32]byte, arg1 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "friendshipByUsers", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FriendshipByUsers is a free data retrieval call binding the contract method 0xfec729e1.
//
// Solidity: function friendshipByUsers(bytes32 , bytes32 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) FriendshipByUsers(arg0 [32]byte, arg1 [32]byte) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendshipByUsers(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// FriendshipByUsers is a free data retrieval call binding the contract method 0xfec729e1.
//
// Solidity: function friendshipByUsers(bytes32 , bytes32 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) FriendshipByUsers(arg0 [32]byte, arg1 [32]byte) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendshipByUsers(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// FriendshipIds is a free data retrieval call binding the contract method 0x687a6884.
//
// Solidity: function friendshipIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) FriendshipIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "friendshipIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FriendshipIds is a free data retrieval call binding the contract method 0x687a6884.
//
// Solidity: function friendshipIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) FriendshipIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendshipIds(&_QuickPicStorage.CallOpts, arg0)
}

// FriendshipIds is a free data retrieval call binding the contract method 0x687a6884.
//
// Solidity: function friendshipIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) FriendshipIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.FriendshipIds(&_QuickPicStorage.CallOpts, arg0)
}

// Friendships is a free data retrieval call binding the contract method 0xd65a128d.
//
// Solidity: function friendships(bytes32 ) view returns(bytes32 id, bytes32 userAId, bytes32 userBId, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCaller) Friendships(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Id        [32]byte
	UserAId   [32]byte
	UserBId   [32]byte
	CreatedAt *big.Int
	Exists    bool
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "friendships", arg0)

	outstruct := new(struct {
		Id        [32]byte
		UserAId   [32]byte
		UserBId   [32]byte
		CreatedAt *big.Int
		Exists    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.UserAId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.UserBId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.CreatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Friendships is a free data retrieval call binding the contract method 0xd65a128d.
//
// Solidity: function friendships(bytes32 ) view returns(bytes32 id, bytes32 userAId, bytes32 userBId, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageSession) Friendships(arg0 [32]byte) (struct {
	Id        [32]byte
	UserAId   [32]byte
	UserBId   [32]byte
	CreatedAt *big.Int
	Exists    bool
}, error) {
	return _QuickPicStorage.Contract.Friendships(&_QuickPicStorage.CallOpts, arg0)
}

// Friendships is a free data retrieval call binding the contract method 0xd65a128d.
//
// Solidity: function friendships(bytes32 ) view returns(bytes32 id, bytes32 userAId, bytes32 userBId, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCallerSession) Friendships(arg0 [32]byte) (struct {
	Id        [32]byte
	UserAId   [32]byte
	UserBId   [32]byte
	CreatedAt *big.Int
	Exists    bool
}, error) {
	return _QuickPicStorage.Contract.Friendships(&_QuickPicStorage.CallOpts, arg0)
}

// GetFriendRequest is a free data retrieval call binding the contract method 0xea65a758.
//
// Solidity: function getFriendRequest(bytes32 id) view returns(bytes32 requestId, bytes32 fromUserId, bytes32 toUserId, uint8 status, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageCaller) GetFriendRequest(opts *bind.CallOpts, id [32]byte) (struct {
	RequestId  [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Status     uint8
	CreatedAt  *big.Int
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getFriendRequest", id)

	outstruct := new(struct {
		RequestId  [32]byte
		FromUserId [32]byte
		ToUserId   [32]byte
		Status     uint8
		CreatedAt  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.FromUserId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ToUserId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Status = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.CreatedAt = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetFriendRequest is a free data retrieval call binding the contract method 0xea65a758.
//
// Solidity: function getFriendRequest(bytes32 id) view returns(bytes32 requestId, bytes32 fromUserId, bytes32 toUserId, uint8 status, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageSession) GetFriendRequest(id [32]byte) (struct {
	RequestId  [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Status     uint8
	CreatedAt  *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetFriendRequest(&_QuickPicStorage.CallOpts, id)
}

// GetFriendRequest is a free data retrieval call binding the contract method 0xea65a758.
//
// Solidity: function getFriendRequest(bytes32 id) view returns(bytes32 requestId, bytes32 fromUserId, bytes32 toUserId, uint8 status, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetFriendRequest(id [32]byte) (struct {
	RequestId  [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Status     uint8
	CreatedAt  *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetFriendRequest(&_QuickPicStorage.CallOpts, id)
}

// GetFriendsOfUser is a free data retrieval call binding the contract method 0x56ad974f.
//
// Solidity: function getFriendsOfUser(bytes32 userId) view returns(bytes32[] friendIds)
func (_QuickPicStorage *QuickPicStorageCaller) GetFriendsOfUser(opts *bind.CallOpts, userId [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getFriendsOfUser", userId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetFriendsOfUser is a free data retrieval call binding the contract method 0x56ad974f.
//
// Solidity: function getFriendsOfUser(bytes32 userId) view returns(bytes32[] friendIds)
func (_QuickPicStorage *QuickPicStorageSession) GetFriendsOfUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetFriendsOfUser(&_QuickPicStorage.CallOpts, userId)
}

// GetFriendsOfUser is a free data retrieval call binding the contract method 0x56ad974f.
//
// Solidity: function getFriendsOfUser(bytes32 userId) view returns(bytes32[] friendIds)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetFriendsOfUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetFriendsOfUser(&_QuickPicStorage.CallOpts, userId)
}

// GetFriendship is a free data retrieval call binding the contract method 0xddf7ef75.
//
// Solidity: function getFriendship(bytes32 id) view returns(bytes32 friendshipId, bytes32 userAId, bytes32 userBId, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageCaller) GetFriendship(opts *bind.CallOpts, id [32]byte) (struct {
	FriendshipId [32]byte
	UserAId      [32]byte
	UserBId      [32]byte
	CreatedAt    *big.Int
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getFriendship", id)

	outstruct := new(struct {
		FriendshipId [32]byte
		UserAId      [32]byte
		UserBId      [32]byte
		CreatedAt    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FriendshipId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.UserAId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.UserBId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.CreatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetFriendship is a free data retrieval call binding the contract method 0xddf7ef75.
//
// Solidity: function getFriendship(bytes32 id) view returns(bytes32 friendshipId, bytes32 userAId, bytes32 userBId, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageSession) GetFriendship(id [32]byte) (struct {
	FriendshipId [32]byte
	UserAId      [32]byte
	UserBId      [32]byte
	CreatedAt    *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetFriendship(&_QuickPicStorage.CallOpts, id)
}

// GetFriendship is a free data retrieval call binding the contract method 0xddf7ef75.
//
// Solidity: function getFriendship(bytes32 id) view returns(bytes32 friendshipId, bytes32 userAId, bytes32 userBId, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetFriendship(id [32]byte) (struct {
	FriendshipId [32]byte
	UserAId      [32]byte
	UserBId      [32]byte
	CreatedAt    *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetFriendship(&_QuickPicStorage.CallOpts, id)
}

// GetMessage is a free data retrieval call binding the contract method 0x0139a221.
//
// Solidity: function getMessage(bytes32 id) view returns(bytes32 messageId, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageCaller) GetMessage(opts *bind.CallOpts, id [32]byte) (struct {
	MessageId        [32]byte
	FromUserId       [32]byte
	ToUserId         [32]byte
	EncryptedContent []byte
	ContentType      uint8
	Signature        string
	CreatedAt        *big.Int
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getMessage", id)

	outstruct := new(struct {
		MessageId        [32]byte
		FromUserId       [32]byte
		ToUserId         [32]byte
		EncryptedContent []byte
		ContentType      uint8
		Signature        string
		CreatedAt        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MessageId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.FromUserId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ToUserId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.EncryptedContent = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.ContentType = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.Signature = *abi.ConvertType(out[5], new(string)).(*string)
	outstruct.CreatedAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetMessage is a free data retrieval call binding the contract method 0x0139a221.
//
// Solidity: function getMessage(bytes32 id) view returns(bytes32 messageId, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageSession) GetMessage(id [32]byte) (struct {
	MessageId        [32]byte
	FromUserId       [32]byte
	ToUserId         [32]byte
	EncryptedContent []byte
	ContentType      uint8
	Signature        string
	CreatedAt        *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetMessage(&_QuickPicStorage.CallOpts, id)
}

// GetMessage is a free data retrieval call binding the contract method 0x0139a221.
//
// Solidity: function getMessage(bytes32 id) view returns(bytes32 messageId, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature, uint256 createdAt)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetMessage(id [32]byte) (struct {
	MessageId        [32]byte
	FromUserId       [32]byte
	ToUserId         [32]byte
	EncryptedContent []byte
	ContentType      uint8
	Signature        string
	CreatedAt        *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetMessage(&_QuickPicStorage.CallOpts, id)
}

// GetMessagesForUser is a free data retrieval call binding the contract method 0x9e2243ef.
//
// Solidity: function getMessagesForUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCaller) GetMessagesForUser(opts *bind.CallOpts, userId [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getMessagesForUser", userId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetMessagesForUser is a free data retrieval call binding the contract method 0x9e2243ef.
//
// Solidity: function getMessagesForUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageSession) GetMessagesForUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetMessagesForUser(&_QuickPicStorage.CallOpts, userId)
}

// GetMessagesForUser is a free data retrieval call binding the contract method 0x9e2243ef.
//
// Solidity: function getMessagesForUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCallerSession) GetMessagesForUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetMessagesForUser(&_QuickPicStorage.CallOpts, userId)
}

// GetMessagesSentByUser is a free data retrieval call binding the contract method 0x563485ed.
//
// Solidity: function getMessagesSentByUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCaller) GetMessagesSentByUser(opts *bind.CallOpts, userId [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getMessagesSentByUser", userId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetMessagesSentByUser is a free data retrieval call binding the contract method 0x563485ed.
//
// Solidity: function getMessagesSentByUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageSession) GetMessagesSentByUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetMessagesSentByUser(&_QuickPicStorage.CallOpts, userId)
}

// GetMessagesSentByUser is a free data retrieval call binding the contract method 0x563485ed.
//
// Solidity: function getMessagesSentByUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCallerSession) GetMessagesSentByUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetMessagesSentByUser(&_QuickPicStorage.CallOpts, userId)
}

// GetPendingRequestsForUser is a free data retrieval call binding the contract method 0xdb0088ac.
//
// Solidity: function getPendingRequestsForUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCaller) GetPendingRequestsForUser(opts *bind.CallOpts, userId [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getPendingRequestsForUser", userId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetPendingRequestsForUser is a free data retrieval call binding the contract method 0xdb0088ac.
//
// Solidity: function getPendingRequestsForUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageSession) GetPendingRequestsForUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetPendingRequestsForUser(&_QuickPicStorage.CallOpts, userId)
}

// GetPendingRequestsForUser is a free data retrieval call binding the contract method 0xdb0088ac.
//
// Solidity: function getPendingRequestsForUser(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCallerSession) GetPendingRequestsForUser(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetPendingRequestsForUser(&_QuickPicStorage.CallOpts, userId)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 id) view returns(bytes32 userId, uint256 userNumber, string username, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt)
func (_QuickPicStorage *QuickPicStorageCaller) GetUser(opts *bind.CallOpts, id [32]byte) (struct {
	UserId       [32]byte
	UserNumber   *big.Int
	Username     string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getUser", id)

	outstruct := new(struct {
		UserId       [32]byte
		UserNumber   *big.Int
		Username     string
		PasswordHash string
		PublicKey    string
		CreatedAt    *big.Int
		UpdatedAt    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UserId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.UserNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Username = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.PasswordHash = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.PublicKey = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.CreatedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 id) view returns(bytes32 userId, uint256 userNumber, string username, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt)
func (_QuickPicStorage *QuickPicStorageSession) GetUser(id [32]byte) (struct {
	UserId       [32]byte
	UserNumber   *big.Int
	Username     string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetUser(&_QuickPicStorage.CallOpts, id)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 id) view returns(bytes32 userId, uint256 userNumber, string username, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetUser(id [32]byte) (struct {
	UserId       [32]byte
	UserNumber   *big.Int
	Username     string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetUser(&_QuickPicStorage.CallOpts, id)
}

// GetUserByUsername is a free data retrieval call binding the contract method 0x9f42af63.
//
// Solidity: function getUserByUsername(string username) view returns(bytes32 userId, uint256 userNumber, string usernameOut, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt)
func (_QuickPicStorage *QuickPicStorageCaller) GetUserByUsername(opts *bind.CallOpts, username string) (struct {
	UserId       [32]byte
	UserNumber   *big.Int
	UsernameOut  string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getUserByUsername", username)

	outstruct := new(struct {
		UserId       [32]byte
		UserNumber   *big.Int
		UsernameOut  string
		PasswordHash string
		PublicKey    string
		CreatedAt    *big.Int
		UpdatedAt    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UserId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.UserNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UsernameOut = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.PasswordHash = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.PublicKey = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.CreatedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUserByUsername is a free data retrieval call binding the contract method 0x9f42af63.
//
// Solidity: function getUserByUsername(string username) view returns(bytes32 userId, uint256 userNumber, string usernameOut, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt)
func (_QuickPicStorage *QuickPicStorageSession) GetUserByUsername(username string) (struct {
	UserId       [32]byte
	UserNumber   *big.Int
	UsernameOut  string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetUserByUsername(&_QuickPicStorage.CallOpts, username)
}

// GetUserByUsername is a free data retrieval call binding the contract method 0x9f42af63.
//
// Solidity: function getUserByUsername(string username) view returns(bytes32 userId, uint256 userNumber, string usernameOut, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetUserByUsername(username string) (struct {
	UserId       [32]byte
	UserNumber   *big.Int
	UsernameOut  string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}, error) {
	return _QuickPicStorage.Contract.GetUserByUsername(&_QuickPicStorage.CallOpts, username)
}

// GetUserCount is a free data retrieval call binding the contract method 0xb5cb15f7.
//
// Solidity: function getUserCount() view returns(uint256)
func (_QuickPicStorage *QuickPicStorageCaller) GetUserCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getUserCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserCount is a free data retrieval call binding the contract method 0xb5cb15f7.
//
// Solidity: function getUserCount() view returns(uint256)
func (_QuickPicStorage *QuickPicStorageSession) GetUserCount() (*big.Int, error) {
	return _QuickPicStorage.Contract.GetUserCount(&_QuickPicStorage.CallOpts)
}

// GetUserCount is a free data retrieval call binding the contract method 0xb5cb15f7.
//
// Solidity: function getUserCount() view returns(uint256)
func (_QuickPicStorage *QuickPicStorageCallerSession) GetUserCount() (*big.Int, error) {
	return _QuickPicStorage.Contract.GetUserCount(&_QuickPicStorage.CallOpts)
}

// GetUserFriendships is a free data retrieval call binding the contract method 0xaf3c0eca.
//
// Solidity: function getUserFriendships(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCaller) GetUserFriendships(opts *bind.CallOpts, userId [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "getUserFriendships", userId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetUserFriendships is a free data retrieval call binding the contract method 0xaf3c0eca.
//
// Solidity: function getUserFriendships(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageSession) GetUserFriendships(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetUserFriendships(&_QuickPicStorage.CallOpts, userId)
}

// GetUserFriendships is a free data retrieval call binding the contract method 0xaf3c0eca.
//
// Solidity: function getUserFriendships(bytes32 userId) view returns(bytes32[])
func (_QuickPicStorage *QuickPicStorageCallerSession) GetUserFriendships(userId [32]byte) ([][32]byte, error) {
	return _QuickPicStorage.Contract.GetUserFriendships(&_QuickPicStorage.CallOpts, userId)
}

// MessageIds is a free data retrieval call binding the contract method 0x9d17b9c5.
//
// Solidity: function messageIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) MessageIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "messageIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MessageIds is a free data retrieval call binding the contract method 0x9d17b9c5.
//
// Solidity: function messageIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) MessageIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.MessageIds(&_QuickPicStorage.CallOpts, arg0)
}

// MessageIds is a free data retrieval call binding the contract method 0x9d17b9c5.
//
// Solidity: function messageIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) MessageIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.MessageIds(&_QuickPicStorage.CallOpts, arg0)
}

// Messages is a free data retrieval call binding the contract method 0x2bbd59ca.
//
// Solidity: function messages(bytes32 ) view returns(bytes32 id, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCaller) Messages(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Id               [32]byte
	FromUserId       [32]byte
	ToUserId         [32]byte
	EncryptedContent []byte
	ContentType      uint8
	Signature        string
	CreatedAt        *big.Int
	Exists           bool
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "messages", arg0)

	outstruct := new(struct {
		Id               [32]byte
		FromUserId       [32]byte
		ToUserId         [32]byte
		EncryptedContent []byte
		ContentType      uint8
		Signature        string
		CreatedAt        *big.Int
		Exists           bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.FromUserId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ToUserId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.EncryptedContent = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.ContentType = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.Signature = *abi.ConvertType(out[5], new(string)).(*string)
	outstruct.CreatedAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[7], new(bool)).(*bool)

	return *outstruct, err

}

// Messages is a free data retrieval call binding the contract method 0x2bbd59ca.
//
// Solidity: function messages(bytes32 ) view returns(bytes32 id, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageSession) Messages(arg0 [32]byte) (struct {
	Id               [32]byte
	FromUserId       [32]byte
	ToUserId         [32]byte
	EncryptedContent []byte
	ContentType      uint8
	Signature        string
	CreatedAt        *big.Int
	Exists           bool
}, error) {
	return _QuickPicStorage.Contract.Messages(&_QuickPicStorage.CallOpts, arg0)
}

// Messages is a free data retrieval call binding the contract method 0x2bbd59ca.
//
// Solidity: function messages(bytes32 ) view returns(bytes32 id, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature, uint256 createdAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCallerSession) Messages(arg0 [32]byte) (struct {
	Id               [32]byte
	FromUserId       [32]byte
	ToUserId         [32]byte
	EncryptedContent []byte
	ContentType      uint8
	Signature        string
	CreatedAt        *big.Int
	Exists           bool
}, error) {
	return _QuickPicStorage.Contract.Messages(&_QuickPicStorage.CallOpts, arg0)
}

// MessagesFromUser is a free data retrieval call binding the contract method 0xe867dc53.
//
// Solidity: function messagesFromUser(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) MessagesFromUser(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "messagesFromUser", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MessagesFromUser is a free data retrieval call binding the contract method 0xe867dc53.
//
// Solidity: function messagesFromUser(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) MessagesFromUser(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.MessagesFromUser(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// MessagesFromUser is a free data retrieval call binding the contract method 0xe867dc53.
//
// Solidity: function messagesFromUser(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) MessagesFromUser(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.MessagesFromUser(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// MessagesToUser is a free data retrieval call binding the contract method 0x7dc45b09.
//
// Solidity: function messagesToUser(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) MessagesToUser(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "messagesToUser", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MessagesToUser is a free data retrieval call binding the contract method 0x7dc45b09.
//
// Solidity: function messagesToUser(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) MessagesToUser(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.MessagesToUser(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// MessagesToUser is a free data retrieval call binding the contract method 0x7dc45b09.
//
// Solidity: function messagesToUser(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) MessagesToUser(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.MessagesToUser(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// NextUserNumber is a free data retrieval call binding the contract method 0x0e7b883d.
//
// Solidity: function nextUserNumber() view returns(uint256)
func (_QuickPicStorage *QuickPicStorageCaller) NextUserNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "nextUserNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextUserNumber is a free data retrieval call binding the contract method 0x0e7b883d.
//
// Solidity: function nextUserNumber() view returns(uint256)
func (_QuickPicStorage *QuickPicStorageSession) NextUserNumber() (*big.Int, error) {
	return _QuickPicStorage.Contract.NextUserNumber(&_QuickPicStorage.CallOpts)
}

// NextUserNumber is a free data retrieval call binding the contract method 0x0e7b883d.
//
// Solidity: function nextUserNumber() view returns(uint256)
func (_QuickPicStorage *QuickPicStorageCallerSession) NextUserNumber() (*big.Int, error) {
	return _QuickPicStorage.Contract.NextUserNumber(&_QuickPicStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_QuickPicStorage *QuickPicStorageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_QuickPicStorage *QuickPicStorageSession) Owner() (common.Address, error) {
	return _QuickPicStorage.Contract.Owner(&_QuickPicStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_QuickPicStorage *QuickPicStorageCallerSession) Owner() (common.Address, error) {
	return _QuickPicStorage.Contract.Owner(&_QuickPicStorage.CallOpts)
}

// PendingRequestsTo is a free data retrieval call binding the contract method 0x74b32b1c.
//
// Solidity: function pendingRequestsTo(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) PendingRequestsTo(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "pendingRequestsTo", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PendingRequestsTo is a free data retrieval call binding the contract method 0x74b32b1c.
//
// Solidity: function pendingRequestsTo(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) PendingRequestsTo(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.PendingRequestsTo(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// PendingRequestsTo is a free data retrieval call binding the contract method 0x74b32b1c.
//
// Solidity: function pendingRequestsTo(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) PendingRequestsTo(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.PendingRequestsTo(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// UserExists is a free data retrieval call binding the contract method 0xa2e8452c.
//
// Solidity: function userExists(bytes32 id) view returns(bool)
func (_QuickPicStorage *QuickPicStorageCaller) UserExists(opts *bind.CallOpts, id [32]byte) (bool, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "userExists", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UserExists is a free data retrieval call binding the contract method 0xa2e8452c.
//
// Solidity: function userExists(bytes32 id) view returns(bool)
func (_QuickPicStorage *QuickPicStorageSession) UserExists(id [32]byte) (bool, error) {
	return _QuickPicStorage.Contract.UserExists(&_QuickPicStorage.CallOpts, id)
}

// UserExists is a free data retrieval call binding the contract method 0xa2e8452c.
//
// Solidity: function userExists(bytes32 id) view returns(bool)
func (_QuickPicStorage *QuickPicStorageCallerSession) UserExists(id [32]byte) (bool, error) {
	return _QuickPicStorage.Contract.UserExists(&_QuickPicStorage.CallOpts, id)
}

// UserFriendships is a free data retrieval call binding the contract method 0x140cc069.
//
// Solidity: function userFriendships(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) UserFriendships(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "userFriendships", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UserFriendships is a free data retrieval call binding the contract method 0x140cc069.
//
// Solidity: function userFriendships(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) UserFriendships(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.UserFriendships(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// UserFriendships is a free data retrieval call binding the contract method 0x140cc069.
//
// Solidity: function userFriendships(bytes32 , uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) UserFriendships(arg0 [32]byte, arg1 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.UserFriendships(&_QuickPicStorage.CallOpts, arg0, arg1)
}

// UserIds is a free data retrieval call binding the contract method 0x4635fd68.
//
// Solidity: function userIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) UserIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "userIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UserIds is a free data retrieval call binding the contract method 0x4635fd68.
//
// Solidity: function userIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) UserIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.UserIds(&_QuickPicStorage.CallOpts, arg0)
}

// UserIds is a free data retrieval call binding the contract method 0x4635fd68.
//
// Solidity: function userIds(uint256 ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) UserIds(arg0 *big.Int) ([32]byte, error) {
	return _QuickPicStorage.Contract.UserIds(&_QuickPicStorage.CallOpts, arg0)
}

// UsernameExists is a free data retrieval call binding the contract method 0xf309e3f9.
//
// Solidity: function usernameExists(string username) view returns(bool)
func (_QuickPicStorage *QuickPicStorageCaller) UsernameExists(opts *bind.CallOpts, username string) (bool, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "usernameExists", username)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsernameExists is a free data retrieval call binding the contract method 0xf309e3f9.
//
// Solidity: function usernameExists(string username) view returns(bool)
func (_QuickPicStorage *QuickPicStorageSession) UsernameExists(username string) (bool, error) {
	return _QuickPicStorage.Contract.UsernameExists(&_QuickPicStorage.CallOpts, username)
}

// UsernameExists is a free data retrieval call binding the contract method 0xf309e3f9.
//
// Solidity: function usernameExists(string username) view returns(bool)
func (_QuickPicStorage *QuickPicStorageCallerSession) UsernameExists(username string) (bool, error) {
	return _QuickPicStorage.Contract.UsernameExists(&_QuickPicStorage.CallOpts, username)
}

// UsernameToId is a free data retrieval call binding the contract method 0x5e1fe1db.
//
// Solidity: function usernameToId(string ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCaller) UsernameToId(opts *bind.CallOpts, arg0 string) ([32]byte, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "usernameToId", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UsernameToId is a free data retrieval call binding the contract method 0x5e1fe1db.
//
// Solidity: function usernameToId(string ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageSession) UsernameToId(arg0 string) ([32]byte, error) {
	return _QuickPicStorage.Contract.UsernameToId(&_QuickPicStorage.CallOpts, arg0)
}

// UsernameToId is a free data retrieval call binding the contract method 0x5e1fe1db.
//
// Solidity: function usernameToId(string ) view returns(bytes32)
func (_QuickPicStorage *QuickPicStorageCallerSession) UsernameToId(arg0 string) ([32]byte, error) {
	return _QuickPicStorage.Contract.UsernameToId(&_QuickPicStorage.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xcea6ab98.
//
// Solidity: function users(bytes32 ) view returns(bytes32 id, uint256 userNumber, string username, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCaller) Users(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Id           [32]byte
	UserNumber   *big.Int
	Username     string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
	Exists       bool
}, error) {
	var out []interface{}
	err := _QuickPicStorage.contract.Call(opts, &out, "users", arg0)

	outstruct := new(struct {
		Id           [32]byte
		UserNumber   *big.Int
		Username     string
		PasswordHash string
		PublicKey    string
		CreatedAt    *big.Int
		UpdatedAt    *big.Int
		Exists       bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.UserNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Username = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.PasswordHash = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.PublicKey = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.CreatedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[7], new(bool)).(*bool)

	return *outstruct, err

}

// Users is a free data retrieval call binding the contract method 0xcea6ab98.
//
// Solidity: function users(bytes32 ) view returns(bytes32 id, uint256 userNumber, string username, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt, bool exists)
func (_QuickPicStorage *QuickPicStorageSession) Users(arg0 [32]byte) (struct {
	Id           [32]byte
	UserNumber   *big.Int
	Username     string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
	Exists       bool
}, error) {
	return _QuickPicStorage.Contract.Users(&_QuickPicStorage.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xcea6ab98.
//
// Solidity: function users(bytes32 ) view returns(bytes32 id, uint256 userNumber, string username, string passwordHash, string publicKey, uint256 createdAt, uint256 updatedAt, bool exists)
func (_QuickPicStorage *QuickPicStorageCallerSession) Users(arg0 [32]byte) (struct {
	Id           [32]byte
	UserNumber   *big.Int
	Username     string
	PasswordHash string
	PublicKey    string
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
	Exists       bool
}, error) {
	return _QuickPicStorage.Contract.Users(&_QuickPicStorage.CallOpts, arg0)
}

// CreateFriendRequest is a paid mutator transaction binding the contract method 0xe6924351.
//
// Solidity: function createFriendRequest(bytes32 id, bytes32 fromUserId, bytes32 toUserId) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) CreateFriendRequest(opts *bind.TransactOpts, id [32]byte, fromUserId [32]byte, toUserId [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "createFriendRequest", id, fromUserId, toUserId)
}

// CreateFriendRequest is a paid mutator transaction binding the contract method 0xe6924351.
//
// Solidity: function createFriendRequest(bytes32 id, bytes32 fromUserId, bytes32 toUserId) returns()
func (_QuickPicStorage *QuickPicStorageSession) CreateFriendRequest(id [32]byte, fromUserId [32]byte, toUserId [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateFriendRequest(&_QuickPicStorage.TransactOpts, id, fromUserId, toUserId)
}

// CreateFriendRequest is a paid mutator transaction binding the contract method 0xe6924351.
//
// Solidity: function createFriendRequest(bytes32 id, bytes32 fromUserId, bytes32 toUserId) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) CreateFriendRequest(id [32]byte, fromUserId [32]byte, toUserId [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateFriendRequest(&_QuickPicStorage.TransactOpts, id, fromUserId, toUserId)
}

// CreateFriendship is a paid mutator transaction binding the contract method 0x60b7ef48.
//
// Solidity: function createFriendship(bytes32 id, bytes32 userAId, bytes32 userBId) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) CreateFriendship(opts *bind.TransactOpts, id [32]byte, userAId [32]byte, userBId [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "createFriendship", id, userAId, userBId)
}

// CreateFriendship is a paid mutator transaction binding the contract method 0x60b7ef48.
//
// Solidity: function createFriendship(bytes32 id, bytes32 userAId, bytes32 userBId) returns()
func (_QuickPicStorage *QuickPicStorageSession) CreateFriendship(id [32]byte, userAId [32]byte, userBId [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateFriendship(&_QuickPicStorage.TransactOpts, id, userAId, userBId)
}

// CreateFriendship is a paid mutator transaction binding the contract method 0x60b7ef48.
//
// Solidity: function createFriendship(bytes32 id, bytes32 userAId, bytes32 userBId) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) CreateFriendship(id [32]byte, userAId [32]byte, userBId [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateFriendship(&_QuickPicStorage.TransactOpts, id, userAId, userBId)
}

// CreateMessage is a paid mutator transaction binding the contract method 0x14782c61.
//
// Solidity: function createMessage(bytes32 id, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) CreateMessage(opts *bind.TransactOpts, id [32]byte, fromUserId [32]byte, toUserId [32]byte, encryptedContent []byte, contentType uint8, signature string) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "createMessage", id, fromUserId, toUserId, encryptedContent, contentType, signature)
}

// CreateMessage is a paid mutator transaction binding the contract method 0x14782c61.
//
// Solidity: function createMessage(bytes32 id, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature) returns()
func (_QuickPicStorage *QuickPicStorageSession) CreateMessage(id [32]byte, fromUserId [32]byte, toUserId [32]byte, encryptedContent []byte, contentType uint8, signature string) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateMessage(&_QuickPicStorage.TransactOpts, id, fromUserId, toUserId, encryptedContent, contentType, signature)
}

// CreateMessage is a paid mutator transaction binding the contract method 0x14782c61.
//
// Solidity: function createMessage(bytes32 id, bytes32 fromUserId, bytes32 toUserId, bytes encryptedContent, uint8 contentType, string signature) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) CreateMessage(id [32]byte, fromUserId [32]byte, toUserId [32]byte, encryptedContent []byte, contentType uint8, signature string) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateMessage(&_QuickPicStorage.TransactOpts, id, fromUserId, toUserId, encryptedContent, contentType, signature)
}

// CreateUser is a paid mutator transaction binding the contract method 0x765f5d35.
//
// Solidity: function createUser(bytes32 id, string username, string passwordHash, string publicKey) returns(uint256 userNumber)
func (_QuickPicStorage *QuickPicStorageTransactor) CreateUser(opts *bind.TransactOpts, id [32]byte, username string, passwordHash string, publicKey string) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "createUser", id, username, passwordHash, publicKey)
}

// CreateUser is a paid mutator transaction binding the contract method 0x765f5d35.
//
// Solidity: function createUser(bytes32 id, string username, string passwordHash, string publicKey) returns(uint256 userNumber)
func (_QuickPicStorage *QuickPicStorageSession) CreateUser(id [32]byte, username string, passwordHash string, publicKey string) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateUser(&_QuickPicStorage.TransactOpts, id, username, passwordHash, publicKey)
}

// CreateUser is a paid mutator transaction binding the contract method 0x765f5d35.
//
// Solidity: function createUser(bytes32 id, string username, string passwordHash, string publicKey) returns(uint256 userNumber)
func (_QuickPicStorage *QuickPicStorageTransactorSession) CreateUser(id [32]byte, username string, passwordHash string, publicKey string) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.CreateUser(&_QuickPicStorage.TransactOpts, id, username, passwordHash, publicKey)
}

// DeleteMessage is a paid mutator transaction binding the contract method 0xfe1e3eca.
//
// Solidity: function deleteMessage(bytes32 id) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) DeleteMessage(opts *bind.TransactOpts, id [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "deleteMessage", id)
}

// DeleteMessage is a paid mutator transaction binding the contract method 0xfe1e3eca.
//
// Solidity: function deleteMessage(bytes32 id) returns()
func (_QuickPicStorage *QuickPicStorageSession) DeleteMessage(id [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.DeleteMessage(&_QuickPicStorage.TransactOpts, id)
}

// DeleteMessage is a paid mutator transaction binding the contract method 0xfe1e3eca.
//
// Solidity: function deleteMessage(bytes32 id) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) DeleteMessage(id [32]byte) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.DeleteMessage(&_QuickPicStorage.TransactOpts, id)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_QuickPicStorage *QuickPicStorageSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.TransferOwnership(&_QuickPicStorage.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.TransferOwnership(&_QuickPicStorage.TransactOpts, newOwner)
}

// UpdateFriendRequestStatus is a paid mutator transaction binding the contract method 0x4018c162.
//
// Solidity: function updateFriendRequestStatus(bytes32 id, uint8 status) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) UpdateFriendRequestStatus(opts *bind.TransactOpts, id [32]byte, status uint8) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "updateFriendRequestStatus", id, status)
}

// UpdateFriendRequestStatus is a paid mutator transaction binding the contract method 0x4018c162.
//
// Solidity: function updateFriendRequestStatus(bytes32 id, uint8 status) returns()
func (_QuickPicStorage *QuickPicStorageSession) UpdateFriendRequestStatus(id [32]byte, status uint8) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.UpdateFriendRequestStatus(&_QuickPicStorage.TransactOpts, id, status)
}

// UpdateFriendRequestStatus is a paid mutator transaction binding the contract method 0x4018c162.
//
// Solidity: function updateFriendRequestStatus(bytes32 id, uint8 status) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) UpdateFriendRequestStatus(id [32]byte, status uint8) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.UpdateFriendRequestStatus(&_QuickPicStorage.TransactOpts, id, status)
}

// UpdateUser is a paid mutator transaction binding the contract method 0x494fb8e7.
//
// Solidity: function updateUser(bytes32 id, string passwordHash, string publicKey) returns()
func (_QuickPicStorage *QuickPicStorageTransactor) UpdateUser(opts *bind.TransactOpts, id [32]byte, passwordHash string, publicKey string) (*types.Transaction, error) {
	return _QuickPicStorage.contract.Transact(opts, "updateUser", id, passwordHash, publicKey)
}

// UpdateUser is a paid mutator transaction binding the contract method 0x494fb8e7.
//
// Solidity: function updateUser(bytes32 id, string passwordHash, string publicKey) returns()
func (_QuickPicStorage *QuickPicStorageSession) UpdateUser(id [32]byte, passwordHash string, publicKey string) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.UpdateUser(&_QuickPicStorage.TransactOpts, id, passwordHash, publicKey)
}

// UpdateUser is a paid mutator transaction binding the contract method 0x494fb8e7.
//
// Solidity: function updateUser(bytes32 id, string passwordHash, string publicKey) returns()
func (_QuickPicStorage *QuickPicStorageTransactorSession) UpdateUser(id [32]byte, passwordHash string, publicKey string) (*types.Transaction, error) {
	return _QuickPicStorage.Contract.UpdateUser(&_QuickPicStorage.TransactOpts, id, passwordHash, publicKey)
}

// QuickPicStorageFriendRequestCreatedIterator is returned from FilterFriendRequestCreated and is used to iterate over the raw logs and unpacked data for FriendRequestCreated events raised by the QuickPicStorage contract.
type QuickPicStorageFriendRequestCreatedIterator struct {
	Event *QuickPicStorageFriendRequestCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageFriendRequestCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageFriendRequestCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageFriendRequestCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageFriendRequestCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageFriendRequestCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageFriendRequestCreated represents a FriendRequestCreated event raised by the QuickPicStorage contract.
type QuickPicStorageFriendRequestCreated struct {
	Id         [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterFriendRequestCreated is a free log retrieval operation binding the contract event 0x9ac14503d7f332c5dbec9975c7e73dc5a8473298dbc7abc78c5997f410d2777c.
//
// Solidity: event FriendRequestCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterFriendRequestCreated(opts *bind.FilterOpts, id [][32]byte, fromUserId [][32]byte, toUserId [][32]byte) (*QuickPicStorageFriendRequestCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromUserIdRule []interface{}
	for _, fromUserIdItem := range fromUserId {
		fromUserIdRule = append(fromUserIdRule, fromUserIdItem)
	}
	var toUserIdRule []interface{}
	for _, toUserIdItem := range toUserId {
		toUserIdRule = append(toUserIdRule, toUserIdItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "FriendRequestCreated", idRule, fromUserIdRule, toUserIdRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageFriendRequestCreatedIterator{contract: _QuickPicStorage.contract, event: "FriendRequestCreated", logs: logs, sub: sub}, nil
}

// WatchFriendRequestCreated is a free log subscription operation binding the contract event 0x9ac14503d7f332c5dbec9975c7e73dc5a8473298dbc7abc78c5997f410d2777c.
//
// Solidity: event FriendRequestCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchFriendRequestCreated(opts *bind.WatchOpts, sink chan<- *QuickPicStorageFriendRequestCreated, id [][32]byte, fromUserId [][32]byte, toUserId [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromUserIdRule []interface{}
	for _, fromUserIdItem := range fromUserId {
		fromUserIdRule = append(fromUserIdRule, fromUserIdItem)
	}
	var toUserIdRule []interface{}
	for _, toUserIdItem := range toUserId {
		toUserIdRule = append(toUserIdRule, toUserIdItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "FriendRequestCreated", idRule, fromUserIdRule, toUserIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageFriendRequestCreated)
				if err := _QuickPicStorage.contract.UnpackLog(event, "FriendRequestCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFriendRequestCreated is a log parse operation binding the contract event 0x9ac14503d7f332c5dbec9975c7e73dc5a8473298dbc7abc78c5997f410d2777c.
//
// Solidity: event FriendRequestCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseFriendRequestCreated(log types.Log) (*QuickPicStorageFriendRequestCreated, error) {
	event := new(QuickPicStorageFriendRequestCreated)
	if err := _QuickPicStorage.contract.UnpackLog(event, "FriendRequestCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuickPicStorageFriendRequestUpdatedIterator is returned from FilterFriendRequestUpdated and is used to iterate over the raw logs and unpacked data for FriendRequestUpdated events raised by the QuickPicStorage contract.
type QuickPicStorageFriendRequestUpdatedIterator struct {
	Event *QuickPicStorageFriendRequestUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageFriendRequestUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageFriendRequestUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageFriendRequestUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageFriendRequestUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageFriendRequestUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageFriendRequestUpdated represents a FriendRequestUpdated event raised by the QuickPicStorage contract.
type QuickPicStorageFriendRequestUpdated struct {
	Id     [32]byte
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFriendRequestUpdated is a free log retrieval operation binding the contract event 0x8a2bf39e9087c33765bcca9a2c61caaa3e82956b669d0f3914a1d95e50f19ff2.
//
// Solidity: event FriendRequestUpdated(bytes32 indexed id, uint8 status)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterFriendRequestUpdated(opts *bind.FilterOpts, id [][32]byte) (*QuickPicStorageFriendRequestUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "FriendRequestUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageFriendRequestUpdatedIterator{contract: _QuickPicStorage.contract, event: "FriendRequestUpdated", logs: logs, sub: sub}, nil
}

// WatchFriendRequestUpdated is a free log subscription operation binding the contract event 0x8a2bf39e9087c33765bcca9a2c61caaa3e82956b669d0f3914a1d95e50f19ff2.
//
// Solidity: event FriendRequestUpdated(bytes32 indexed id, uint8 status)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchFriendRequestUpdated(opts *bind.WatchOpts, sink chan<- *QuickPicStorageFriendRequestUpdated, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "FriendRequestUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageFriendRequestUpdated)
				if err := _QuickPicStorage.contract.UnpackLog(event, "FriendRequestUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFriendRequestUpdated is a log parse operation binding the contract event 0x8a2bf39e9087c33765bcca9a2c61caaa3e82956b669d0f3914a1d95e50f19ff2.
//
// Solidity: event FriendRequestUpdated(bytes32 indexed id, uint8 status)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseFriendRequestUpdated(log types.Log) (*QuickPicStorageFriendRequestUpdated, error) {
	event := new(QuickPicStorageFriendRequestUpdated)
	if err := _QuickPicStorage.contract.UnpackLog(event, "FriendRequestUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuickPicStorageFriendshipCreatedIterator is returned from FilterFriendshipCreated and is used to iterate over the raw logs and unpacked data for FriendshipCreated events raised by the QuickPicStorage contract.
type QuickPicStorageFriendshipCreatedIterator struct {
	Event *QuickPicStorageFriendshipCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageFriendshipCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageFriendshipCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageFriendshipCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageFriendshipCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageFriendshipCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageFriendshipCreated represents a FriendshipCreated event raised by the QuickPicStorage contract.
type QuickPicStorageFriendshipCreated struct {
	Id      [32]byte
	UserAId [32]byte
	UserBId [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFriendshipCreated is a free log retrieval operation binding the contract event 0x63bb48a2197f36fe8120ecc4f504ccefdd2d712a3bf2d5645c0e6dac33e0301b.
//
// Solidity: event FriendshipCreated(bytes32 indexed id, bytes32 indexed userAId, bytes32 indexed userBId)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterFriendshipCreated(opts *bind.FilterOpts, id [][32]byte, userAId [][32]byte, userBId [][32]byte) (*QuickPicStorageFriendshipCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var userAIdRule []interface{}
	for _, userAIdItem := range userAId {
		userAIdRule = append(userAIdRule, userAIdItem)
	}
	var userBIdRule []interface{}
	for _, userBIdItem := range userBId {
		userBIdRule = append(userBIdRule, userBIdItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "FriendshipCreated", idRule, userAIdRule, userBIdRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageFriendshipCreatedIterator{contract: _QuickPicStorage.contract, event: "FriendshipCreated", logs: logs, sub: sub}, nil
}

// WatchFriendshipCreated is a free log subscription operation binding the contract event 0x63bb48a2197f36fe8120ecc4f504ccefdd2d712a3bf2d5645c0e6dac33e0301b.
//
// Solidity: event FriendshipCreated(bytes32 indexed id, bytes32 indexed userAId, bytes32 indexed userBId)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchFriendshipCreated(opts *bind.WatchOpts, sink chan<- *QuickPicStorageFriendshipCreated, id [][32]byte, userAId [][32]byte, userBId [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var userAIdRule []interface{}
	for _, userAIdItem := range userAId {
		userAIdRule = append(userAIdRule, userAIdItem)
	}
	var userBIdRule []interface{}
	for _, userBIdItem := range userBId {
		userBIdRule = append(userBIdRule, userBIdItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "FriendshipCreated", idRule, userAIdRule, userBIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageFriendshipCreated)
				if err := _QuickPicStorage.contract.UnpackLog(event, "FriendshipCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFriendshipCreated is a log parse operation binding the contract event 0x63bb48a2197f36fe8120ecc4f504ccefdd2d712a3bf2d5645c0e6dac33e0301b.
//
// Solidity: event FriendshipCreated(bytes32 indexed id, bytes32 indexed userAId, bytes32 indexed userBId)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseFriendshipCreated(log types.Log) (*QuickPicStorageFriendshipCreated, error) {
	event := new(QuickPicStorageFriendshipCreated)
	if err := _QuickPicStorage.contract.UnpackLog(event, "FriendshipCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuickPicStorageMessageCreatedIterator is returned from FilterMessageCreated and is used to iterate over the raw logs and unpacked data for MessageCreated events raised by the QuickPicStorage contract.
type QuickPicStorageMessageCreatedIterator struct {
	Event *QuickPicStorageMessageCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageMessageCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageMessageCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageMessageCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageMessageCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageMessageCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageMessageCreated represents a MessageCreated event raised by the QuickPicStorage contract.
type QuickPicStorageMessageCreated struct {
	Id         [32]byte
	FromUserId [32]byte
	ToUserId   [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMessageCreated is a free log retrieval operation binding the contract event 0x4f2ba766b8c822d461311c87d079178a55231792db9f142af3ecd797a5cc5d28.
//
// Solidity: event MessageCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterMessageCreated(opts *bind.FilterOpts, id [][32]byte, fromUserId [][32]byte, toUserId [][32]byte) (*QuickPicStorageMessageCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromUserIdRule []interface{}
	for _, fromUserIdItem := range fromUserId {
		fromUserIdRule = append(fromUserIdRule, fromUserIdItem)
	}
	var toUserIdRule []interface{}
	for _, toUserIdItem := range toUserId {
		toUserIdRule = append(toUserIdRule, toUserIdItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "MessageCreated", idRule, fromUserIdRule, toUserIdRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageMessageCreatedIterator{contract: _QuickPicStorage.contract, event: "MessageCreated", logs: logs, sub: sub}, nil
}

// WatchMessageCreated is a free log subscription operation binding the contract event 0x4f2ba766b8c822d461311c87d079178a55231792db9f142af3ecd797a5cc5d28.
//
// Solidity: event MessageCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchMessageCreated(opts *bind.WatchOpts, sink chan<- *QuickPicStorageMessageCreated, id [][32]byte, fromUserId [][32]byte, toUserId [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromUserIdRule []interface{}
	for _, fromUserIdItem := range fromUserId {
		fromUserIdRule = append(fromUserIdRule, fromUserIdItem)
	}
	var toUserIdRule []interface{}
	for _, toUserIdItem := range toUserId {
		toUserIdRule = append(toUserIdRule, toUserIdItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "MessageCreated", idRule, fromUserIdRule, toUserIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageMessageCreated)
				if err := _QuickPicStorage.contract.UnpackLog(event, "MessageCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMessageCreated is a log parse operation binding the contract event 0x4f2ba766b8c822d461311c87d079178a55231792db9f142af3ecd797a5cc5d28.
//
// Solidity: event MessageCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseMessageCreated(log types.Log) (*QuickPicStorageMessageCreated, error) {
	event := new(QuickPicStorageMessageCreated)
	if err := _QuickPicStorage.contract.UnpackLog(event, "MessageCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuickPicStorageMessageDeletedIterator is returned from FilterMessageDeleted and is used to iterate over the raw logs and unpacked data for MessageDeleted events raised by the QuickPicStorage contract.
type QuickPicStorageMessageDeletedIterator struct {
	Event *QuickPicStorageMessageDeleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageMessageDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageMessageDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageMessageDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageMessageDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageMessageDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageMessageDeleted represents a MessageDeleted event raised by the QuickPicStorage contract.
type QuickPicStorageMessageDeleted struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMessageDeleted is a free log retrieval operation binding the contract event 0xf4ef3cbd1d2cff45998b1fb60a13a6505834c570a5a9a44f55c3a3c2b5c852fd.
//
// Solidity: event MessageDeleted(bytes32 indexed id)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterMessageDeleted(opts *bind.FilterOpts, id [][32]byte) (*QuickPicStorageMessageDeletedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "MessageDeleted", idRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageMessageDeletedIterator{contract: _QuickPicStorage.contract, event: "MessageDeleted", logs: logs, sub: sub}, nil
}

// WatchMessageDeleted is a free log subscription operation binding the contract event 0xf4ef3cbd1d2cff45998b1fb60a13a6505834c570a5a9a44f55c3a3c2b5c852fd.
//
// Solidity: event MessageDeleted(bytes32 indexed id)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchMessageDeleted(opts *bind.WatchOpts, sink chan<- *QuickPicStorageMessageDeleted, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "MessageDeleted", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageMessageDeleted)
				if err := _QuickPicStorage.contract.UnpackLog(event, "MessageDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMessageDeleted is a log parse operation binding the contract event 0xf4ef3cbd1d2cff45998b1fb60a13a6505834c570a5a9a44f55c3a3c2b5c852fd.
//
// Solidity: event MessageDeleted(bytes32 indexed id)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseMessageDeleted(log types.Log) (*QuickPicStorageMessageDeleted, error) {
	event := new(QuickPicStorageMessageDeleted)
	if err := _QuickPicStorage.contract.UnpackLog(event, "MessageDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuickPicStorageUserCreatedIterator is returned from FilterUserCreated and is used to iterate over the raw logs and unpacked data for UserCreated events raised by the QuickPicStorage contract.
type QuickPicStorageUserCreatedIterator struct {
	Event *QuickPicStorageUserCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageUserCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageUserCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageUserCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageUserCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageUserCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageUserCreated represents a UserCreated event raised by the QuickPicStorage contract.
type QuickPicStorageUserCreated struct {
	Id         [32]byte
	UserNumber *big.Int
	Username   string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUserCreated is a free log retrieval operation binding the contract event 0xfa1f973aa4e6a4cd973d33775be141ea493df90f6cba4b68b4fcc7fe352522f0.
//
// Solidity: event UserCreated(bytes32 indexed id, uint256 userNumber, string username)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterUserCreated(opts *bind.FilterOpts, id [][32]byte) (*QuickPicStorageUserCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "UserCreated", idRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageUserCreatedIterator{contract: _QuickPicStorage.contract, event: "UserCreated", logs: logs, sub: sub}, nil
}

// WatchUserCreated is a free log subscription operation binding the contract event 0xfa1f973aa4e6a4cd973d33775be141ea493df90f6cba4b68b4fcc7fe352522f0.
//
// Solidity: event UserCreated(bytes32 indexed id, uint256 userNumber, string username)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchUserCreated(opts *bind.WatchOpts, sink chan<- *QuickPicStorageUserCreated, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "UserCreated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageUserCreated)
				if err := _QuickPicStorage.contract.UnpackLog(event, "UserCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUserCreated is a log parse operation binding the contract event 0xfa1f973aa4e6a4cd973d33775be141ea493df90f6cba4b68b4fcc7fe352522f0.
//
// Solidity: event UserCreated(bytes32 indexed id, uint256 userNumber, string username)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseUserCreated(log types.Log) (*QuickPicStorageUserCreated, error) {
	event := new(QuickPicStorageUserCreated)
	if err := _QuickPicStorage.contract.UnpackLog(event, "UserCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuickPicStorageUserUpdatedIterator is returned from FilterUserUpdated and is used to iterate over the raw logs and unpacked data for UserUpdated events raised by the QuickPicStorage contract.
type QuickPicStorageUserUpdatedIterator struct {
	Event *QuickPicStorageUserUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *QuickPicStorageUserUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuickPicStorageUserUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(QuickPicStorageUserUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *QuickPicStorageUserUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuickPicStorageUserUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuickPicStorageUserUpdated represents a UserUpdated event raised by the QuickPicStorage contract.
type QuickPicStorageUserUpdated struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUserUpdated is a free log retrieval operation binding the contract event 0xdf3a752736081f0b857124cd76b543d42ea93c426e4ae00b28927d4de82ee102.
//
// Solidity: event UserUpdated(bytes32 indexed id)
func (_QuickPicStorage *QuickPicStorageFilterer) FilterUserUpdated(opts *bind.FilterOpts, id [][32]byte) (*QuickPicStorageUserUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.FilterLogs(opts, "UserUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &QuickPicStorageUserUpdatedIterator{contract: _QuickPicStorage.contract, event: "UserUpdated", logs: logs, sub: sub}, nil
}

// WatchUserUpdated is a free log subscription operation binding the contract event 0xdf3a752736081f0b857124cd76b543d42ea93c426e4ae00b28927d4de82ee102.
//
// Solidity: event UserUpdated(bytes32 indexed id)
func (_QuickPicStorage *QuickPicStorageFilterer) WatchUserUpdated(opts *bind.WatchOpts, sink chan<- *QuickPicStorageUserUpdated, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _QuickPicStorage.contract.WatchLogs(opts, "UserUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuickPicStorageUserUpdated)
				if err := _QuickPicStorage.contract.UnpackLog(event, "UserUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUserUpdated is a log parse operation binding the contract event 0xdf3a752736081f0b857124cd76b543d42ea93c426e4ae00b28927d4de82ee102.
//
// Solidity: event UserUpdated(bytes32 indexed id)
func (_QuickPicStorage *QuickPicStorageFilterer) ParseUserUpdated(log types.Log) (*QuickPicStorageUserUpdated, error) {
	event := new(QuickPicStorageUserUpdated)
	if err := _QuickPicStorage.contract.UnpackLog(event, "UserUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
