package backend

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/expandorg/requester-service/pkg/nulls"
)

var addressRegex = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
var txHashRegex = regexp.MustCompile("^0x[a-fA-F0-9]{64}$")

type Address nulls.String

func NewAddress(a string) Address {
	return Address(nulls.NewString(a))
}

func (a Address) String() string {
	return a.NullString.String
}

func (a Address) Normalize() string {
	return strings.ToLower(a.String())
}

func (a Address) MarshalJSON() ([]byte, error) {
	return nulls.String(a).MarshalJSON()
}

func (a *Address) UnmarshalJSON(data []byte) error {
	s := new(nulls.String)
	err := s.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	*a = Address(*s)
	return nil
}

func (a Address) IsValid() bool {
	return addressRegex.MatchString(a.NullString.String)
}

func (a Address) WithoutPrefix() string {
	return a.NullString.String[2:]
}

type TXHash nulls.String

func NewTXHash(tx string) TXHash {
	return TXHash(nulls.NewString(tx))
}

func (tx TXHash) String() string {
	return tx.NullString.String
}

func (tx TXHash) Normalize() string {
	return strings.ToLower(tx.String())
}

func (tx TXHash) MarshalJSON() ([]byte, error) {
	return nulls.String(tx).MarshalJSON()
}

func (tx *TXHash) UnmarshalJSON(data []byte) error {
	s := new(nulls.String)
	err := s.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	*tx = TXHash(*s)
	return nil
}

func (tx *TXHash) IsValid() bool {
	return txHashRegex.MatchString(tx.NullString.String)
}

type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
)

func (t TransactionType) IsValid() bool {
	return t == TransactionDeposit || t == TransactionWithdrawal
}

type UnableToAuthErr struct{}

func (err UnableToAuthErr) Error() string {
	return "Unable to authenticate user"
}

type UserResponse struct {
	User              User         `json:"user"`
	AssignedTasks     []Task       `json:"assignedTasks"`
	AssignedResponses []Response   `json:"assignedResponses"`
	Assignments       []Assignment `json:"assignments"`
}

type User struct {
	ID               uint64        `json:"id"`
	Email            nulls.String  `json:"email"`
	EmailConfirmed   bool          `json:"emailConfirmed"`
	Address          Address       `json:"address"`
	AddressConfirmed bool          `json:"addressConfirmed"`
	Score            nulls.Float64 `json:"score"`
	Stats            UserStats     `json:"stats"`
	Gems             UserGems      `json:"gems"`
	PendingTX        *Transaction  `json:"pendingTx"`
}

type Assignment struct {
	ID         uint64      `json:"id"`
	AssignedAt time.Time   `json:"assignedAt"`
	ExpiresAt  nulls.Time  `json:"expiresAt"`
	JobID      uint64      `json:"jobId"`
	TaskID     uint64      `json:"taskId"`
	ResponseID nulls.Int64 `json:"responseId"`
}

type UserStats struct {
	Pending  uint64 `json:"pending"`
	Accepted uint64 `json:"accepted"`
	Rejected uint64 `json:"rejected"`
}

type UserGems struct {
	Balance  float64 `json:"balance"`
	Reserved float64 `json:"reserved"`
}

type Transaction struct {
	Hash TXHash          `json:"hash"`
	Type TransactionType `json:"type"`
}

type TaskStats struct {
	Assignment   uint64 `json:"assignment"`
	Pending      uint64 `json:"pending"`
	Accepted     uint64 `json:"accepted"`
	Verification uint64 `json:"verification"`
}

type Task struct {
	ID       uint64          `json:"id"`
	JobID    uint64          `json:"jobId"`
	TaskData json.RawMessage `json:"taskData"`
	IsActive bool            `json:"isActive"`
	Stats    TaskStats       `json:"stats"`
}

type Response struct {
	ID         uint64          `json:"id"`
	WorkerID   uint64          `json:"workerId"`
	JobID      uint64          `json:"jobId"`
	TaskID     uint64          `json:"taskId"`
	Value      json.RawMessage `json:"value"`
	IsAccepted nulls.Bool      `json:"isAccepted"`
}
