{ buildGoModule
, writeShellScriptBin
, go-tools
, ...
}:
let
  drv = buildGoModule {
    pname = "gobar";
    version = "0.1.3";
    CGO_ENABLED = 0;
    src = ./.;
    vendorSha256 = "sha256-+x1zKaT4IiTMHvKXfGn8RFAwYJxBf3OArShOXsQh2cM=";
    preCheck = ''
      HOME=/tmp ${go-tools}/bin/staticcheck ./...
    '';
    passthru.update = writeShellScriptBin "update" ''
      if [[ $(${drv.go}/bin/go get -u all 2>&1) != "" ]]; then
        ${drv.go}/bin/go mod tidy
        sed -i 's/vendorSha256\ =.*;/vendorSha256="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";/' default.nix
      fi
    '';
  };
in
drv
