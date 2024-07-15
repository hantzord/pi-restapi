package transaction

import (
	"capstone/constants"
	"capstone/controllers/transaction/response"
	"capstone/entities"
	transactionEntities "capstone/entities/transaction"
	"capstone/repositories/mysql/consultation"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) transactionEntities.TransactionRepository {
	return &TransactionRepo{db}
}

func (repository *TransactionRepo) Insert(transaction *transactionEntities.Transaction) (*transactionEntities.Transaction, error) {
	transactionDb := ToTransactionModel(transaction)
	if err := repository.db.Create(&transactionDb).Error; err != nil {
		return nil, constants.ErrInsertDatabase
	}

	if err := repository.db.Preload("Consultation").Preload("Consultation.Doctor").First(&transactionDb, transactionDb.ID).Error; err != nil {
		return nil, constants.ErrInsertDatabase
	}

	return transactionDb.ToEntities(), nil
}

func (repository *TransactionRepo) FindByID(ID string) (*transactionEntities.Transaction, error) {
	transactionDB := new(Transaction)
	if err := repository.db.Preload("Consultation").Preload("Consultation.Doctor").First(&transactionDB, "id LIKE ?", ID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	return transactionDB.ToEntities(), nil
}

func (repository *TransactionRepo) FindByConsultationID(consultationID uint) (*transactionEntities.Transaction, error) {
	transactionDB := new(Transaction)
	if err := repository.db.Preload("Consultation").Preload("Consultation.Doctor").First(&transactionDB, "consultation_id LIKE ?", consultationID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	return transactionDB.ToEntities(), nil
}

func (repository *TransactionRepo) FindAllByUserID(metadata *entities.Metadata, userID uint, status string) (*[]transactionEntities.Transaction, error) {
	transactionDB := new([]Transaction)
	if err := repository.db.
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Joins("JOIN users ON consultations.user_id = users.id").
		Where("transactions.status LIKE ?", "%"+status+"%").
		Where("users.id LIKE ?", userID).
		Preload("Consultation").
		Preload("Consultation.Doctor").
		Limit(metadata.Limit).
		Offset(metadata.Offset()).
		Find(&transactionDB).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var transactions []transactionEntities.Transaction
	for _, transaction := range *transactionDB {
		transactions = append(transactions, *transaction.ToEntities())
	}
	return &transactions, nil
}

func (repository *TransactionRepo) Update(transaction *transactionEntities.Transaction) (*transactionEntities.Transaction, error) {
	transactionDB := ToTransactionModel(transaction)
	if err := repository.db.Model(&Transaction{}).Where("id LIKE ?", transactionDB.ID).Update("status", transactionDB.Status).Update("payment_status", transactionDB.Status).Error; err != nil {
		return nil, constants.ErrInsertDatabase
	}
	transactionDB.Consultation = *consultation.ToConsultationModel(&transaction.Consultation)
	return transactionDB.ToEntities(), nil
}

func (repository *TransactionRepo) Delete(ID string) error {
	_, err := repository.FindByID(ID)
	if err != nil {
		return constants.ErrDataNotFound
	}
	if err = repository.db.Delete(&Transaction{}, "id LIKE ?", ID).Error; err != nil {
		return constants.ErrDeleteDatabase
	}
	return nil
}

func (repository *TransactionRepo) FindAllByDoctorID(metadata *entities.Metadata, doctorID uint, status string) (*[]transactionEntities.Transaction, error) {
	transactionDB := new([]Transaction)
	if err := repository.db.
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Joins("JOIN doctors ON consultations.doctor_id = doctors.id").
		Where("doctors.id LIKE ?", doctorID).
		Where("transactions.status LIKE ?", "%"+status+"%").
		Preload("Consultation").
		Preload("Consultation.User").
		Preload("Consultation.Complaint").
		Limit(metadata.Limit).
		Offset(metadata.Offset()).
		Find(&transactionDB).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var transactions []transactionEntities.Transaction
	for _, transaction := range *transactionDB {
		transactions = append(transactions, *transaction.ToEntities())
	}
	return &transactions, nil
}

func (repository *TransactionRepo) CountTransactionByDoctorID(doctorID uint) (*response.TransactionCount, error) {
	var transactionTotal, transactionSuccess, transactionFailed, transactionToday, transactionWeek, transactionMonth, transactionYear int64
	var err error

	// Count All Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status NOT LIKE ?", doctorID, constants.PENDING).
		Count(&transactionTotal).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	// Count Success Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status LIKE ?", doctorID, constants.Success).
		Count(&transactionSuccess).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	// Count Failed Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status LIKE ?", doctorID, constants.Failed).
		Count(&transactionFailed).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	// Count Today Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status NOT LIKE ? AND DATE(transactions.created_at) = CURDATE()", doctorID, constants.PENDING).
		Count(&transactionToday).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	// Count Week Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status NOT LIKE ? AND YEARWEEK(transactions.created_at, 1) = YEARWEEK(CURDATE(), 1)", doctorID, constants.PENDING).
		Count(&transactionWeek).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	// Count Month Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status NOT LIKE ? AND MONTH(transactions.created_at) = MONTH(CURDATE())", doctorID, constants.PENDING).
		Count(&transactionMonth).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	// Count Year Transaction
	if err = repository.db.
		Model(&Transaction{}).
		Joins("JOIN consultations ON consultations.id = transactions.consultation_id").
		Where("consultations.doctor_id LIKE ? AND transactions.status NOT LIKE ? AND YEAR(transactions.created_at) = YEAR(CURDATE())", doctorID, constants.PENDING).
		Count(&transactionYear).
		Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	countTransaction := response.ToTransactionCount(transactionTotal, transactionSuccess, transactionFailed, transactionToday, transactionWeek, transactionMonth, transactionYear)
	return countTransaction, nil
}
