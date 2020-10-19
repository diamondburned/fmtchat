# fmtchat

Discord chat log extractor used mainly for archiving.

## Usage

```sh
export TOKEN="token goes here"
go run . -channelID 123123 [-beforeID 123123|-afterID 123123] -limit 500
```

A limit of 0 will fetch everything.
