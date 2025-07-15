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

### Sur NixOS

1. Dans le fichier `flake.nix`, ajoutez le répertoire `socme` dans la section
   `inputs` et importez le module `socme.nixosModules.socme-backend`:

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

2. Pour le frontend, activez le service `nginx` avec la configuration suivante :

```nix
  services.nginx = {
    virtualHosts."localhost" = { # Changez "localhost" par votre nom de domaine si nécessaire
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

3. Activez le service backend :

```nix
services.socme-backend = {
  enable = true;
  domain = "socme.wiki";
  githubClientId = "...";
  githubClientSecret = "...";
  githubRedirectUrl = "...";
};
```
