package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ============================================
// WebSearchSkill - 生产级网页搜索技能
// ============================================

// WebSearchParameters 定义了网页搜索技能的参数结构
type WebSearchParameters struct {
	Type       string                   `json:"type"`
	Properties WebSearchParamProperties `json:"properties"`
	Required   []string                 `json:"required"`
}

// WebSearchParamProperties 定义参数属性
type WebSearchParamProperties struct {
	Query      WebSearchParamProperty `json:"query"`
	Num        WebSearchParamProperty `json:"num"`
	Engine     WebSearchParamProperty `json:"engine"`
	Language   WebSearchParamProperty `json:"language"`
	Region     WebSearchParamProperty `json:"region"`
	SafeSearch WebSearchParamProperty `json:"safe_search"`
}

// WebSearchParamProperty 定义单个参数属性
type WebSearchParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// WebSearchArgs 定义执行搜索时的参数
type WebSearchArgs struct {
	Query      string `json:"query"`
	Num        int    `json:"num,omitempty"`
	Engine     string `json:"engine,omitempty"`
	Language   string `json:"language,omitempty"`
	Region     string `json:"region,omitempty"`
	SafeSearch bool   `json:"safe_search,omitempty"`
}

// SearchResult 定义单个搜索结果
type SearchResult struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Snippet     string `json:"snippet"`
	Description string `json:"description,omitempty"`
	HostName    string `json:"host_name"`
	Rank        int    `json:"rank"`
	Date        string `json:"date,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
}

// WebSearchResult 定义搜索结果结构
type WebSearchResult struct {
	Query       string         `json:"query"`
	Total       int            `json:"total"`
	Engine      string         `json:"engine"`
	SearchTime  float64        `json:"search_time_ms"`
	Results     []SearchResult `json:"results"`
	Suggestions []string       `json:"suggestions,omitempty"`
	Error       string         `json:"error,omitempty"`
}

// SearchEngine 搜索引擎接口
type SearchEngine interface {
	Search(ctx context.Context, args WebSearchArgs) (*WebSearchResult, error)
	Name() string
}

// WebSearchSkill 实现了网页搜索技能
type WebSearchSkill struct {
	name          string
	descName      string
	description   string
	parameters    WebSearchParameters
	httpClient    *http.Client
	engines       map[string]SearchEngine
	defaultEngine string
}

// WebSearchConfig 搜索技能配置
type WebSearchConfig struct {
	// HTTP超时时间
	Timeout time.Duration
	// 默认搜索引擎
	DefaultEngine string
	// SerpAPI密钥（可选）
	SerpAPIKey string
	// Google Custom Search API密钥（可选）
	GoogleAPIKey string
	// Google Custom Search引擎ID（可选）
	GoogleEngineID string
	// Bing API密钥（可选）
	BingAPIKey string
	// 代理URL（可选）
	ProxyURL string
	// User-Agent
	UserAgent string
}

// NewWebSearchSkill 创建一个新的网页搜索技能实例
func NewWebSearchSkill(config WebSearchConfig) *WebSearchSkill {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.DefaultEngine == "" {
		config.DefaultEngine = "duckduckgo"
	}
	if config.UserAgent == "" {
		config.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}

	// 创建HTTP客户端
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}

	if config.ProxyURL != "" {
		proxyURL, _ := url.Parse(config.ProxyURL)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	// 注册搜索引擎
	engines := make(map[string]SearchEngine)

	// DuckDuckGo - 无需API密钥
	engines["duckduckgo"] = NewDuckDuckGoEngine(httpClient, config.UserAgent)

	// SerpAPI - 需要API密钥
	if config.SerpAPIKey != "" {
		engines["serpapi"] = NewSerpAPIEngine(httpClient, config.SerpAPIKey)
	}

	// Google Custom Search - 需要API密钥和引擎ID
	if config.GoogleAPIKey != "" && config.GoogleEngineID != "" {
		engines["google"] = NewGoogleSearchEngine(httpClient, config.GoogleAPIKey, config.GoogleEngineID)
	}

	// Bing Search - 需要API密钥
	if config.BingAPIKey != "" {
		engines["bing"] = NewBingSearchEngine(httpClient, config.BingAPIKey)
	}

	return &WebSearchSkill{
		name:     "web_search",
		descName: "网页搜索",
		description: `网页搜索技能允许大模型通过互联网搜索获取实时信息。

这个技能可以帮助大模型：
- 搜索最新的新闻和事件
- 获取实时数据和统计信息
- 查找特定主题的相关资料
- 验证信息的准确性

支持的搜索引擎：
- duckduckgo: DuckDuckGo搜索（无需API密钥）
- google: Google自定义搜索（需要API密钥）
- bing: Bing搜索（需要API密钥）
- serpapi: SerpAPI聚合搜索（需要API密钥）`,
		parameters: WebSearchParameters{
			Type: "object",
			Properties: WebSearchParamProperties{
				Query: WebSearchParamProperty{
					Type:        "string",
					Description: "搜索查询关键词或问题",
				},
				Num: WebSearchParamProperty{
					Type:        "integer",
					Description: "返回结果数量，默认为10",
					Default:     10,
				},
				Engine: WebSearchParamProperty{
					Type:        "string",
					Description: "搜索引擎选择",
					Default:     config.DefaultEngine,
					Enum:        []string{"duckduckgo", "google", "bing", "serpapi"},
				},
				Language: WebSearchParamProperty{
					Type:        "string",
					Description: "搜索语言代码",
					Default:     "zh-CN",
				},
				Region: WebSearchParamProperty{
					Type:        "string",
					Description: "搜索地区代码",
					Default:     "cn",
				},
				SafeSearch: WebSearchParamProperty{
					Type:        "boolean",
					Description: "是否启用安全搜索",
					Default:     true,
				},
			},
			Required: []string{"query"},
		},
		httpClient:    httpClient,
		engines:       engines,
		defaultEngine: config.DefaultEngine,
	}
}

// GetName 返回技能名称
func (s *WebSearchSkill) GetName() string {
	return s.name
}

// GetDescName 返回技能的描述性名称
func (s *WebSearchSkill) GetDescName() string {
	return s.descName
}

// GetDescription 返回技能的详细描述
func (s *WebSearchSkill) GetDescription() string {
	return s.description
}

// GetParameters 返回技能的参数定义
func (s *WebSearchSkill) GetParameters() any {
	return s.parameters
}

// Execute 执行网页搜索
func (s *WebSearchSkill) Execute(ctx context.Context, args string) (string, error) {
	// 解析参数
	var searchArgs WebSearchArgs
	if err := json.Unmarshal([]byte(args), &searchArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	// 验证必要参数
	if searchArgs.Query == "" {
		return "", fmt.Errorf("搜索查询不能为空")
	}

	// 设置默认值
	if searchArgs.Num <= 0 {
		searchArgs.Num = 10
	}
	if searchArgs.Engine == "" {
		searchArgs.Engine = s.defaultEngine
	}
	if searchArgs.Language == "" {
		searchArgs.Language = "zh-CN"
	}
	if searchArgs.Region == "" {
		searchArgs.Region = "cn"
	}

	// 获取搜索引擎
	engine, exists := s.engines[searchArgs.Engine]
	if !exists {
		return "", fmt.Errorf("不支持的搜索引擎: %s，可用引擎: %v",
			searchArgs.Engine, s.getAvailableEngines())
	}

	// 执行搜索
	result, err := engine.Search(ctx, searchArgs)
	if err != nil {
		return "", fmt.Errorf("搜索失败: %w", err)
	}

	// 限制结果数量
	if len(result.Results) > searchArgs.Num {
		result.Results = result.Results[:searchArgs.Num]
	}
	result.Total = len(result.Results)

	// 将结果序列化为JSON
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(resultJSON), nil
}

// getAvailableEngines 获取可用的搜索引擎列表
func (s *WebSearchSkill) getAvailableEngines() []string {
	engines := make([]string, 0, len(s.engines))
	for name := range s.engines {
		engines = append(engines, name)
	}
	return engines
}

// ToJSON 将技能信息序列化为JSON
func (s *WebSearchSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// FromJSON 从JSON反序列化技能信息
func (s *WebSearchSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string              `json:"name"`
		DescName    string              `json:"descName"`
		Description string              `json:"description"`
		Parameters  WebSearchParameters `json:"parameters"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}

	s.name = data.Name
	s.descName = data.DescName
	s.description = data.Description
	s.parameters = data.Parameters

	return nil
}

// ============================================
// DuckDuckGo 搜索引擎实现（无需API密钥）
// ============================================

// DuckDuckGoEngine DuckDuckGo搜索引擎
type DuckDuckGoEngine struct {
	httpClient *http.Client
	userAgent  string
}

// NewDuckDuckGoEngine 创建DuckDuckGo搜索引擎
func NewDuckDuckGoEngine(httpClient *http.Client, userAgent string) *DuckDuckGoEngine {
	return &DuckDuckGoEngine{
		httpClient: httpClient,
		userAgent:  userAgent,
	}
}

// Name 返回引擎名称
func (e *DuckDuckGoEngine) Name() string {
	return "duckduckgo"
}

// Search 执行搜索
func (e *DuckDuckGoEngine) Search(ctx context.Context, args WebSearchArgs) (*WebSearchResult, error) {
	startTime := time.Now()
	result := &WebSearchResult{
		Query:  args.Query,
		Engine: "duckduckgo",
	}

	// 方式1: 使用DuckDuckGo HTML版本爬取
	results, err := e.searchHTML(ctx, args)
	if err == nil && len(results) > 0 {
		result.Results = results
		result.SearchTime = float64(time.Since(startTime).Milliseconds())
		return result, nil
	}

	// 方式2: 使用DuckDuckGo Instant Answer API
	results, err = e.searchInstantAnswer(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("DuckDuckGo搜索失败: %w", err)
	}

	result.Results = results
	result.SearchTime = float64(time.Since(startTime).Milliseconds())
	return result, nil
}

// searchHTML 通过爬取HTML页面获取搜索结果
func (e *DuckDuckGoEngine) searchHTML(ctx context.Context, args WebSearchArgs) ([]SearchResult, error) {
	// 构建搜索URL
	queryURL := fmt.Sprintf(
		"https://html.duckduckgo.com/html/?q=%s&kl=%s",
		url.QueryEscape(args.Query),
		url.QueryEscape(args.Region+"-"+args.Language),
	)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("User-Agent", e.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")

	// 发送请求
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 提取搜索结果
	var results []SearchResult
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		if i >= args.Num {
			return
		}

		// 提取标题和链接
		titleElem := s.Find(".result__title a")
		title := strings.TrimSpace(titleElem.Text())
		href, _ := titleElem.Attr("href")

		// DuckDuckGo使用重定向URL，需要提取真实URL
		if strings.Contains(href, "uddg=") {
			if u, err := url.Parse(href); err == nil {
				if realURL := u.Query().Get("uddg"); realURL != "" {
					href = realURL
				}
			}
		}

		// 提取摘要
		snippet := strings.TrimSpace(s.Find(".result__snippet").Text())

		// 提取主机名
		hostname := ""
		if parsedURL, err := url.Parse(href); err == nil {
			hostname = parsedURL.Hostname()
		}

		if title != "" && href != "" {
			results = append(results, SearchResult{
				URL:      href,
				Title:    title,
				Snippet:  snippet,
				HostName: hostname,
				Rank:     i + 1,
			})
		}
	})

	return results, nil
}

// searchInstantAnswer 使用DuckDuckGo Instant Answer API
func (e *DuckDuckGoEngine) searchInstantAnswer(ctx context.Context, args WebSearchArgs) ([]SearchResult, error) {
	// 构建API URL
	apiURL := fmt.Sprintf(
		"https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1",
		url.QueryEscape(args.Query),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", e.userAgent)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析JSON响应
	var ddgResp struct {
		AbstractText   string `json:"AbstractText"`
		AbstractURL    string `json:"AbstractURL"`
		AbstractSource string `json:"AbstractSource"`
		Heading        string `json:"Heading"`
		RelatedTopics  []struct {
			Text string `json:"Text"`
			URL  string `json:"FirstURL"`
		} `json:"RelatedTopics"`
		Results []struct {
			Text string `json:"Text"`
			URL  string `json:"FirstURL"`
		} `json:"Results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ddgResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	var results []SearchResult

	// 添加摘要结果
	if ddgResp.AbstractText != "" && ddgResp.AbstractURL != "" {
		hostname := ""
		if parsedURL, err := url.Parse(ddgResp.AbstractURL); err == nil {
			hostname = parsedURL.Hostname()
		}
		results = append(results, SearchResult{
			URL:      ddgResp.AbstractURL,
			Title:    ddgResp.Heading,
			Snippet:  ddgResp.AbstractText,
			HostName: hostname,
			Rank:     1,
		})
	}

	// 添加相关主题
	for i, topic := range ddgResp.RelatedTopics {
		if i >= args.Num-1 {
			break
		}
		if topic.URL != "" && topic.Text != "" {
			hostname := ""
			if parsedURL, err := url.Parse(topic.URL); err == nil {
				hostname = parsedURL.Hostname()
			}
			results = append(results, SearchResult{
				URL:      topic.URL,
				Title:    extractTitleFromText(topic.Text),
				Snippet:  topic.Text,
				HostName: hostname,
				Rank:     i + 2,
			})
		}
	}

	return results, nil
}

// extractTitleFromText 从文本中提取标题
func extractTitleFromText(text string) string {
	// 通常格式是 "Title - Description"
	parts := strings.SplitN(text, " - ", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return text
}

// ============================================
// SerpAPI 搜索引擎实现
// ============================================

// SerpAPIEngine SerpAPI搜索引擎
type SerpAPIEngine struct {
	httpClient *http.Client
	apiKey     string
}

// NewSerpAPIEngine 创建SerpAPI搜索引擎
func NewSerpAPIEngine(httpClient *http.Client, apiKey string) *SerpAPIEngine {
	return &SerpAPIEngine{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

// Name 返回引擎名称
func (e *SerpAPIEngine) Name() string {
	return "serpapi"
}

// Search 执行搜索
func (e *SerpAPIEngine) Search(ctx context.Context, args WebSearchArgs) (*WebSearchResult, error) {
	startTime := time.Now()
	result := &WebSearchResult{
		Query:  args.Query,
		Engine: "serpapi",
	}

	// 构建API URL
	apiURL := fmt.Sprintf(
		"https://serpapi.com/search.json?q=%s&api_key=%s&num=%d&hl=%s&gl=%s",
		url.QueryEscape(args.Query),
		url.QueryEscape(e.apiKey),
		args.Num,
		url.QueryEscape(args.Language),
		url.QueryEscape(args.Region),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var serpResp struct {
		OrganicResults []struct {
			Position int    `json:"position"`
			Title    string `json:"title"`
			Link     string `json:"link"`
			Snippet  string `json:"snippet"`
			Date     string `json:"date"`
		} `json:"organic_results"`
		SearchInformation struct {
			QueryDisplayed string `json:"query_displayed"`
		} `json:"search_information"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &serpResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	if serpResp.Error != "" {
		return nil, fmt.Errorf("SerpAPI错误: %s", serpResp.Error)
	}

	// 转换结果
	for i, r := range serpResp.OrganicResults {
		hostname := ""
		if parsedURL, err := url.Parse(r.Link); err == nil {
			hostname = parsedURL.Hostname()
		}
		result.Results = append(result.Results, SearchResult{
			URL:      r.Link,
			Title:    r.Title,
			Snippet:  r.Snippet,
			HostName: hostname,
			Rank:     i + 1,
			Date:     r.Date,
		})
	}

	result.SearchTime = float64(time.Since(startTime).Milliseconds())
	return result, nil
}

// ============================================
// Google Custom Search 实现
// ============================================

// GoogleSearchEngine Google自定义搜索引擎
type GoogleSearchEngine struct {
	httpClient *http.Client
	apiKey     string
	engineID   string
}

// NewGoogleSearchEngine 创建Google搜索引擎
func NewGoogleSearchEngine(httpClient *http.Client, apiKey, engineID string) *GoogleSearchEngine {
	return &GoogleSearchEngine{
		httpClient: httpClient,
		apiKey:     apiKey,
		engineID:   engineID,
	}
}

// Name 返回引擎名称
func (e *GoogleSearchEngine) Name() string {
	return "google"
}

// Search 执行搜索
func (e *GoogleSearchEngine) Search(ctx context.Context, args WebSearchArgs) (*WebSearchResult, error) {
	startTime := time.Now()
	result := &WebSearchResult{
		Query:  args.Query,
		Engine: "google",
	}

	// 构建API URL
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d&lr=lang_%s",
		url.QueryEscape(e.apiKey),
		url.QueryEscape(e.engineID),
		url.QueryEscape(args.Query),
		args.Num,
		url.QueryEscape(strings.Split(args.Language, "-")[0]),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var googleResp struct {
		Items []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
			Pagemap struct {
				Metatags []struct {
					ArticlePublishedTime string `json:"article:published_time"`
				} `json:"metatags"`
			} `json:"pagemap"`
		} `json:"items"`
		SearchInformation struct {
			SearchTime float64 `json:"searchTime"`
		} `json:"searchInformation"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &googleResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	if googleResp.Error != nil {
		return nil, fmt.Errorf("Google API错误 %d: %s", googleResp.Error.Code, googleResp.Error.Message)
	}

	// 转换结果
	for i, item := range googleResp.Items {
		hostname := ""
		if parsedURL, err := url.Parse(item.Link); err == nil {
			hostname = parsedURL.Hostname()
		}

		date := ""
		if len(item.Pagemap.Metatags) > 0 {
			date = item.Pagemap.Metatags[0].ArticlePublishedTime
		}

		result.Results = append(result.Results, SearchResult{
			URL:      item.Link,
			Title:    item.Title,
			Snippet:  item.Snippet,
			HostName: hostname,
			Rank:     i + 1,
			Date:     date,
		})
	}

	result.SearchTime = float64(time.Since(startTime).Milliseconds())
	return result, nil
}

// ============================================
// Bing Search API 实现
// ============================================

// BingSearchEngine Bing搜索引擎
type BingSearchEngine struct {
	httpClient *http.Client
	apiKey     string
}

// NewBingSearchEngine 创建Bing搜索引擎
func NewBingSearchEngine(httpClient *http.Client, apiKey string) *BingSearchEngine {
	return &BingSearchEngine{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

// Name 返回引擎名称
func (e *BingSearchEngine) Name() string {
	return "bing"
}

// Search 执行搜索
func (e *BingSearchEngine) Search(ctx context.Context, args WebSearchArgs) (*WebSearchResult, error) {
	startTime := time.Now()
	result := &WebSearchResult{
		Query:  args.Query,
		Engine: "bing",
	}

	// 构建API URL
	apiURL := fmt.Sprintf(
		"https://api.bing.microsoft.com/v7.0/search?q=%s&count=%d&mkt=%s&setLang=%s",
		url.QueryEscape(args.Query),
		args.Num,
		url.QueryEscape(args.Region+"-"+strings.ToUpper(args.Language)),
		url.QueryEscape(args.Language),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", e.apiKey)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var bingResp struct {
		WebPages struct {
			Value []struct {
				ID              string `json:"id"`
				Name            string `json:"name"`
				URL             string `json:"url"`
				Snippet         string `json:"snippet"`
				DateLastCrawled string `json:"dateLastCrawled"`
				DisplayURL      string `json:"displayUrl"`
			} `json:"value"`
		} `json:"webPages"`
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &bingResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	if bingResp.Error != nil {
		return nil, fmt.Errorf("Bing API错误 %s: %s", bingResp.Error.Code, bingResp.Error.Message)
	}

	// 转换结果
	for i, item := range bingResp.WebPages.Value {
		hostname := ""
		if parsedURL, err := url.Parse(item.URL); err == nil {
			hostname = parsedURL.Hostname()
		}

		result.Results = append(result.Results, SearchResult{
			URL:      item.URL,
			Title:    item.Name,
			Snippet:  item.Snippet,
			HostName: hostname,
			Rank:     i + 1,
			Date:     item.DateLastCrawled,
		})
	}

	result.SearchTime = float64(time.Since(startTime).Milliseconds())
	return result, nil
}
