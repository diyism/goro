package core

import (
	"io"
	"log"

	"git.atonline.com/tristantech/gophp/core/tokenizer"
)

var itemTypeHandler map[tokenizer.ItemType]func(i *tokenizer.Item, c *compileCtx) (runnable, error)
var itemSingleHandler map[rune]func(i *tokenizer.Item, c *compileCtx) (runnable, error)

func init() {
	itemTypeHandler = map[tokenizer.ItemType]func(i *tokenizer.Item, c *compileCtx) (runnable, error){
		tokenizer.T_OPEN_TAG:    nil,
		tokenizer.T_INLINE_HTML: compileInlineHtml,
		tokenizer.T_FUNCTION:    compileFunction,
		tokenizer.T_RETURN:      compileReturn,
	}

	itemSingleHandler = map[rune]func(i *tokenizer.Item, c *compileCtx) (runnable, error){
		'{': compileBase,
		//'}': return compileBase (hidden)
	}
}

// compileIgnore will ignore a given token
func compileIgnore(i *tokenizer.Item, c *compileCtx) (runnable, error) {
	return nil, nil
}

func compileBase(i *tokenizer.Item, c *compileCtx) (runnable, error) {
	var res runnables

	for {
		i, err := c.NextItem()
		if err != nil {
			if err == io.EOF {
				break
			}
			return res, err
		}

		log.Printf("compileBase: %s:%d %s %q", i.Filename, i.Line, i.Type, i.Data)

		// is it a single char item?
		if i.Type == tokenizer.ItemSingleChar {
			ch := []rune(i.Data)[0]

			if ch == '}' {
				// end of block
				return res, nil
			}

			h, ok := itemSingleHandler[ch]
			if !ok {
				return nil, i.Unexpected()
			}
			if h == nil {
				// ignore this tag
				continue
			}

			r, err := h(i, c)
			if err != nil {
				return res, err
			}

			if r != nil {
				res = append(res, r)
			}
		} else {
			// is it a token?
			h, ok := itemTypeHandler[i.Type]
			if !ok {
				return nil, i.Unexpected()
			}
			if h == nil {
				// ignore this tag
				continue
			}

			r, err := h(i, c)
			if err != nil {
				return res, err
			}

			if r != nil {
				res = append(res, r)
			}
		}

		// check for ';'
		i, err = c.NextItem()
		if err != nil {
			if err == io.EOF {
				break
			}
			return res, err
		}

		if !i.IsSingle(';') {
			// expecting a ';' after a var
			return nil, i.Unexpected()
		}
	}

	return res, nil
}

func compileReturn(i *tokenizer.Item, c *compileCtx) (runnable, error) {
	i, err := c.NextItem()
	c.backup()
	if err != nil {
		return &runReturn{}, err
	}

	if i.IsSingle(';') {
		return &runReturn{}, err
	}

	v, err := compileExpr(c)
	if err != nil {
		return nil, err
	}

	return &runReturn{v}, nil
}

type runReturn struct {
	v runnable
}

func (r *runReturn) run(ctx Context) (*ZVal, error) {
	return r.v.run(ctx) // TODO
}