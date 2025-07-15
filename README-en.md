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

### Nix installation

1. In the `flake.nix` file, add the `socme` directory in the `inputs` section:

```nix
{
  inputs = {
    socme.url = "github:socme-project/socme";
  };
}
```

2. Add in the same file in the `outputs` section the backend in the NixOS configuration:

```nix
{
  outputs = inputs @ {nixpkgs, ...}: {
    nixosConfigurations = {
      core = nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
          {_module.args = {inherit inputs;};}
          inputs.socme.nixosModules.socme-backend
          ./hosts/core/configuration.nix
        ];
      };
    };
  };
}
```

3. For the frontend, enable the `nginx` service with the following configuration:

```nix
{ config, inputs, pkgs, ... }: {
  services.nginx = {
    enable = true;
    virtualHosts."localhost" = { # Change "localhost" to your domain name if necessary
      root =
        "${inputs.socme.packages.${pkgs.system}.socme-frontend}/socme-frontend";

      locations."/api/" = {
        proxyPass =
          "http://127.0.0.1:${toString config.services.socme-backend.port}/";
        recommendedProxySettings = true;
        extraConfig = ''
          rewrite ^/api/(.*) /$1 break;
          add_header 'Access-Control-Allow-Origin' '*' always;
          add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, OPTIONS' always;
          add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type' always;
          add_header 'Access-Control-Allow-Credentials' 'true' always;
          if ($request_method = 'OPTIONS') {
            add_header 'Access-Control-Allow-Origin' '*';
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, OPTIONS';
            add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type';
            add_header 'Access-Control-Max-Age' 1728000;
            add_header 'Content-Type' 'text/plain charset=UTF-8';
            add_header 'Content-Length' 0;
            return 204;
          }
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

4. Enable the backend service:

```nix
_: {
  services.socme-backend = {
    enable = true;
    ... # Other configurations if necessary
  };
}
```

5. Rebuild the NixOS configuration:

```bash
sudo nixos-rebuild switch --flake .#core
```

### Installation manuelle

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
