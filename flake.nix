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
        
        # Configuration Python avec Textual pré-installé
        pythonEnv = pkgs.python3.withPackages (ps: with ps; [
          textual
          requests # Probablement utile pour parler à ton API Go
          # Tu pourras ajouter d'autres libs ici (ex: numpy pour l'IA ?)
        ]);

      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # --- BACKEND (GO) ---
            go              # Le langage Go
            gopls           # Le serveur de langage (pour l'autocomplétion VSCode/Neovim)
            air             # Outil génial pour le hot-reload en Go (optionnel mais recommandé)
            golangci-lint   # Linter standard pour Go

            # --- CLIENT TUI (PYTHON) ---
            pythonEnv       # Notre Python custom avec Textual
            # poetry        # (Optionnel) Si tu préfères gérer les dépendances avec Poetry au lieu de Nix

            # --- CLIENT WEB (ANGULAR) ---
            nodejs_20       # Node.js (version LTS recommandée)
            nodePackages.npm 
            nodePackages.angular-cli # La commande 'ng' pour créer/gérer le projet Angular
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
