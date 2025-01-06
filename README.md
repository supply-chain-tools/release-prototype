# Release Prototype

Tie together external tools with those in https://github.com/supply-chain-tools for an end-to-end release pipeline.

## Verify

Features
 - [Bit-by-bit verified reproducible build](https://google.github.io/building-secure-and-reliable-systems/raw/ch14.html#verifiable_build_architectures) that is signed by the third-party SLSA builder and a maintainer.
This should increase the confidence that the build process itself has not been compromised.
   - Further maintainer or verifier signatures can be added.
   - Others can verify the build as long as the dependencies continue to be available. It is possible to include the dependencies
 in this repository using using `go mod vendor`.
 - The attestation is added a public transparency log, so it's possible check that the
signing keys have not been used to create another `v0.0.4` release.

In the future this should be a single command and there should be a better way to boostrap things.

### Get binary and attestation
```sh
wget https://github.com/supply-chain-tools/release-prototype/releases/download/v0.0.4/binary-linux-amd64
wget https://github.com/supply-chain-tools/release-prototype/releases/download/v0.0.4/binary-linux-amd64.intoto.jsonl
```

### SLSA
`binary-linux-amd64.intoto.jsonl` was created by the [SLSA workflow](https://github.com/slsa-framework/slsa-github-generator/blob/v2.0.0/internal/builders/go/README.md), which achieves achieves [SLSA build security level 3](https://slsa.dev/spec/v1.0/levels). It runs i a separate VM
which is not controlled by this project. The workflow's OIDC token is used to sign the attestation (via [Fulcio](https://github.com/sigstore/fulcio)) which creates an [entry](https://search.sigstore.dev/?uuid=108e9186e8c5677ab6b091326f7d8c447d1a69572e7c94e560b51ade929719f643b891942371d7af)
the transparency log [Rekor](https://github.com/sigstore/rekor).

Install `slsa-verifier`
```sh
go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.6.0
```

```sh
slsa-verifier verify-artifact binary-linux-amd64 \
--provenance-path binary-linux-amd64.intoto.jsonl \
--source-uri github.com/supply-chain-tools/release-prototype \
--source-tag v0.0.4
```

This should return `PASSED: SLSA verification passed`

### Countersigning

`binary-linux-amd64` was reproduced bit-for-bit locally by a maintainer using [reproducible-build.sh](reproducible-build.sh). After validating
that the attestation `binary-linux-amd64.intoto.jsonl` it was countersigned, which can be verified as follows.

Install `dsse`
```sh
go install github.com/supply-chain-tools/go-sandbox/cmd/dsse@latest
```

```sh
curl https://api.github.com/users/stiankri-telenor/ssh_signing_keys | jq '.[] | select(.id==343845)' | jq -r '.key' > key.pub
cat binary-linux-amd64.intoto.jsonl | jq -r '.signatures[1].sig' | base64 --decode  > sig
dsse pae binary-linux-amd64.intoto.jsonl |  ssh-keygen -Y check-novalidate -n dsse -f key.pub -s sig
```

This should return
```sh
Good "dsse" signature with ED25519-SK key SHA256:SWd1HXZXtEr5ohh7P12awcTiSnc/W7MnbY6VWvR0zfk
```
