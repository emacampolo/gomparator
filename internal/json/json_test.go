package json

import (
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name    string
		isEqual bool
		b1      []byte
		b2      []byte
	}{
		{
			name:    "empty array",
			isEqual: true,
			b1:      []byte(`[]`),
			b2:      []byte(`[]`),
		},
		{
			name:    "empty object",
			isEqual: true,
			b1:      []byte(`{}`),
			b2:      []byte(`{}`),
		},
		{
			name:    "different basic types",
			isEqual: false,
			b1:      []byte(`[]`),
			b2:      []byte(`{}`),
		},
		{
			name:    "different first name",
			isEqual: false,
			b1:      []byte(`{"FirstName":"Alan"}`),
			b2:      []byte(`{"FirstName":"Galileo"}`),
		},
		{
			name:    "equal first name",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan"}`),
			b2:      []byte(`{"FirstName":"Alan"}`),
		},
		{
			name:    "equal fields different order",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan", "LastName": "Turing"}`),
			b2:      []byte(`{"LastName": "Turing", "FirstName":"Alan"}`),
		},
		{
			name:    "equal fields different order",
			isEqual: true,
			b1:      []byte(`{"x": {"t": 1, "s": 2}, "z": 1}`),
			b2:      []byte(`{"z": 1, "x": {"s": 2, "t": 1}}`),
		},
		{
			name:    "equal array objects",
			isEqual: true,
			b1:      []byte(`[{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}]`),
			b2:      []byte(`[{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}]`),
		},
		{
			name:    "equal exact same objects",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
		},
		{
			name:    "equal flip one field at the root",
			isEqual: true,
			b1:      []byte(`{"LastName": "Turing", "FirstName":"Alan", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
		},
		{
			name:    "flip one field inside a node",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Martin Fowler", "Rob Pike"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
		},
		{
			name:    "Flip one field inside inner node",
			isEqual: true,
			b1:      []byte(`{"LastName": "Turing", "FirstName":"Alan", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"LastName": "Turing", "FirstName":"Alan", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Number": 111, "Country":  "US"}]}`),
		},
		{
			name:    "flip object in root array",
			isEqual: true,
			b1:      []byte(`[{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}, {"FirstName":"Martin", "LastName": "Fowler", "Age" : 30, "Friends" : []}]`),
			b2:      []byte(`[{"FirstName":"Martin", "LastName": "Fowler", "Age" : 30, "Friends" : []}, {"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}]`),
		},
		{
			name:    "some complex example",
			isEqual: true,
			b1:      []byte(`[{"a": 1, "b": [{"c": [1,5,2,4]}, {"d": [1]}]}]`),
			b2:      []byte(`[{"b": [{"d": [1]}, {"c": [1,2,4,5]}], "a": 1}]`),
		},
		{
			name:    "some complex example with slices of different types",
			isEqual: true,
			b1:      []byte(`{"paging":{"total":2,"limit":30,"offset":0},"results":[{"financial_institutions":[],"secure_thumbnail":"https://www.mercadopago.com/org-img/MP3/API/logos/visa.gif","payer_costs":[{"installment_reduced_tea":null,"installment_reduced_cft":null,"installments":1,"installment_rate":0,"id":67696444,"installment_full_cft":null,"discount_rate":0,"min_allowed_amount":0,"installment_full_tea":null,"labels":["CFT_0,00%|TEA_0,00%"],"max_allowed_amount":250000,"base_installment_rate":0,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":150.35,"installment_reduced_cft":199.44,"installments":3,"installment_rate":19.72,"id":67696441,"installment_full_cft":199.44,"discount_rate":0,"min_allowed_amount":2,"installment_full_tea":150.35,"labels":["CFT_199,44%|TEA_150,35%"],"max_allowed_amount":250000,"base_installment_rate":19.72,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":142.79,"installment_reduced_cft":187.01,"installments":6,"installment_rate":34.49,"id":67696442,"installment_full_cft":187.01,"discount_rate":0,"min_allowed_amount":3,"installment_full_tea":142.79,"labels":["CFT_187,01%|TEA_142,79%","recommended_interest_installment_with_some_banks"],"max_allowed_amount":250000,"base_installment_rate":34.49,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":136.21,"installment_reduced_cft":176.57,"installments":9,"installment_rate":49.19,"id":441302368,"installment_full_cft":176.57,"discount_rate":0,"min_allowed_amount":2,"installment_full_tea":136.21,"labels":["CFT_176,57%|TEA_136,21%"],"max_allowed_amount":250000,"base_installment_rate":49.19,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":130.33,"installment_reduced_cft":167.52,"installments":12,"installment_rate":63.77,"id":67696440,"installment_full_cft":167.52,"discount_rate":0,"min_allowed_amount":6,"installment_full_tea":130.33,"labels":["recommended_installment","CFT_167,52%|TEA_130,33%"],"max_allowed_amount":250000,"base_installment_rate":63.77,"installment_rate_collector":["MERCADOPAGO"]}],"issuer":{"default":true,"name":"Visa Argentina S.A.","id":1},"total_financial_cost":null,"min_accreditation_days":0,"max_accreditation_days":2,"merchant_account_id":null,"id":"visa","payment_type_id":"credit_card","accreditation_time":2880,"owner":"site","settings":[{"security_code":{"mode":"mandatory","card_location":"back","length":3},"card_number":{"length":16,"validation":"standard"},"bin":{"pattern":"^4","installments_pattern":"^4","exclusion_pattern":"^(476520|473713|473713|473227|444493|410122|405517|402789|417856|448712|453770|434541|411199|423465|434540|434542|434538|423018|488241|489634|434537|434539|434536|427156|427157|434535|434534|434533|423077|434532|434586|423001|434531|411197|443264|400276|400615|402914|404625|405069|434543|416679|405515|405516|405755|405896|405897|406290|406291|406375|406652|406998|406999|408515|410082|410083|410121|410123|410853|411849|417309|421738|423623|428062|428063|428064|434795|437996|439818|442371|442548|444060|446343|446344|446347|450412|450799|451377|451701|451751|451756|451757|451758|451761|451763|451764|451765|451766|451767|451768|451769|451770|451772|451773|457596|457665|462815|463465|468508|473710|473711|473712|473714|473715|473716|473717|473718|473719|473720|473721|473722|473725|477051|477053|481397|481501|481502|481550|483002|483020|483188|489412|492528|499859|446344|446345|446346|400448)"},"id":67696446}],"thumbnail":"http://img.mlstatic.com/org-img/MP3/API/logos/visa.gif","bins":[],"marketplace":"MELI","deferred_capture":"supported","labels":["recommended_method"],"financing_deals":{"legals":null,"installments":null,"expiration_date":null,"start_date":null,"status":"deactive"},"name":"Visa","site_id":"MLA","processing_mode":"aggregator","additional_info_needed":["cardholder_name","cardholder_identification_type","cardholder_identification_number"],"status":"active"},{"financial_institutions":[],"secure_thumbnail":"https://www.mercadopago.com/org-img/MP3/API/logos/visa.gif","payer_costs":[{"installment_reduced_tea":69,"installment_reduced_cft":87.81,"installments":1,"installment_rate":0,"id":67696449,"installment_full_cft":87.81,"discount_rate":0,"min_allowed_amount":2,"installment_full_tea":69,"labels":["CFT_0,00%|TEA_0,00%"],"max_allowed_amount":250000,"base_installment_rate":0,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":150.35,"installment_reduced_cft":199.44,"installments":3,"installment_rate":19.72,"id":67696450,"installment_full_cft":199.44,"discount_rate":0,"min_allowed_amount":3,"installment_full_tea":150.35,"labels":["CFT_199,44%|TEA_150,35%"],"max_allowed_amount":250000,"base_installment_rate":19.72,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":142.79,"installment_reduced_cft":187.01,"installments":6,"installment_rate":34.49,"id":67696451,"installment_full_cft":187.01,"discount_rate":0,"min_allowed_amount":5,"installment_full_tea":142.79,"labels":["CFT_187,01%|TEA_142,79%","recommended_interest_installment_with_some_banks"],"max_allowed_amount":250000,"base_installment_rate":34.49,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":165.77,"installment_reduced_cft":217.13,"installments":9,"installment_rate":56.9,"id":67696452,"installment_full_cft":217.13,"discount_rate":0,"min_allowed_amount":6,"installment_full_tea":165.77,"labels":["CFT_217,13%|TEA_165,77%"],"max_allowed_amount":250000,"base_installment_rate":56.9,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":168.35,"installment_reduced_cft":219.13,"installments":12,"installment_rate":77.4,"id":67696453,"installment_full_cft":219.13,"discount_rate":0,"min_allowed_amount":6,"installment_full_tea":168.35,"labels":["recommended_installment","CFT_219,13%|TEA_168,35%"],"max_allowed_amount":250000,"base_installment_rate":77.4,"installment_rate_collector":["MERCADOPAGO"]}],"issuer":{"default":true,"name":"Visa Argentina S.A.","id":1},"total_financial_cost":null,"min_accreditation_days":0,"max_accreditation_days":2,"merchant_account_id":null,"id":"visa","payment_type_id":"credit_card","accreditation_time":2880,"owner":"site","settings":[{"security_code":{"mode":"mandatory","card_location":"back","length":3},"card_number":{"length":16,"validation":"standard"},"bin":{"pattern":"^4","installments_pattern":"^4","exclusion_pattern":"^(476520|473713|473713|473227|444493|410122|405517|402789|417856|448712|453770|434541|411199|423465|434540|434542|434538|423018|488241|489634|434537|434539|434536|427156|427157|434535|434534|434533|423077|434532|434586|423001|434531|411197|443264|400276|400615|402914|404625|405069|434543|416679|405515|405516|405755|405896|405897|406290|406291|406375|406652|406998|406999|408515|410082|410083|410121|410123|410853|411849|417309|421738|423623|428062|428063|428064|434795|437996|439818|442371|442548|444060|446343|446344|446347|450412|450799|451377|451701|451751|451756|451757|451758|451761|451763|451764|451765|451766|451767|451768|451769|451770|451772|451773|457596|457665|462815|463465|468508|473710|473711|473712|473714|473715|473716|473717|473718|473719|473720|473721|473722|473725|477051|477053|481397|481501|481502|481550|483002|483020|483188|489412|492528|499859|446344|446345|446346|400448)"},"id":67696455}],"thumbnail":"http://img.mlstatic.com/org-img/MP3/API/logos/visa.gif","bins":[],"marketplace":"NONE","deferred_capture":"supported","labels":["recommended_method"],"financing_deals":{"legals":null,"installments":null,"expiration_date":null,"start_date":null,"status":"deactive"},"name":"Visa","site_id":"MLA","processing_mode":"aggregator","additional_info_needed":["cardholder_name","cardholder_identification_type","cardholder_identification_number"],"status":"active"}]}`),
			b2:      []byte(`{"paging":{"total":2,"limit":30,"offset":0},"results":[{"financial_institutions":[],"secure_thumbnail":"https://www.mercadopago.com/org-img/MP3/API/logos/visa.gif","payer_costs":[{"installment_reduced_tea":null,"installment_reduced_cft":null,"installments":1,"installment_rate":0,"id":67696444,"installment_full_cft":null,"discount_rate":0,"min_allowed_amount":0,"installment_full_tea":null,"labels":["CFT_0,00%|TEA_0,00%"],"max_allowed_amount":250000,"base_installment_rate":0,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":150.35,"installment_reduced_cft":199.44,"installments":3,"installment_rate":19.72,"id":67696441,"installment_full_cft":199.44,"discount_rate":0,"min_allowed_amount":2,"installment_full_tea":150.35,"labels":["CFT_199,44%|TEA_150,35%"],"max_allowed_amount":250000,"base_installment_rate":19.72,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":142.79,"installment_reduced_cft":187.01,"installments":6,"installment_rate":34.49,"id":67696442,"installment_full_cft":187.01,"discount_rate":0,"min_allowed_amount":3,"installment_full_tea":142.79,"labels":["CFT_187,01%|TEA_142,79%","recommended_interest_installment_with_some_banks"],"max_allowed_amount":250000,"base_installment_rate":34.49,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":136.21,"installment_reduced_cft":176.57,"installments":9,"installment_rate":49.19,"id":441302368,"installment_full_cft":176.57,"discount_rate":0,"min_allowed_amount":2,"installment_full_tea":136.21,"labels":["CFT_176,57%|TEA_136,21%"],"max_allowed_amount":250000,"base_installment_rate":49.19,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":130.33,"installment_reduced_cft":167.52,"installments":12,"installment_rate":63.77,"id":67696440,"installment_full_cft":167.52,"discount_rate":0,"min_allowed_amount":6,"installment_full_tea":130.33,"labels":["recommended_installment","CFT_167,52%|TEA_130,33%"],"max_allowed_amount":250000,"base_installment_rate":63.77,"installment_rate_collector":["MERCADOPAGO"]}],"issuer":{"default":true,"name":"Visa Argentina S.A.","id":1},"total_financial_cost":null,"min_accreditation_days":0,"max_accreditation_days":2,"merchant_account_id":null,"id":"visa","payment_type_id":"credit_card","accreditation_time":2880,"owner":"site","settings":[{"security_code":{"mode":"mandatory","card_location":"back","length":3},"card_number":{"length":16,"validation":"standard"},"bin":{"pattern":"^4","installments_pattern":"^4","exclusion_pattern":"^(476520|473713|473713|473227|444493|410122|405517|402789|417856|448712|453770|434541|411199|423465|434540|434542|434538|423018|488241|489634|434537|434539|434536|427156|427157|434535|434534|434533|423077|434532|434586|423001|434531|411197|443264|400276|400615|402914|404625|405069|434543|416679|405515|405516|405755|405896|405897|406290|406291|406375|406652|406998|406999|408515|410082|410083|410121|410123|410853|411849|417309|421738|423623|428062|428063|428064|434795|437996|439818|442371|442548|444060|446343|446344|446347|450412|450799|451377|451701|451751|451756|451757|451758|451761|451763|451764|451765|451766|451767|451768|451769|451770|451772|451773|457596|457665|462815|463465|468508|473710|473711|473712|473714|473715|473716|473717|473718|473719|473720|473721|473722|473725|477051|477053|481397|481501|481502|481550|483002|483020|483188|489412|492528|499859|446344|446345|446346|400448)"},"id":67696446}],"thumbnail":"http://img.mlstatic.com/org-img/MP3/API/logos/visa.gif","bins":[],"marketplace":"MELI","deferred_capture":"supported","labels":["recommended_method"],"financing_deals":{"legals":null,"installments":null,"expiration_date":null,"start_date":null,"status":"deactive"},"name":"Visa","site_id":"MLA","processing_mode":"aggregator","additional_info_needed":["cardholder_name","cardholder_identification_type","cardholder_identification_number"],"status":"active"},{"financial_institutions":[],"secure_thumbnail":"https://www.mercadopago.com/org-img/MP3/API/logos/visa.gif","payer_costs":[{"installment_reduced_tea":69,"installment_reduced_cft":87.81,"installments":1,"installment_rate":0,"id":67696449,"installment_full_cft":87.81,"discount_rate":0,"min_allowed_amount":2,"installment_full_tea":69,"labels":["CFT_0,00%|TEA_0,00%"],"max_allowed_amount":250000,"base_installment_rate":0,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":150.35,"installment_reduced_cft":199.44,"installments":3,"installment_rate":19.72,"id":67696450,"installment_full_cft":199.44,"discount_rate":0,"min_allowed_amount":3,"installment_full_tea":150.35,"labels":["CFT_199,44%|TEA_150,35%"],"max_allowed_amount":250000,"base_installment_rate":19.72,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":142.79,"installment_reduced_cft":187.01,"installments":6,"installment_rate":34.49,"id":67696451,"installment_full_cft":187.01,"discount_rate":0,"min_allowed_amount":5,"installment_full_tea":142.79,"labels":["CFT_187,01%|TEA_142,79%","recommended_interest_installment_with_some_banks"],"max_allowed_amount":250000,"base_installment_rate":34.49,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":165.77,"installment_reduced_cft":217.13,"installments":9,"installment_rate":56.9,"id":67696452,"installment_full_cft":217.13,"discount_rate":0,"min_allowed_amount":6,"installment_full_tea":165.77,"labels":["CFT_217,13%|TEA_165,77%"],"max_allowed_amount":250000,"base_installment_rate":56.9,"installment_rate_collector":["MERCADOPAGO"]},{"installment_reduced_tea":168.35,"installment_reduced_cft":219.13,"installments":12,"installment_rate":77.4,"id":67696453,"installment_full_cft":219.13,"discount_rate":0,"min_allowed_amount":6,"installment_full_tea":168.35,"labels":["recommended_installment","CFT_219,13%|TEA_168,35%"],"max_allowed_amount":250000,"base_installment_rate":77.4,"installment_rate_collector":["MERCADOPAGO"]}],"issuer":{"default":true,"name":"Visa Argentina S.A.","id":1},"total_financial_cost":null,"min_accreditation_days":0,"max_accreditation_days":2,"merchant_account_id":null,"id":"visa","payment_type_id":"credit_card","accreditation_time":2880,"owner":"site","settings":[{"security_code":{"mode":"mandatory","card_location":"back","length":3},"card_number":{"length":16,"validation":"standard"},"bin":{"pattern":"^4","installments_pattern":"^4","exclusion_pattern":"^(476520|473713|473713|473227|444493|410122|405517|402789|417856|448712|453770|434541|411199|423465|434540|434542|434538|423018|488241|489634|434537|434539|434536|427156|427157|434535|434534|434533|423077|434532|434586|423001|434531|411197|443264|400276|400615|402914|404625|405069|434543|416679|405515|405516|405755|405896|405897|406290|406291|406375|406652|406998|406999|408515|410082|410083|410121|410123|410853|411849|417309|421738|423623|428062|428063|428064|434795|437996|439818|442371|442548|444060|446343|446344|446347|450412|450799|451377|451701|451751|451756|451757|451758|451761|451763|451764|451765|451766|451767|451768|451769|451770|451772|451773|457596|457665|462815|463465|468508|473710|473711|473712|473714|473715|473716|473717|473718|473719|473720|473721|473722|473725|477051|477053|481397|481501|481502|481550|483002|483020|483188|489412|492528|499859|446344|446345|446346|400448)"},"id":67696455}],"thumbnail":"http://img.mlstatic.com/org-img/MP3/API/logos/visa.gif","bins":[],"marketplace":"NONE","deferred_capture":"supported","labels":["recommended_method"],"financing_deals":{"legals":null,"installments":null,"expiration_date":null,"start_date":null,"status":"deactive"},"name":"Visa","site_id":"MLA","processing_mode":"aggregator","additional_info_needed":["cardholder_name","cardholder_identification_type","cardholder_identification_number"],"status":"active"}]}`),
		},
		{
			name:    "trying to reproduce error",
			isEqual: false,
			b1:      []byte(`{"nums":[1,3,3]}`),
			b2:      []byte(`{"nums":[1,1,3]}`),
		},
		{
			name:    "slice with different default things",
			isEqual: false,
			b1:      []byte(`[1,3,3]`),
			b2:      []byte(`[1,1,"casa"]`),
		},
		{
			name:    "slice with different complex things",
			isEqual: false,
			b1:      []byte(`[{"nums": 2}, 3, 3]`),
			b2:      []byte(`[1,1,"casa"]`),
		},
		{
			name:    "slice with different very complex things",
			isEqual: true,
			b1:      []byte(`{"a": [1,"a",{"a":"a","b":"b"}]}`),
			b2:      []byte(`{"a": [1,{"b":"b","a":"a"},"a"]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j1, j2, err := Unmarshal(test.b1, test.b2)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, test.isEqual, Equal(j1, j2))
		})

	}
}
