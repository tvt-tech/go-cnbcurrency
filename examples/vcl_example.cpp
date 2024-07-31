//---------------------------------------------------------------------------
#include <vcl.h>
#include "currency.hpp"

#pragma hdrstop

#include "Unit1.h"
//---------------------------------------------------------------------------
#pragma package(smart_init)
#pragma resource "*.dfm"
TForm1 *Form1;
//---------------------------------------------------------------------------
__fastcall TForm1::TForm1(TComponent* Owner)
	: TForm(Owner)
{

}
//---------------------------------------------------------------------------

void __fastcall TForm1::Button1Click(TObject *Sender)
{
	int cacheUpdateErrorCode;
	long long int cacheUpdateTime;
	double usdRate;

	char bankId[] = "NBU";
	char code[] = "USD";

	UnicodeString dllPath = ExtractFileDir(Application->ExeName) + L"\\currency32.dll";

	std::basic_string dllPathStr = std::wstring(dllPath.c_str(), dllPath.Length());

	CurrencyDLLWrapper currencyWrapper(dllPathStr);
	cacheUpdateErrorCode = currencyWrapper.UpdateCurrencies();
	cacheUpdateTime = currencyWrapper.GetCacheUpdateTime();
	usdRate =currencyWrapper.GetCurrency(bankId, code);
	ShowMessage(L"Update result: " + IntToStr(cacheUpdateErrorCode));
	ShowMessage(L"Cache update time (Epoch): " + IntToStr(cacheUpdateTime));
	ShowMessage(L"USD rate: " + FloatToStr(usdRate));
}