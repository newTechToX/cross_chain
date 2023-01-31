package etherscan

import (
	"app/utils"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/time/rate"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
}

func TestGetLogs(t *testing.T) {
	p := NewEtherScanProvider("https://api.etherscan.io/", []string{"Y5CIXMXJ23Y6H8JSRAUQ5T8SMT2VV9W17Z", "4RYCK1WU1W2NBCGDNVEV36GHSZTF6CGW2M"}, "http://192.168.3.59:10809", 20)
	ret, err := p.GetLogs([]string{"0xbd3531da5cf5857e7cfaa92426877b022e612cf8"}, []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}, 12878196, 12878196)
	utils.PrintPretty(ret)
	fmt.Println(len(ret), err)
}

func TestFirstCalled(t *testing.T) {
	p := NewEtherScanProvider("https://api.etherscan.io/", []string{"Y5CIXMXJ23Y6H8JSRAUQ5T8SMT2VV9W17Z", "4RYCK1WU1W2NBCGDNVEV36GHSZTF6CGW2M"}, "http://192.168.3.59:10809", 20)
	fmt.Println(p.GetContractFirstInvocation("0x8765b1a0eb57ca49be7eacd35b24a574d0203656"))
}

func TestOptiscanFirstCalled(t *testing.T) {
	p := NewEtherScanProvider("https://api-optimistic.etherscan.io/", []string{"TX5FYFU9QWEMCQ9UGP865H74VTBWEBVW8X", "TV4RKAHHUXRVKYDJJ1ZPDXR75QJGG75WRB"}, "http://192.168.3.59:10809", 20)
	fmt.Println(p.GetContractFirstInvocation("0xe7798f023fc62146e8aa1b36da45fb70855a77ea"))
}

func TestArbiscanFirstCalled(t *testing.T) {
	p := NewEtherScanProvider("https://api.arbiscan.io/", []string{"NHQD5YRM8PMU76UB6GX7UNKQC7W48E4Y8M", "ZUIBQZJVP5W5RCTU8ZVWPHSWMGTMHR5ZHD"}, "http://192.168.3.59:10809", 20)
	fmt.Println(p.GetContractFirstInvocation("0x6c68eb45d5c2019136c8362cc928fb4f13f5385d"))
}

func TestGetCall(t *testing.T) {
	p := NewEtherScanProvider("https://api.etherscan.io/", []string{"Y5CIXMXJ23Y6H8JSRAUQ5T8SMT2VV9W17Z", "4RYCK1WU1W2NBCGDNVEV36GHSZTF6CGW2M"}, "http://192.168.3.59:10809", 20)
	ret, err := p.GetCalls([]string{"0x5006192340d83bfa47ee2f28edd0fd16a56d5b5e"}, []string{"transfer"}, 15826384, 15826384)
	fmt.Println(len(ret), err)
	utils.PrintPretty(ret)
}

func TestLatest(t *testing.T) {
	p := NewEtherScanProvider("https://api-optimistic.etherscan.io/", []string{"GU95IB2QWKF5THUW5QDZ9WGI8FRDCFPFXA"}, "http://192.168.3.59:10809", 20)
	fmt.Println(p.LatestNumber())
}

func TestNum(t *testing.T) {
	fmt.Println(utils.EtherScanMaxResult * 0.8)
}

func TestContains(t *testing.T) {
	if strings.Contains("asd aaa", "") {
		fmt.Println(1)
	}
}

func TestEtherscanRate(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	keys := []string{"SHR5J5UQC5JQ2GBPSWTBESK94GUVDX38ER", "PWB8MP823SWFPZF3F84NA6N6N7A9G2S68P", "QRG3QFE7RAPY4BX666MVPJ4AFI56QQHIQE", "2VJXZ2GN1WMH22V8R7KB62MJNFQGQSGNNH", "Q7APJ9KF4R63FVBKSC3WZPKR25I4RCAIEY", "UE1WC2EBPQ5I3G7IM3X4TZA19ICZY9EYUJ", "YFT6V4HFIFGQHXWW4R9Z5G7N5N29S9JA24", "JXPD8V2Y25XPV7621KZM73R28KWT7TJMNT", "NBFMYYT5XVWA4HCHJ41IEJJDDD41X4B8MB"}
	p := NewEtherScanProvider("https://api.etherscan.io/", keys, "http://192.168.3.59:10809", 45)
	limiter := rate.NewLimiter(45, 1)

	for {
		limiter.Wait(context.Background())
		fmt.Println("call")
		go func() {
			_, err := p.GetLogs([]string{"0xbd3531da5cf5857e7cfaa92426877b022e612caa"}, []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}, 12878196, 12878196)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}
