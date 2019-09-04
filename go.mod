module github.com/yannh/terraform-provider-statuspage

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

require (
	github.com/hashicorp/terraform v0.12.3
	github.com/yannh/statuspage-go-sdk v0.0.0-20190706125613-73f6fed15b1a
)
