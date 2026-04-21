// Пакет model содержит доменные типы системы лояльности Gophermart.
package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus — статус обработки заказа в системе начислений.
type OrderStatus string

const (
	// OrderStatusNew — заказ загружен, ещё не обработан.
	OrderStatusNew OrderStatus = "NEW"
	// OrderStatusProcessing — заказ принят в обработку сервисом начислений.
	OrderStatusProcessing OrderStatus = "PROCESSING"
	// OrderStatusInvalid — заказ не принят системой начислений (не будет начислено баллов).
	OrderStatusInvalid OrderStatus = "INVALID"
	// OrderStatusProcessed — заказ обработан, баллы начислены.
	OrderStatusProcessed OrderStatus = "PROCESSED"
)

// User представляет зарегистрированного пользователя системы.
type User struct {
	// ID — уникальный идентификатор пользователя.
	ID uuid.UUID
	// Login — уникальный логин пользователя.
	Login string
	// PasswordHash — хэш пароля пользователя.
	PasswordHash string
}

// Order представляет заказ, загруженный пользователем для начисления баллов.
type Order struct {
	// Number — номер заказа (проходит проверку алгоритмом Луна).
	Number string
	// UserId — идентификатор пользователя, загрузившего заказ.
	UserId uuid.UUID
	// Status — текущий статус обработки заказа.
	Status OrderStatus
	// Accrual — количество начисленных баллов; nil, если начисление ещё не произошло.
	Accrual *float64
	// UploadedAt — время загрузки заказа.
	UploadedAt time.Time
}

// Withdrawal представляет операцию списания баллов пользователем.
type Withdrawal struct {
	// ID — уникальный идентификатор операции списания.
	ID uuid.UUID
	// UserId — идентификатор пользователя, выполнившего списание.
	UserId uuid.UUID
	// Order — номер заказа, в счёт которого произведено списание.
	Order string
	// Sum — количество списанных баллов.
	Sum float64
	// ProcessedAt — время выполнения операции списания.
	ProcessedAt time.Time
}
