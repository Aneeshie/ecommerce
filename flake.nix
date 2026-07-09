{
  description = "Development environment for Ecommerce";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-26.05-darwin";
  };

  outputs = { self, nixpkgs }:
  let
    system = "aarch64-darwin";

    pkgs = import nixpkgs {
      inherit system;
    };

    go-migrate-pg = pkgs.go-migrate.overrideAttrs (old: {
      tags = [ "postgres" ];
    });

  in {
    devShells.${system}.default = pkgs.mkShell {
      packages = with pkgs; [
        go
        go-migrate-pg

        git
        gnumake

        postgresql
      ];

      shellHook = ''
        echo " Ecommerce development environment"
        echo "Go: $(go version)"
      '';
    };
  };
}
