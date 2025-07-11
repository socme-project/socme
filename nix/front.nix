{
  pkgs,
  lib,
  bun2nix,
}: let
  packageJson = lib.importJSON ../front/package.json;

  pname =
    if lib.isAttrs packageJson && lib.hasAttr "name" packageJson
    then packageJson.name
    else "socme-frontend-default";

  version =
    if lib.isAttrs packageJson && lib.hasAttr "version" packageJson
    then packageJson.version
    else "0.0.0";
in
  bun2nix.lib.${pkgs.system}.mkBunDerivation {
    inherit pname version;
    src = ../front;

    bunNix = ../front/bun.nix;

    outputHash = "sha256-etiek3sHnDpc+eSdwwA3Qvc2kkAimc4YMF2+56lXas4=";
    outputHashAlgo = "sha256";
    outputHashMode = "recursive";

    index = null;

    nativeBuildInputs = [pkgs.bun pkgs.vite pkgs.typescript pkgs.rsync];

    buildPhase = ''
      ${pkgs.bun}/bin/bun run build
    '';

    installPhase = ''
      mkdir -p $out/${pname}
      cp -r dist/* $out/${pname}/
    '';
  }
