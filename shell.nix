#...
{} :
  let
    # Pinned nixpkgs
    pkgs = import
      (builtins.fetchGit {
        name = "nixpkgs-22.05";
        url = "https://github.com/NixOS/nixpkgs/";
        ref = "refs/tags/22.05";
        # rev = "fefb0df7d2ab2e1cabde7312238026dcdc972441";
      }) {};
    buildNodejs = with pkgs; callPackage <nixpkgs/pkgs/development/web/nodejs/nodejs.nix> {
      inherit openssl;
      python = python3;
    };
    nodejs-16_7_0 = buildNodejs {
      enableNpm = true;
      version = "16.7.0"; 
      sha256 = "0drd7zyadjrhng9k0mspz456j3pmr7kli5dd0kx8grbqsgxzv1gs";
      # Not sure what the callPackage anolog is for referencing patches in a nix package. ¯\_(ツ)_/¯
      # patches = [
      #    pkgs.callPackage <nixpkgs/pkgs/development/web/nodejs/disable-darwin-v8-system-instrumentation.patch>
      # ];
    };
    vim-html5 = pkgs.vimUtils.buildVimPluginFrom2Nix {
      pname = "html5.vim";
      version = "0.27";
      src = pkgs.fetchFromGitHub {
        owner = "othree";
        repo = "html5.vim";
        rev = "60545015ff66cf38581cda89aea386cbbc0b9861";
        sha256 = "0f4h7fxb8jpk7cpv8a9bbi53m406f62a2r4mm3m09x6n8rx20ky7";
      };
      meta.homepage = "https://github.com/othree/html5.vim/";
    };
    this-neovim = pkgs.neovim.override {
      vimAlias = true;
      configure = {
        customRC = builtins.readFile (./.vimrc);
        packages.myplugins = with pkgs.vimPlugins; {
          start = [
            vim-nix
            vim-svelte
            vim-html5
            vim-javascript
          ];
        };
      };
    };
  in
    pkgs.mkShell {
      nativeBuildInputs = with pkgs.buildPackages; [ 
        # minio
        go
        nodejs-16_7_0
        entr
        ffmpeg
        pkg-config
        this-neovim
      ];
    }
