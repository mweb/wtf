package devto

import (
	"context"
	"fmt"

	"github.com/VictorAvelar/devto-api-go/devto"
	"github.com/rivo/tview"

	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget
	articles []devto.Article
	settings *Settings
	err      error
}

func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := &Widget{
		KeyboardWidget:   view.NewKeyboardWidget(app, pages, settings.common),
		ScrollableWidget: view.NewScrollableWidget(app, settings.common),

		settings: settings,
	}

	widget.SetRenderFunction(widget.Render)
	widget.View.SetScrollable(true)
	widget.initializeKeyboardControls()
	widget.View.SetInputCapture(widget.InputCapture)

	widget.KeyboardWidget.SetView(widget.View)

	return widget
}

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	ctx := context.Background()
	wCfg, _ := devto.NewConfig(false, "")

	c, _ := devto.NewClient(ctx, wCfg, nil, devto.BaseURL)

	options := devto.ArticleListOptions{
		Tags:     widget.settings.contentTag,
		Username: widget.settings.contentUsername,
		State:    widget.settings.contentState,
	}

	articles, err := c.Articles.List(ctx, options)
	if err != nil {
		widget.err = err
		widget.articles = nil
		widget.SetItemCount(0)
	} else {
		var displayArticles []devto.Article
		var l int
		if len(articles) < widget.settings.numberOfArticles {
			l = len(articles)
		} else {
			l = widget.settings.numberOfArticles - 1
		}
		for i, art := range articles {
			if i > l {
				break
			}
			displayArticles = append(displayArticles, art)
		}
		widget.articles = displayArticles
		widget.SetItemCount(len(displayArticles))
	}

	widget.Render()
}

// Render sets up the widget data for redrawing to the screen
func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() (string, string, bool) {
	title := fmt.Sprintf("%s - %s stories", widget.CommonSettings().Title, widget.settings.contentTag)

	if widget.err != nil {
		return title, widget.err.Error(), true
	}

	articles := widget.articles
	if articles == nil || len(articles) == 0 {
		return title, "No stories to display", false
	}
	var str string

	for idx, article := range articles {
		row := fmt.Sprintf(
			`[%s]%2d. %s [lightblue](%s)[white]`,
			widget.RowColor(idx),
			idx+1,
			article.Title,
			article.User.Username,
		)

		str += utils.HighlightableHelper(widget.View, row, idx, len(article.Title))
	}

	return title, str, false
}

func (widget *Widget) openStory() {
	sel := widget.GetSelected()
	if sel >= 0 && widget.articles != nil && sel < len(widget.articles) {
		article := &widget.articles[sel]
		utils.OpenFile(article.URL.String())
	}
}
