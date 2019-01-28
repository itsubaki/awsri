package ec2

import "testing"

func TestGetPriceAPN1(t *testing.T) {
	_, err := GetPrice("ap-northeast-1")
	if err != nil {
		t.Error(err)
	}
}

/*
https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/region_index.json
{ap-northeast-1 4QTYZ3AAJFVFSW2D APN1-BoxUsage:r3.2xlarge Windows 16 {4QTYZ3AAJFVFSW2D 1.1770000000 7147}}
{ap-northeast-1 Q2ZJCCY86AZE98CP APN1-BoxUsage:m4.16xlarge Windows 128 {Q2ZJCCY86AZE98CP 7.0720000000 48591}}
{ap-northeast-1 HAB58CVXVW943UVQ APN1-BoxUsage:r4.xlarge Windows 8 {HAB58CVXVW943UVQ 0.5040000000 3260}}
{ap-northeast-1 FBB5W5WTFXJSNGPN APN1-BoxUsage:t2.2xlarge Linux 16 {FBB5W5WTFXJSNGPN 0.4864000000 2687}}
{ap-northeast-1 2DCU7R2BA29KKVRK APN1-BoxUsage:c5d.9xlarge Linux 72 {2DCU7R2BA29KKVRK 2.1960000000 12770}}
{ap-northeast-1 T3WJ2QREPM29JX83 APN1-BoxUsage:c5.9xlarge Windows 72 {T3WJ2QREPM29JX83 3.5820000000 25826}}
{ap-northeast-1 N9RJNUBUK9KF2CJD APN1-BoxUsage:m5d.4xlarge Windows 32 {N9RJNUBUK9KF2CJD 1.9040000000 13100}}
{ap-northeast-1 86BX29K37BK2Z2E2 APN1-BoxUsage:c3.xlarge Windows 8 {86BX29K37BK2Z2E2 0.4620000000 2745}}
{ap-northeast-1 8P7XSPEWNQ9EDW27 APN1-BoxUsage:c3.2xlarge Windows 16 {8P7XSPEWNQ9EDW27 0.9250000000 5459}}
{ap-northeast-1 EZCSGZJ8PMXA2QF2 APN1-BoxUsage:i2.xlarge Linux 8 {EZCSGZJ8PMXA2QF2 1.0010000000 4573}}
{ap-northeast-1 JKWD3YB3EB4C3ZXN APN1-BoxUsage:t2.large Windows 4 {JKWD3YB3EB4C3ZXN 0.1496000000 917}}
{ap-northeast-1 3WMXZQY28ZJRFU7T APN1-BoxUsage:t2.medium Windows 2 {3WMXZQY28ZJRFU7T 0.0788000000 494}}
{ap-northeast-1 XBTA6RG7KN7YMKKP APN1-BoxUsage:x1e.8xlarge Linux 64 {XBTA6RG7KN7YMKKP 9.6720000000 48726}}
{ap-northeast-1 ESYFZEM6DQKDVAH7 APN1-BoxUsage:z1d.2xlarge Linux 16 {ESYFZEM6DQKDVAH7 0.9080000000 4739}}
{ap-northeast-1 ERVWZ4V3UBYH4NQH APN1-BoxUsage:t1.micro Linux 0.5 {ERVWZ4V3UBYH4NQH 0.0260000000 138}}
{ap-northeast-1 AG62W8FCAWTS4BU6 APN1-BoxUsage:i2.8xlarge Windows 64 {AG62W8FCAWTS4BU6 8.9030000000 48656}}
{ap-northeast-1 3KYNGSJGJJ4XEFVU APN1-BoxUsage:m5.xlarge Linux 8 {3KYNGSJGJJ4XEFVU 0.2480000000 1412}}
{ap-northeast-1 DDX2JPPMM28BXD7D APN1-BoxUsage:r3.8xlarge Linux 64 {DDX2JPPMM28BXD7D 3.1920000000 17837}}
{ap-northeast-1 HQ8U9TFKGY2BFMZK APN1-BoxUsage:m4.2xlarge Windows 16 {HQ8U9TFKGY2BFMZK 0.8840000000 6074}}
{ap-northeast-1 N8CWR4CSAJWBBE57 APN1-BoxUsage:i3.xlarge Windows 8 {N8CWR4CSAJWBBE57 0.5500000000 3654}}
{ap-northeast-1 Y74PTFH3JBHER34P APN1-BoxUsage:z1d.6xlarge Linux 48 {Y74PTFH3JBHER34P 2.7240000000 14218}}
{ap-northeast-1 HZCFGWJFUGVWKJCT APN1-BoxUsage:m5.large Windows 4 {HZCFGWJFUGVWKJCT 0.2160000000 1512}}
{ap-northeast-1 KXWF29RB7W5NHUYZ APN1-BoxUsage:p3.16xlarge Windows 128 {KXWF29RB7W5NHUYZ 44.8880000000 265963}}
{ap-northeast-1 E5ZC2EJP47JC4Y2A APN1-BoxUsage:x1.16xlarge Linux 128 {E5ZC2EJP47JC4Y2A 9.6710000000 48719}}
{ap-northeast-1 HTNXMK8Z5YHMU737 APN1-BoxUsage:c3.xlarge Linux 8 {HTNXMK8Z5YHMU737 0.2550000000 1505}}
{ap-northeast-1 V6AP4NM9BG9JX6S8 APN1-BoxUsage:d2.8xlarge Windows 64 {V6AP4NM9BG9JX6S8 7.4300000000 31629}}
{ap-northeast-1 FCC4C43QD9KUHD2X APN1-BoxUsage:c5.2xlarge Linux 16 {FCC4C43QD9KUHD2X 0.4280000000 2515}}
{ap-northeast-1 7A24VVDQEZ54FYXU APN1-BoxUsage:d2.2xlarge Linux 16 {7A24VVDQEZ54FYXU 1.6880000000 7760}}
{ap-northeast-1 5YM2V38VPGM8RC7S APN1-BoxUsage:g3.8xlarge Windows 64 {5YM2V38VPGM8RC7S 4.6320000000 32931}}
{ap-northeast-1 JM4YNKP9HA3MAVA2 APN1-BoxUsage:c5.4xlarge Windows 32 {JM4YNKP9HA3MAVA2 1.5920000000 11478}}
{ap-northeast-1 DAFNCUF5NU3EK9D3 APN1-BoxUsage:z1d.xlarge Linux 8 {DAFNCUF5NU3EK9D3 0.4540000000 2370}}
{ap-northeast-1 URZU4GXQC7AT6RE9 APN1-BoxUsage:c1.xlarge Linux 8 {URZU4GXQC7AT6RE9 0.6320000000 3736}}
{ap-northeast-1 XJZSRKDB2CMA66QR APN1-BoxUsage:c5d.2xlarge Windows 16 {XJZSRKDB2CMA66QR 0.8560000000 6061}}
{ap-northeast-1 VTKKFBR2Z5YZ5U2E APN1-BoxUsage:c5d.xlarge Linux 8 {VTKKFBR2Z5YZ5U2E 0.2440000000 1419}}
{ap-northeast-1 CR3WZ76QUCUUVDB3 APN1-BoxUsage:c5d.18xlarge Linux 144 {CR3WZ76QUCUUVDB3 4.3920000000 25540}}
{ap-northeast-1 HHB3QTXEKWVTMUJB APN1-BoxUsage:c5.18xlarge Windows 144 {HHB3QTXEKWVTMUJB 7.1640000000 51652}}
{ap-northeast-1 YJ2E4JTYGN2FMNQM APN1-BoxUsage:cc2.8xlarge Linux 64 {YJ2E4JTYGN2FMNQM 2.3490000000 11097}}
{ap-northeast-1 7TNGFE9SMYR8JYNZ APN1-BoxUsage:m3.large Windows 4 {7TNGFE9SMYR8JYNZ 0.2920000000 1676}}
{ap-northeast-1 EP8EJMA4GKSUCMU6 APN1-BoxUsage:c5.xlarge Linux 8 {EP8EJMA4GKSUCMU6 0.2140000000 1258}}
{ap-northeast-1 KB5V9BJ77S8AV7TK APN1-BoxUsage:c5.large Linux 4 {KB5V9BJ77S8AV7TK 0.1070000000 629}}
{ap-northeast-1 NRPXFBPFDSHUN7HQ APN1-BoxUsage:m5d.4xlarge Linux 32 {NRPXFBPFDSHUN7HQ 1.1680000000 6653}}
{ap-northeast-1 77PTRYZ5MAUP8HU6 APN1-BoxUsage:t3.medium Linux 2 {77PTRYZ5MAUP8HU6 0.0544000000 304}}
{ap-northeast-1 DZXQUR7N3Y4FB3Q4 APN1-BoxUsage:x1e.8xlarge Windows 64 {DZXQUR7N3Y4FB3Q4 11.1440000000 61621}}
{ap-northeast-1 EZQXE684FX263WYR APN1-BoxUsage:cc2.8xlarge Windows 64 {EZQXE684FX263WYR 2.9190000000 15759}}
{ap-northeast-1 FPXP8QM9DMXHP6QP APN1-BoxUsage:c5.9xlarge Linux 72 {FPXP8QM9DMXHP6QP 1.9260000000 11319}}
{ap-northeast-1 AWVVMPYS5HFBJQ2E APN1-BoxUsage:g2.8xlarge Windows 64 {AWVVMPYS5HFBJQ2E 3.8700000000 21171}}
{ap-northeast-1 3P87DYEBFP6F4BQT APN1-BoxUsage:m4.10xlarge Windows 80 {3P87DYEBFP6F4BQT 4.4200000000 30369}}
{ap-northeast-1 K6RT5YHWMXE2Q7NX APN1-BoxUsage:t3.small Windows 1 {K6RT5YHWMXE2Q7NX 0.0456000000 313}}
{ap-northeast-1 36A5UXP5XPUNYGKZ APN1-BoxUsage:x1e.16xlarge Windows 128 {36A5UXP5XPUNYGKZ 22.2880000000 123242}}
{ap-northeast-1 NWAMNVT9VTV64ZM6 APN1-BoxUsage:c1.xlarge Windows 8 {NWAMNVT9VTV64ZM6 1.0320000000 6102}}
{ap-northeast-1 KR2NPP9N7R68MJVX APN1-BoxUsage:c5.xlarge Windows 8 {KR2NPP9N7R68MJVX 0.3980000000 2870}}
{ap-northeast-1 A8RFGEGVBUQQDWD4 APN1-BoxUsage:m4.large Windows 4 {A8RFGEGVBUQQDWD4 0.2210000000 1518}}
{ap-northeast-1 DAPC5MD4ZQN9K67N APN1-BoxUsage:t3.micro Linux 0.5 {DAPC5MD4ZQN9K67N 0.0136000000 76}}
{ap-northeast-1 ZTWDHNAC35QT5ZQZ APN1-BoxUsage:i3.4xlarge Windows 32 {ZTWDHNAC35QT5ZQZ 2.2000000000 14617}}
{ap-northeast-1 6BH4MPBPGY9986DC APN1-BoxUsage:c5d.xlarge Windows 8 {6BH4MPBPGY9986DC 0.4280000000 3031}}
{ap-northeast-1 GNEKD47PUMN4FP4J APN1-BoxUsage:g3.8xlarge Linux 64 {GNEKD47PUMN4FP4J 3.1600000000 20037}}
{ap-northeast-1 2XF9NDPWBAEXYY6S APN1-BoxUsage:x1e.32xlarge Linux 256 {2XF9NDPWBAEXYY6S 38.6880000000 194905}}
{ap-northeast-1 U8JUARJS4SHG5W54 APN1-BoxUsage:r4.xlarge Linux 8 {U8JUARJS4SHG5W54 0.3200000000 1648}}
{ap-northeast-1 UMV7384WFS5N9T5F APN1-BoxUsage:m2.2xlarge Linux 16 {UMV7384WFS5N9T5F 0.5750000000 2231}}
{ap-northeast-1 FGDPWUNAJG8Y9UQR APN1-BoxUsage:m5d.12xlarge Windows 96 {FGDPWUNAJG8Y9UQR 5.7120000000 39300}}
{ap-northeast-1 N72TAGFJTQXEK4GH APN1-BoxUsage:g3.4xlarge Windows 32 {N72TAGFJTQXEK4GH 2.3160000000 16466}}
{ap-northeast-1 DKTG88F6VHT8E4PM APN1-BoxUsage:z1d.12xlarge Windows 96 {DKTG88F6VHT8E4PM 7.6560000000 47778}}
{ap-northeast-1 R5HC3S3BJM43SW7V APN1-BoxUsage:m5.xlarge Windows 8 {R5HC3S3BJM43SW7V 0.4320000000 3024}}
{ap-northeast-1 2GK4E7T3QWMC7QHC APN1-BoxUsage:r3.large Windows 4 {2GK4E7T3QWMC7QHC 0.3000000000 1913}}
{ap-northeast-1 RPSNHYM8M88X8DF5 APN1-BoxUsage:r4.4xlarge Linux 32 {RPSNHYM8M88X8DF5 1.2800000000 6593}}
{ap-northeast-1 6A63VD4RDBFRY4JK APN1-BoxUsage:x1e.16xlarge Linux 128 {6A63VD4RDBFRY4JK 19.3440000000 97453}}
{ap-northeast-1 6TMC6UD2UCCDAMNP APN1-BoxUsage:m1.large Linux 4 {6TMC6UD2UCCDAMNP 0.2430000000 1128}}
{ap-northeast-1 SR8U7JF2N76XJGXV APN1-BoxUsage:r3.8xlarge Windows 64 {SR8U7JF2N76XJGXV 4.2500000000 19623}}
{ap-northeast-1 BHVGUPJAB9MQKXCD APN1-BoxUsage:p3.8xlarge Windows 64 {BHVGUPJAB9MQKXCD 22.4440000000 132982}}
{ap-northeast-1 CMZ72PGCABYSQMPU APN1-BoxUsage:r4.2xlarge Windows 16 {CMZ72PGCABYSQMPU 1.0080000000 6520}}
{ap-northeast-1 7VD689A8AZD59K29 APN1-BoxUsage:z1d.xlarge Windows 8 {7VD689A8AZD59K29 0.6380000000 3982}}
{ap-northeast-1 VQP5R9BSB4AJR4CR APN1-BoxUsage:g2.2xlarge Windows 16 {VQP5R9BSB4AJR4CR 1.0100000000 6098}}
{ap-northeast-1 4J9CTGS8W7B7SFTP APN1-BoxUsage:t3.nano Windows 0.25 {4J9CTGS8W7B7SFTP 0.0114000000 78}}
{ap-northeast-1 SA9UW2TC8EGBE7NW APN1-BoxUsage:i3.2xlarge Linux 16 {SA9UW2TC8EGBE7NW 0.7320000000 4085}}
{ap-northeast-1 5TNGT5CHGYHEM47D APN1-BoxUsage:x1.32xlarge Windows 256 {5TNGT5CHGYHEM47D 25.2290000000 149017}}
{ap-northeast-1 R49K2Y7KZ6527C35 APN1-BoxUsage:x1e.2xlarge Linux 16 {R49K2Y7KZ6527C35 2.4180000000 12182}}
{ap-northeast-1 6MPP3K5KV3SKV8Q9 APN1-BoxUsage:m3.2xlarge Windows 16 {6MPP3K5KV3SKV8Q9 1.1660000000 6676}}
{ap-northeast-1 BURRP7APFUCZFSZK APN1-BoxUsage:m4.xlarge Linux 8 {BURRP7APFUCZFSZK 0.2580000000 1425}}
{ap-northeast-1 YRQT9KKZFWN77DXT APN1-BoxUsage:t2.small Windows 1 {YRQT9KKZFWN77DXT 0.0396000000 249}}
{ap-northeast-1 MPQ46C9QZE4FXJ4R APN1-BoxUsage:m5d.2xlarge Windows 16 {MPQ46C9QZE4FXJ4R 0.9520000000 6550}}
{ap-northeast-1 6JP9PA73B58NZHUN APN1-BoxUsage:d2.4xlarge Linux 32 {6JP9PA73B58NZHUN 3.3760000000 15520}}
{ap-northeast-1 E2747P4S267E56HR APN1-BoxUsage:t2.2xlarge Windows 16 {E2747P4S267E56HR 0.5484000000 3230}}
{ap-northeast-1 TZG97WFA265PFBMW APN1-BoxUsage:t3.nano Linux 0.25 {TZG97WFA265PFBMW 0.0068000000 38}}
{ap-northeast-1 NZHXGSV3KSMEQT45 APN1-BoxUsage:d2.xlarge Windows 8 {NZHXGSV3KSMEQT45 0.9750000000 4394}}
{ap-northeast-1 Y5V29EJZ67GJPZMS APN1-BoxUsage:m2.xlarge Windows 8 {Y5V29EJZ67GJPZMS 0.3520000000 2288}}
{ap-northeast-1 5YHRKH4DFNQ4XWHZ APN1-BoxUsage:i3.4xlarge Linux 32 {5YHRKH4DFNQ4XWHZ 1.4640000000 8169}}
{ap-northeast-1 3CJSUV6SJ9TG2J2F APN1-BoxUsage:c5d.4xlarge Linux 32 {3CJSUV6SJ9TG2J2F 0.9760000000 5676}}
{ap-northeast-1 XSS7VX8UVDYEA2F4 APN1-BoxUsage:m3.medium Windows 2 {XSS7VX8UVDYEA2F4 0.1460000000 834}}
{ap-northeast-1 R64MZQ7UD9UPA4ZW APN1-BoxUsage:r3.xlarge Windows 8 {R64MZQ7UD9UPA4ZW 0.5990000000 3707}}
{ap-northeast-1 UJB452HW969DQZFB APN1-BoxUsage:c4.xlarge Linux 8 {UJB452HW969DQZFB 0.2520000000 1477}}
{ap-northeast-1 ZAEMVYU798AFMXPQ APN1-BoxUsage:z1d.large Linux 4 {ZAEMVYU798AFMXPQ 0.2270000000 1185}}
{ap-northeast-1 46WBDTB7HP3F5A8Z APN1-BoxUsage Windows 1 {46WBDTB7HP3F5A8Z 0.0880000000 483}}
{ap-northeast-1 TPZKPCAQBPBS7CF8 APN1-BoxUsage:i3.8xlarge Linux 64 {TPZKPCAQBPBS7CF8 2.9280000000 16339}}
{ap-northeast-1 AG6SXPB68M2AKYAV APN1-BoxUsage:c5d.large Linux 4 {AG6SXPB68M2AKYAV 0.1220000000 709}}
{ap-northeast-1 3AHMWXGBEX6HSCMQ APN1-BoxUsage:t3.large Windows 4 {3AHMWXGBEX6HSCMQ 0.1364000000 850}}
{ap-northeast-1 4GHFAT5CNS2FEKB2 APN1-BoxUsage:m4.2xlarge Linux 16 {4GHFAT5CNS2FEKB2 0.5160000000 2850}}
{ap-northeast-1 C8EY297E42DNDAKF APN1-BoxUsage:c4.large Windows 4 {C8EY297E42DNDAKF 0.2180000000 1481}}
{ap-northeast-1 S3Y9ND5T3E2PMPBK APN1-BoxUsage:r4.8xlarge Windows 64 {S3Y9ND5T3E2PMPBK 4.0320000000 26081}}
{ap-northeast-1 9E3A3M4ACM2BKD5B APN1-BoxUsage:m5d.xlarge Windows 8 {9E3A3M4ACM2BKD5B 0.4760000000 3275}}
{ap-northeast-1 9NSP3EV3G593P35X APN1-BoxUsage:t2.micro Linux 0.5 {9NSP3EV3G593P35X 0.0152000000 84}}
{ap-northeast-1 AY6XZ64M22QQJCXE APN1-BoxUsage:m3.large Linux 4 {AY6XZ64M22QQJCXE 0.1930000000 950}}
{ap-northeast-1 MZQS5X6TXEPY4CG9 APN1-BoxUsage:m2.2xlarge Windows 16 {MZQS5X6TXEPY4CG9 0.7050000000 4578}}
{ap-northeast-1 UQNVRJQ7NEHJMYXU APN1-BoxUsage:m2.4xlarge Windows 32 {UQNVRJQ7NEHJMYXU 1.4100000000 9155}}
{ap-northeast-1 GJHUHQSUU37VCQ5A APN1-BoxUsage:r3.xlarge Linux 8 {GJHUHQSUU37VCQ5A 0.3990000000 2230}}
{ap-northeast-1 UUEJKASUECCGFNVR APN1-BoxUsage:g3.16xlarge Windows 128 {UUEJKASUECCGFNVR 9.2640000000 65863}}
{ap-northeast-1 E5MWNHYU3BAVZCRP APN1-BoxUsage:c4.4xlarge Linux 32 {E5MWNHYU3BAVZCRP 1.0080000000 5915}}
{ap-northeast-1 Q4QTSF7H37JFW9ER APN1-BoxUsage:c3.large Linux 4 {Q4QTSF7H37JFW9ER 0.1280000000 753}}
{ap-northeast-1 J85A5X44TT267EH8 APN1-BoxUsage:m3.xlarge Linux 8 {J85A5X44TT267EH8 0.3850000000 1909}}
{ap-northeast-1 THKGMUKKFXV9CKUW APN1-BoxUsage:p2.16xlarge Linux 128 {THKGMUKKFXV9CKUW 24.6720000000 141314}}
{ap-northeast-1 X5M9TWU789QSW4YG APN1-BoxUsage:z1d.3xlarge Windows 24 {X5M9TWU789QSW4YG 1.9140000000 11945}}
{ap-northeast-1 8VN3HX7E6Z8JVZ78 APN1-BoxUsage:t3.xlarge Linux 8 {8VN3HX7E6Z8JVZ78 0.2176000000 1216}}
{ap-northeast-1 6M27QQ6HYCNA5KGA APN1-BoxUsage:m3.medium Linux 2 {6M27QQ6HYCNA5KGA 0.0960000000 471}}
{ap-northeast-1 BYV8H4R4VJNAH42Q APN1-BoxUsage:r3.4xlarge Linux 32 {BYV8H4R4VJNAH42Q 1.5960000000 8919}}
{ap-northeast-1 EES7AKKAGBY837K3 APN1-BoxUsage:p2.xlarge Windows 8 {EES7AKKAGBY837K3 1.7260000000 10444}}
{ap-northeast-1 SFUYMZQAV538QWXK APN1-BoxUsage:c5.18xlarge Linux 144 {SFUYMZQAV538QWXK 3.8520000000 22639}}
{ap-northeast-1 PCB5ARVZ6TNS7A96 APN1-BoxUsage:m3.2xlarge Linux 16 {PCB5ARVZ6TNS7A96 0.7700000000 3819}}
{ap-northeast-1 RCJ9VNKFJCUCGU3W APN1-BoxUsage:p2.8xlarge Linux 64 {RCJ9VNKFJCUCGU3W 12.3360000000 70657}}
{ap-northeast-1 E5R2MT5HUM6D2JM3 APN1-BoxUsage:i3.16xlarge Windows 128 {E5R2MT5HUM6D2JM3 8.8000000000 58467}}
{ap-northeast-1 PN5FSE6BV774Z4CN APN1-BoxUsage:x1e.2xlarge Windows 16 {PN5FSE6BV774Z4CN 2.7860000000 15405}}
{ap-northeast-1 9UBMZYZ6SXZ5JQGV APN1-BoxUsage:i3.xlarge Linux 8 {9UBMZYZ6SXZ5JQGV 0.3660000000 2042}}
{ap-northeast-1 KWUMHT4YYUHYMCEV APN1-BoxUsage:m5d.24xlarge Linux 192 {KWUMHT4YYUHYMCEV 7.0080000000 39915}}
{ap-northeast-1 J466PGV3UVD449YQ APN1-BoxUsage:r4.4xlarge Windows 32 {J466PGV3UVD449YQ 2.0160000000 13040}}
{ap-northeast-1 E9E26HA2R4KDRBC2 APN1-BoxUsage:z1d.6xlarge Windows 48 {E9E26HA2R4KDRBC2 3.8280000000 23889}}
{ap-northeast-1 8JAU22FMXQBNYDG6 APN1-BoxUsage:c5d.large Windows 4 {8JAU22FMXQBNYDG6 0.2140000000 1515}}
{ap-northeast-1 F7XCNBBYFKX42QF3 APN1-BoxUsage:t2.nano Linux 0.25 {F7XCNBBYFKX42QF3 0.0076000000 42}}
{ap-northeast-1 72F7QWUGHW3NB9PS APN1-BoxUsage:t3.xlarge Windows 8 {72F7QWUGHW3NB9PS 0.2912000000 1860}}
{ap-northeast-1 G6G6ZNFBYMW2V8BH APN1-BoxUsage:m2.xlarge Linux 8 {G6G6ZNFBYMW2V8BH 0.2870000000 1110}}
{ap-northeast-1 AKQ89V8E78T6H534 APN1-BoxUsage Linux 1 {AKQ89V8E78T6H534 0.0610000000 278}}
{ap-northeast-1 8YT9TXPXJ6KCKS3Z APN1-BoxUsage:m5d.2xlarge Linux 16 {8YT9TXPXJ6KCKS3Z 0.5840000000 3326}}
{ap-northeast-1 BQQUCAM9PFTSUNQX APN1-BoxUsage:m2.4xlarge Linux 32 {BQQUCAM9PFTSUNQX 1.1500000000 4444}}
{ap-northeast-1 BZ22G7HWGF9WW7EF APN1-BoxUsage:t3.2xlarge Windows 16 {BZ22G7HWGF9WW7EF 0.5824000000 3721}}
{ap-northeast-1 WHR37BGS9EYEPVKT APN1-BoxUsage:p3.2xlarge Linux 16 {WHR37BGS9EYEPVKT 5.2430000000 30022}}
{ap-northeast-1 BDT5TR37G2FCDKV2 APN1-BoxUsage:t3.large Linux 4 {BDT5TR37G2FCDKV2 0.1088000000 608}}
{ap-northeast-1 E3PNQB9BY2K5GSXP APN1-BoxUsage:t1.micro Windows 0.5 {E3PNQB9BY2K5GSXP 0.0330000000 215}}
{ap-northeast-1 26QPRVEP3SD3YZR7 APN1-BoxUsage:i2.4xlarge Windows 32 {26QPRVEP3SD3YZR7 4.4510000000 24328}}
{ap-northeast-1 2XSTZ4JVJ7XG7Y23 APN1-BoxUsage:i3.2xlarge Windows 16 {2XSTZ4JVJ7XG7Y23 1.1000000000 7308}}
{ap-northeast-1 9KMZWGZVTXKAQXNM APN1-BoxUsage:r3.2xlarge Linux 16 {9KMZWGZVTXKAQXNM 0.7980000000 4459}}
{ap-northeast-1 4REMK3MMXCZ55ZX3 APN1-BoxUsage:i2.8xlarge Linux 64 {4REMK3MMXCZ55ZX3 8.0040000000 36568}}
{ap-northeast-1 6KTQUBH8Y7R2NVBW APN1-BoxUsage:m5d.12xlarge Linux 96 {6KTQUBH8Y7R2NVBW 3.5040000000 19958}}
{ap-northeast-1 S4EWKNDHYM7FSPG6 APN1-BoxUsage:c5.4xlarge Linux 32 {S4EWKNDHYM7FSPG6 0.8560000000 5031}}
{ap-northeast-1 W5295P7RR98PKJFH APN1-BoxUsage:m1.xlarge Windows 8 {W5295P7RR98PKJFH 0.7060000000 3887}}
{ap-northeast-1 G55JJ7CXZ5E2QE8H APN1-BoxUsage:r4.large Linux 4 {G55JJ7CXZ5E2QE8H 0.1600000000 824}}
{ap-northeast-1 2JSMK4YRVSAHV4RW APN1-BoxUsage:m4.16xlarge Linux 128 {2JSMK4YRVSAHV4RW 4.1280000000 22801}}
{ap-northeast-1 RNTGYXD2G78F5ZVB APN1-BoxUsage:z1d.large Windows 4 {RNTGYXD2G78F5ZVB 0.3190000000 1991}}
{ap-northeast-1 PAHJQFQFMDZ4DMFV APN1-BoxUsage:z1d.2xlarge Windows 16 {PAHJQFQFMDZ4DMFV 1.2760000000 7963}}
{ap-northeast-1 7MJR4RD25PP93ENY APN1-BoxUsage:t3.small Linux 1 {7MJR4RD25PP93ENY 0.0272000000 152}}
{ap-northeast-1 4BJPFU3PAZJ4AKMM APN1-BoxUsage:m1.xlarge Linux 8 {4BJPFU3PAZJ4AKMM 0.4860000000 2239}}
{ap-northeast-1 MKVJ4C4XUPQ3657J APN1-BoxUsage:i3.16xlarge Linux 128 {MKVJ4C4XUPQ3657J 5.8560000000 32677}}
{ap-northeast-1 CX79CXQ739SJJJ6P APN1-BoxUsage:t3.2xlarge Linux 16 {CX79CXQ739SJJJ6P 0.4352000000 2431}}
{ap-northeast-1 68DZA4NDNQQUJY3E APN1-BoxUsage:d2.2xlarge Windows 16 {68DZA4NDNQQUJY3E 1.9090000000 8352}}
{ap-northeast-1 MYX88QW5HYQW9KS4 APN1-BoxUsage:r4.8xlarge Linux 64 {MYX88QW5HYQW9KS4 2.5600000000 13186}}
{ap-northeast-1 RAJXW44F9F9QMXYK APN1-BoxUsage:t3.medium Windows 2 {RAJXW44F9F9QMXYK 0.0728000000 465}}
{ap-northeast-1 Q85F79PK8VHHZT6X APN1-BoxUsage:c4.2xlarge Linux 16 {Q85F79PK8VHHZT6X 0.5040000000 2962}}
{ap-northeast-1 ST43JST6BQDZGU9F APN1-BoxUsage:m1.medium Windows 2 {ST43JST6BQDZGU9F 0.1770000000 978}}
{ap-northeast-1 5VR9TEU2W5PW2JAN APN1-BoxUsage:c4.xlarge Windows 8 {5VR9TEU2W5PW2JAN 0.4360000000 2954}}
{ap-northeast-1 5R5TRXGCZSGH3HQ4 APN1-BoxUsage:d2.4xlarge Windows 32 {5R5TRXGCZSGH3HQ4 3.6780000000 16111}}
{ap-northeast-1 N4H3ZWBACT9JMMZS APN1-BoxUsage:x1e.xlarge Windows 8 {N4H3ZWBACT9JMMZS 1.3930000000 7703}}
{ap-northeast-1 GSXGEKMNEVYDGNSF APN1-BoxUsage:c3.4xlarge Windows 32 {GSXGEKMNEVYDGNSF 1.8490000000 10927}}
{ap-northeast-1 APGTJ6NBJA89PCKZ APN1-BoxUsage:i3.large Linux 4 {APGTJ6NBJA89PCKZ 0.1830000000 1021}}
{ap-northeast-1 UDHFRPKESN82BQYQ APN1-BoxUsage:g3.16xlarge Linux 128 {UDHFRPKESN82BQYQ 6.3200000000 40073}}
{ap-northeast-1 VWK982QWFHEZBQPG APN1-BoxUsage:p2.16xlarge Windows 128 {VWK982QWFHEZBQPG 27.6160000000 167103}}
{ap-northeast-1 AJ7VEHQACPH6WCUW APN1-BoxUsage:c5.2xlarge Windows 16 {AJ7VEHQACPH6WCUW 0.7960000000 5739}}
{ap-northeast-1 5ETHUFF38RM2SEXX APN1-BoxUsage:c1.medium Windows 2 {5ETHUFF38RM2SEXX 0.2580000000 1527}}
{ap-northeast-1 PNWEYBEMH3JBNDGR APN1-BoxUsage:x1e.32xlarge Windows 256 {PNWEYBEMH3JBNDGR 44.5760000000 246484}}
{ap-northeast-1 REQJXNHA97H376TZ APN1-BoxUsage:m5d.xlarge Linux 8 {REQJXNHA97H376TZ 0.2920000000 1663}}
{ap-northeast-1 WZ88G3NJFBWUCY5P APN1-BoxUsage:m4.xlarge Windows 8 {WZ88G3NJFBWUCY5P 0.4420000000 3037}}
{ap-northeast-1 QD2F9REPXN5QEMC7 APN1-BoxUsage:c4.2xlarge Windows 16 {QD2F9REPXN5QEMC7 0.8720000000 5925}}
{ap-northeast-1 Y6EKUV29QTSRMM7Y APN1-BoxUsage:m5d.large Linux 4 {Y6EKUV29QTSRMM7Y 0.1460000000 832}}
{ap-northeast-1 FE9CPXTKGGW59Q7V APN1-BoxUsage:x1e.xlarge Linux 8 {FE9CPXTKGGW59Q7V 1.2090000000 6091}}
{ap-northeast-1 N75CY8HAQEC34CR6 APN1-BoxUsage:m5.4xlarge Windows 32 {N75CY8HAQEC34CR6 1.7280000000 12097}}
{ap-northeast-1 5AEKU56K9NUEPTXT APN1-BoxUsage:c4.8xlarge Windows 64 {5AEKU56K9NUEPTXT 3.6720000000 23725}}
{ap-northeast-1 USE46NRHP4S7UT6J APN1-BoxUsage:m5.large Linux 4 {USE46NRHP4S7UT6J 0.1240000000 706}}
{ap-northeast-1 B3RQCPARF6B7RHW5 APN1-BoxUsage:m5.2xlarge Windows 16 {B3RQCPARF6B7RHW5 0.8640000000 6048}}
{ap-northeast-1 5RFFBCKXAV4J29B7 APN1-BoxUsage:m5.24xlarge Linux 192 {5RFFBCKXAV4J29B7 5.9520000000 33896}}
{ap-northeast-1 R3WKUHQUFW8Q9SAW APN1-BoxUsage:r4.16xlarge Windows 128 {R3WKUHQUFW8Q9SAW 8.0640000000 52162}}
{ap-northeast-1 YR67H6NVBRN37HRZ APN1-BoxUsage:c3.2xlarge Linux 16 {YR67H6NVBRN37HRZ 0.5110000000 3012}}
{ap-northeast-1 PZSZYGTCM8M3CAC3 APN1-BoxUsage:p3.2xlarge Windows 16 {PZSZYGTCM8M3CAC3 5.6110000000 33245}}
{ap-northeast-1 YPF8479KWUC4WWRC APN1-BoxUsage:t2.xlarge Windows 8 {YPF8479KWUC4WWRC 0.2842000000 1702}}
{ap-northeast-1 JDHFUYQU2RVEVNEX APN1-BoxUsage:z1d.3xlarge Linux 24 {JDHFUYQU2RVEVNEX 1.3620000000 7109}}
{ap-northeast-1 9XQ29FEKBNY683HQ APN1-BoxUsage:p2.8xlarge Windows 64 {9XQ29FEKBNY683HQ 13.8080000000 83552}}
{ap-northeast-1 CV9PY4ZHFE7HDJKY APN1-BoxUsage:m5d.24xlarge Windows 192 {CV9PY4ZHFE7HDJKY 11.4240000000 78599}}
{ap-northeast-1 ZV2DS4C98AB8SS7J APN1-BoxUsage:t2.xlarge Linux 8 {ZV2DS4C98AB8SS7J 0.2432000000 1343}}
{ap-northeast-1 7MYWT7Y96UT3NJ2D APN1-BoxUsage:m4.large Linux 4 {7MYWT7Y96UT3NJ2D 0.1290000000 713}}
{ap-northeast-1 FBUWUPNC8FXRUS5W APN1-BoxUsage:i2.4xlarge Linux 32 {FBUWUPNC8FXRUS5W 4.0020000000 18275}}
{ap-northeast-1 SD97GUGBCUND24YK APN1-BoxUsage:m5.4xlarge Linux 32 {SD97GUGBCUND24YK 0.9920000000 5649}}
{ap-northeast-1 T7CGRZ4XENPHVK6D APN1-BoxUsage:p2.xlarge Linux 8 {T7CGRZ4XENPHVK6D 1.5420000000 8832}}
{ap-northeast-1 F2RRJYX33EGMBSFR APN1-BoxUsage:m1.medium Linux 2 {F2RRJYX33EGMBSFR 0.1220000000 565}}
{ap-northeast-1 KM8DYQWHEC32CGGX APN1-BoxUsage:i2.2xlarge Linux 16 {KM8DYQWHEC32CGGX 2.0010000000 9146}}
{ap-northeast-1 2NK92W5SRKRY46GS APN1-BoxUsage:p3.16xlarge Linux 128 {2NK92W5SRKRY46GS 41.9440000000 240174}}
{ap-northeast-1 U9EUE7H4E7G5TZN2 APN1-BoxUsage:z1d.12xlarge Linux 96 {U9EUE7H4E7G5TZN2 5.4480000000 28436}}
{ap-northeast-1 CFGR4ZZWZE9QJFME APN1-BoxUsage:m5.12xlarge Windows 96 {CFGR4ZZWZE9QJFME 5.1840000000 36290}}
{ap-northeast-1 SE3ZWED6JFUSH56R APN1-BoxUsage:i2.xlarge Windows 8 {SE3ZWED6JFUSH56R 1.1130000000 6085}}
{ap-northeast-1 SKJWZ9DAN9BMXCEN APN1-BoxUsage:x1.16xlarge Windows 128 {SKJWZ9DAN9BMXCEN 12.6150000000 74508}}
{ap-northeast-1 DVHBU5E2WXYKER69 APN1-BoxUsage:r4.large Windows 4 {DVHBU5E2WXYKER69 0.2520000000 1630}}
{ap-northeast-1 XEQVYBYPC2ZXKZ9H APN1-BoxUsage:m1.large Windows 4 {XEQVYBYPC2ZXKZ9H 0.3530000000 1944}}
{ap-northeast-1 SHF9TA3MCU6W2BRA APN1-BoxUsage:c5d.2xlarge Linux 16 {SHF9TA3MCU6W2BRA 0.4880000000 2838}}
{ap-northeast-1 6J4947QU7ZJUKHRB APN1-BoxUsage:c5d.4xlarge Windows 32 {6J4947QU7ZJUKHRB 1.7120000000 12123}}
{ap-northeast-1 QXQ5UG8SEH8X3RN3 APN1-BoxUsage:t3.micro Windows 0.5 {QXQ5UG8SEH8X3RN3 0.0228000000 157}}
{ap-northeast-1 MFRDGSP29SNABK4D APN1-BoxUsage:c4.4xlarge Windows 32 {MFRDGSP29SNABK4D 1.7440000000 11867}}
{ap-northeast-1 C8MHVPEYQG6UHPS4 APN1-BoxUsage:r4.2xlarge Linux 16 {C8MHVPEYQG6UHPS4 0.6400000000 3297}}
{ap-northeast-1 77K4UJJUNGQ6UXXN APN1-BoxUsage:g2.2xlarge Linux 16 {77K4UJJUNGQ6UXXN 0.8980000000 5140}}
{ap-northeast-1 VWWQ7S3GZ9J8TF77 APN1-BoxUsage:d2.xlarge Linux 8 {VWWQ7S3GZ9J8TF77 0.8440000000 3880}}
{ap-northeast-1 EBSRPCDMGT2V87YG APN1-BoxUsage:c3.8xlarge Windows 64 {EBSRPCDMGT2V87YG 3.6990000000 21849}}
{ap-northeast-1 8V2J26TGMFGDSDKE APN1-BoxUsage:m5.24xlarge Windows 192 {8V2J26TGMFGDSDKE 10.3680000000 72581}}
{ap-northeast-1 XJ88E6MSR3AYHFXA APN1-BoxUsage:c3.4xlarge Linux 32 {XJ88E6MSR3AYHFXA 1.0210000000 6016}}
{ap-northeast-1 Q5HVB8NUA7UMHV4Y APN1-BoxUsage:t2.large Linux 4 {Q5HVB8NUA7UMHV4Y 0.1216000000 672}}
{ap-northeast-1 GCVKVN4MVRGS7UK3 APN1-BoxUsage:t2.nano Windows 0.25 {GCVKVN4MVRGS7UK3 0.0099000000 62}}
{ap-northeast-1 N6SGMGNN8CA3TG6Q APN1-BoxUsage:g3.4xlarge Linux 32 {N6SGMGNN8CA3TG6Q 1.5800000000 10018}}
{ap-northeast-1 YUYNTU8AZ9MKK68V APN1-BoxUsage:t2.small Linux 1 {YUYNTU8AZ9MKK68V 0.0304000000 168}}
{ap-northeast-1 PCNBVATW49APFGZQ APN1-BoxUsage:c4.8xlarge Linux 64 {PCNBVATW49APFGZQ 2.0160000000 11830}}
{ap-northeast-1 JTQKHD7ZTEEM4DC5 APN1-BoxUsage:m4.10xlarge Linux 80 {JTQKHD7ZTEEM4DC5 2.5800000000 14251}}
{ap-northeast-1 424M8CDE42M2JPGD APN1-BoxUsage:r3.4xlarge Windows 32 {424M8CDE42M2JPGD 2.2760000000 12586}}
{ap-northeast-1 YCYU3NQCQRYQ2TU7 APN1-BoxUsage:r3.large Linux 4 {YCYU3NQCQRYQ2TU7 0.2000000000 1114}}
{ap-northeast-1 EWRM596KUQ2YH8ER APN1-BoxUsage:r4.16xlarge Linux 128 {EWRM596KUQ2YH8ER 5.1200000000 26373}}
{ap-northeast-1 WEMS88BYNHHUKWC8 APN1-BoxUsage:m5.2xlarge Linux 16 {WEMS88BYNHHUKWC8 0.4960000000 2825}}
{ap-northeast-1 CTK39QJHQN4CZ3PC APN1-BoxUsage:g2.8xlarge Linux 64 {CTK39QJHQN4CZ3PC 3.5920000000 20560}}
{ap-northeast-1 MJ7YVW9J2WD856AC APN1-BoxUsage:x1.32xlarge Linux 256 {MJ7YVW9J2WD856AC 19.3410000000 97438}}
{ap-northeast-1 WCVRN7FCX9MS9SAD APN1-BoxUsage:i3.8xlarge Windows 64 {WCVRN7FCX9MS9SAD 4.4000000000 29233}}
{ap-northeast-1 E3J2T7B8EQDFXWDR APN1-BoxUsage:c3.8xlarge Linux 64 {E3J2T7B8EQDFXWDR 2.0430000000 12032}}
{ap-northeast-1 XU2NYYPCRTK4T7CN APN1-BoxUsage:m4.4xlarge Linux 32 {XU2NYYPCRTK4T7CN 1.0320000000 5700}}
{ap-northeast-1 PTSCWYT4DGMHMSYG APN1-BoxUsage:c1.medium Linux 2 {PTSCWYT4DGMHMSYG 0.1580000000 932}}
{ap-northeast-1 6986XC33S6FFMJGG APN1-BoxUsage:m5.12xlarge Linux 96 {6986XC33S6FFMJGG 2.9760000000 16948}}
{ap-northeast-1 GHWKPR98HVZHT4FD APN1-BoxUsage:c5.large Windows 4 {GHWKPR98HVZHT4FD 0.1990000000 1435}}
{ap-northeast-1 KNVQZWZRBTHCFMS5 APN1-BoxUsage:p3.8xlarge Linux 64 {KNVQZWZRBTHCFMS5 20.9720000000 120087}}
{ap-northeast-1 3REN7JCBCCW7XYHC APN1-BoxUsage:c5d.18xlarge Windows 144 {3REN7JCBCCW7XYHC 7.7040000000 54553}}
{ap-northeast-1 FFMZ4PMY5YBPUU9F APN1-BoxUsage:m3.xlarge Windows 8 {FFMZ4PMY5YBPUU9F 0.5830000000 3342}}
{ap-northeast-1 4VCVKYYH4EGXEE23 APN1-BoxUsage:c3.large Windows 4 {4VCVKYYH4EGXEE23 0.2310000000 1369}}
{ap-northeast-1 PHXN6XU2RQN4NYNJ APN1-BoxUsage:x1e.4xlarge Windows 32 {PHXN6XU2RQN4NYNJ 5.5720000000 30811}}
{ap-northeast-1 4CSYQ7HB95TUZPRR APN1-BoxUsage:m5d.large Windows 4 {4CSYQ7HB95TUZPRR 0.2380000000 1637}}
{ap-northeast-1 KA565JRTVNZB5VF2 APN1-BoxUsage:x1e.4xlarge Linux 32 {KA565JRTVNZB5VF2 4.8360000000 24363}}
{ap-northeast-1 H2HQURMVXXGUYJ8G APN1-BoxUsage:i2.2xlarge Windows 16 {H2HQURMVXXGUYJ8G 2.2260000000 12177}}
{ap-northeast-1 H89SS44PSUZRUVP9 APN1-BoxUsage:c5d.9xlarge Windows 72 {H89SS44PSUZRUVP9 3.8520000000 27277}}
{ap-northeast-1 R2NRCU4EJAHR98QB APN1-BoxUsage:t2.micro Windows 0.5 {R2NRCU4EJAHR98QB 0.0198000000 124}}
{ap-northeast-1 XDDRDMN5QWVPB9FG APN1-BoxUsage:i3.large Windows 4 {XDDRDMN5QWVPB9FG 0.2750000000 1827}}
{ap-northeast-1 E6F66FZ47YZNXAJ2 APN1-BoxUsage:t2.medium Linux 2 {E6F66FZ47YZNXAJ2 0.0608000000 336}}
{ap-northeast-1 FEBPNB3R5SDBR8DM APN1-BoxUsage:m4.4xlarge Windows 32 {FEBPNB3R5SDBR8DM 1.7680000000 12148}}
{ap-northeast-1 PSF2TQK8WMUGUPYK APN1-BoxUsage:d2.8xlarge Linux 64 {PSF2TQK8WMUGUPYK 6.7520000000 31040}}
{ap-northeast-1 7KXQBZSKETPTG6QZ APN1-BoxUsage:c4.large Linux 4 {7KXQBZSKETPTG6QZ 0.1260000000 738}}
*/
