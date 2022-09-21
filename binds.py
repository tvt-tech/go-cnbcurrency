from ctypes import *
from ctypes import cdll
from pathlib import Path


class GoString(Structure):
    _fields_ = [
        ("p", c_char_p),
        ("n", c_int)]


lib_path = Path("cnbcurrency.dll").absolute().__str__()
lib = cdll.LoadLibrary(lib_path)

lib.GetUsd.restype = c_double
lib.GetEur.restype = c_double


def get_currency(code: str):
    char = c_char_p(code.encode('utf-8'))
    gostring = GoString(char, len(code))
    lib.GetCurrency.restype = c_double
    return lib.GetCurrency(gostring)


if __name__ == '__main__':
    get_currency("ils")
    lib.GetUsd()
