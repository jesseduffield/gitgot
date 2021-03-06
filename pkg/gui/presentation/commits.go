package presentation

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
)

func GetCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitShaMap map[string]bool, diffName string, parseEmoji bool) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*models.Commit, map[string]bool, bool, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForCommit
	} else {
		displayFunc = getDisplayStringsForCommit
	}

	for i := range commits {
		diffed := commits[i].Sha == diffName
		lines[i] = displayFunc(commits[i], cherryPickedCommitShaMap, diffed, parseEmoji)
	}

	return lines
}

func getFullDescriptionDisplayStringsForCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed, parseEmoji bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	defaultColor := color.New(theme.DefaultTextColor)
	diffedColor := color.New(theme.DiffTerminalColor)

	// for some reason, setting the background to blue pads out the other commits
	// horizontally. For the sake of accessibility I'm considering this a feature,
	// not a bug
	copied := color.New(color.FgCyan, color.BgBlue)

	var shaColor *color.Color
	switch c.Status {
	case "unpushed":
		shaColor = red
	case "pushed":
		shaColor = yellow
	case "merged":
		shaColor = green
	case "rebasing":
		shaColor = blue
	case "reflog":
		shaColor = blue
	default:
		shaColor = defaultColor
	}

	if diffed {
		shaColor = diffedColor
	} else if cherryPickedCommitShaMap[c.Sha] {
		shaColor = copied
	}

	tagString := ""
	secondColumnString := blue.Sprint(utils.UnixToDate(c.UnixTimestamp))
	if c.Action != "" {
		secondColumnString = color.New(actionColorMap(c.Action)).Sprint(c.Action)
	} else if c.ExtraInfo != "" {
		tagColor := color.New(color.FgMagenta, color.Bold)
		tagString = utils.ColoredStringDirect(c.ExtraInfo, tagColor) + " "
	}

	truncatedAuthor := utils.TruncateWithEllipsis(c.Author, 17)

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{shaColor.Sprint(c.ShortSha()), secondColumnString, yellow.Sprint(truncatedAuthor), tagString + defaultColor.Sprint(name)}
}

func getDisplayStringsForCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed, parseEmoji bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	defaultColor := color.New(theme.DefaultTextColor)
	diffedColor := color.New(theme.DiffTerminalColor)

	// for some reason, setting the background to blue pads out the other commits
	// horizontally. For the sake of accessibility I'm considering this a feature,
	// not a bug
	copied := color.New(color.FgCyan, color.BgBlue)

	var shaColor *color.Color
	switch c.Status {
	case "unpushed":
		shaColor = red
	case "pushed":
		shaColor = yellow
	case "merged":
		shaColor = green
	case "rebasing":
		shaColor = blue
	case "reflog":
		shaColor = blue
	default:
		shaColor = defaultColor
	}

	if diffed {
		shaColor = diffedColor
	} else if cherryPickedCommitShaMap[c.Sha] {
		shaColor = copied
	}

	actionString := ""
	tagString := ""
	if c.Action != "" {
		actionString = color.New(actionColorMap(c.Action)).Sprint(utils.WithPadding(c.Action, 7)) + " "
	} else if len(c.Tags) > 0 {
		tagColor := color.New(color.FgMagenta, color.Bold)
		tagString = utils.ColoredStringDirect(strings.Join(c.Tags, " "), tagColor) + " "
	}

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{shaColor.Sprint(c.ShortSha()), actionString + tagString + defaultColor.Sprint(name)}
}

func actionColorMap(str string) color.Attribute {
	switch str {
	case "pick":
		return color.FgCyan
	case "drop":
		return color.FgRed
	case "edit":
		return color.FgGreen
	case "fixup":
		return color.FgMagenta
	default:
		return color.FgYellow
	}
}
