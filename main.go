package main

import (
	"fmt"
	"strings"

	"github.com/YumaFuu/ssm-tui/aws/ssm"
	component "github.com/YumaFuu/ssm-tui/node"
	"github.com/YumaFuu/ssm-tui/ui"
	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

func main() {
	// ctx := context.Background()

	client, err := ssm.NewClient()
	if err != nil {
		panic(err)
	}

	params, err := ssm.List(client, "/")
	if err != nil {
		panic(err)
	}
	app := tview.NewApplication()

	parameters := buildMapFromPaths(params)
	infoView := component.BuildInfoView()
	valueView := component.BuildValueView()
	tree := component.BuildParameterTree(
		parameters,
		infoView,
		valueView,
	)

	// // Function to display information of the selected node
	displayNodeInfo := func(node *tview.TreeNode) {
		if node == nil {
			infoView.SetText("")
			return
		}
		if len(node.GetChildren()) != 0 {
			infoView.SetText("Not a parameter")
			return
		}
		param := node.GetReference().(types.Parameter)
		info := fmt.Sprintf(
			`Version:          %d
Name:             %s
Type:             %s
LastModifiedDate: %s`,
			param.Version,
			*param.Name,
			param.Type,
			param.LastModifiedDate,
		)
		infoView.SetText(info)
		valueView.SetText(*param.Value)
	}

	// Capture the cursor movement
	tree.SetChangedFunc(func(node *tview.TreeNode) {
		displayNodeInfo(node)
	})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if len(node.GetChildren()) == 0 {
			// If the node has no children, it might be a file, display its value
			displayNodeInfo(node)
		}
		if node.IsExpanded() {
			node.CollapseAll()
		} else {
			node.ExpandAll()
		}
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'c' {
			node := tree.GetCurrentNode()
			if node != nil && len(node.GetChildren()) == 0 {
				param := node.GetReference().(types.Parameter)
				if err := clipboard.WriteAll(*param.Value); err != nil {
					infoView.SetText(fmt.Sprintf("[red]Error copying to clipboard: %s", err))
				} else {
					infoView.SetText("[green]Value copied to clipboard")
				}
			}
			return nil
		}
		return event
	})

	layout := ui.BuildLayout(tree, infoView, valueView)

	app.SetRoot(layout, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

// buildMapFromPaths function to build the map from the given paths
func buildMapFromPaths(params []types.Parameter) map[string]any {
	paramTree := make(map[string]any)
	for _, param := range params {
		key := *param.Name
		parts := strings.Split(strings.TrimPrefix(key, "/"), "/")
		insertPath(paramTree, parts, param)
	}
	return paramTree
}

// insertPath function to insert the parts into the map
func insertPath(m map[string]any, parts []string, param types.Parameter) {
	if len(parts) == 0 {
		return
	}
	if len(parts) == 1 {
		m[parts[0]] = param
		return
	}
	if _, ok := m[parts[0]]; !ok {
		m[parts[0]] = make(map[string]any)
	}
	insertPath(m[parts[0]].(map[string]any), parts[1:], param)
}
