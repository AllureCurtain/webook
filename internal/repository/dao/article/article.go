package article

//import (
//	"context"
//	"errors"
//	"gorm.io/gorm"
//	"gorm.io/gorm/clause"
//	"time"
//)
//
//var ErrPossibleIncorrectAuthor = errors.New("用户在尝试操作非本人数据")
//
//type ArticleDAO interface {
//	Insert(ctx context.Context, art Article) (int64, error)
//	UpdateById(ctx context.Context, article Article) error
//	Sync(ctx context.Context, art Article) (int64, error)
//	SyncStatus(ctx context.Context, author, id int64, status uint8) error
//}
//
//func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
//	return &GORMArticleDAO{
//		db: db,
//	}
//}
//
//type GORMArticleDAO struct {
//	db *gorm.DB
//}
//
//func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
//	now := time.Now().UnixMilli()
//	art.Ctime = now
//	art.Utime = now
//	err := dao.db.WithContext(ctx).Create(&art).Error
//	return art.Id, err
//}
//
//func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
//	now := time.Now().UnixMilli()
//	art.Utime = now
//	err := dao.db.WithContext(ctx).Model(&art).Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
//		Updates(map[string]any{
//			"title":   art.Title,
//			"content": art.Content,
//			"status":  art.Status,
//			"utime":   art.Utime,
//		}).Error
//	return err
//}
//
//func (dao *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
//	tx := dao.db.WithContext(ctx).Begin()
//	now := time.Now().UnixMilli()
//	defer tx.Rollback()
//	txDAO := NewGORMArticleDAO(tx)
//	var (
//		id  = art.Id
//		err error
//	)
//	if id == 0 {
//		id, err = txDAO.Insert(ctx, art)
//	} else {
//		err = txDAO.UpdateById(ctx, art)
//	}
//	if err != nil {
//		return 0, err
//	}
//	art.Id = id
//	publishArt := PublishedArticle(art)
//	publishArt.Utime = now
//	publishArt.Ctime = now
//	err = tx.Clauses(clause.OnConflict{
//		// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
//		Columns: []clause.Column{{Name: "id"}},
//		DoUpdates: clause.Assignments(map[string]interface{}{
//			"title":   art.Title,
//			"content": art.Content,
//			"status":  art.Status,
//			"utime":   now,
//		}),
//	}).Create(&publishArt).Error
//	if err != nil {
//		return 0, err
//	}
//	tx.Commit()
//	return id, tx.Error
//}
//
//func (dao *GORMArticleDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
//	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		res := tx.Model(&Article{}).
//			Where("id=? AND author_id = ?", id, author).
//			Update("status", status)
//		if res.Error != nil {
//			return res.Error
//		}
//		if res.RowsAffected != 1 {
//			return ErrPossibleIncorrectAuthor
//		}
//
//		res = tx.Model(&PublishedArticle{}).
//			Where("id=? AND author_id = ?", id, author).Update("status", status)
//		if res.Error != nil {
//			return res.Error
//		}
//		if res.RowsAffected != 1 {
//			return ErrPossibleIncorrectAuthor
//		}
//		return nil
//	})
//}
//
////type Article struct {
////	Id      int64  `gorm:"primaryKey,autoIncrement"`
////	Title   string `gorm:"type=varchar(1024)"`
////	Content string `gorm:"BLOB"`
////	// 如何设计索引
////	AuthorId int64 `gorm:"index"`
////	Ctime    int64
////	Utime    int64
////}
