{
  description = "Flake for SOCme";

  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs, ... }:
    let
      supportedSystems =
        [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f:
        nixpkgs.lib.genAttrs supportedSystems
        (system: f system (import nixpkgs { inherit system; }));
    in {
      packages = forAllSystems (system: pkgs: {
        socme-backend = pkgs.buildGoModule {
          pname = "socme-backend";
          version = "0.1.0";
          src = ./backend;
          vendorHash = "sha256-axt8YiVCGHnlXh174vUEKjCVWrR5YXYT+JPaNkQrE+0=";
        };
      });

      nixosModules.socme-backend = { config, lib, pkgs, ... }: {
        options.services.socme-backend = {
          enable = lib.mkEnableOption "Exposed backend for SOCme";
          port = lib.mkOption {
            type = lib.types.port;
            default = 8090;
            description = "Port on which the SOCme backend will listen.";
          };
          user = lib.mkOption {
            type = lib.types.str;
            default = "socme-backend";
            description =
              "User under which the SOCme backend service will run.";
          };
          group = lib.mkOption {
            type = lib.types.str;
            default = "socme-backend";
            description =
              "Group under which the SOCme backend service will run.";
          };
          dbPath = lib.mkOption {
            type = lib.types.str;
            default = "/var/lib/socme-backend/socme.db";
            description =
              "Path to the database used by the SOCme backend service.";
          };
          alertRetrievalInterval = lib.mkOption {
            type = lib.types.str;
            default = "5m";
            description =
              "Interval for retrieving alerts from the SOCme backend.";
          };
          githubClientId = lib.mkOption {
            type = lib.types.str;
            default = "";
            description =
              "Client ID for GitHub OAuth integration in SOCme backend.";
          };
          githubClientSecret = lib.mkOption {
            type = lib.types.str;
            default = "";
            description = "Client Secret for GitHub in SOCme backend.";
          };
          githubRedirectUrl = lib.mkOption {
            type = lib.types.str;
            default = "";
            description = "Redirect URL for GitHub OAuth in SOCme backend.";
          };
          developmentMode = lib.mkOption {
            type = lib.types.bool;
            default = false;
            description =
              "Enable development mode for the SOCme backend service.";
          };
        };
        config = lib.mkIf config.services.socme-backend.enable {
          systemd.services.socme-backend = {
            description = "SOCme Backend Service";
            after = [ "network.target" ];
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              ExecStart =
                "${self.packages.${pkgs.system}.socme-backend}/bin/backend";
              Restart = "always";
              User = config.services.socme-backend.user;
              Group = config.services.socme-backend.group;
              DynamicUser = true;
              StateDirectory = "socme-backend";
              ReadWritePaths = [ "/var/lib/socme-backend" ];
              Environment = [
                "BACKEND_PORT=${toString config.services.socme-backend.port}"
                "IS_PROD=${
                  if config.services.socme-backend.developmentMode then
                    "false"
                  else
                    "true"
                }"
                "DB_PATH=${config.services.socme-backend.dbPath}"
                "ALERT_RETRIEVAL_INTERVAL=${config.services.socme-backend.alertRetrievalInterval}"
                "GITHUB_CLIENT_ID=${config.services.socme-backend.githubClientId}"
                "GITHUB_CLIENT_SECRET=${config.services.socme-backend.githubClientSecret}"
                "GITHUB_REDIRECT_URL=${config.services.socme-backend.githubRedirectUrl}"
              ];
            };
          };
        };
      };
    };
}
