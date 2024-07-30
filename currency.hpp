#include <Windows.h>
#include <string>
#include <stdexcept>


class CurrencyDLLWrapper {
public:
    typedef int (*UpdateCurrenciesC)();
	typedef double (*GetCurrencyC)(const char*, const char*);

    CurrencyDLLWrapper(const std::wstring& dllPath) {
        // Load the DLL
        dllHandle = LoadLibrary(dllPath.c_str());
        if (!dllHandle) {
            throw std::runtime_error("Error loading DLL");
        }

        // Get the function addresses
        updateCurrencies = (UpdateCurrenciesC)GetProcAddress(dllHandle, "UpdateCurrenciesC");
        getCurrency = (GetCurrencyC)GetProcAddress(dllHandle, "GetCurrencyC");

        if (!updateCurrencies || !getCurrency) {
            FreeLibrary(dllHandle);
            throw std::runtime_error("Failed to get function address");
        }
    }

    ~CurrencyDLLWrapper() {
        // Unload the DLL
        if (dllHandle) {
            //FreeLibrary(dllHandle);
        }
    }

    int UpdateCurrencies() {
        if (updateCurrencies) {
            return updateCurrencies();
        }
        throw std::runtime_error("UpdateCurrencies function not loaded");
    }

    double GetCurrency(const std::string& bankId, const std::string& code) {
        if (getCurrency) {
            return getCurrency(bankId.c_str(), code.c_str());
        }
        throw std::runtime_error("GetCurrency function not loaded");
    }

private:
    HMODULE dllHandle;
    UpdateCurrenciesC updateCurrencies;
    GetCurrencyC getCurrency;
};


/* USAGE EXAMPLE:
	CurrencyDLLWrapper currencyWrapper(std::wstring(dllPath.c_str(), dllPath.Length()));
	int err = currencyWrapper.UpdateCurrencies();
	ShowMessage(IntToStr(err));

	ShowMessage(FloatToStr(currencyWrapper.GetCurrency(bank_id, code)));
*/