package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/smallnest/langgraphgo/tool"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type logKey struct{}

func logf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	// Always print to stdout
	fmt.Print(msg)

	// If log channel exists in context, send it there too
	if ch, ok := ctx.Value(logKey{}).(chan string); ok {
		// Non-blocking send to avoid stalling if channel is full or no one listening
		select {
		case ch <- msg:
		default:
		}
	}
}

// PlannerNode generates a research plan based on the query.
func PlannerNode(ctx context.Context, state interface{}) (interface{}, error) {
	s := state.(*State)
	logf(ctx, "--- 规划节点：正在为查询 '%s' 进行规划 ---\n", s.Request.Query)

	llm, err := getLLM()
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf("你是一名研究规划师。请为以下查询创建一个分步研究计划：%s。仅返回编号列表形式的计划。必须使用中文回复。", s.Request.Query)
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		return nil, err
	}

	// Simple parsing of the plan (splitting by newlines)
	lines := strings.Split(completion, "\n")
	var plan []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			plan = append(plan, trimmed)
		}
	}
	s.Plan = plan
	s.Plan = plan

	// Format plan for better readability
	formattedPlan := "生成的计划：\n"
	for _, step := range s.Plan {
		formattedPlan += fmt.Sprintf("%s\n", step)
	}
	logf(ctx, "%s", formattedPlan)

	return s, nil
}

// ResearcherNode executes the research plan.
func ResearcherNode(ctx context.Context, state interface{}) (interface{}, error) {
	s := state.(*State)
	logf(ctx, "--- 研究节点：正在执行计划（并发） ---\n")

	// Create Tavily search tool
	tavilyTool, err := tool.NewTavilySearch("", tool.WithTavilySearchDepth("advanced"))
	if err != nil {
		logf(ctx, "警告：无法初始化 Tavily 搜索工具 (%v)，将使用 LLM 模拟研究\n", err)
		// Fallback to LLM-based research
		return researchWithLLM(ctx, s)
	}

	type stepResult struct {
		summary string
		images  []string
	}

	results := make([]stepResult, len(s.Plan))
	var wg sync.WaitGroup
	// Limit concurrency to avoid rate limits
	sem := make(chan struct{}, 5)

	for i, step := range s.Plan {
		wg.Add(1)
		go func(i int, step string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			logf(ctx, "正在研究步骤：%s\n", step)

			// Use Tavily to search for information with images
			searchResult, err := tavilyTool.CallWithImages(ctx, step)
			if err != nil {
				logf(ctx, "搜索失败 (%v)，跳过此步骤\n", err)
				return
			}

			// Collect images (limit to first 1 per step to avoid too many images)
			var stepImages []string
			imageCount := 0
			for _, imgURL := range searchResult.Images {
				if imageCount >= 1 {
					break
				}
				stepImages = append(stepImages, imgURL)
				imageCount++
			}

			// Use LLM to summarize the search results
			llm, err := getLLM()
			if err != nil {
				logf(ctx, "LLM 初始化失败 (%v)\n", err)
				return
			}

			prompt := fmt.Sprintf("你是一名研究员。请根据以下搜索结果为研究步骤 '%s' 提供详细摘要。必须使用中文回复。\n\n搜索结果：\n%s", step, searchResult.Text)
			summary, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
			if err != nil {
				logf(ctx, "摘要生成失败 (%v)，使用原始搜索结果\n", err)
				summary = searchResult.Text
			}

			results[i] = stepResult{
				summary: summary,
				images:  stepImages,
			}
		}(i, step)
	}
	wg.Wait()

	// Aggregate results
	var finalResults []string
	var allImages []string

	for i, res := range results {
		if res.summary == "" {
			continue
		}

		// Calculate image indices for this step to help LLM reference them
		var imgIndices []string
		startIdx := len(allImages) + 1
		for j, img := range res.images {
			allImages = append(allImages, img)
			// We don't strictly need to store indices here if we just append,
			// but passing the ID to LLM helps it know which image is which.
			imgIndices = append(imgIndices, fmt.Sprintf("IMAGE_%d", startIdx+j))
		}

		imgNote := ""
		if len(imgIndices) > 0 {
			imgNote = fmt.Sprintf("\n(可用图片: %s)", strings.Join(imgIndices, ", "))
		}

		finalResults = append(finalResults, fmt.Sprintf("Step: %s\nFindings: %s%s", s.Plan[i], res.summary, imgNote))
	}

	s.ResearchResults = finalResults
	s.Images = allImages
	logf(ctx, "收集到 %d 张图片\n", len(allImages))
	return s, nil
}

// researchWithLLM is a fallback function that uses LLM to simulate research
func researchWithLLM(ctx context.Context, s *State) (interface{}, error) {
	llm, err := getLLM()
	if err != nil {
		return nil, err
	}

	var results []string
	for _, step := range s.Plan {
		logf(ctx, "正在研究步骤（使用 LLM）：%s\n", step)
		prompt := fmt.Sprintf("你是一名研究员。请为这个研究步骤查找详细信息：%s。提供发现摘要。必须使用中文回复。", step)
		completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
		if err != nil {
			return nil, err
		}
		results = append(results, fmt.Sprintf("Step: %s\nFindings: %s", step, completion))
	}

	s.ResearchResults = results
	return s, nil
}

// Replace image placeholders with actual image tags
// Regex matches [IMAGE_X：Title] or [IMAGE_X:Title]
var imgRe = regexp.MustCompile(`\[IMAGE_(\d+)[：:]([^\]]+)\]`)

// ReporterNode compiles the final report.
func ReporterNode(ctx context.Context, state interface{}) (interface{}, error) {
	s := state.(*State)
	logf(ctx, "--- 报告节点：正在生成最终报告 ---\n")

	llm, err := getLLM()
	if err != nil {
		return nil, err
	}

	researchData := strings.Join(s.ResearchResults, "\n\n")

	// Inform LLM about available images
	imageInfo := ""
	if len(s.Images) > 0 {
		imageInfo = fmt.Sprintf("\n\n注意：研究过程中收集到 %d 张相关图片。在报告中适当的位置，你可以使用 [IMAGE_X：图片标题] 占位符来标记应该插入图片的位置（X 为 1 到 %d，图片标题为你为该图片起的标题）。例如：[IMAGE_1：某某图表]。请务必确保引用的图片与周围的文字内容高度相关，如果图片与当前段落无关，请不要强行插入。", len(s.Images), len(s.Images))
	}

	prompt := fmt.Sprintf("你是一名资深报告撰写员。请根据以下研究结果撰写一份全面的最终报告。使用 Markdown 格式，包含清晰的标题、要点，并在适当的地方使用代码块。数学公式请使用 ```math 代码块包裹，或者使用 $$...$$ (块级) 和 $...$ (行内) 包裹。不要透漏撰写人信息。%s必须使用中文撰写报告：\n\n%s\n\n原始查询是：%s", imageInfo, researchData, s.Request.Query)

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		return nil, err
	}

	// Convert Markdown to HTML
	// Clean up markdown code blocks if present
	completion = strings.TrimPrefix(completion, "```markdown")
	completion = strings.TrimPrefix(completion, "```")
	completion = strings.TrimSuffix(completion, "```")

	completion = imgRe.ReplaceAllStringFunc(completion, func(match string) string {
		parts := imgRe.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		idxStr := parts[1]
		title := strings.TrimSpace(parts[2])

		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx < 1 || idx > len(s.Images) {
			return match
		}

		imgURL := s.Images[idx-1]
		return fmt.Sprintf("\n\n<img src=\"%s\" alt=\"%s\" style=\"max-width: 90%%; display: block; margin: 10px auto;\" />\n\n", imgURL, title)
	})

	// If LLM didn't use placeholders, append images at the end
	if len(s.Images) > 0 && !strings.Contains(completion, "<img") {
		completion += "\n\n## 相关图片\n\n"
		for i, imgURL := range s.Images {
			completion += fmt.Sprintf("<img src=\"%s\" alt=\"图片 %d\" style=\"max-width: 90%%; display: block; margin: 10px auto;\" />\n\n", imgURL, i+1)
		}
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(completion))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	s.FinalReport = string(markdown.Render(doc, renderer))
	logf(ctx, "最终报告已生成（包含 %d 张图片）。\n", len(s.Images))
	return s, nil
}

func getLLM() (llms.Model, error) {
	// Use DeepSeek as per user preference
	// Ensure OPENAI_API_KEY and OPENAI_API_BASE are set in the environment
	return openai.New()
}
