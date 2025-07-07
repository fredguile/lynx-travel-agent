package gwt

import (
	"testing"
)

func TestParseGWTRetrieveItineraryResponseBody(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		expectedError  bool
		expectedResult *RetrieveItineraryResponseArray
	}{
		{
			name:          "Test Case 1 - Complex itinerary with multiple transactions",
			responseBody:  "//OK[125,124,123,0,122,'$xOpT',0,121,'7LC',0,0,120,0,119,'Bfz5A',23,0,51,'zR',6,21,20,114,0,0,0,0,118,'U5Ko',0,117,18,0,0,0,3,6,0,0,116,16,115,61,114,0,0,1,113,0,0,0,1,0,112,13,0,12,0,0,111,'BgIpk',23,0,22,'B',6,21,20,14,0,0,0,0,0,'D',0,110,18,0,1,0,3,0,0,0,109,16,15,0,14,0,0,1,0,0,0,0,1,0,14,13,0,12,108,0,107,'BgBFz',23,0,106,'BAjc',6,21,20,99,105,0,0,0,104,'U5n_',0,103,18,0,0,0,3,102,0,0,101,16,100,61,99,0,1,1,98,0,0,0,1,0,97,13,0,12,96,0,95,'BgExP',23,0,76,'tU',6,75,20,93,74,0,-33,73,0,71,0,72,71,0,-32,-31,1,67,66,0,65,'DMx',0,64,18,0,0,0,3,6,0,0,63,16,94,61,93,0,0,1,92,0,0,0,1,0,91,13,0,12,90,0,89,'BgFGt',23,0,88,'BTY',6,21,20,81,0,0,0,0,87,'WqD',0,86,18,0,0,0,3,85,0,0,84,83,82,61,81,0,0,1,80,0,0,0,1,0,79,13,0,12,78,0,77,'BgExO',23,0,76,'tU',6,75,20,60,74,0,'P5',23,73,0,71,0,72,71,0,70,68,69,68,1,67,66,0,65,'DMx',0,64,18,0,0,0,3,6,0,0,63,16,62,61,60,0,0,1,59,0,0,0,1,0,58,13,0,12,57,0,56,'BgBFx',23,0,34,'u_',6,21,20,27,33,0,0,0,32,'BSwI',0,31,18,0,0,0,3,6,0,0,55,16,29,28,27,0,0,1,54,0,0,0,1,0,25,13,0,12,53,0,52,-24,0,51,'zR',6,21,20,39,0,0,0,0,50,'BPNU',0,49,18,0,0,0,0,0,0,'Bf$n4',23,0,0,35.0,45,48,43,0,0,'A',2,0,'U5Kf',0,0,0,0,0,0,47,0,0,0,0,0,0,35.0,45,46,43,0,0,0,0,0,1,0,0,0,0,0.0,0,0,0,0,31.5,45,44,43,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42,1,3,6,0,0,41,16,40,28,39,0,0,1,38,0,0,0,1,0,37,13,0,12,36,0,35,'BgBFw',23,0,34,'u_',6,21,20,27,33,0,0,0,32,'BSwI',0,31,18,0,0,0,3,6,0,0,30,16,29,28,27,0,0,1,26,0,0,0,1,0,25,13,0,12,0,0,24,'BgFLy',23,0,22,'B',6,21,20,14,0,0,0,0,0,'D',0,19,18,0,1,0,3,0,0,0,17,16,15,0,14,0,0,1,0,0,0,0,1,0,14,13,0,12,10,3,0,0,0,9,4,1,3,863,0,0,1,11,1,0,200,10,0,1,1,9,4,1,3,0,8,0,0,0,0,21,7,0,6,0,5,4,1,3,0,2,0,2,1,[\"com.lynxtraveltech.client.shared.model.FileSummary/2582189904\",\"10016.0000\",\"java.util.ArrayList/4159755760\",\"java.lang.String/2004016611\",\"AUD 9,203.46\",\"\",\"$812.54\",\"AUD\",\"AUD 10,016.00\",\"8.11%\",\"0\",\"com.lynxtraveltech.client.shared.model.TransactionSummary/2232640144\",\"-     \",\"-\",\"N/A\",\"Kathrin Biberstein\",\"31 Dec 2024 12:44 PM\",\"2 adults\",\"*** Mobile Number Lucy NOTA B: +41 79 361 40 30 *** No allergies, dietary requirements or health issues. ***\",\"Confirmed\",\"S\",\"Own Arrangements\",\"java.lang.Long/4227064769\",\"02 Oct 2025\",\"AUD 162.00\",\"2194C574\",\"AUD 184.00\",\"ALICE SPRINGS, NT\",\"11.96%\",\"31 Dec 2024 06:05 PM\",\"1 X Standard Twin Room (1-2 Pax, Max 3A or 2A2C)\",\"Standard Twin Room (1-2 Pax, Max 3A or 2A2C)\",\"2 x Queen;  \",\"CROWNE PLAZA ALICE SPRINGS LASSETERS\",\"04/05 Oct 2025\",\"16454569-4\",\"AUD 2,242.50\",\"PPA3CC057402/XX124156\",\"AUD 2,422.00\",\"7.48%\",\"16 Jan 2025 10:49 AM\",\"com.lynxtraveltech.common.gui.shared.model.OptionalService/1553404829\",\"java.math.BigDecimal/8151472\",\"63.00\",\"java.lang.Double/858496421\",\"70.0000\",\"Sleeping Bag Hire\",\"70.00\",\"2 X 5 Day Outback Camping Adventure - Alice Springs to Alice Springs - Powered Tent/Swag Camping (AR5)\",\"5 Day Outback Camping Adventure - Alice Springs to Alice Springs - Powered Tent/Swag Camping (AR5)\",\"ADVENTURE TOURS AUSTRALIA\",\"05/09 Oct 2025\",\"16454569-3\",\"223447575\",\"31 Dec 2024 06:04 PM\",\"09/10 Oct 2025\",\"16454569-5\",\"AUD 549.82\",\"1773Z45J9500\",\"AUD 622.00\",\"PERTH, WA\",\"11.60%\",\"30 Dec 2024 06:01 PM\",\"1 X Standard Room (1-2 Pax, Max 2)\",\"Standard Room (1-2 Pax, Max 2)\",\"DH\",\"com.lynxtraveltech.common.gui.shared.model.SimpleRatePlanBean/240540241\",\"java.sql.Date/730999118\",\"120-0-1-0-0-0\",\"150-11-31-0-0-0\",\"Best Flexible\",\"RMONY\",\"P1P\",\"1 x Double;  \",\"C\",\"IBIS PERTH\",\"10/13 Oct 2025\",\"16454569-12\",\"AUD 2,317.50\",\"RUE234MPTA\",\"AUD 2,504.00\",\"7.45%\",\"Becky Caiels\",\"03 Jul 2025 03:05 PM\",\"******* MINIMUM NUMBERS NOT MET YET ******* *** A 40% deposit is due at time of booking. *** Hotel pick-up is not provided, clients are asked to make their own way to the departure location. *** ********* SLEEPING BAGS CANNOT BE BOOKED *********\",\"2 X Adventures To Awaken - 6 Day Esperance \x26 Margaret River Adventure - Camping (PE6)\",\"Adventures To Awaken - 6 Day Esperance \x26 Margaret River Adventure - Camping (PE6)\",\"UNTAMED ESCAPES\",\"13/18 Oct 2025\",\"16454569-16\",\"AUD 199.64\",\"1773ZJH546y\",\"AUD 226.00\",\"11.66%\",\"18/19 Oct 2025\",\"16454569-13\",\"AUD 814.50\",\"2C4WH6XX\",\"AUD 896.00\",\"9.10%\",\"08 Jan 2025 11:01 PM\",\"*** Departure location: Pier 3, Barrack St Jetty, Perth - 08.30h *** Boarding for all passengers closes 10 minutes prior to departure. ***\",\"1 X 2 Night Package - Standard Tent Including Breakfast \x26 Return Ferry (1-2 Pax, Max 2) - Ex Perth\",\"2 Night Package - Standard Tent Including Breakfast \x26 Return Ferry (1-2 Pax, Max 2) - Ex Perth\",\"1 x King;  \",\"SEALINK ROTTNEST ISLAND\",\"19/21 Oct 2025\",\"16454569-7\",\"15 Jan 2025 12:13 PM\",\"DoubleTree by Hilton in Perth booked directly through bedbanks by agent.\",\"21/25 Oct 2025\",\"AUD 2,692.50\",\"PPA3XXX503-778\",\"AUD 2,908.00\",\"7.41%\",\"08 Jan 2025 10:44 PM\",\"2 X 7 Day Perth to Exmouth \x26 Return - Hostel Dorm (PX7)\",\"7 Day Perth to Exmouth \x26 Return - Hostel Dorm (PX7)\",\"25/31 Oct 2025\",\"16454569-2\",\"1061848\",\"FTSWA230184\",\"STIRFRY, Mrs / NOTA B, Lucy Mrs\",\"Partial\",\"08 Oct 2025\"],0,7]",
			expectedError: false,
			expectedResult: &RetrieveItineraryResponseArray{
				Type:             "Partial",
				PartyName:        "STIRFRY, Mrs / NOTA B, Lucy Mrs",
				FileReference:    "FTSWA230184",
				FileIdentifier:   "$xOpT",
				ClientIdentifier: "7LC",
				AgentReference:   "1061848",
				ItineraryCount:   10,
				Itineraries: []ItineraryTransactionSummary{
					{
						VoucherIdentifier:     "",
						Date:                  "02 Oct 2025",
						TransactionIdentifier: "BgFLy",
						Supplier:              "Own Arrangements",
						Status:                "Confirmed",
						ConfirmationNumber:    "",
						Location:              "",
					},
					{
						VoucherIdentifier:     "16454569-4",
						Date:                  "04/05 Oct 2025",
						TransactionIdentifier: "BgBFw",
						Supplier:              "CROWNE PLAZA ALICE SPRINGS LASSETERS",
						Status:                "Confirmed",
						ConfirmationNumber:    "2194C574",
						Location:              "ALICE SPRINGS, NT",
					},
					{
						VoucherIdentifier:     "16454569-3",
						Date:                  "05/09 Oct 2025",
						TransactionIdentifier: "Bf$n4",
						Supplier:              "ADVENTURE TOURS AUSTRALIA",
						Status:                "Confirmed",
						ConfirmationNumber:    "PPA3CC057402/XX124156",
						Location:              "ALICE SPRINGS, NT",
					},
					{
						VoucherIdentifier:     "16454569-5",
						Date:                  "09/10 Oct 2025",
						TransactionIdentifier: "BgBFx",
						Supplier:              "CROWNE PLAZA ALICE SPRINGS LASSETERS",
						Status:                "Confirmed",
						ConfirmationNumber:    "223447575",
						Location:              "ALICE SPRINGS, NT",
					},
					{
						VoucherIdentifier:     "16454569-12",
						Date:                  "10/13 Oct 2025",
						TransactionIdentifier: "BgExO",
						Supplier:              "IBIS PERTH",
						Status:                "Confirmed",
						ConfirmationNumber:    "1773Z45J9500",
						Location:              "PERTH, WA",
					},
					{
						VoucherIdentifier:     "16454569-16",
						Date:                  "13/18 Oct 2025",
						TransactionIdentifier: "BgFGt",
						Supplier:              "UNTAMED ESCAPES",
						Status:                "Confirmed",
						ConfirmationNumber:    "RUE234MPTA",
						Location:              "PERTH, WA",
					},
					{
						VoucherIdentifier:     "16454569-13",
						Date:                  "18/19 Oct 2025",
						TransactionIdentifier: "BgExP",
						Supplier:              "IBIS PERTH",
						Status:                "Confirmed",
						ConfirmationNumber:    "1773ZJH546y",
						Location:              "PERTH, WA",
					},
					{
						VoucherIdentifier:     "16454569-7",
						Date:                  "19/21 Oct 2025",
						TransactionIdentifier: "BgBFz",
						Supplier:              "SEALINK ROTTNEST ISLAND",
						Status:                "Confirmed",
						ConfirmationNumber:    "2C4WH6XX",
						Location:              "PERTH, WA",
					},
					{
						VoucherIdentifier:     "",
						Date:                  "21/25 Oct 2025",
						TransactionIdentifier: "BgIpk",
						Supplier:              "Own Arrangements",
						Status:                "Confirmed",
						ConfirmationNumber:    "",
						Location:              "",
					},
					{
						VoucherIdentifier:     "16454569-2",
						Date:                  "25/31 Oct 2025",
						TransactionIdentifier: "Bfz5A",
						Supplier:              "ADVENTURE TOURS AUSTRALIA",
						Status:                "Confirmed",
						ConfirmationNumber:    "PPA3XXX503-778",
						Location:              "PERTH, WA",
					},
				},
			},
		},
		{
			name:          "Test Case 2 - Simpler itinerary",
			responseBody:  "//OK[54,53,52,0,51,'$2s7',0,50,'BAkC',0,0,49,0,48,'Bg2yK',26,0,47,'tXw',6,36,35,41,0,0,0,0,46,'BKTt',0,45,31,0,0,0,3,6,0,0,44,18,43,42,41,0,1,0,0,0,0,0,1,0,40,13,0,12,39,0,38,-24,0,37,'sw',6,36,35,15,34,0,0,0,33,'UxU1',0,32,31,0,0,0,0,0,0,-24,0,0,0.0,22,-26,0,0,'A',1,0,'8C9',30,0,0,0,0,0,29,0,0,0,0,0,0,0.0,22,28,23,0,0,0,0,0,1,0,0,0,0,0.0,0,0,0,0,0.0,22,27,23,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,21,0,0,0,0,'BgsOD',26,0,0,30.0,22,24,23,0.0,22,0,'A',1,1,'DwL',6,0,0,0,0,0,25,0,0,0,0,0,0,30.0,22,24,23,0.0,22,0,0,0,0,1,0,0,0,0,0.0,0,11,23,11,23,11,23,30.0,22,24,23,0.0,22,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,21,2,3,20,0,0,19,18,17,16,15,0,1,0,6,0,0,0,1,0,14,13,0,12,2,3,0,0,0,9,4,1,3,863,1,0,1,11,0,0,200,10,0,1,0,9,4,1,3,0,8,0,0,0,0,1,7,0,6,0,5,4,1,3,0,2,0,2,1,[\"com.lynxtraveltech.client.shared.model.FileSummary/2582189904\",\"430.0000\",\"java.util.ArrayList/4159755760\",\"java.lang.String/2004016611\",\"AUD 414.00\",\"\",\"$16.00\",\"AUD\",\"AUD 430.00\",\"3.72%\",\"0\",\"com.lynxtraveltech.client.shared.model.TransactionSummary/2232640144\",\"-     \",\"AUD 300.00\",\"AUD 310.00\",\"SYDNEY, NSW\",\"2.94%\",\"Claire-Marie Bray\",\"07 Jul 2025 10:02 PM\",\"notes\",\"com.lynxtraveltech.common.gui.shared.model.OptionalService/1553404829\",\"java.lang.Double/858496421\",\"java.math.BigDecimal/8151472\",\"30.00\",\"Full Breakfast\",\"java.lang.Long/4227064769\",\"0.00\",\"0.0000\",\"Honeymoon Offer: Bottle Of Sparkling Wine - Combinable\",\"HMOON\",\"2 adults\",\"1 X Premium City View Room (1-2 Pax, Max 4)\",\"Premium City View Room (1-2 Pax, Max 4)\",\"2 x Queen;  \",\"On Request\",\"A\",\"HOLIDAY INN SYDNEY POTTS POINT\",\"01/02 Nov 2025\",\"16476987-1\",\"AUD 84.00\",\"AUD 90.00\",\"CAIRNS, QLD\",\"6.67%\",\"06 Jul 2025 06:35 AM\",\"2 X Cairns to Fitzroy Island - one way Transfer with The Fitzroy Flyer\",\"Cairns to Fitzroy Island - one way Transfer with The Fitzroy Flyer\",\"FITZROY ISLAND RESORT\",\"02 Nov 2025\",\"16476987-2\",\"25052025\",\"FTAUB252039\",\"BRAY Claire-Marie Mrs \x26 LAGIRAFE Sophie Mrs QUOTE ONLY\",\"Quote\",\"01 Nov 2025\"],0,7]",
			expectedError: false,
			expectedResult: &RetrieveItineraryResponseArray{
				Type:             "Quote",
				PartyName:        "BRAY Claire-Marie Mrs & LAGIRAFE Sophie Mrs QUOTE ONLY",
				FileReference:    "FTAUB252039",
				FileIdentifier:   "$2s7",
				ClientIdentifier: "BAkC",
				AgentReference:   "25052025",
				ItineraryCount:   2,
				Itineraries: []ItineraryTransactionSummary{
					{
						VoucherIdentifier:     "16476987-1",
						Date:                  "01/02 Nov 2025",
						TransactionIdentifier: "BgsOD",
						Supplier:              "HOLIDAY INN SYDNEY POTTS POINT",
						Status:                "On Request",
						ConfirmationNumber:    "",
						Location:              "SYDNEY, NSW",
					},
					{
						VoucherIdentifier:     "16476987-2",
						Date:                  "02 Nov 2025",
						TransactionIdentifier: "Bg2yK",
						Supplier:              "FITZROY ISLAND RESORT",
						Status:                "On Request",
						ConfirmationNumber:    "",
						Location:              "CAIRNS, QLD",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseGWTRetrieveItineraryResponseBody(tt.responseBody)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			// Check basic fields
			if result.Type != tt.expectedResult.Type {
				t.Errorf("Type mismatch: got %s, want %s", result.Type, tt.expectedResult.Type)
			}

			if result.PartyName != tt.expectedResult.PartyName {
				t.Errorf("PartyName mismatch: got %s, want %s", result.PartyName, tt.expectedResult.PartyName)
			}

			if result.FileReference != tt.expectedResult.FileReference {
				t.Errorf("FileReference mismatch: got %s, want %s", result.FileReference, tt.expectedResult.FileReference)
			}

			if result.FileIdentifier != tt.expectedResult.FileIdentifier {
				t.Errorf("FileIdentifier mismatch: got %s, want %s", result.FileIdentifier, tt.expectedResult.FileIdentifier)
			}

			if result.ClientIdentifier != tt.expectedResult.ClientIdentifier {
				t.Errorf("ClientIdentifier mismatch: got %s, want %s", result.ClientIdentifier, tt.expectedResult.ClientIdentifier)
			}

			if result.AgentReference != tt.expectedResult.AgentReference {
				t.Errorf("AgentReference mismatch: got %s, want %s", result.AgentReference, tt.expectedResult.AgentReference)
			}

			if result.ItineraryCount != tt.expectedResult.ItineraryCount {
				t.Errorf("ItineraryCount mismatch: got %d, want %d", result.ItineraryCount, tt.expectedResult.ItineraryCount)
			}

			// Check that itineraries slice is not nil
			if result.Itineraries == nil {
				t.Errorf("Itineraries slice is nil")
			}

			for i, itinerary := range result.Itineraries {
				if itinerary.VoucherIdentifier != tt.expectedResult.Itineraries[i].VoucherIdentifier {
					t.Errorf("Itinerary[%d].VoucherIdentifier mismatch: got %s, want %s", i, itinerary.VoucherIdentifier, tt.expectedResult.Itineraries[i].VoucherIdentifier)
				}

				if itinerary.Date != tt.expectedResult.Itineraries[i].Date {
					t.Errorf("Itinerary[%d].Date mismatch: got %s, want %s", i, itinerary.Date, tt.expectedResult.Itineraries[i].Date)
				}

				if itinerary.TransactionIdentifier != tt.expectedResult.Itineraries[i].TransactionIdentifier {
					t.Errorf("Itinerary[%d].TransactionIdentifier mismatch: got %s, want %s", i, itinerary.TransactionIdentifier, tt.expectedResult.Itineraries[i].TransactionIdentifier)
				}

				if itinerary.Supplier != tt.expectedResult.Itineraries[i].Supplier {
					t.Errorf("Itinerary[%d].Supplier mismatch: got %s, want %s", i, itinerary.Supplier, tt.expectedResult.Itineraries[i].Supplier)
				}

				if itinerary.Status != tt.expectedResult.Itineraries[i].Status {
					t.Errorf("Itinerary[%d].Status mismatch: got %s, want %s", i, itinerary.Status, tt.expectedResult.Itineraries[i].Status)
				}

				if itinerary.ConfirmationNumber != tt.expectedResult.Itineraries[i].ConfirmationNumber {
					t.Errorf("Itinerary[%d].ConfirmationNumber mismatch: got %s, want %s", i, itinerary.ConfirmationNumber, tt.expectedResult.Itineraries[i].ConfirmationNumber)
				}

				if itinerary.Location != tt.expectedResult.Itineraries[i].Location {
					t.Errorf("Itinerary[%d].Location mismatch: got %s, want %s", i, itinerary.Location, tt.expectedResult.Itineraries[i].Location)
				}
			}
		})
	}
}

func TestParseGWTRetrieveItineraryResponseBody_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		responseBody  string
		expectedError bool
	}{
		{
			name:          "Empty response body",
			responseBody:  "",
			expectedError: true,
		},
		{
			name:          "Invalid GWT format",
			responseBody:  "invalid format",
			expectedError: true,
		},
		{
			name:          "Missing protocol version",
			responseBody:  "//OK[1,2,3]",
			expectedError: true,
		},
		{
			name:          "Too few array elements",
			responseBody:  "//OK[1,2,3,4]",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseGWTRetrieveItineraryResponseBody(tt.responseBody)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
