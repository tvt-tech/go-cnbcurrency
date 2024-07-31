//---------------------------------------------------------------------------

#include <vcl.h>
#include <debugapi.h>
#pragma hdrstop

#include "Unit1.h"
//---------------------------------------------------------------------------
#pragma package(smart_init)
#pragma resource "*.dfm"
TForm1* Form1;
//---------------------------------------------------------------------------
__fastcall TForm1::TForm1(TComponent* Owner) : TForm(Owner) {}
//---------------------------------------------------------------------------

typedef double(__cdecl* GetCurrencyC)(char*, char*);
typedef int(__cdecl* UpdateCurrenciesC)();

void __fastcall TForm1::Button1Click(TObject* Sender)
{
    char bank_id[] = "NBU";
    char code[] = "USD";

    UnicodeString dllPath =
        ExtractFileDir(Application->ExeName) + L"\\currency32.dll";

    HMODULE currency32 = LoadLibrary(dllPath.c_str());
    if (!currency32) {
        OutputDebugString(L"Error loading currency32.dll: ");
        return;
    }

    GetCurrencyC GetCurrency =
        (GetCurrencyC)GetProcAddress(currency32, "GetCurrencyC");
    UpdateCurrenciesC UpdateCurrencies =
        (UpdateCurrenciesC)GetProcAddress(currency32, "UpdateCurrenciesC");

    if (!GetCurrency || !UpdateCurrencies) {
        OutputDebugString(L"Failed to get function address: ");
        FreeLibrary(currency32);
        return;
    }

    try {
        int err = UpdateCurrencies();
        if (err < 0) {
            OutputDebugString(
                L"Error: UpdateCurrencies. Trying to load from cache");
        } else {
			ShowMessage(L"Currencies cache successfully updated");
        }

        double result = GetCurrency(bank_id, code);

		//double result = 1;
        if (result <= 0) {
        } else {
            ShowMessage(L"RESULT: " + FloatToStr(result));
        }
    } catch (const Exception &e) {
        OutputDebugString(L"Exception: ");
    } catch (...) {
        OutputDebugString(L"Unknown error occurred.");
    }
    //FreeLibrary(currency32);       // optional
}
//---------------------------------------------------------------------------
