package blockchain

import "errors"

var (
	ErrMinPercent               = errors.New("not enought value")
	ErrNilTx                    = errors.New("nil tx")
	ErrStorageRewardPass        = errors.New("storage reward pass")
	ErrNotEnoghtMoney           = errors.New("not enought mney")
	ErrNilBlock                 = errors.New("Nil block")
	ErrEqualRandBytes           = errors.New("Equal rand bytes")
	ErrSecondStorageSender      = errors.New("May be only one storage sender")
	ErrIncorrectStorageReceiver = errors.New("Incorrect storage receiver")

	ErrTxHash            = errors.New("not valid hash tx")
	ErrTxSign            = errors.New("not valid sign tx")
	ErrTxBalanceSender   = errors.New("not valid sender balance tx")
	ErrTxBalanceReceiver = errors.New("not valid receiver balance tx")

	ErrNothaveStorage        = errors.New("not have sttorage in block")
	ErrMissingAddressInBlock = errors.New("missing address in block")
	ErrIncorrectBalanceBlock = errors.New("incorrect balance block")
	ErrIncorrectTimeBlock    = errors.New("incorrect time block")

	ErrNotProof = errors.New("Not proof")
)
