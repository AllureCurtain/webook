package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository/article"
	repomocks "webook/internal/repository/mocks"
	"webook/pkg/logger"
)

func TestArticleService_PublishV1(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)

		art domain.Article

		wantErr error
		wantId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用量了制作库的 ID
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return ar, rr
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 1,
		},
		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用量了制作库的 ID
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return ar, rr
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 2,
		},
		{
			name: "修改保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("保存失败"))
				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				return ar, rr
			},
			art: domain.Article{
				Id:      7,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("保存失败"),
		},
		{
			name: "修改保存到线上库失败-重试都失败了",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用量了制作库的 ID
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					//	AnyTimes 确保一直都失败，最终引起重试次数耗尽，都失败了
				}).AnyTimes().Return(errors.New("保存到线上库失败"))
				return ar, rr
			},
			art: domain.Article{
				Id:      7,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("保存到线上库失败"),
		},
		{
			name: "修改保存到线上库失败-重试成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				ar := repomocks.NewMockArticleAuthorRepository(ctrl)
				ar.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				rr := repomocks.NewMockArticleReaderRepository(ctrl)
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用量了制作库的 ID
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("保存到线上库失败"))
				rr.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用量了制作库的 ID
					Id:      7,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
					//	重试成功了
				}).Return(nil)
				return ar, rr
			},
			art: domain.Article{
				Id:      7,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 7,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleServiceV1(author, reader, logger.NewNoOpLogger())
			id, err := svc.PublishV1(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
