{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
    pre-commit-hooks.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs: with inputs; {
    overlays.default = _: prev: { gobar = prev.callPackage ./. { }; };
  } // flake-utils.lib.eachSystem [ "aarch64-linux" "x86_64-linux" ] (system:
    let
      pkgs = import nixpkgs {
        overlays = [ self.overlays.default ];
        inherit system;
      };
      preCommitCheck = pre-commit-hooks.lib.${system}.run {
        src = ./.;
        hooks = {
          nixpkgs-fmt.enable = true;
          govet.enable = true;
          gofmt = {
            enable = true;
            entry = "${pkgs.gobar.go}/bin/gofmt -w";
            types = [ "go" ];
          };
        };
      };
    in
    rec {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [ just go-tools nix-prefetch ];
        inherit (pkgs.gobar) CGO_ENABLED nativeBuildInputs;
        inherit (preCommitCheck) shellHook;
      };
      packages.default = pkgs.gobar;
      apps.default = { type = "app"; program = "${pkgs.gobar}/bin/gobar"; };
    });
}
