package vida

const v rune = 'V'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

func (c *Compiler) appendHeader() {
	c.module.Code = append(c.module.Code, byte(v))
	c.module.Code = append(c.module.Code, byte(i))
	c.module.Code = append(c.module.Code, byte(d))
	c.module.Code = append(c.module.Code, byte(a))
	c.module.Code = append(c.module.Code, byte(major))
	c.module.Code = append(c.module.Code, byte(minor))
	c.module.Code = append(c.module.Code, byte(patch))
	c.module.Code = append(c.module.Code, byte(inception))
}

func (c *Compiler) appendStop() {
	c.module.Code = append(c.module.Code, end)
}
