package symc

import (
	"reflect"
	"testing"
)

func TestModules(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  string
	}{
		{
			"test module 1 bootpack.c",
			`
# 1 "examples/Kitax/bootpack.c"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 366 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./bootpack.h" 1





typedef struct {
    char cyls, leds, vmode, reserve;
    short scrnx, scrny;
    char *vram;
} BOOTINFO;
# 2 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./dsctbl/dsctbl.h" 1



void init_gdtidt(void);
# 3 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./fifo/fifo.h" 1



typedef struct {
    unsigned char *buf;
    int w;
    int r;
    int size;
    int free_num;
    int flags;
} FIFO8;

void fifo8_init(FIFO8 *fifo, int size, unsigned char *buf);
int fifo8_put(FIFO8 *fifo, unsigned char data);
int fifo8_get(FIFO8 *fifo);
int fifo8_data_count(FIFO8 *fifo);
# 4 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./graphic/graphic.h" 1
# 21 "examples/Kitax/./graphic/graphic.h"
void init_palette(void);
void init_screen(char *vram, int xsize, int ysize);
void init_mouse_cursor8(char *mouse, char bc);
void putblock8_8(char vram[], short vxsize, int pxsize, int pysize, int px,
                 int py, char buf[], int bxsize);
void boxfill8(char *vram, int xsize, unsigned char color, int x_s, int y_s,
              int x_e, int y_e);
void putfonts8_asc(char *vram, short xsize, int x, int y, char color, char s[]);
# 5 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./int/int.h" 1
# 21 "examples/Kitax/./int/int.h"
extern FIFO8 keyfifo;
extern FIFO8 mousefifo;

void init_pic(void);
# 6 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./lib/lib.h" 1



void mysprintf(char *str, char *fmt, ...);
# 7 "examples/Kitax/bootpack.c" 2
# 1 "examples/Kitax/./naskfunc/naskfunc.h" 1




void io_hlt(void);
void io_stihlt(void);
int io_load_eflags(void);
void io_store_eflags(int eflags);
void io_out8(int port, char data);
int io_in8(int port);
void io_cli(void);
void io_sti(void);
void load_gdtr(int limit, int addr);
void load_idtr(int limit, int addr);

void asm_inthandler21(void);
void asm_inthandler27(void);
void asm_inthandler2c(void);
# 8 "examples/Kitax/bootpack.c" 2
# 21 "examples/Kitax/bootpack.c"
static void wait_KBC_sendready(void);
static void init_keyboard(void);
static void enable_mouse(void);

void HariMain(void) {
    BOOTINFO *binfo = (BOOTINFO *)(0x00000ff0);
    char s[64];
    int mx, my;
    char mcursor[16 * 16];
    unsigned char data;
    unsigned char keybuf[(32)];
    unsigned char mousebuf[(128)];

    init_gdtidt();
    init_pic();

    io_sti();

    fifo8_init(&keyfifo, (32), keybuf);
    fifo8_init(&mousefifo, (128), mousebuf);
    io_out8((0x0021), 0xf9);
    io_out8((0x00a1), 0xef);

    init_keyboard();

    init_palette();
    init_screen(binfo->vram, binfo->scrnx, binfo->scrny);
    mx = (binfo->scrnx - 16) / 2;
    my = (binfo->scrny - 28 - 16) / 2;
    init_mouse_cursor8(mcursor, (14));
    putblock8_8(binfo->vram, binfo->scrnx, 16, 16, mx, my, mcursor, 16);

    mysprintf(s, "(%d, %d)", mx, my);
    putfonts8_asc(binfo->vram, binfo->scrnx, 0, 0, (7), s);

    enable_mouse();

    for (;;) {
        io_cli();
        if (fifo8_data_count(&keyfifo) + fifo8_data_count(&mousefifo) == 0) {
            io_stihlt();
        } else {
            if (fifo8_data_count(&keyfifo) != 0) {
                data = fifo8_get(&keyfifo);
                io_sti();
                mysprintf(s, "%x", data);
                boxfill8(binfo->vram, binfo->scrnx, (14), 0, 16, 15, 31);
                putfonts8_asc(binfo->vram, binfo->scrnx, 0, 16, (7), s);
            } else if (fifo8_data_count(&mousefifo) != 0) {
                data = fifo8_get(&mousefifo);
                io_sti();
                mysprintf(s, "%x", data);
                boxfill8(binfo->vram, binfo->scrnx, (14), 32, 16, 47,
                         31);
                putfonts8_asc(binfo->vram, binfo->scrnx, 32, 16, (7),
                              s);
            }
        }
    }
}

static void wait_KBC_sendready(void) {

    while (1) {
        if ((io_in8((0x0064)) & (0x0002)) == 0x0000) {
            break;
        }
    }
}

static void init_keyboard(void) {


    wait_KBC_sendready();
    io_out8((0x0064), (0x60));
    wait_KBC_sendready();
    io_out8((0x0060), (0x47));
}

static void enable_mouse(void) {

    wait_KBC_sendready();
    io_out8((0x0064), (0xd4));
    wait_KBC_sendready();
    io_out8((0x0060), (0xf4));


}
`,
			`PROTOTYPE init_gdtidt
PROTOTYPE fifo8_init
PROTOTYPE fifo8_put
PROTOTYPE fifo8_get
PROTOTYPE fifo8_data_count
PROTOTYPE init_palette
PROTOTYPE init_screen
PROTOTYPE init_mouse_cursor8
PROTOTYPE putblock8_8
PROTOTYPE boxfill8
PROTOTYPE putfonts8_asc
DECLARE keyfifo
DECLARE mousefifo
PROTOTYPE init_pic
PROTOTYPE mysprintf
PROTOTYPE io_hlt
PROTOTYPE io_stihlt
PROTOTYPE io_load_eflags
PROTOTYPE io_store_eflags
PROTOTYPE io_out8
PROTOTYPE io_in8
PROTOTYPE io_cli
PROTOTYPE io_sti
PROTOTYPE load_gdtr
PROTOTYPE load_idtr
PROTOTYPE asm_inthandler21
PROTOTYPE asm_inthandler27
PROTOTYPE asm_inthandler2c
PROTOTYPE wait_KBC_sendready
PROTOTYPE init_keyboard
PROTOTYPE enable_mouse
FUNC HariMain() {
    DEFINITION binfo
    DEFINITION s
    DEFINITION mx
    DEFINITION my
    DEFINITION mcursor
    DEFINITION data
    DEFINITION keybuf
    DEFINITION mousebuf
    init_gdtidt()
    init_pic()
    io_sti()
    fifo8_init(keyfifo, keybuf)
    fifo8_init(mousefifo, mousebuf)
    io_out8()
    io_out8()
    init_keyboard()
    init_palette()
    init_screen(binfo->vram, binfo->scrnx, binfo->scrny)
    ASSIGNE mx
    binfo->scrnx
    ASSIGNE my
    binfo->scrny
    init_mouse_cursor8(mcursor)
    putblock8_8(binfo->vram, binfo->scrnx, mx, my, mcursor)
    mysprintf(s, mx, my)
    putfonts8_asc(binfo->vram, binfo->scrnx, s)
    enable_mouse()
    io_cli()
    fifo8_data_count(keyfifo)
    fifo8_data_count(mousefifo)
    io_stihlt()
    fifo8_data_count(keyfifo)
    ASSIGNE data
    fifo8_get(keyfifo)
    io_sti()
    mysprintf(s, data)
    boxfill8(binfo->vram, binfo->scrnx)
    putfonts8_asc(binfo->vram, binfo->scrnx, s)
    fifo8_data_count(mousefifo)
    ASSIGNE data
    fifo8_get(mousefifo)
    io_sti()
    mysprintf(s, data)
    boxfill8(binfo->vram, binfo->scrnx)
    putfonts8_asc(binfo->vram, binfo->scrnx, s)
}

FUNC wait_KBC_sendready() {
    io_in8()
}

FUNC init_keyboard() {
    wait_KBC_sendready()
    io_out8()
    wait_KBC_sendready()
    io_out8()
}

FUNC enable_mouse() {
    wait_KBC_sendready()
    io_out8()
    wait_KBC_sendready()
    io_out8()
}

`,
		},
		{
			"test module 2 mystdio.c",
			`
# 1 "examples/Kitax/lib/mystdio.c"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 366 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "examples/Kitax/lib/mystdio.c" 2
# 1 "/Library/Developer/CommandLineTools/usr/lib/clang/12.0.0/include/stdarg.h" 1 3
# 14 "/Library/Developer/CommandLineTools/usr/lib/clang/12.0.0/include/stdarg.h" 3
typedef __builtin_va_list va_list;
# 32 "/Library/Developer/CommandLineTools/usr/lib/clang/12.0.0/include/stdarg.h" 3
typedef __builtin_va_list __gnuc_va_list;
# 2 "examples/Kitax/lib/mystdio.c" 2


int dec2asc(char *str, int dec) {
    int len = 0, len_buf;
    int buf[10];
    while (1) {
        buf[len++] = dec % 10;
        if (dec < 10)
            break;
        dec /= 10;
    }
}
`,
			`FUNC dec2asc(DEFINITION str, DEFINITION dec) {
    DEFINITION len
    DEFINITION len_buf
    DEFINITION buf
    ASSIGNE buf
    len
    dec
    dec
    ASSIGNE dec
}

`,
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse().PrettyString()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("\ngot=\n%v\nexpect=\n%v\n", got, tt.expect)
		}
	}
}
