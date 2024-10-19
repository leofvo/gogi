# Getting started

Github profile inspector. Grab all email associated to the github profile.

## Scan your github account

```bash
go run main.go scan <username> --output result.json --exclude-mail <public_email> --token <github_api_token>
```

### Parameters

```txt
Flags:
  -f, --exclude-forks              Exclude forked repositories (default true)
  -e, --exclude-mail stringArray   Email(s) to exclude from the scan
  -h, --help                       help for scan
  -o, --output string              File to output results to
  -r, --private                    Include private repositories (default true)
  -p, --public                     Include public repositories (default true)
  -t, --token string               GitHub token for authentication
```

## How to change your email and username from repo history

If you want to rewrite the history and change the email or username, you can use the following:

```bash
pip install git-filter-repo
cd <your_git_repository>
git-filter-repo --email-callback 'return email.replace(b"old@email.com", b"new@email.com")' --name-callback 'return name.replace(b"old-username", b"new-username")' --force
```
