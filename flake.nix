{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-21.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }@inputs: {
    overlay = final: prev: {
      gobar = nixpkgs.legacyPackages.${prev.system}.buildGo117Module {
        pname = "gobar";
        version = "0.1.0";
        CGO_ENABLED = 0;
        src = builtins.path { path = ./.; };
        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
      };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { overlays = [ self.overlay ]; inherit system; };
    in
    rec {
      devShell = pkgs.mkShell {
        buildInputs = [ pkgs.go_1_17 ];
      };
      packages.gobar = pkgs.gobar;
      defaultPackage = packages.gobar;
      apps.gobar = flake-utils.lib.mkApp { drv = pkgs.gobar; name = "gobar"; };
      defaultApp = apps.gobar;
    });
}
