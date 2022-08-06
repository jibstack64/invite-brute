## Invite Brute

### What is it?
A program that generates the provided amount of valid Discord invite codes - these codes are then tested, through proxy or raw http.
This allows somebody to create a database of valid invite codes. Given enough computing power, this could be used to index and search for Discord invites given a query, such as server name, or inviter user id.

### Flags
- `-codes #` : the number of codes you wish to generate and subsequently validify.
- `-proxies #,#,#,...` : the proxies you require to be used in GETting invite data.
- `-proxy_selection #` : can be either 'random', 'in_order' or 'reverse'. Indicates the order in which the proxies are selected for each request.
- `-url #` : the base Discord invite url to be used. Make sure to include '%s' in placement for the code's position in the url, or issues will arise.
- `-timeout_delay #` : the number of seconds to wait after recieving a 429 (too many requests) error. The default is 5.
- `-out_path #/#/...` : the path of the JSON file in which will be used to store all successfully grabbed invite data.
