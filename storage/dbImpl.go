///////////////////////////////////////////////////////////////////////////////
// Copyright © 2021 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package storage

import (
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/gorm/clause"
)

func (db *DatabaseImpl) InsertMembers(members []Member) error {
	return db.db.Create(&members).Error
}

func (db *DatabaseImpl) InsertCommitment(commitment Commitment) error {
	return db.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&commitment).Error
}

func (db *DatabaseImpl) GetMember(id string) (*Member, error) {
	jww.INFO.Printf("Getting member with id %+v", id)
	m := Member{}
	return &m, db.db.First(&m, "id = ?", id).Error
}
