# azure-env

Azure env is a CLI teal built to help me solve a specific problem.
Sometimes while testing certain CLI tools (e.g.: terraform import) locally I can't automatically use all credentials in my shell.

This tool allows me to use env vars and still use azure key vault to store my secrets.
The credentials are, by default, not stored on the file system. Later on support for locally (temporary) caching
and azure app configuration might get added.