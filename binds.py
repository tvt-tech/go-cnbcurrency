from ctypes import cdll, c_char_p, c_double, c_int
from pathlib import Path


lib_path = Path("currency32.dll").absolute().__str__()
lib = cdll.LoadLibrary(lib_path)
handle = lib._handle

lib.GetCurrencyC.argtypes = [c_char_p, c_char_p]
lib.GetCurrencyC.restype = c_double
lib.UpdateCurrenciesC.restype = c_int


def main():
    if lib.UpdateCurrenciesC() >= 0:
        print(lib.GetCurrencyC("NBU".encode('utf-8'), "USD".encode('utf-8')))


if __name__ == '__main__':
    main()