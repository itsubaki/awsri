package usage

/*
Region returns fullname (ap-northeast-1) from UsageType style (APN1)
https://docs.aws.amazon.com/ja_jp/AmazonS3/latest/dev/aws-usage-report-understand.html
https://docs.aws.amazon.com/ja_jp/general/latest/gr/rande.html
*/
var region = map[string]string{
	"APN1": "ap-northeast-1",
	"APN2": "ap-northeast-2",
	"APS1": "ap-southeast-1",
	"APS2": "ap-southeast-2",
	"APS3": "ap-southeast-3",
	"CAN1": "ca-central-1",
	"EUC1": "eu-central-1",
	"EU":   "eu-west-1",
	"EUW2": "eu-west-2",
	"EUW3": "eu-west-3",
	"SAE1": "sa-east-1",
	"USE1": "us-east-1",
	"USE2": "us-east-2",
	"USW1": "us-west-1",
	"USW2": "us-west-2",
}
