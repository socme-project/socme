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

<<<<<<< HEAD
    outputHash = "sha256-2tWytM5rAqzDmoFdvJ6oOWfSFrypbQrGvoU3b5/eF4U=";
=======
    outputHash = "sha256-ybuje0uSUAk5IMmIPY6K5XN1ja4racEvKrJXirF1ah8=";
>>>>>>> c80b852e863036717eb3c0f284ce28608df5097f
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
