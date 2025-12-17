{
  description = "Environnement de dev Gomoku : Go (Gin) + Python (Textual) + Angular";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        
        pythonEnv = pkgs.python3.withPackages (ps: with ps; [
          textual
          requests 
        ]);

      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # --- BACKEND (GO) ---
            go
            gopls
            air
            golangci-lint

            # --- CLIENT TUI (PYTHON) ---
            pythonEnv

            # --- CLIENT WEB (ANGULAR) ---
            nodejs_20
            nodePackages.npm 
            nodePackages."@angular/cli" # <--- LIGNE CORRIGÉE
          ];

          shellHook = ''
            echo "🚀 Environnement Gomoku chargé !"
            echo "--------------------------------"
            echo "Backend : $(go version)"
            echo "Client TUI : Python $(python --version) (Textual inclus)"
            echo "Client Web : Node $(node --version) + Angular CLI"
            echo "--------------------------------"
          '';
        };
      }
    );
}