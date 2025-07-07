package gwt

import "fmt"

type GWTLoginArgs struct {
	RemoteHost  string
	CompanyCode string
	Username    string
	Password    string
}

// BuildGWTLoginBody constructs the GWT-RPC login body with the given company code.
func BuildGWTLoginBody(args *GWTLoginArgs) string {
	return fmt.Sprintf("7|0|9|https://%s/lynx/lynx/|4775EB021C85EC0B04470837F40FC64A|com.lynxtraveltech.common.gui.client.rpc.SecurityService|login|java.lang.String/2004016611|Z|%s|%s|%s|1|2|3|4|4|5|5|5|6|7|8|9|0|", args.RemoteHost, args.CompanyCode, args.Username, args.Password)
}
