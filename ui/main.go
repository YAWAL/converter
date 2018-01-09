package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
	"github.com/google/gxui/themes/light"
	"github.com/KitlerUa/csvparser/parser"
)

var (
	fileColor = gxui.Color{R: 0.7, G: 0.8, B: 1.0, A: 1}
	//fileColor = gxui.Black
	directoryColor = gxui.Color{R: 0.8, G: 1.0, B: 0.7, A: 1}
	//directoryColor = gxui.Black
	filename string
)

// filesAt returns a list of all immediate files in the given directory path.
func filesAt(path string) []string {
	files := []string{}
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err == nil && path != subpath {
			files = append(files, subpath)
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return files
}

// filesAdapter is an implementation of the gxui.ListAdapter interface.
// The AdapterItems returned by this adapter are absolute file path strings.
type filesAdapter struct {
	gxui.AdapterBase
	files []string // The absolute file paths
}

// SetFiles assigns the specified list of absolute-path files to this adapter.
func (a *filesAdapter) SetFiles(files []string) {
	a.files = files
	a.DataChanged(false)
}

func (a *filesAdapter) Count() int {
	return len(a.files)
}

func (a *filesAdapter) ItemAt(index int) gxui.AdapterItem {
	return a.files[index]
}

func (a *filesAdapter) ItemIndex(item gxui.AdapterItem) int {
	path := item.(string)
	for i, f := range a.files {
		if f == path {
			return i
		}
	}
	return -1 // Not found
}

func (a *filesAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	path := a.files[index]
	_, name := filepath.Split(path)
	label := theme.CreateLabel()
	label.SetText(name)
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		label.SetColor(directoryColor)
	} else {
		label.SetColor(fileColor)
	}
	return label
}

func (a *filesAdapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}

// directory implements the gxui.TreeNode interface to represent a directory
// node in a file-system.
type directory struct {
	path    string   // The absolute path of this directory.
	subdirs []string // The absolute paths of all immediate sub-directories.
}

// directoryAt returns a directory structure populated with the immediate
// subdirectories at the given path.
func directoryAt(path string) directory {
	directory := directory{path: path}
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err == nil && path != subpath {
			if info.IsDir() {
				directory.subdirs = append(directory.subdirs, subpath)
				return filepath.SkipDir
			}
		}
		return nil
	})
	return directory
}

// Count implements gxui.TreeNodeContainer.
func (d directory) Count() int {
	return len(d.subdirs)
}

// NodeAt implements gxui.TreeNodeContainer.
func (d directory) NodeAt(index int) gxui.TreeNode {
	return directoryAt(d.subdirs[index])
}

// ItemIndex implements gxui.TreeNodeContainer.
func (d directory) ItemIndex(item gxui.AdapterItem) int {
	path := item.(string)
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	for i, subpath := range d.subdirs {
		subpath += string(filepath.Separator)
		if strings.HasPrefix(path, subpath) {
			return i
		}
	}
	return -1
}

// Item implements gxui.TreeNode.
func (d directory) Item() gxui.AdapterItem {
	return d.path
}

// Create implements gxui.TreeNode.
func (d directory) Create(theme gxui.Theme) gxui.Control {
	_, name := filepath.Split(d.path)
	if name == "" {
		name = d.path
	}
	l := theme.CreateLabel()
	l.SetText(name)
	l.SetColor(directoryColor)
	return l
}

// directoryAdapter is an implementation of the gxui.TreeAdapter interface.
// The AdapterItems returned by this adapter are absolute file path strings.
type directoryAdapter struct {
	gxui.AdapterBase
	directory
}

func (a directoryAdapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}

// Override directory.Create so that the full root is shown, unaltered.
func (a directoryAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	l := theme.CreateLabel()
	l.SetText(a.subdirs[index])
	l.SetColor(directoryColor)
	return l
}

// CreateCustomTheme for white background of windows
func CreateCustomTheme(driver gxui.Driver) gxui.Theme {
	return light.CreateTheme(driver)
}

func showSuccessWindow(driver gxui.Driver) {
	theme := CreateCustomTheme(driver)
	wnd := theme.CreateWindow(200, 50, "Success")
	okButton := theme.CreateButton()
	okButton.SetText("Ok")
	okButton.OnClick(func(gxui.MouseEvent) {
		wnd.Close()
	})
	wndLayout := theme.CreateLinearLayout()
	//wndLayout.SetDirection(gxui.TopToBottom)
	wndLayout.SetHorizontalAlignment(gxui.AlignCenter)
	//wndLayout.SetVerticalAlignment(gxui.AlignMiddle)
	wndLayout.AddChild(okButton)
	//okButton.SetVerticalAlignment(gxui.AlignTop)
	//okButton.SetHorizontalAlignment(gxui.AlignRight)
	wnd.AddChild(wndLayout)
}

// Roots returns the list of drives avaliable on this machine.
func Roots() []string {
	roots := []string{}
	for drive := 'A'; drive <= 'Z'; drive++ {
		path := string(drive) + ":"
		if _, err := os.Stat(path); err == nil {
			roots = append(roots, path)
		}
	}
	return roots
}

func appMain(driver gxui.Driver) {
	//dark background
	theme := flags.CreateTheme(driver)

	//white background
	//theme := CreateCustomTheme(driver)

	window := theme.CreateWindow(800, 600, "Open file...")
	window.SetScale(flags.DefaultScaleFactor)

	// fullpath is the textbox at the top of the window holding the current
	// selection's absolute file path.
	fullpath := theme.CreateTextBox()
	fullpath.SetDesiredWidth(math.MaxSize.W)

	outputDir := theme.CreateTextBox()
	outputDir.SetDesiredWidth(math.MaxSize.W)
	outputDir.SetText("Choose output directory")

	// directories is the Tree of directories on the left of the window.
	// It uses the directoryAdapter to show the entire system's directory
	// hierarchy.
	directories := theme.CreateTree()

	directories.SetAdapter(&directoryAdapter{
		directory: directory{
			path:    "C:",
			subdirs: Roots(), // A: , c: ,  E: ....
		},
	})

	// filesAdapter is the adapter used to show the currently selected directory's
	// content. The adapter has its data changed whenever the selected directory
	// changes.
	filesAdapter := &filesAdapter{}

	// files is the List of files in the selected directory to the right of the
	// window.
	files := theme.CreateList()
	files.SetAdapter(filesAdapter)

	open := theme.CreateButton()
	open.SetText("Open")
	open.OnClick(func(gxui.MouseEvent) {

		filename := fullpath.Text()
		err := parser.Parse(filename, "")
		if err != nil {
			fmt.Println("err : ", err)
		}

		showSuccessWindow(driver)
		fmt.Printf("File '%s' selected!\n", files.Selected())
		//window.Close()
	})

	// If the user hits the enter key while the fullpath control has focus,
	// attempt to select the directory.
	fullpath.OnKeyDown(func(ev gxui.KeyboardEvent) {
		if ev.Key == gxui.KeyEnter || ev.Key == gxui.KeyKpEnter {
			path := fullpath.Text()
			if directories.Select(path) {
				directories.Show(path)
			}
		}
	})

	// When the directory selection changes, update the files list
	directories.OnSelectionChanged(func(item gxui.AdapterItem) {
		dir := item.(string)
		filesAdapter.SetFiles(filesAt(dir))
		fullpath.SetText(dir)
	})

	// ????????
	directories.OnClick(func(gxui.MouseEvent) {
		directories.Selected()

		fullpath.SetText("work")
	})

	// When the file selection changes, update the fullpath text
	files.OnSelectionChanged(func(item gxui.AdapterItem) {
		fullpath.SetText(item.(string))
	})

	// When the user double-clicks a directory in the file list, select it in the
	// directories tree view.
	files.OnDoubleClick(func(gxui.MouseEvent) {
		if path, ok := files.Selected().(string); ok {
			if fi, err := os.Stat(path); err == nil && fi.IsDir() {
				if directories.Select(path) {
					directories.Show(path)
				}
			} else {
				filename := fullpath.Text()
				err := parser.Parse(filename, "")
				if err != nil {
					fmt.Println("err : ", err)
				}
				fmt.Printf("File '%s' selected!\n", path)
				//window.Close()
			}
		}
	})

	//Start with the CWD selected and visible.
	if cwd, err := os.Getwd(); err == nil {
		if directories.Select(cwd) {
			directories.Show(directories.Selected())
		}
	}

	splitter := theme.CreateSplitterLayout()
	splitter.SetOrientation(gxui.Horizontal)
	splitter.AddChild(directories)
	splitter.AddChild(files)

	topLayout := theme.CreateLinearLayout()
	topLayout.SetDirection(gxui.TopToBottom)
	topLayout.AddChild(fullpath)

	topLayout.AddChild(outputDir) // ?????????????

	topLayout.AddChild(splitter)

	btmLayout := theme.CreateLinearLayout()
	btmLayout.SetDirection(gxui.BottomToTop)
	btmLayout.SetHorizontalAlignment(gxui.AlignRight)
	btmLayout.AddChild(open)
	btmLayout.AddChild(topLayout)

	window.AddChild(btmLayout)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func main() {
	gl.StartDriver(appMain)
}
