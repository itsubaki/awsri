package reservation

/*
Region returns fullname (ap-northeast-1) from UsageType style (APN1)
https://docs.aws.amazon.com/ja_jp/AmazonS3/latest/dev/aws-usage-report-understand.html
https://docs.aws.amazon.com/ja_jp/general/latest/gr/rande.html
*/
var region = map[string]string{
	"ap-northeast-1": "APN1",
	"ap-northeast-2": "APN2",
	"ap-southeast-1": "APS1",
	"ap-southeast-2": "APS2",
	"ap-southeast-3": "APS3",
	"ca-central-1":   "CAN1",
	"eu-central-1":   "EUC1",
	"eu-west-1":      "EU",
	"eu-west-2":      "EUW2",
	"eu-west-3":      "EUW3",
	"sa-east-1":      "SAE1",
	"us-east-1":      "USE1",
	"us-east-2":      "USE2",
	"us-west-1":      "USW1",
	"us-west-2":      "USW2",
}
