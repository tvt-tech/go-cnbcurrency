#include <windows.h>

// Declare functions to import from Go DLL
extern double __stdcall GetCurrencyC(char* bankID, char* code);
extern int __stdcall UpdateCurrenciesC(void);

// Export functions with __stdcall calling convention
__declspec(dllexport) double __stdcall _GetCurrencyC(char* bankID, char* code) {
    return GetCurrencyC(bankID, code);
}

__declspec(dllexport) int __stdcall _UpdateCurrenciesC(void) {
    return UpdateCurrenciesC();
}