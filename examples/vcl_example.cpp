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
	char bank_id[] = "NBU";
	char code[] = "USD";

	UnicodeString dllPath = ExtractFileDir(Application->ExeName) + L"\\currency32.dll";

	CurrencyDLLWrapper currencyWrapper(std::wstring(dllPath.c_str(), dllPath.Length()));
	int err = currencyWrapper.UpdateCurrencies();
	ShowMessage(L"Update result: " + IntToStr(err));
	ShowMessage("USD rate: " + FloatToStr(currencyWrapper.GetCurrency(bank_id, code)));
}