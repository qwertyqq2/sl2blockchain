package blockchain

import "errors"

var (
	ErrMinPercent               = errors.New("not enought value")
	ErrNilTx                    = errors.New("nil tx")
	ErrStorageRewardPass        = errors.New("storage reward pass")
	ErrNotEnoghtMoney           = errors.New("not enought mney")
	ErrNilBlock                 = errors.New("nil block")
	ErrEqualRandBytes           = errors.New("equal rand bytes")
	ErrSecondStorageSender      = errors.New("may be only one storage sender")
	ErrIncorrectStorageReceiver = errors.New("incorrect storage receiver")
	ErrIncorrectDiff            = errors.New("incorrect difficalty")
	ErrIncorrectHash            = errors.New("incorrect hash")
	ErrIncorrectTrasnactions    = errors.New("incorrect transactions")
	ErrIncorrectSign            = errors.New("incorrect sign")
	ErrIncorrectProof           = errors.New("incorrect proof")
	ErrIncorrectTime            = errors.New("incorrect time")

	ErrTxHash            = errors.New("not valid hash tx")
	ErrTxSign            = errors.New("not valid sign tx")
	ErrTxBalanceSender   = errors.New("not valid sender balance tx")
	ErrTxBalanceReceiver = errors.New("not valid receiver balance tx")

	ErrNothaveStorage        = errors.New("not have sttorage in block")
	ErrMissingAddressInBlock = errors.New("missing address in block")
	ErrIncorrectBalanceBlock = errors.New("incorrect balance block")
	ErrIncorrectTimeBlock    = errors.New("incorrect time block")

	ErrNotProof = errors.New("not proof")

	ErrChainAlreadyExist = errors.New("chain already exist")

	ErrNotFileExist = errors.New("not file exist")
	ErrNotRows      = errors.New("not rows")
	ErrIsLastBlock  = errors.New("is last block")
)
