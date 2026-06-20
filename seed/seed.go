package seed

import (
	"blog/core"
	"blog/internal/utils"
	"blog/models"
	"time"

	"go.uber.org/zap"
)

func Run() {
	zap.L().Info("开始插入种子数据...")

	// 1. 分类
	categories := []models.Category{
		{ID: 1, Name: "前端技术", Description: "HTML、CSS、JavaScript、Vue、React 等前端技术分享", Sort: 10},
		{ID: 2, Name: "开源项目", Description: "开源项目开发记录与经验总结", Sort: 9},
		{ID: 3, Name: "工具效率", Description: "开发工具、效率方法和日常工作流分享", Sort: 8},
		{ID: 4, Name: "周刊", Description: "每周技术动态与思考", Sort: 7},
		{ID: 5, Name: "设计", Description: "UI设计、排版、视觉美学相关", Sort: 6},
		{ID: 6, Name: "年度总结", Description: "年度回顾与技术趋势观察", Sort: 5},
	}
	for _, c := range categories {
		if err := core.DB.FirstOrCreate(&models.Category{}, models.Category{Name: c.Name}).Error; err != nil {
			zap.L().Error("插入分类失败: " + c.Name + " " + err.Error())
		}
	}

	// 2. 管理员 (password: admin123)
	hash, _ := utils.HashPassword("admin123")
	if err := core.DB.FirstOrCreate(&models.Admin{}, models.Admin{
		Username: "admin",
		Password: hash,
		Nickname: "站长",
		Email:    "admin@blog.com",
	}).Error; err != nil {
		zap.L().Error("插入管理员失败: " + err.Error())
	} else {
		core.DB.Model(&models.Admin{}).Where("username = ?", "admin").Update("password", hash)
	}

	// 3. 示例用户
	userHash, _ := utils.HashPassword("123456")
	user := models.User{
		Email:    "demo@blog.com",
		Password: userHash,
		Nickname: "Demo用户",
		Status:   1,
	}
	core.DB.Where(models.User{Email: "demo@blog.com"}).FirstOrCreate(&user)

	// 4. 文章
	now := time.Now()
	articles := []models.Article{
		{
			Title:       "Vue 3 组合式 API 实战指南",
			Summary:     "深入理解 Vue 3 Composition API 的设计理念与最佳实践。",
			Content:     "<h2>为什么选择 Composition API？</h2><p>Vue 3 引入了 Composition API，它为开发者提供了更灵活的代码组织方式。</p><p>与 Options API 相比，Composition API 可以更好地提取和复用逻辑。</p><h2>核心概念</h2><ul><li>ref 和 reactive</li><li>computed 和 watch</li><li>生命周期钩子</li></ul><p>通过组合式函数（composables），我们可以把组件逻辑提取到独立的函数中，实现更好的代码复用。</p>",
			ContentType: 1, Cover: "", CategoryID: 1, AuthorID: 1,
			ViewCount: 1250, LikeCount: 86, CommentCount: 12,
			Status: 2, Tags: "Vue3,前端,JavaScript",
		},
		{
			Title:       "Go + Gin 构建 RESTful API 最佳实践",
			Summary:     "从项目结构到中间件设计，全面解析 Go Web 开发。",
			Content:     "## 项目结构\n\n良好的项目结构是可维护性的基础：\n\n- config/    配置加载\n- core/      核心初始化\n- internal/  业务逻辑\n  - dao/     数据访问层\n  - service/ 服务层\n  - router/  路由层\n- models/    数据模型\n- pkg/       公共包\n\n## 分层架构\n\n采用 Router -> Service -> DAO -> Model 四层架构，各层职责清晰。",
			ContentType: 2, Cover: "", CategoryID: 2, AuthorID: 1,
			ViewCount: 2350, LikeCount: 156, CommentCount: 28,
			Status: 2, Tags: "Go,Gin,后端,API",
		},
		{
			Title:       "我的 VS Code 插件推荐 2025",
			Summary:     "分享日常开发中最实用的 VS Code 插件和配置技巧。",
			Content:     "<p>工欲善其事，必先利其器。下面分享我每天都在用的 VS Code 插件：</p><h3>必装插件</h3><ol><li><strong>Volar</strong> - Vue 3 开发必备</li><li><strong>Go</strong> - Go 语言官方插件</li><li><strong>GitHub Copilot</strong> - AI 辅助编程</li><li><strong>Prettier</strong> - 代码格式化</li></ol><p>这些插件能极大提升开发效率。</p>",
			ContentType: 1, Cover: "", CategoryID: 3, AuthorID: 1,
			ViewCount: 980, LikeCount: 72, CommentCount: 8,
			Status: 2, Tags: "VS Code,效率,工具",
		},
		{
			Title:       "潮流周刊第 99 期 — AI 编程工具全面对比",
			Summary:     "GitHub Copilot vs Claude Code vs Cursor，谁是 AI 编程之王？",
			Content:     "## AI 编程工具对比\n\n| 工具 | 优势 | 劣势 |\n| Copilot | 深度 IDE 集成 | 上下文有限 |\n| Claude Code | 强大推理能力 | CLI 操作 |\n| Cursor | 全功能 IDE | 需付费 |\n\n## 结论\n\n不同场景选择不同工具，三者配合使用效果最佳。",
			ContentType: 2, Cover: "", CategoryID: 4, AuthorID: 1,
			ViewCount: 5200, LikeCount: 320, CommentCount: 56,
			Status: 2, Tags: "AI,编程工具,周刊",
		},
		{
			Title:       "极简主义在 UI 设计中的应用",
			Summary:     "探讨如何在界面设计中追求简洁与功能的完美平衡。",
			Content:     "<p>极简设计不是越少越好，而是恰到好处。</p><h3>核心原则</h3><ul><li>去掉多余的装饰</li><li>保持留白空间</li><li>统一的间距系统</li><li>克制的色彩使用</li></ul><p>好的设计让人感觉「本该如此」。</p>",
			ContentType: 1, Cover: "", CategoryID: 5, AuthorID: 1,
			ViewCount: 760, LikeCount: 48, CommentCount: 6,
			Status: 2, Tags: "UI设计,极简主义",
		},
		{
			Title:       "2025 年度技术回顾与展望",
			Summary:     "回顾过去一年的技术发展，展望新一年的技术趋势。",
			Content:     "## 2025 年度回顾\n\n### 大前端\n- Vue 3 生态持续完善\n- React 19 带来全新特性\n- Rust 工具链在前端领域加速\n\n### 后端\n- Go 继续在企业级应用中占据重要地位\n- Rust 后端框架逐渐成熟\n- AI 原生应用成为新趋势\n\n## 2026 展望\n\nAI 将继续深刻改变软件开发的方式，但基础技术的重要性不会减弱。",
			ContentType: 2, Cover: "", CategoryID: 6, AuthorID: 1,
			ViewCount: 3800, LikeCount: 245, CommentCount: 42,
			Status: 2, Tags: "年度总结,技术趋势,2025",
		},
		{
			Title:       "从零搭建个人博客全栈项目",
			Summary:     "使用 Go + Vue 3 从零开始搭建一个完整的个人博客。",
			Content:     "<h2>技术选型</h2><ul><li>后端：Go + Gin + GORM + MySQL</li><li>前端：Vue 3 + Vite + Vue Router</li><li>编辑器：wangEditor + md-editor-v3</li></ul><h2>核心功能</h2><p>支持文章管理、评论系统、分类管理、图片上传、数据面板等功能。</p><p>前后端分离架构，JWT 认证，统一响应格式。</p>",
			ContentType: 1, Cover: "", CategoryID: 2, AuthorID: 1,
			ViewCount: 650, LikeCount: 35, CommentCount: 5,
			Status: 2, Tags: "全栈,Go,Vue3,博客",
		},
		{
			Title:       "[草稿] Docker 容器化部署实践",
			Summary:     "记录将 Go 应用打包成 Docker 镜像并部署到服务器的全过程。",
			Content:     "# Docker 部署实践\n\n内容完善中...",
			ContentType: 2, Cover: "", CategoryID: 2, AuthorID: 1,
			ViewCount: 0, LikeCount: 0, CommentCount: 0,
			Status: 1, Tags: "Docker,部署",
		},
	}

	for i := range articles {
		a := &articles[i]
		if a.PublishTime == nil && a.Status == 2 {
			t := now.Add(-time.Duration(7-i) * 24 * time.Hour)
			a.PublishTime = &t
		}
		a.CreatedAt = now.Add(-time.Duration(8-i) * 24 * time.Hour)
		a.UpdatedAt = now.Add(-time.Duration(1) * 24 * time.Hour)
		if err := core.DB.Where(models.Article{Title: a.Title}).FirstOrCreate(a).Error; err != nil {
			zap.L().Error("插入文章失败: " + a.Title + " " + err.Error())
		}
	}

	zap.L().Info("种子数据插入完成")
}
