package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Не смог подключиться к базе, ошибка")
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err, "Не смог добавить посылку, ошибка")
	assert.NotZero(t, id, "Посылку положил, но id у нее 0")
	parcel.Number = id

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	secondParcel, err := store.Get(id)
	assert.NoError(t, err, "Не смог получить посылку, ошибка")
	assert.Equal(t, parcel, secondParcel, "Отправленная и полученная посылки не совпадают")

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	assert.NoError(t, err, "Не смог удалить посылку, ошибка")
	_, err = store.Get(id)
	assert.ErrorIs(t, err, sql.ErrNoRows, "При запросе удаленной посылки не получена ошибка sql.ErrNoRows")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Не смог подключиться к базе, ошибка")
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err, "Не смог добавить посылку, ошибка")
	assert.NotZero(t, id, "Посылку положил, но id у нее 0")

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.NoError(t, err, "Не смог изменить адрес, ошибка")

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	secondParcel, err := store.Get(id)
	assert.NoError(t, err, "Не смог получить посылку, ошибка")
	assert.Equal(t, newAddress, secondParcel.Address, "Адрес у полученной из базы посылки не верный")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepareы
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Не смог подключиться к базе, ошибка")
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err, "Не смог добавить посылку, ошибка")
	assert.NotZero(t, id, "Посылку положил, но id у нее 0")

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := ParcelStatusDelivered
	err = store.SetStatus(id, newStatus)
	assert.NoError(t, err, "Не смог изменить статус посылки, ошибка")

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	secondParcel, err := store.Get(id)
	assert.NoError(t, err, "Не смог получить посылку, ошибка")
	assert.Equal(t, newStatus, secondParcel.Status, "Статус у полученной из базы посылки не верный")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Не смог подключиться к базе, ошибка")
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		assert.NoError(t, err, "Не смог добавить посылку, ошибка")
		assert.NotZero(t, id, "Посылку положил, но id у нее 0")

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	assert.NoError(t, err, "Не смог получить посылку для клиента, ошибка")
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	assert.Len(t, storedParcels, 3, "Загрузили одно количество посылок, а получили другое, ошибка")
	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		expectedParcel, exists := parcelMap[parcel.Number]
		assert.True(t, exists, "Посылка с ID %d не найдена в parcelMap", parcel.Number)
		assert.Equal(t, expectedParcel, parcel, "Посылка с ID %d не совпадает с ожидаемой", parcel.Number)
	}
}
