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
  in {
    devShells.${system}.default = pkgs.mkShell {
      packages = with pkgs; [
        go
        go-migrate

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
