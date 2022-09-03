{ buildGoModule, writeShellScriptBin }:
let
  gobar = buildGoModule {
    pname = "gobar";
    version = "0.1.2";
    CGO_ENABLED = 0;
    src = ./.;
    vendorSha256 = "sha256-FZxD24dKT3HZU2oEA8F4dmIxajMqlEUc3ZS1qQHjdWU=";
    passthru.update = writeShellScriptBin "update" ''
      if [[ $(${gobar.go}/bin/go get -u ./...) != "" ]]; then
        sed -i 's/vendorSha256\ =\ "sha256-.*";/vendorSha256="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";/' default.nix
        echo "run 'nix build' then update the vendorSha256 field with the correct value"
      fi
    '';
  };
in
gobar
