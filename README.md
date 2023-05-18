# onlygithub

GitHub Sponsors as a content subscription service provider. Kind of like
OnlyFans but you own the content and it uses GitHub Sponsors.

<div align="center">
  <img src=".github/screenshot01.png" alt="screenshot" width="650px">
</div>

## Building

```sh
make
```

## Running

### Setting Up a Database

A database is required to run this. Currently, only SQLite and PostgreSQL
are supported.

#### SQLite

No setup is required. The database will be created automatically.

#### PostgreSQL

_TODO_

### Locally or on a Server

You will first need to make a GitHub OAuth app. The callback URL should be
set to whatever the domain is that you are running this on.

Then, set these as environment variables in whatever environment you are
running this on, for example:

```sh
export GITHUB_CLIENT_ID=     # github oauth app ID
export GITHUB_CLIENT_SECRET= # github oauth app secret
export GITHUB_REDIRECT_URL=  # github oauth app callback URL, e.g. https://example.com
```

Then, run the server:

```sh
./dist/onlygithub
```

The server will be listening on port 8080 by default. You can change this
by setting the `--http` flag. Unix sockets are supported using `--http unix://...`.

### Netlify

_TODO_
