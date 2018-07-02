# Kuttlefish

_A simple reporter for PR reviews_

## Usage

```bash
$ kuttlefish -r org/repos -f index.html --token 12345678 --id 42 -t 0
````

* -r : Repository where you need to post a comment
* -f : The filename where the comment body is stored (prefered HTML)
* --token : Github Token
* --id : PR or Issue number
* -t : Thread type where to put comment (0 for pull request and 1 for issues)
