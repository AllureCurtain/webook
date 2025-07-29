package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/integration/startup"
	"webook/internal/repository/dao"
	ijwt "webook/internal/web/jwt"
)

// ArticleTestSuite 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 在所有测试执行之前，初始化一些内容
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	s.db = startup.InitTestDB()
	artHdl := startup.InitArticleHandler()
	// 注册好了路由
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownTest() {
	// 情况所有数据，并且自增主键恢复到 1
	s.db.Exec("TRUNCATE TABLE articles")
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		// 集成测试准备数据
		before func(t *testing.T)
		// 集成测试验证数据
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子-保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				// 验证数据库
				var art dao.Article
				err := s.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "修改已有帖子，并保存",
			before: func(t *testing.T) {
				// 提前准备数据
				err := s.db.Create(dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Ctime:    123,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art dao.Article
				err := s.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)
				// 为了确保更新 Utime
				assert.True(t, art.Utime > 234)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Ctime:    123,
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "OK",
			},
		},
		{
			name: "修改别人的帖子",
			before: func(t *testing.T) {
				// 提前准备数据
				err := s.db.Create(dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 789,
					Ctime:    123,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art dao.Article
				err := s.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Ctime:    123,
					Utime:    234,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 0,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 构造请求

			// 执行

			// 验证结果

			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里就可以继续使用 req
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 相应写会到 resp 里
			s.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				var webRes Result[int64]
				err = json.NewDecoder(resp.Body).Decode(&webRes)
				require.NoError(t, err)
				assert.Equal(t, tc.wantRes, webRes)
			}
			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
