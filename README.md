# Lima City + Dynamic DNS = 💖

This tool automatically obtains public IPv4 and IPv6 addresses of the private internet connection and stores them in selected DNS records on Lima City's name servers. This enables DDNS without the use of classic providers like DuckDNS, NoIP etc.

Similar to [github.com/timothymiller/cloudflare-ddns](https://github.com/timothymiller/cloudflare-ddns).

# Quickstart Guide

## 1. Buy or transfer domain on lima-city.de

Free domains do not have the option to change DNS records. Hence a "real" domain is needed, e.g. `.de` or `.com`.

## 2. Create A and AAAA records

If not already provided by default. The TTL should be set to a low value (e.g. 300 s = 5 min) so the old IP is not cached for too long when it is changed.

If necessary create subdomains and/or wildcards at this point.

## 3. Identify domain ID and all record IDs

- Separate the record IDs by type (A and AAAA).
- [Lima City ➜ Domains](https://www.lima-city.de/usercp/domains) ➜ Manage DNS ➜ Change entry ➜ Extract values from the URL.

## 4. Obtain an API key

[Lima City ➜ API Keys](https://www.lima-city.de/usercp/api_keys)

Requires the following permissions:

- `dns.admin`
- `dns.editor`
- `dns.reader`

## 5. Regularly run the script with cron or systemd

    lima_ddns -a [api-key] -d 13573 -4 23727 -4 23728 -6 23811

Sets two DNS records (23727 and 23728) for the domain with the ID 13573 to the IPv4 address and one record (23811) to the IPv6 address.
