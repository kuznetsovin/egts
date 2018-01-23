package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type EgtsPkgHeader struct {
	//Параметр определяет версию используемой структуры заголовка и должен содержать значение 0x01.
	//Значение данного параметра инкрементируется каждый раз при внесении изменений в структуру заголовка.
	PRV byte

	//Параметр определяет идентификатор ключа, используемый при шифровании.
	SKID byte

	// Данный параметр определяет префикс заголовка Транспортного Уровня и для данной версии
	// должен содержать значение 00.
	PRF uint8

	// Битовое поле определяет необходимость дальнейшей маршрутизации данного пакета на удалённую телематическую
	// платформу, а также наличие опциональных параметров PRA, RCA, TTL, необходимых для маршрутизации данного пакета.
	// Если поле имеет значение 1, то необходима маршрутизация, и поля PRA, RCA, TTL присутствуют в пакете.
	// Данное поле устанавливает Диспетчер той ТП, на которой сгенерирован пакет, или АТ,
	// сгенерировавший пакет для отправки на ТП, в случае установки в нём параметра «HOME_DISPATCHER_ID»,
	// определяющего адрес ТП, на которой данный АТ зарегистрирован.
	RTE uint8

	// // Битовое поле определяет код алгоритма, используемый для шифрования данных из поля SFRD.
	// // Если поле имеет значение 0 0 , то данные в поле SFRD не шифруются.
	// // Состав и коды алгоритмов не определены в данной версии Протокола
	ENA uint8

	// // Битовое поле определяет, используется ли сжатие данных из поля SFRD. Если поле имеет значение 1,
	// // то данные в поле SFRD считаются сжатыми. Алгоритм сжатия не определен в данной версии Протокола.
	CMP uint8

	// // Битовое поле определяет приоритет маршрутизации данного пакета и может принимать следующие значения:
	// // 0 0 – наивысший
	// // 0 1 – высокий
	// // 1 0 – средний
	// // 1 1 – низкий
	// // Установка большего приоритета позволяет передавать пакеты, содержащие срочные данные, такие, например,
	// // как пакет с минимальным набором данных услуги «ЭРА ГЛОНАСС» или данные о срабатывании сигнализации на ТС.
	// // При получении пакета Диспетчер, анализируя данное поле, производит маршрутизацию пакета с более высоким
	// // приоритетом быстрее, чем пакетов с низким приоритетом, тем самым достигается более оперативная обработка
	// // информации при наступлении критически важных событий.
	PR uint8

	// Длина заголовка Транспортного Уровня в байтах с учётом байта контрольной суммы (поля HCS).
	HL byte

	// Определяет применяемый метод кодирования следующей за данным параметром части заголовка Транспортного Уровня.
	HE byte

	// Определяет размер в байтах поля данных SFRD, содержащего информацию Протокола Уровня Поддержки Услуг.
	FDL uint16

	// Содержит номер пакета Транспортного Уровня, увеличивающийся на 1 при отправке каждого нового
	// пакета на стороне отправителя. Значения в данном поле изменяются по правилам циклического счётчика в
	// диапазоне от 0 до 65535, т.е. при достижении значения 65535, следующее значение должно быть 0.
	PID uint16

	// Тип пакета Транспортного Уровня. Поле PT может принимать следующие значения:
	// 0 – EGTS_PT_RESPONSE (подтверждение на пакет Транспортного Уровня);
	// 1 – EGTS_PT_APPDATA (пакет, содержащий данные Протокола Уровня Поддержки Услуг);
	// 2 – EGTS_PT_SIGNED_APPDATA (пакет, содержащий данные Протокола Уровня Поддержки Услуг с цифровой подписью);
	PT byte

	// Адрес ТП, на которой данный пакет сгенерирован. Данный адрес является уникальным в рамках связной сети и
	// используется для создания пакета-подтверждения на принимающей стороне.
	PRA uint16

	// Адрес ТП, для которой данный пакет предназначен. По данному адресу производится идентификация
	// принадлежности пакета определённой ТП и его маршрутизация при использовании промежуточных ТП.
	RCA uint16

	// Время жизни пакета при его маршрутизации между ТП.
	TTL byte

	// Контрольная сумма заголовка Транспортного Уровня (начиная с поля «PRV» до поля «HCS», не включая последнего).
	// Для подсчёта значения поля HCS ко всем байтам указанной последовательности применяется алгоритм CRC-8.
	// Пример программного кода расчета CRC-8 приведен в Приложении 3.
	HCS byte
}

// метод преобразования структуры в строку байт
func (h *EgtsPkgHeader) ToBytes() ([]byte, error) {
	result := []byte{}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, h.PRV); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.SKID); err != nil {
		return result, err
	}

	//составной байт
	flagsByte := fmt.Sprintf("%02b%01b%02b%01b%02b", h.PRF, h.RTE, h.ENA, h.CMP, h.PR)
	flagByte, err := bitsToByte(flagsByte)
	if err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, flagByte); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.HL); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.HE); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.FDL); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.PID); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.PT); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.PRA); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.RCA); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.TTL); err != nil {
		return result, err
	}

	if err := h.CalcCRC8(); err != nil {
		return result, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.HCS); err != nil {
		return result, err
	}

	result = buf.Bytes()
	return result, nil
}

func (h *EgtsPkgHeader) CalcCRC8() error {
	// ЭТО ЗАГЛУШКА ЗАМЕНИТЬ НА НОРМАЛЬНЫЙ АЛГОРИТМ!!!!!
	h.HCS = 202

	return nil
}

type EgtsPkg struct {
	EgtsPkgHeader

	// Структура данных, зависящая от типа Пакета и содержащая информацию Протокола Уровня Поддержки Услуг.
	// Формат структуры данных в зависимости от типа Пакета описан в п.8.2.
	SFRD []byte

	// Контрольная сумма поля уровня Протокола Поддержки Услуг. Для подсчёта контрольной суммы по данным из поля SFRD,
	// используется алгоритм CRC-16. Данное поле присутствует только в том случае, если есть поле SFRD.
	// Пример программного кода расчета CRC-16 приведен в Приложении 2.
	// Блок схема алгоритма разбора пакета Протокола Транспортного Уровня при приеме представлена на рисунке 3.
	SFRCS uint16
}