cookie = 'PHPSESSID=78no3vorv613akdfmbsrh1o1nq'

import requests

resp = requests.get('https://free.vps.vc/vps-info', headers=
{
    'Cookie': cookie
})
print(resp.text)
