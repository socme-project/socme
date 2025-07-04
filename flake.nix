{
  description = "Flake for SOCme";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    bun2nix = {
      url = "github:baileyluTCD/bun2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { systems, self, nixpkgs, bun2nix, ... }:
    let
      eachSystem = nixpkgs.lib.genAttrs (import systems);

      importFrontend = system:
        import ./nix/front.nix {
          pkgs = nixpkgs.legacyPackages.${system};
          lib = nixpkgs.lib;
          bun2nix = bun2nix.packages.${system};
        };

      importBackend = system:
        import ./nix/back.nix {
          pkgs = nixpkgs.legacyPackages.${system};
          lib = nixpkgs.lib;
          inherit self;
        };
    in {
      packages = eachSystem (system: {
        socme-frontend = importFrontend system;
        socme-backend = (importBackend system).package;
      });

      nixosModules.socme-backend = { config, lib, pkgs, ... }:
        (importBackend pkgs.system).nixosModule {
          inherit config lib pkgs;
          inherit self;
        };
    };
}
