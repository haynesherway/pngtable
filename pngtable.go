package pngtable

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "io/ioutil"
    "os"
    
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
)

const (
    DEFAULT_COLWIDTH = 50
    DEFAULT_ROWHEIGHT = 25
    
    DEFAULT_WIDTH = 500
    DEFAULT_HEIGHT = 250
)

var (
    FONT_FILE = os.Getenv("GOPATH") + "/src/github.com/golang/freetype/testdata/luxisr.ttf"
    font *truetype.Font
)

type Table struct {
    ColCount int // Number of columns in table
    RowCount int // Number of rows in table
    Rows []*Row
    Header *Row
    Image *image.RGBA
    FileName string
    Options *TableOptions
}

type Row struct {
    Values []string
    Background color.Color
    Color color.Color
    bounds image.Rectangle
}

type TableOptions struct {
    Width int
    Height int
    ColWidth int
    ColWidths []int
    RowHeight int
    Color color.Color
    Background color.Color
    BorderWidth int
    BorderColor color.Color
    FontSize int
}

func New() *Table {
    options := DefaultOptions()
    rows := []*Row{}
    return &Table{Header: &Row{}, Rows: rows, Options: options}
}

func DefaultOptions() *TableOptions {
    options := &TableOptions{
        Width: 0,
        Height: 0,
        RowHeight: 25,
        ColWidth: 50,
        Color: color.White,
        Background: color.Black,
        BorderWidth: 1,
        BorderColor: color.White,
        FontSize: 25,
    }
    return options
}

func (to *TableOptions) SetWidth(w int) {
    to.Width = w
    to.ColWidth = 0
    return
}

func (to *TableOptions) SetHeight(h int) {
    to.Height = h
    to.RowHeight = 0
    return
}

func (to *TableOptions) SetColWidth(w int) {
    to.ColWidth = w
    to.Width = 0
    return
}

func (to *TableOptions) SetColWidths(w []int) {
    to.ColWidths = w
    to.Width = 0
    return
}

func (to *TableOptions) SetRowHeight(h int) {
    to.RowHeight = h
    to.Height = 0
    return
}

func (to *TableOptions) SetBorder(b int) {
    to.BorderWidth = b
    return
}

func (to *TableOptions) SetFontSize(s int) {
    to.FontSize = s
    return
}

func (r *Row) SetBackground(c color.Color) (*Row) {
    r.Background = c
    return r
}

func (r *Row) SetColor(c color.Color) (*Row) {
    r.Color = c
    return r
}

func (t *Table) SetHeaders(h []string) (*Row) {
    var row Row
    row.Values = h
    t.ColCount = len(h)
    t.Header = &row
    t.RowCount++
    return &row
}

func (t *Table) AddRow(r []string) (*Row) {
    if len(r) > t.ColCount {
        t.ColCount = len(r)
    }
    var row Row
    row.Values = r
    row.Background = color.Black
    row.Color = color.White
    t.Rows = append(t.Rows, &row)
    t.RowCount++
    return &row
}

func (t *Table) Draw() {
    t.drawTable()
    
    return
    
}

func (t *Table) drawTable() {
    //Options
    fg := image.NewUniform(t.Options.Color)
    bg := image.NewUniform(t.Options.Background)
    colWidths := t.getColWidths()
    rowHeight := t.getRowHeight()
    borderWidth := t.Options.BorderWidth
    
    width := t.getWidth()
    height := t.getHeight()
    boxWidth := width - borderWidth*2
    boxHeight := height - borderWidth*2
    
    fmt.Println(boxWidth, boxHeight)
    
    // Start Drawing Table 
    t.Image = image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(t.Image, t.Image.Bounds(), bg, image.ZP, draw.Src)
    
    var x,y,x1,y1,bw,bx,by int
     //Draw Rows
    x,y = 0,0
    rows := append([]*Row{t.Header}, t.Rows...)
    for _, row := range rows {
        top := y+borderWidth
        y += rowHeight
        row.bounds = image.Rect(0, top, width, y)
        
        //Set background
        if row.Background != t.Options.Background && row.Background != nil {
            bg = image.NewUniform(row.Background)
            draw.Draw(t.Image, row.bounds, bg, image.ZP, draw.Src)
        }
        
        for x1 = 0 ; x1 <= width; x1++ {
            for bw = 0; bw < borderWidth; bw++ {
                t.Image.Set(x1, y+bw, t.Options.Color)
            }
        }
    }
    
    //Draw Border
    // Top
    for bx,by=0,0; bx <= width; bx++ {
        for bw = 0; bw < borderWidth; bw++ {
            t.Image.Set(bx, by+bw, t.Options.BorderColor)
        }
    }
    // Right
    for bx,by=width-1,0; by <= height; by++ {
        for bw = 0; bw < borderWidth; bw++ {
            t.Image.Set(bx-bw, by, t.Options.BorderColor)
        }
    }
    // Bottom 
    by = height - 1
    for bx,by=0,height-1; bx <= width; bx++ {
        for bw = 0; bw < borderWidth; bw++ {
            t.Image.Set(bx, by-bw, t.Options.BorderColor)
        }
    }
    // Left
    for bx,by=0,0; by <= height; by++ {
        for bw = 0; bw < borderWidth; bw++ {
            t.Image.Set(bx+bw, by, t.Options.BorderColor)
        }
    }
    
    /*for i, _ := range t.Header.Values {
        x += colWidths[i]
        for y1 = 0; y1 < height; y1++ {
            for bw = 0; bw < borderWidth; bw++ {
                t.Image.Set(x+bw, y1, t.Options.Color)
            }
        }
    }*/
    
    // Set up font
        c := freetype.NewContext()
    	c.SetDPI(72)
    	c.SetFont(font)
    	c.SetFontSize(float64(t.Options.FontSize))
    	c.SetClip(t.Image.Bounds())
    	c.SetDst(t.Image)
    	c.SetSrc(fg)
    	opts := truetype.Options{}
        opts.Size = float64(t.Options.FontSize)
        face := truetype.NewFace(font, &opts)
    
    
    
    //for  := range t.Header.Values {
    // Draw Header Row 
   /*y += rowHeight
    for x1 = 0 ; x1 <= width; x1++ {
        for bw = 0; bw < borderWidth; bw++ {
            t.Image.Set(x1, y+bw, t.Options.Color)
        }
    }*/
    x,y = 0,0
    y += rowHeight
    for i, h := range t.Header.Values {
        x += colWidths[i]
        for y1 = 0; y1 < height; y1++ {
            for bw = 0; bw < borderWidth; bw++ {
                t.Image.Set(x+bw, y1, t.Options.Color)
            }
        }
        
        //Print
        hWidth := 0
        for _, l := range(h) {
            awidth, ok := face.GlyphAdvance(rune(l))
            if ok != true {
                return
            }
            fmt.Println(awidth)
            hWidth += int(float64(awidth) / 64)
        }
        fmt.Println(colWidths[i])
        pt := freetype.Pt((x-colWidths[i])+((colWidths[i]/2)-hWidth/2), y-(rowHeight-t.Options.FontSize)+2)
        fmt.Println(pt)
        c.DrawString(string(h), pt)
    }
    
    for _, r := range t.Rows {
        y += rowHeight
        /*for x1 = 0; x1 <= width; x1++ {
            bw = 0
            for ; bw < borderWidth; bw++ {
                t.Image.Set(x1, y+bw, t.Options.Color)
            }
        }
        */
        x := 0
        for i, v := range r.Values {
            x += colWidths[i]
            //Print
            vWidth := 0
            for _, l := range(v) {
                awidth, ok := face.GlyphAdvance(rune(l))
                if ok != true {
                    return
                }
                vWidth += int(float64(awidth) / 64)
            }
            
            pt := freetype.Pt((x-colWidths[i])+((colWidths[i]/2)-vWidth/2), y-(rowHeight-t.Options.FontSize))
            c.DrawString(string(v), pt)
        }
    }
    
    
    //t.draw()
}

func (t *Table) getColWidths() []int {
    if len(t.Options.ColWidths) != 0 {
        return t.Options.ColWidths
    }
    if t.Options.ColWidth != 0 {
        cw := []int{}
        for c := 0; c < t.ColCount; c++ {
            cw = append(cw, t.Options.ColWidth)
        }
        t.Options.ColWidths = cw
        return cw
    } else if t.Options.Width != 0 {
        cw := []int{}
        for c := 0; c < t.ColCount; c++ {
            cw = append(cw, (t.Options.Width - t.Options.BorderWidth) / t.ColCount)
        }
        t.Options.ColWidths = cw
        return cw
    }
    
    cw := []int{}
        for c := 0; c < t.ColCount; c++ {
            cw = append(cw, DEFAULT_COLWIDTH)
        }
        t.Options.ColWidths = cw
        return cw
}

func (t *Table) getRowHeight() int {
    if t.Options.RowHeight != 0 {
        return t.Options.RowHeight
    } else if t.Options.Height != 0 {
        return (t.Options.Height - t.Options.BorderWidth) / t.RowCount
    }
    
    return DEFAULT_ROWHEIGHT
}

func (t *Table) getHeight() int {
    if t.Options.Height != 0 {
        return t.Options.Height
    } else if t.Options.RowHeight != 0 {
        return (t.Options.RowHeight * t.RowCount) + t.Options.BorderWidth
    }
    
    return DEFAULT_HEIGHT
}

func (t *Table) getWidth() int {
    if len(t.Options.ColWidths) != 0 {
        w := 0
        for _, cw := range t.Options.ColWidths {
            w += cw
        }
        w += t.Options.BorderWidth
        t.Options.Width = w
        return w
    } else if t.Options.Width != 0 {
        return t.Options.Width 
    } else if t.Options.ColWidth != 0 {
        return t.Options.ColWidth * t.ColCount + t.Options.BorderWidth
    }
    
    return DEFAULT_WIDTH
}

func (t *Table) draw() {
    f, err := os.Create("draw.png")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    err = png.Encode(f, t.Image)
    if err != nil {
        panic(err)
    }
    
    return 
}

func init() {
    // Read the font data.
	fontBytes, err := ioutil.ReadFile(FONT_FILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
}