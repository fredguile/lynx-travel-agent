package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnescapeGWTErrorString(t *testing.T) {
	// Test the example provided by the user
	testInput := `"Type \x27com.lynxtraveltech.client.shared.model.FileSearchCriteria\x27 was not assignable to \x27com.google.gwt.user.client.rpc.IsSerializable\x27 and did not have a custom field serializer. For security purposes"`

	expected := `Type 'com.lynxtraveltech.client.shared.model.FileSearchCriteria' was not assignable to 'com.google.gwt.user.client.rpc.IsSerializable' and did not have a custom field serializer. For security purposes`

	result := unescapeGWTErrorString(testInput)

	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
		fmt.Printf("Expected length: %d, Got length: %d\n", len(expected), len(result))
	}
}

func TestParseResponseError(t *testing.T) {
	// Test case with the example error provided
	testError := `//EX[2,1,["com.google.gwt.user.client.rpc.IncompatibleRemoteServiceException/3936916533","Type \x27com.lynxtraveltech.client.shared.model.FileSearchCriteria\x27 was not assignable to \x27com.google.gwt.user.client.rpc.IsSerializable\x27 and did not have a custom field serializer. For security purposes, this type will not be deserialized."],0,7]`

	expectedError := `Type 'com.lynxtraveltech.client.shared.model.FileSearchCriteria' was not assignable to 'com.google.gwt.user.client.rpc.IsSerializable' and did not have a custom field serializer. For security purposes, this type will not be deserialized.`

	result, err := ParseResponseError(testError)
	if err != nil {
		t.Fatalf("ParseResponseError failed: %v", err)
	}

	actual := strings.TrimSpace(result)
	expected := strings.TrimSpace(expectedError)

	if actual != expected {
		t.Errorf("Expected error message: %s\nGot: %s", expected, actual)
		fmt.Printf("Expected length: %d, Got length: %d\n", len(expected), len(actual))
	}
}

func TestParseFileSummaryResponse(t *testing.T) {
	responseBody := `//OK[29,28,27,0,26,'$2s7',0,25,'BAkC',0,0,24,0,23,'BgsOD',22,0,21,'sw',6,20,19,5,18,0,0,0,17,'UxU1',0,16,15,0,0,0,3,6,0,0,14,13,12,11,5,0,0,0,6,0,0,0,1,0,5,10,0,9,1,3,0,0,0,5,4,1,3,863,1,0,1,8,0,0,200,0,0,1,0,5,4,1,3,0,7,0,0,0,0,0,0,0,6,0,5,4,1,3,0,2,0,2,1,["com.lynxtraveltech.client.shared.model.FileSummary/2582189904","295.0000","java.util.ArrayList/4159755760","java.lang.String/2004016611","AUD 295.00","","AUD","0","com.lynxtraveltech.client.shared.model.TransactionSummary/2232640144","-     ","SYDNEY, NSW","0.00%","Claire-Marie Bray","04 Jul 2025 07:13 PM","2 adults","1 X Premium City View Room (1-2 Pax, Max 4)","Premium City View Room (1-2 Pax, Max 4)","2 x Queen;  ","On Request","A","HOLIDAY INN SYDNEY POTTS POINT","java.lang.Long/4227064769","01/02 Nov 2025","16476987-1","25052025","FTAUB252039","BRAY Claire-Marie Mrs \x26 LAGIRAFE Sophie Mrs QUOTE ONLY","Quote","01 Nov 2025"],0,7]`

	result, err := ParseFileSummaryResponse(responseBody)
	if err != nil {
		t.Fatalf("ParseFileSummaryResponse failed: %v", err)
	}

	fileSummary, ok := result.(GWTFileSummary)
	if !ok {
		t.Fatalf("Expected result to be GWTFileSummary, got %T", result)
	}

	if fileSummary.TotalBuyPrice != "295.0000" {
		t.Errorf("Expected TotalBuyPrice '295.0000', got '%s'", fileSummary.TotalBuyPrice)
	}
	if fileSummary.Currency != "AUD" {
		t.Errorf("Expected Currency 'AUD', got '%s'", fileSummary.Currency)
	}
	if fileSummary.HotelName != "HOLIDAY INN SYDNEY POTTS POINT" {
		t.Errorf("Expected HotelName 'HOLIDAY INN SYDNEY POTTS POINT', got '%s'", fileSummary.HotelName)
	}
	if fileSummary.Status != "Quote" {
		t.Errorf("Expected Status 'Quote', got '%s'", fileSummary.Status)
	}
	if fileSummary.ClientName != "BRAY Claire-Marie Mrs & LAGIRAFE Sophie Mrs QUOTE ONLY" {
		t.Errorf("Expected ClientName 'BRAY Claire-Marie Mrs & LAGIRAFE Sophie Mrs QUOTE ONLY', got '%s'", fileSummary.ClientName)
	}
}

func TestParseFileSummaryResponse_Array(t *testing.T) {
	responseBody := `//OK[103,102,101,0,100,'$3Jl',0,6,'7MK',0,0,99,0,98,'Bgzo9',33,0,97,'93Y',6,50,35,89,96,0,0,0,95,'BLza',0,94,22,0,0,0,3,93,0,0,92,45,91,90,89,0,0,1,88,0,0,0,1,0,87,14,0,13,86,0,85,'Bgzo6',33,0,84,'7V',6,36,35,74,83,0,'Bwi',33,82,0,81,0,31,80,0,29,27,63,27,1,26,62,0,79,'U2Wo',0,78,22,0,0,0,3,6,0,0,77,45,76,75,74,0,0,1,73,0,0,0,1,0,72,14,0,13,71,0,70,'Bgzo4',33,0,69,'oC',6,36,35,56,68,0,'BwU',33,67,0,66,0,65,64,0,29,27,63,27,1,26,62,0,61,'U2U5',0,60,22,0,0,0,3,6,0,0,59,45,58,57,56,0,0,1,55,0,0,0,1,0,54,14,0,13,53,0,52,'BgyE4',33,0,51,'8gP',6,50,35,42,0,0,0,0,49,'BTGw',0,48,22,0,0,0,3,47,0,0,46,45,44,43,42,0,0,1,41,1,0,0,1,0,40,14,0,13,39,0,38,'BgzlJ',33,0,37,'BUg8',6,36,35,17,34,0,'dh',33,32,0,30,0,31,30,0,29,27,28,27,1,26,25,0,24,'U0oR',0,23,22,0,0,0,3,6,0,0,21,20,19,18,17,0,0,1,16,0,0,0,1,0,15,14,0,13,5,3,12,1,0,9,4,1,3,863,0,0,1,11,0,0,200,10,0,1,1,9,4,1,3,0,8,0,0,0,0,35,7,0,6,0,5,4,1,3,1,2,1,2,1,["com.lynxtraveltech.client.shared.model.FileSummary/2582189904","3045.9600","java.util.ArrayList/4159755760","java.lang.String/2004016611","AUD 2,775.75","","$270.21","AUD","AUD 3,045.96","8.87%","0","All Transactions Confirmed","com.lynxtraveltech.client.shared.model.TransactionSummary/2232640144","-     ","AUD 159.55","C045ZKS500","AUD 181.00","DEVONPORT, TAS","11.85%","Nathalie da Costa","23 Jun 2025 07:02 PM","2 adults","1 X Standard Room 2 Queen Beds (1-4 Pax)","Standard Room 2 Queen Beds (1-4 Pax)","DH","com.lynxtraveltech.common.gui.shared.model.SimpleRatePlanBean/240540241","java.sql.Date/730999118","120-0-1-0-0-0","150-11-31-0-0-0","Best Flexible","RMONY","P1P","java.lang.Long/4227064769","2 x Queen;  ","Confirmed","C","NOVOTEL DEVONPORT","29/30 Nov 2025","16478821-9","AUD 475.20","3509-5490-AU-3","AUD 498.96","NATIONAL","4.76%","Claire-Marie Bray","24 Jun 2025 05:08 PM","Devonport Airport\x27s operating hours:\n\n    Mon - Fri: 8am - 5pm\n    Sat: 9am - 12pm\n    Sun: 1pm-5pm\n\nReservation No: 3509-5490-AU-3 \nCustomer Name: Mrs Muriel Jacquemoud Jaquier.\nPick Up Location: Devonport Airport \nPick Up Date / Time: 30NOV25/1300 \nDrop Off Location: Hobart Airport \nDrop Off Date / Time: 05DEC25/2000 \nVehicle Type: Compact SUV - Mitsubishi ASX or similar - S \n","1 X Group S - CFAR Auto - Compact SUV - On Airport - Excluding Northern Territory (Toyota CH-R o.s.) (T131170-74)","Group S - CFAR Auto - Compact SUV - On Airport - Excluding Northern Territory (Toyota CH-R o.s.) (T131170-74)","S","AVIS RENTAL CAR SWITZERLAND (Wizard# DA948V) - Knecht \x26 Voyageplan only","30 Nov/05 Dec 2025","16478821-6","AUD 770.00","PPN-FTSWB252257-10-1SM","AUD 871.00","CRADLE MOUNTAIN, TAS","11.60%","24 Jun 2025 05:05 PM","1 X Deluxe Spa Room (1-2 Pax, Max 4)","Deluxe Spa Room (1-2 Pax, Max 4)","SM","123-0-1-0-0-0","Best Available Rate - Including Full Breakfast in Altitude Restaurant.","BFAST","Nett Rates (TD1)","TD1","2 x Double;  ","CRADLE MOUNTAIN HOTEL","30 Nov/02 Dec 2025","16478821-10","AUD 215.00","PPN-FTSWB252257-11-1SM","AUD 233.00","STRAHAN, TAS","7.73%","24 Jun 2025 05:06 PM","1 X Village Waterfront Terrace Room (1-2 Pax, Max 2) (VWTQ)","Village Waterfront Terrace Room (1-2 Pax, Max 2) (VWTQ)","Best Available Rate - Room Only.","DTrade Nett (DTRADE)","DTRADE","1 x Queen;  ","STRAHAN VILLAGE","02/03 Jan 2026","16478821-11","AUD 1,156.00","LH25062043961957","AUD 1,262.00","COLES BAY, TAS","8.40%","24 Jun 2025 05:10 PM","Mount Paul Lounge Restaurant offers a Seasonal Omakase set menu that seamlessly fuses Tasmanian produce with Japanese cooking techniques and presentation. Dinner reservations are mandatory and must be made a minimum of 3 days in advance, no walk-ins accepted. Due to the nature of the dining experience, the restaurant is unable to cater for any major dietary requests.","1 X Ocean View Studio Including Continental Breakfast Provisions (1-2 Pax, Max 2)","Ocean View Studio Including Continental Breakfast Provisions (1-2 Pax, Max 2)","1 x King;  ","FREYCINET RESORT","03/05 Jan 2026","16478821-13","FTSWB252257","JACQUEMOUD","Open","29 Nov 2025"],0,7]`

	result, err := ParseFileSummaryResponse(responseBody)
	if err != nil {
		t.Fatalf("ParseFileSummaryResponse failed: %v", err)
	}

	// Debug output
	t.Logf("Result type: %T", result)
	if summaries, ok := result.([]GWTFileSummary); ok {
		t.Logf("Found %d summaries", len(summaries))
		for i, summary := range summaries {
			t.Logf("Summary %d: TotalBuyPrice=%s, HotelName=%s, Status=%s",
				i+1, summary.TotalBuyPrice, summary.HotelName, summary.Status)
		}
	} else if summary, ok := result.(GWTFileSummary); ok {
		t.Logf("Single summary: TotalBuyPrice=%s, HotelName=%s, Status=%s",
			summary.TotalBuyPrice, summary.HotelName, summary.Status)
	}

	summaries, ok := result.([]GWTFileSummary)
	if !ok {
		t.Fatalf("Expected result to be []GWTFileSummary, got %T", result)
	}

	if len(summaries) != 5 {
		t.Errorf("Expected 5 file summaries, got %d", len(summaries))
	}

	// Check a few fields from the first summary
	first := summaries[0]
	if first.TotalBuyPrice != "3045.9600" {
		t.Errorf("First summary: expected TotalBuyPrice '3045.9600', got '%s'", first.TotalBuyPrice)
	}
	if first.Currency != "AUD" {
		t.Errorf("First summary: expected Currency 'AUD', got '%s'", first.Currency)
	}
	if first.HotelName != "NOVOTEL DEVONPORT" {
		t.Errorf("First summary: expected HotelName 'NOVOTEL DEVONPORT', got '%s'", first.HotelName)
	}
	if first.Status != "Confirmed" {
		t.Errorf("First summary: expected Status 'Confirmed', got '%s'", first.Status)
	}

	// Check a few fields from the last summary
	last := summaries[4]
	if last.HotelName != "FREYCINET RESORT" {
		t.Errorf("Last summary: expected HotelName 'FREYCINET RESORT', got '%s'", last.HotelName)
	}
	if last.Status != "Open" {
		t.Errorf("Last summary: expected Status 'Open', got '%s'", last.Status)
	}
}
