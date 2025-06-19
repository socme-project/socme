{
  description = "Flake for SOCme";

  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      socmeBackendPackage = pkgs.buildGoModule {
        pname = "socme-backend";
        version = "0.1.0";
        src = ./backend;
        vendorHash = "sha256-axt8YiVCGHnlXh174vUEKjCVWrR5YXYT+JPaNkQrE+0=";
      };
    in {
      packages.${system}.socme-backend = socmeBackendPackage;

      nixosModules.socme-backend = { config, lib, pkgs, ... }: {
        options.services.socme-backend = {
          enable = lib.mkEnableOption "Exposed backend for SOCme";
          port = lib.mkOption {
            type = lib.types.port;
            default = 8090;
            description = "Port on which the SOCme backend will listen.";
          };
        };
        config = lib.mkIf config.services.socme-backend.enable {
          systemd.services.socme-backend = {
            description = "SOCme Backend Service";
            after = [ "network.target" ];
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              ExecStart = "${socmeBackendPackage}/bin/socme-backend";
              Restart = "always";
              User = config.services.socme-backend.user;
              Group = config.services.socme-backend.group;
              DynamicUser = true;
              StateDirectory = "socme-backend";
              ReadWritePaths = [ "/var/lib/socme-backend" ];
              Environment = [
                "SOCME_BACKEND_PORT=${
                  toString config.services.socme-backend.port
                }"
              ];
            };
          };
        };
      };
    };
}
