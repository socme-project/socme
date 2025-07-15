# SOCme

- [English version 🇬🇧](./README-en.md)
- [Frontend (FR)](./front/README.md)
- [Frontend (EN)](./front/README-en.md)
- [Backend (FR)](./back/README.md)
- [Backend (EN)](./back/README-en.md)

## Introduction

This repository contains the main application, composed of:

- A **frontend** (Svelte) using:
  - [`sv-router`](https://sv-router.vercel.app/)
  - [`shadcn-svelte`](https://shadcn-svelte.com/)
  - [`tailwindcss`](https://tailwindcss.com/)
- A **backend** (Go) using:
  - [`gin`](https://gin-gonic.com/)
  - [`gorm`](https://gorm.io/index.html)

## Installation

### Manual Installation

1. Clone the repository:

```bash
git clone https://github.com/socme-project/socme.git
```

2. Install the necessary dependencies:
   - [Just](https://github.com/casey/just)
   - [Golang](https://go.dev/doc/install)
   - [Bun](https://bun.sh/)

3. Build the application:

```bash
just build
```

### On NixOS

1. In the `flake.nix` file, add the `socme` directory in the `inputs` section and import the `socme.nixosModules.socme-backend` module:

```nix
{
  inputs = {
    socme.url = "github:socme-project/socme";
  };
  outputs = { 
    # ...
    modules = [
      inputs.socme.nixosModules.socme-backend
    ];
    # ...
  }
}
```

2. For the frontend, enable the `nginx` service with the following configuration:

```nix
  services.nginx = {
    virtualHosts."localhost" = { # Change "localhost" by your domain name if necessary
      root =
        "${inputs.socme.packages.${pkgs.system}.socme-frontend}/socme-frontend";

      locations."/api/" = {
        proxyPass =
          "http://127.0.0.1:${toString config.services.socme-backend.port}/";
        extraConfig = ''
          rewrite ^/api/(.*) /$1 break;
        '';
      };

      locations."/" = {
        extraConfig = ''
          try_files $uri $uri/ /index.html =404;
        '';
      };
    };
  };
}
```

3. Enable the backend service:

```nix
services.socme-backend = {
  enable = true;
  domain = "socme.wiki";
  githubClientId = "...";
  githubClientSecret = "...";
  githubRedirectUrl = "...";
};
```
