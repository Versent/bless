# versent bless client

This project is a port of the [netflix/bless/bless_client](https://github.com/Netflix/bless/tree/master/bless_client).

# usage

```
usage: versent-bless [<flags>] <command> [<args> ...]

A command line client for netflix bless.

Flags:
  --help   Show context-sensitive help (also try --help-long and --help-man).
  --debug  Enable debug mode.

Commands:
  help [<command>...]
    Show help.

  login <region> <lambda_function_name> <bastion_user> <bastion_user_ip> <remote_usernames> <bastion_ips> <bastion_command> <public_key_to_sign> <certificate_filename> [<kmsauth_token>]
    Login and retrieve a key.

```

# license

This code is released under MIT License.
