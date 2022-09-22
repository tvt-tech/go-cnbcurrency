from ctypes import *
from pathlib import Path


lib_path = Path("cnbcurrency.dll").absolute().__str__()
lib = cdll.LoadLibrary(lib_path)


cnblib.GetUsd.restype = c_double
cnblib.GetEur.restype = c_double
cnblib.GetCurrency.argtypes = [c_char_p]
cnblib.GetCurrency.restype = c_double


if __name__ == '__main__':
    cnblib.GetUsd()
    cnblib.GetCurrency("ils".encode('utf-8'))
