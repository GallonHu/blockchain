package core

// TxInput represents a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
