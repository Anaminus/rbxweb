// The rbxweb package provides an interface to many of ROBLOX's web-based
// services.
package rbxweb

// BaseDomain is the URL domain to which all requests will be sent.
//
// Subdomains are handled automatically as a part of API requests. Alternative
// domains, such as gametest, follow a scheme that makes switching domains
// easier:
//
//     BaseDomain:                  With subdomain:
//     roblox.com               --> www.roblox.com
//     gametest.robloxlabs.com  --> www.gametest.robloxlabs.com
var BaseDomain = `roblox.com`
