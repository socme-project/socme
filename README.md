# SOCme

- [English version 🇬🇧](./README-en.md)
- [Frontend (FR)](./front/README.md)
- [Frontend (EN)](./front/README-en.md)
- [Backend (FR)](./back/README.md)
- [Backend (EN)](./back/README-en.md)

## Introduction

Ce dépôt contient l’application principale, composée de :

- Un **frontend** (Svelte) utilisant :
  - [`sv-router`](https://sv-router.vercel.app/)
  - [`shadcn-svelte`](https://shadcn-svelte.com/)
  - [`tailwindcss`](https://tailwindcss.com/)
- Un **backend** (Go) utilisant :
  - [`gin`](https://gin-gonic.com/)
  - [`gorm`](https://gorm.io/index.html)

## Installation

### Nix installation

1. Dans le fichier `flake.nix`, ajoutez le répertoire `socme` dans la section `inputs` :

```nix
{
  inputs = {
    socme.url = "github:socme-project/socme";
  };
}
```

2. Ajoutez dans le même fichier dans la section `outputs` le backend dans la configuration NixOS :

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

3. Pour le frontend, activez le service `nginx` avec la configuration suivante :

```nix
{ config, inputs, pkgs, ... }: {
  services.nginx = {
    enable = true;
    virtualHosts."localhost" = { # Changez "localhost" par votre nom de domaine si nécessaire
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

4. Activez le service backend :

```nix
_: {
  services.socme-backend = {
    enable = true;
    ... # Autres configurations si nécessaire
  };
}
```

5. Rebuildez la configuration NixOS :

```bash
sudo nixos-rebuild switch --flake .#core
```

### Installation manuelle

1. Clonez le répertoire :

```bash
git clone https://github.com/socme-project/socme.git
```

2. Installez les dépendances nécessaires :
   - [Just](https://github.com/casey/just)
   - [Golang](https://go.dev/doc/install)
   - [Bun](https://bun.sh/)

3. Buildez l'application :

```bash
just build
```
