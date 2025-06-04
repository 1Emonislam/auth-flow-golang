package services

import (
	"errors"
	"own-paynet/models"
	"own-paynet/repository"
	"own-paynet/services/bitcoin"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	walletService   *PayoutWalletService
	bitcoinService  *bitcoin.BitcoinService
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, walletService *PayoutWalletService, bitcoinService *bitcoin.BitcoinService) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		walletService:   walletService,
		bitcoinService:  bitcoinService,
	}
}

// CreateTransaction creates a new transaction for sender only
func (s *TransactionService) CreateTransaction(walletID uint, transactionType models.TransactionType, amount float64, priceCurrency, payCurrency, comment string, senderID uint, receiverWallet string) (*models.Transaction, error) {
	// Validate inputs
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// Get the sender's wallet
	senderWallet, err := s.walletService.GetPayoutWallet(walletID)
	if err != nil {
		return nil, err
	}

	// Validate sender owns the wallet
	if senderWallet.UserID != senderID {
		return nil, errors.New("wallet does not belong to sender")
	}

	// Find receiver's wallet by address
	receiverWalletFound, err := s.walletService.GetPayoutWalletByAddress(receiverWallet)
	if err != nil {
		return nil, err
	}

	if receiverWalletFound == nil {
		return nil, errors.New("receiver wallet not found")
	}

	// Check sender's balance
	if senderWallet.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	// Create transaction record with pending status
	transaction := &models.Transaction{
		PayoutWalletID: walletID,
		Type:           models.TransactionTypeDebit,
		Amount:         amount,
		PriceCurrency:  priceCurrency,
		PayCurrency:    payCurrency,
		Comment:        comment,
		Status:         "pending",
		SenderID:       senderID,
		ReceiverWallet: receiverWallet,
		IsProcessed:    false, // Track if transaction has been processed
	}

	// Create transaction record
	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Start monitoring blockchain confirmation
	go func() {
		// Check if transaction is already processed
		if transaction.IsProcessed {
			return
		}

		// Monitor blockchain confirmation
		confirmations, err := s.bitcoinService.GetTransactionConfirmations(transaction.TransactionID)
		if err != nil {
			// Update transaction status to failed
			_ = s.transactionRepo.UpdateStatus(transaction.ID, "failed")
			_ = s.transactionRepo.MarkAsProcessed(transaction.ID)
			return
		}

		// If we have enough confirmations (e.g., 6 or more)
		if confirmations >= 6 {
			// Mark as processed before fund transfer to prevent duplicate processing
			err = s.transactionRepo.MarkAsProcessed(transaction.ID)
			if err != nil {
				return
			}

			// Deduct amount from sender's wallet
			err = s.walletService.WithdrawFunds(walletID, amount)
			if err != nil {
				_ = s.transactionRepo.UpdateStatus(transaction.ID, "failed")
				return
			}

			// Add amount to receiver's wallet
			err = s.walletService.AddFunds(receiverWalletFound.ID, amount)
			if err != nil {
				// If adding to receiver fails, refund sender
				_ = s.walletService.AddFunds(walletID, amount)
				_ = s.transactionRepo.UpdateStatus(transaction.ID, "failed")
				return
			}

			// Update transaction status to completed
			_ = s.transactionRepo.UpdateStatus(transaction.ID, "completed")
		} else {
			// Not enough confirmations, mark as failed and processed
			_ = s.transactionRepo.UpdateStatus(transaction.ID, "failed")
			_ = s.transactionRepo.MarkAsProcessed(transaction.ID)
		}
	}()

	return transaction, nil
}

// GetTransaction retrieves a transaction by ID
func (s *TransactionService) GetTransaction(id uint) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(id)
}

// GetWalletTransactions retrieves all transactions for a wallet
func (s *TransactionService) GetWalletTransactions(walletID uint) ([]models.Transaction, error) {
	return s.transactionRepo.FindByPayoutWalletID(walletID)
}

// GetUserTransactions retrieves all transactions where the user is sender
func (s *TransactionService) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	return s.transactionRepo.FindByUserID(userID)
}
