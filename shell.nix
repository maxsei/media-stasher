#...
{ } :
  let
    # Pinned nixpkgs
    pkgs = import (builtins.fetchGit {
      name = "nixpkg-21.11";                         
      url = "https://github.com/NixOS/nixpkgs/";             
      ref = "refs/tags/21.05";           
      rev = "fefb0df7d2ab2e1cabde7312238026dcdc972441";                       
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
  in
    pkgs.mkShell {
      nativeBuildInputs = with pkgs.buildPackages; [ 
        # minio
	docker
        go
        nodejs-16_7_0
	entr
	ffmpeg
	pkg-config
      ];
    }
