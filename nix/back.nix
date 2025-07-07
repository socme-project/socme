{
  pkgs,
  lib,
  self,
}: let
  socmeBackend = pkgs.buildGoModule {
    pname = "socme-backend";
    version = "0.1.0";
    src = ../back;
    vendorHash = "sha256-iZ0wVOCu6te35W+A9/JmYvxtW9ZoXFyYsU5k6oMmWAM=";
  };
in {
  package = socmeBackend;

  nixosModule = {config, ...}: {
    options.services.socme-backend = {
      enable = lib.mkEnableOption "Exposed backend for SOCme";
      domain = lib.mkOption {
        type = lib.types.str;
        default = "localhost";
        description = "Domain name for the SOCme backend service.";
      };
      port = lib.mkOption {
        type = lib.types.port;
        default = 8080;
        description = "Port on which the SOCme backend will listen.";
      };
      user = lib.mkOption {
        type = lib.types.str;
        default = "socme-backend";
        description = "User under which the SOCme backend service will run.";
      };
      group = lib.mkOption {
        type = lib.types.str;
        default = "socme-backend";
        description = "Group under which the SOCme backend service will run.";
      };
      dbPath = lib.mkOption {
        type = lib.types.str;
        default = "/var/lib/socme-backend/socme.db";
        description = "Path to the database used by the SOCme backend service.";
      };
      alertRetrievalInterval = lib.mkOption {
        type = lib.types.str;
        default = "5m";
        description = "Interval for retrieving alerts from the SOCme backend.";
      };
      githubClientId = lib.mkOption {
        type = lib.types.str;
        default = "";
        description = "Client ID for GitHub OAuth integration in SOCme backend.";
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
    };

    config = lib.mkIf config.services.socme-backend.enable {
      systemd.services.socme-backend = {
        description = "SOCme Backend Service";
        after = ["network.target"];
        wantedBy = ["multi-user.target"];
        serviceConfig = {
          ExecStart = "${self.packages.${pkgs.system}.socme-backend}/bin/cmd";
          Restart = "always";
          User = config.services.socme-backend.user;
          Group = config.services.socme-backend.group;
          DynamicUser = true;
          StateDirectory = "socme-backend";
          ReadWritePaths = ["/var/lib/socme-backend"];
          Environment = [
            "DOMAIN=${toString config.services.socme-backend.domain}"
            "BACKEND_PORT=${toString config.services.socme-backend.port}"
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
}
