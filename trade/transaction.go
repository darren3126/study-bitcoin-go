package trade

import (
	"encoding/gob"
	"log"
	"crypto/sha256"
	"bytes"
	"fmt"
	"github.com/study-bitcion-go/block"
	"encoding/hex"
)

const subsidy  = 10 //初始化补助为10


//交易事物
type Transaction struct {
	ID   []byte     //交易hash
	Vin  []TXInput  //事物输入
	Vout []TXOutput //事物输出
}
//一个事物输入
type TXInput struct {
	Txid []byte //交易ID的hash
	Vout int    //交易输出
	ScriptSig string //解锁脚本
}
//一个事物输出
type TXOutput struct {
	Value int
	ScriptPubkey string
}


func (tx *Transaction)IsCoinbase()bool  {
    return len(tx.Vin)==1&&len(tx.Vin[0].Txid)==0&&tx.Vin[0].Vout==-1
}
//设置交易ID hash
func (tx *Transaction) SetID(){
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	//将交易信息sha256
	hash = sha256.Sum256(encoded.Bytes())
	//生成hash
	tx.ID = hash[:]
}

func NewCoinbaseTX(to,data string) *Transaction  {
	if data==""{
		data=fmt.Sprintf("Reward to '%s'",to)
	}

	txin :=TXInput{[]byte{},-1,data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx

}

func NewUTXOTransaction(from, to string, amount int, bc *block.Blockchain)  *Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
	
}