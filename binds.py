# from ctypes import *
# from pathlib import Path


# lib_path = Path("cnbcurrency.dll").absolute().__str__()
# lib = cdll.LoadLibrary(lib_path)


# cnblib.GetUsd.restype = c_double
# cnblib.GetEur.restype = c_double
# cnblib.GetCurrency.argtypes = [c_char_p]
# cnblib.GetCurrency.restype = c_double


# if __name__ == '__main__':
#     cnblib.GetUsd()
#     cnblib.GetCurrency("ils".encode('utf-8'))


from ctypes import *
from pathlib import Path


lib_path = Path("currency64.dll").absolute().__str__()
lib = cdll.LoadLibrary(lib_path)
handle = lib._handle

lib.GetCurrencyC.argtypes = [c_char_p, c_char_p]
lib.GetCurrencyC.restype = c_double




if __name__ == '__main__':
    print(lib.GetCurrencyC("NBU".encode('utf-8'), "USD".encode('utf-8')))
    del lib
    while True:
        ...