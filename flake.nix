{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }@inputs: {
    overlays.default = final: prev: {
      gobar = prev.buildGo118Module {
        pname = "gobar";
        version = "0.1.0";
        CGO_ENABLED = 0;
        src = builtins.path { path = ./.; };
        vendorSha256 = "sha256-jfe7Zge16YkQRdrkxrlh/RTADgawQ+mhdXW/KC2Y3d0=";
      };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { overlays = [ self.overlays.default ]; inherit system; };
    in
    rec {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [ go_1_18 go-tools ];
      };
      packages.default = pkgs.gobar;
      apps.default = flake-utils.lib.mkApp { drv = pkgs.gobar; name = "gobar"; };
    });
}
