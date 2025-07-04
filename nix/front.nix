{ pkgs, lib, bun2nix }:

let
  packageJson = lib.importJSON ../front/package.json;

  pname = if lib.isAttrs packageJson && lib.hasAttr "name" packageJson
          then packageJson.name
          else "socme-frontend-default";

  version = if lib.isAttrs packageJson && lib.hasAttr "version" packageJson
           then packageJson.version
           else "0.0.0";

in
bun2nix.lib.${pkgs.system}.mkBunDerivation {
  inherit pname version;
  src = ../front;

  bunNix = ../front/bun.nix;

  outputHash = "sha256-yeIqS+QKMwMrdb3SAhNKoINBVE8AjFAXNAIoEDZZ/GI="; 
  outputHashAlgo = "sha256";
  outputHashMode = "recursive"; 

  index = null;

  nativeBuildInputs = [ pkgs.bun pkgs.vite pkgs.typescript pkgs.rsync ];

  buildPhase = ''

    echo "--- Contents of current directory before build: ---"
    ls -lA .
    echo "---------------------------------------------------"

    ${pkgs.bun}/bin/bun run build

    echo "--- Contents of 'dist' directory after build: ---"
    ls -lA dist/ || echo "dist directory not found or empty after build!"
    echo "---------------------------------------------------"
  '';

  installPhase = ''
    echo "Installing frontend artifacts..."
    mkdir -p $out/share/${pname}

    echo "--- Contents of 'dist' directory before copy in installPhase: ---"
    ls -lA dist/ || echo "dist directory not found or empty before copy!"
    echo "----------------------------------------------------------------"

    cp -r dist/* $out/share/${pname}/

    echo "--- Contents of $out/share/${pname}/ after copy: ---"
    ls -lA $out/share/${pname}/ || echo "Output directory not found or empty after copy!"
    echo "----------------------------------------------------"
  '';
}

