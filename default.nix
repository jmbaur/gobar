{ buildGoModule, writeShellScriptBin }:
let
  drv = buildGoModule {
    pname = "gobar";
    version = "0.1.2";
    CGO_ENABLED = 0;
    src = ./.;
    vendorSha256 = "sha256-QnOWFMbzzzL3ZWdnPqMPU+5/YgEOyErVDvrhEx8SJfI=";
    passthru.update = writeShellScriptBin "update" ''
      if [[ $(${drv.go}/bin/go get -u all 2>&1) != "" ]]; then
        sed -i 's/vendorSha256\ =.*;/vendorSha256="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";/' default.nix
        ${drv.go}/bin/go mod tidy
      fi
    '';
  };
in
drv
