package backend

import (
	"fmt"

	"github.com/quickpic/server/internal/repository/blockchain"
	"github.com/quickpic/server/internal/repository/sqlite"
	"github.com/quickpic/server/internal/storage"
)

// Type represents the type of storage backend
type Type string

const (
	TypeSQLite     Type = "sqlite"
	TypeBlockchain Type = "blockchain"
)

// Config holds the configuration for creating a backend
type Config struct {
	Type Type

	// SQLite configuration
	SQLitePath string

	// Blockchain configuration
	BlockchainRPCURL          string
	BlockchainPrivateKey      string
	BlockchainContractAddress string
}

// Result contains the backend and its repositories
type Result struct {
	Backend storage.Backend
	Repos   storage.Repositories
}

// New creates a new storage backend based on the configuration
func New(cfg Config) (*Result, error) {
	switch cfg.Type {
	case TypeSQLite:
		if cfg.SQLitePath == "" {
			cfg.SQLitePath = "./quickpic.db"
		}
		backend, err := sqlite.NewBackend(cfg.SQLitePath)
		if err != nil {
			return nil, err
		}
		return &Result{
			Backend: backend,
			Repos: storage.Repositories{
				Users:    backend.Users(),
				Friends:  backend.Friends(),
				Messages: backend.Messages(),
			},
		}, nil

	case TypeBlockchain:
		if cfg.BlockchainRPCURL == "" {
			return nil, fmt.Errorf("blockchain RPC URL is required")
		}
		if cfg.BlockchainPrivateKey == "" {
			return nil, fmt.Errorf("blockchain private key is required")
		}
		if cfg.BlockchainContractAddress == "" {
			return nil, fmt.Errorf("blockchain contract address is required")
		}
		backend, err := blockchain.NewBackend(blockchain.Config{
			RPCURL:          cfg.BlockchainRPCURL,
			PrivateKey:      cfg.BlockchainPrivateKey,
			ContractAddress: cfg.BlockchainContractAddress,
		})
		if err != nil {
			return nil, err
		}
		return &Result{
			Backend: backend,
			Repos: storage.Repositories{
				Users:    backend.Users(),
				Friends:  backend.Friends(),
				Messages: backend.Messages(),
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown backend type: %s", cfg.Type)
	}
}

// DefaultAnvilConfig returns a Config configured for local Anvil development
// Uses the first Anvil default account
func DefaultAnvilConfig(contractAddress string) Config {
	return Config{
		Type:                      TypeBlockchain,
		BlockchainRPCURL:          "http://localhost:8545",
		BlockchainPrivateKey:      "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		BlockchainContractAddress: contractAddress,
	}
}

// DefaultSQLiteConfig returns a Config configured for SQLite storage
func DefaultSQLiteConfig(dbPath string) Config {
	if dbPath == "" {
		dbPath = "./quickpic.db"
	}
	return Config{
		Type:       TypeSQLite,
		SQLitePath: dbPath,
	}
}
