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
        
        # Passage à Node 22 pour supporter Angular 21 proprement
        node = pkgs.nodejs_22;

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
            gotools
            
            # --- CLIENT TUI (PYTHON) ---
            pythonEnv

            # --- CLIENT WEB (ANGULAR) ---
            node
            nodePackages.npm 
            # Note: on ne met pas le CLI ici pour Angular 21, 
            # on le gère via npm install ou npx pour avoir la version exacte.
          ];

          shellHook = ''
            # 1. Personnalisation du titre du terminal
            echo -ne "\033]0;GOMOKU\007"
            
            # 2. Affichage des infos
            echo "🚀 Environnement Gomoku chargé !"
            echo "--------------------------------"
            echo "Backend : $(go version)"
            echo "Client TUI : Python $(python --version)"
            echo "Client Web : Node $(node --version)"
            echo "--------------------------------"

            # 3. ASTUCE POUR GARDER ZSH / OH MY ZSH
            # Si on est déjà dans zsh, on ne fait rien pour éviter les boucles infinies.
            # Sinon, on lance zsh.
            if [[ $(ps -p $PPID -o comm=) != "zsh" ]]; then
              exec zsh
            fi
          '';
        };
      }
    );
}