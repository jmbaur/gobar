{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs: with inputs; {
    overlays.default = _: prev: { gobar = prev.callPackage ./. { }; };
  } // flake-utils.lib.eachSystem [ "aarch64-linux" "x86_64-linux" ] (system:
    let
      pkgs = import nixpkgs {
        overlays = [ self.overlays.default ];
        inherit system;
      };
    in
    rec {
      devShells.default = pkgs.mkShell {
        inherit (pkgs.gobar) CGO_ENABLED;
        buildInputs = with pkgs; [ pkgs.gobar.go ];
      };
      packages.default = pkgs.gobar;
      apps.default = { type = "app"; program = "${pkgs.gobar}/bin/gobar"; };
    });
}
