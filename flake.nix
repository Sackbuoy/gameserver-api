{
  description = "Go project with Docker image";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    system = "x86_64-linux";
    pkgs = import nixpkgs {inherit system;};
    pname = "gameserver-api";
    author = "Sackbuoy";
  in {
    packages.${system} = {
      default = pkgs.buildGoModule {
        inherit pname;
        version = "0.1.0";
        src = ./.;
        vendorHash = null; # Will be updated on first build
      };

      docker = pkgs.dockerTools.buildImage {
        name = "ghcr.io/${author}/${pname}";
        tag = "latest";
        created = "now";

        copyToRoot = pkgs.buildEnv {
          name = "image-root";
          paths = [
            self.packages.${system}.default
            pkgs.coreutils
            pkgs.shadow
            pkgs.bashInteractive
          ];
          pathsToLink = ["/bin" "/etc" "/home" "/var"];
        };

        config = {
          Cmd = ["/bin/${pname}"];
          ExposedPorts = {
            "8080/tcp" = {};
          };
          WorkingDir = "/app";
          Volumes = {
            "/home/nonroot/.kube" = {};
          };
          User = "nonroot:nonroot";
        };

        runAsRoot = ''
          #!${pkgs.runtimeShell}
          ${pkgs.dockerTools.shadowSetup}
          groupadd -r nonroot
          useradd -m -r -g nonroot nonroot

          mkdir -p /app /tmp
          chmod 1777 /tmp

          mkdir /home/nonroot/.kube
          touch /home/nonroot/.kube/config
          chmod 744 /home/nonroot/.kube
          chmod 644 /home/nonroot/.kube/config
          chown -R nonroot:nonroot /home/nonroot /app
        '';
      };
    };

    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [
        go
        gopls
        gotools
        go-outline
        delve
        docker
      ];
    };
  };
}
