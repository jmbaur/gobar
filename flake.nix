{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }@inputs: {
    overlay = final: prev: {
      gobar = prev.buildGo118Module {
        pname = "gobar";
        version = "0.1.0";
        CGO_ENABLED = 0;
        src = builtins.path { path = ./.; };
        vendorSha256 = "sha256-ulJNWhEqMz2V21vf910DMADCf1UxKIcnZsX51+eYLXo=";
      };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { overlays = [ self.overlay ]; inherit system; };
    in
    rec {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [ go_1_18 entr ];
      };
      packages.default = pkgs.gobar;
      apps.default = flake-utils.lib.mkApp { drv = pkgs.gobar; name = "gobar"; };
    });
}
