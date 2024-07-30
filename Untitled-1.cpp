//---------------------------------------------------------------------------

#include <vcl.h>
#pragma hdrstop

#include "Unit1.h"
#include "currency32.h"
//#include <IdHTTP.hpp>
//#include <IdSSLOpenSSL.hpp>
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

void __fastcall TForm1::FormCreate(TObject *Sender)
{
	TTabSheet *firstTabSheet = PageControl1->Pages[0];
	firstTabSheet->Caption = "First caption";
    TabSheet2->Caption = "Second caption";
}
//---------------------------------------------------------------------------

void __fastcall TForm1::Button1Click(TObject *Sender)
{

	WCHAR exePath[MAX_PATH];

	if (GetModuleFileName(NULL, exePath, MAX_PATH) == 0)
	{
		ShowMessage("Error getting executable path: " + GetLastError());
		return;
	}

	UnicodeString exeDir = ExtractFileDir(exePath) + L"\\";
	HMODULE currency32 = LoadLibrary((exeDir + L"currency32.dll").c_str());

//	WCHAR currency32Path[MAX_PATH];
//    ShowMessage((exeDir + L"currency32.dll").c_str());
//
//	if (GetModuleFileName(currency32, currency32Path, MAX_PATH) == 0)
//	{
//		ShowMessage("Error getting currency32 path: " + GetLastError());
//		return;
//	}
//
//	ShowMessage("currency32.dll loaded from: " + UnicodeString(currency32Path));
//
//	if (!currency32)
//	{
//		ShowMessage("Error loading OpenSSL libraries: " + GetLastError());
//		return;
//	}


	GetCurrencyC GetCurrency = (GetCurrencyC)GetProcAddress(currency32, "GetCurrencyC");
	if (GetCurrency == NULL) {
		ShowMessage("Failed to get function address");
		FreeLibrary(currency32);
		return;
	}

}


////	 Buffer to store the full path of the executable
//	WCHAR exePath[MAX_PATH];
//	// Get the full path of the executable
//	if (GetModuleFileName(NULL, exePath, MAX_PATH) == 0)
//    {
//        ShowMessage("Error getting executable path: " + GetLastError());
//        return;
//	}
//
//	UnicodeString exeDir = ExtractFileDir(exePath) + L"\\";
//
//	HMODULE libeay32 = LoadLibrary((exeDir + "libeay32.dll").c_str());
//	HMODULE ssleay32 = LoadLibrary((exeDir + "ssleay32.dll").c_str());
//
//	if (!libeay32 || !ssleay32)
//	{
//		ShowMessage("Error loading OpenSSL libraries: " + GetLastError());
//		return;
//	}
//
//    // Get the path of the loaded libraries
//    WCHAR libeay32Path[MAX_PATH];
//    WCHAR ssleay32Path[MAX_PATH];
//    if (GetModuleFileName(libeay32, libeay32Path, MAX_PATH) == 0)
//    {
//		ShowMessage("Error getting libeay32 path: " + GetLastError());
//		return;
//	}
//	if (GetModuleFileName(ssleay32, ssleay32Path, MAX_PATH) == 0)
//	{
//        ShowMessage("Error getting ssleay32 path: " + GetLastError());
//        return;
//    }
//    // Display the paths
//    ShowMessage("libeay32.dll loaded from: " + UnicodeString(libeay32Path));
//	ShowMessage("ssleay32.dll loaded from: " + UnicodeString(ssleay32Path));
//
//	 // Create components
//    TIdHTTP *IdHTTP = new TIdHTTP(NULL);
//    TIdSSLIOHandlerSocketOpenSSL *IdSSLIOHandler = new TIdSSLIOHandlerSocketOpenSSL(NULL);
//    try
//    {
//        // Set the IOHandler
//        IdHTTP->IOHandler = IdSSLIOHandler;
//        // Perform the GET request
//        UnicodeString url = L"https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange";
//        UnicodeString response = IdHTTP->Get(url);
//        // Display the response
//        ShowMessage(response);
//    }
//    __finally
//	{
//        // Clean up
//        delete IdHTTP;
//        delete IdSSLIOHandler;
//	}
//}
//---------------------------------------------------------------------------
