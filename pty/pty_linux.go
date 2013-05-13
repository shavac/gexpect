package pty
import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	sys_TIOCGPTN   = 0x80045430
	sys_TIOCSPTLCK = 0x40045431
	sys_TIOCSWINSZ = syscall.TIOCSWINSZ
)

// Opens a pty and its corresponding tty.
func Open() (pty, tty *os.File, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	err = unlockpt(p)
	if err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}


func ptsname(f *os.File) (string, error) {
	var n int
	err := ioctl(f.Fd(), sys_TIOCGPTN, &n)
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(n), nil
}

func unlockpt(f *os.File) error {
	var u int
	return ioctl(f.Fd(), sys_TIOCSPTLCK, &u)
}

func ioctl(fd uintptr, cmd uintptr, data *int) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		cmd,
		uintptr(unsafe.Pointer(data)),
	)
	if e != 0 {
		return syscall.ENOTTY
	}
	return nil
}
type ttySize struct {
	Rows uint16
	Cols uint16
	Xpixel uint16
	Ypixel uint16
}
func SetWinSize(f *os.File, cols uint16, rows uint16) error {
	_, _, e := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&ttySize{rows, cols, 0, 0})),
		0, 0, 0,
	)
	if e != 0 {
		return syscall.ENOTTY
	}
	return nil
}
func GetWinSize(f *os.File) (width, height int, err error) {
    var dimensions ttySize
    if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(f.Fd()), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
        return -1, -1, err
    }
    return int(dimensions.Cols), int(dimensions.Rows), nil
}
type Terminal struct {
	Pty *os.File
    Tty *os.File
}

func (t *Terminal)Write(b []byte) (n int, err error){
	return t.Pty.Write(b)
}

func (t *Terminal)Read(b []byte) (n int, err error) {
	return t.Pty.Read(b)
}

func (t *Terminal)SetWinSize(x, y int) error {
	SetWinSize(t.Pty, uint16(x), uint16(y))
	return nil
}
func (t *Terminal)GetWinSize() (x, y int, err error) {
	return 0,0, nil
}

func NewTerminal() (term *Terminal, err error) {
	pty, tty, err := Open()
	term = &Terminal{ Pty: pty, Tty: tty}
	return
}
