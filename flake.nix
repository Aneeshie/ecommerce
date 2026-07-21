{
  description = "Development environment for Ecommerce";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }:
  let
    supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
  in {
    devShells = forAllSystems (system:
      let
        pkgs = import nixpkgs { inherit system; };

        go-migrate-pg = pkgs.go-migrate.overrideAttrs (old: {
          tags = [ "postgres" ];
        });
      in {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            go-migrate-pg
            git
            gnumake
            postgresql
            bun
          ];

          shellHook = ''
            echo " Ecommerce development environment"
            echo "Go: $(go version)"
          '';
        };
      }
    );
  };
}
