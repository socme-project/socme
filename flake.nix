{
  description = "Flake for SOCme";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs = { systems, self, nixpkgs, ... }:
    let
      eachSystem = nixpkgs.lib.genAttrs (import systems);

      importFrontend = system:
        import ./nix/front.nix {
          pkgs = nixpkgs.legacyPackages.${system};
          lib = nixpkgs.lib;
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
