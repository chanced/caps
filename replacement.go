/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Chance Dinkins <chanceusc@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, Subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or Substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package caps

// DefaultReplacements is the list of Replacements passed to DefaultConverter.
//
//	{"Acl", "ACL"},
//	{"Api", "API"},
//	{"Ascii", "ASCII"},
//	{"Cpu", "CPU"},
//	{"Css", "CSS"},
//	{"Dns", "DNS"},
//	{"Eof", "EOF"},
//	{"Guid", "GUID"},
//	{"Html", "HTML"},
//	{"Http", "HTTP"},
//	{"Https", "HTTPS"},
//	{"Id", "ID"},
//	{"Ip", "IP"},
//	{"Json", "JSON"},
//	{"Lhs", "LHS"},
//	{"Qps", "QPS"},
//	{"Ram", "RAM"},
//	{"Rhs", "RHS"},
//	{"Rpc", "RPC"},
//	{"Sla", "SLA"},
//	{"Smtp", "SMTP"},
//	{"Sql", "SQL"},
//	{"Ssh", "SSH"},
//	{"Tcp", "TCP"},
//	{"Tls", "TLS"},
//	{"Ttl", "TTL"},
//	{"Udp", "UDP"},
//	{"Ui", "UI"},
//	{"Uid", "UID"},
//	{"Uuid", "UUID"},
//	{"Uri", "URI"},
//	{"Url", "URL"},
//	{"Utf8", "UTF8"},
//	{"Vm", "VM"},
//	{"Xml", "XML"},
//	{"Xmpp", "XMPP"},
//	{"Xsrf", "XSRF"},
//	{"Xss", "XSS"},
var DefaultReplacements []Replacement = []Replacement{
	{"Acl", "ACL"},
	{"Api", "API"},
	{"Ascii", "ASCII"},
	{"Cpu", "CPU"},
	{"Css", "CSS"},
	{"Dns", "DNS"},
	{"Eof", "EOF"},
	{"Guid", "GUID"},
	{"Html", "HTML"},
	{"Http", "HTTP"},
	{"Https", "HTTPS"},
	{"Id", "ID"},
	{"Ip", "IP"},
	{"Json", "JSON"},
	{"Lhs", "LHS"},
	{"Qps", "QPS"},
	{"Ram", "RAM"},
	{"Rhs", "RHS"},
	{"Rpc", "RPC"},
	{"Sla", "SLA"},
	{"Smtp", "SMTP"},
	{"Sql", "SQL"},
	{"Ssh", "SSH"},
	{"Tcp", "TCP"},
	{"Tls", "TLS"},
	{"Ttl", "TTL"},
	{"Udp", "UDP"},
	{"Ui", "UI"},
	{"Uid", "UID"},
	{"Uuid", "UUID"},
	{"Uri", "URI"},
	{"Url", "URL"},
	{"Utf8", "UTF8"},
	{"Vm", "VM"},
	{"Xml", "XML"},
	{"Xmpp", "XMPP"},
	{"Xsrf", "XSRF"},
	{"Xss", "XSS"},
}

type (
	Replacement struct {
		// Camel case variant of the word which should be replaced.
		// e.g. "Http"
		Camel string
		// Screaming (all upper case) representation of the word to replace.
		// e.g. "HTTP"
		Screaming string
	}
)
