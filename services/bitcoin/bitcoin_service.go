package bitcoin

import (
	"log"
	"own-paynet/config"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

type BitcoinService struct {
	client *rpcclient.Client
}

func NewBitcoinService(cfg *config.Config) (*BitcoinService, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         cfg.BitcoinRPCURL,
		User:         cfg.BitcoinRPCUser,
		Pass:         cfg.BitcoinRPCPass,
		HTTPPostMode: true,
		DisableTLS:   true, // Use TLS in production
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return &BitcoinService{client: client}, nil
}

func (s *BitcoinService) GenerateAddress() (string, error) {
	address, err := s.client.GetNewAddress("")
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// MonitorAddress polls transactions and calls the callback when transactions are found.
func (s *BitcoinService) MonitorAddress(address string, callback func(txID string, status string), netParams *chaincfg.Params) error {
	go func() {
		for {
			addr, err := btcutil.DecodeAddress(address, netParams)
			if err != nil {
				log.Printf("failed to decode address: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			txs, err := s.client.SearchRawTransactions(addr, 0, 10, true, nil)
			if err != nil {
				log.Printf("failed to search transactions: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for _, tx := range txs {
				txID := tx.TxID()
				status := "confirmed"

				txHash, err := chainhash.NewHashFromStr(txID)
				if err != nil {
					log.Printf("invalid tx hash: %v", err)
					continue
				}

				txDetails, err := s.client.GetTransaction(txHash)
				if err != nil {
					log.Printf("failed to get transaction details: %v", err)
					continue
				}

				if txDetails.Confirmations == 0 {
					status = "pending"
				}

				// Execute user-defined callback
				callback(txID, status)
			}

			// Sleep to avoid excessive polling
			time.Sleep(60 * time.Second)
		}
	}()

	return nil
}
