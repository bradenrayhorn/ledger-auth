#!/bin/bash

# configures vault for local development
# this configuration is insecure and should not be used in deployed environments

vault secrets enable -path=pki pki

vault secrets tune -max-lease-ttl=8760h pki
vault write pki/root/generate/internal common_name=test_root ttl=8760h
vault write pki/roles/ledger_auth allow_bare_domains=true allow_subdomains=true allow_glob_domains=true allow_any_name=true enforce_hostnames=false key_type=any

vault read -field=data -field=certificate pki/cert/ca > ca.pem

cat << HCL | vault policy write ledger_auth -
path "pki/sign/ledger_auth" {
  capabilities = ["read", "update"]
}
HCL

vault token create -field=token -policy=ledger_auth -period=30m > token
