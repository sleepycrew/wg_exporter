{
  description = "Wireguard exporter for prometheus";

  inputs = {
      flake-utils.url = "github:numtide/flake-utils";
    };

    outputs = { self, flake-utils, nixpkgs, ... }: flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs { inherit system; };
    in {
      packages = flake-utils.lib.flattenTree {
        default = pkgs.buildGoModule (let
          version = "0.0.${nixpkgs.lib.substring 0 8 self.lastModifiedDate}.${self.shortRev or "dirty"}";
        in {
          pname = "wg_exporter";
          inherit version;

          src = ./.;
          vendorHash = "sha256-uM2peb9s3lDdFay2GTY1+U6PZVd0DDT4OxHgATfW2Dw=";

        });
      };

      # Set up flake module
      nixosModules.default = { options, config, lib, pkgs, ... }: let
        cfg = config.services.wg_exporter;
        pkg = self.packages.${system}.default;
      in {
        # Set up module options
        options.services.wg_exporter.enable = lib.mkEnableOption
          "Wireguard exporter for prometheus";

        # Set up module implementation
        config = lib.mkIf cfg.enable {
          systemd.services.wg_exporter = {
            description = "Wireguard exporter for prometheus";
            after = [ "network.target" ];
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              Type = "simple";
              User = "nobody";
              ExecStart = "${pkg}/bin/wg_exporter";
              Restart = "always";
            };
          };
        };
      };
    });
}
