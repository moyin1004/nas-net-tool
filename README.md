# nas_net_tool
Nas ipv6 工具

- 由于使用slaac，每次ipv6前缀变更会增加nas的ipv6地址，ddnsgo没办法使用一个正确的ipv6地址
- 本工具通过获取路由器RA，提供api返回nas和路由器通告前缀匹配的ipv6地址