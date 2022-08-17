{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs: with inputs; {
    overlays.default = _: prev: {
      gobar = prev.buildGoModule {
        pname = "gobar";
        version = "0.1.2";
        CGO_ENABLED = 0;
        src = ./.;
        vendorSha256 = "sha256-5+BYPHVyGDmTbiSjqRWpJzZOc82KwCjv2RKA75Oz4EI=";
      };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { overlays = [ self.overlays.default ]; inherit system; };
    in
    rec {
      devShells.default = pkgs.mkShell {
        CGO_ENABLED = 0;
        buildInputs = with pkgs; [ go_1_18 go-tools ];
      };
      packages.default = pkgs.gobar;
      apps.default = flake-utils.lib.mkApp { drv = pkgs.gobar; name = "gobar"; };
    });
}
