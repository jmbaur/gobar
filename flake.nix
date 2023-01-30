{
  description = "gobar";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    pre-commit.url = "github:cachix/pre-commit-hooks.nix";
    pre-commit.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs: with inputs;
    let
      forAllSystems = f: nixpkgs.lib.genAttrs [ "aarch64-linux" "x86_64-linux" ] (system: f {
        inherit system;
        pkgs = import nixpkgs { inherit system; overlays = [ self.overlays.default ]; };
      });
    in
    {
      overlays.default = _: prev: { gobar = prev.callPackage ./. { }; };
      devShells = forAllSystems ({ pkgs, system, ... }: {
        ci = pkgs.mkShell {
          buildInputs = with pkgs; [ go-tools just nix-prefetch revive ];
          inherit (pkgs.gobar) nativeBuildInputs;
        };
        default = self.devShells.${system}.ci.overrideAttrs (old: {
          inherit (pre-commit.lib.${system}.run {
            src = ./.;
            hooks = {
              nixpkgs-fmt.enable = true;
              govet.enable = true;
              revive.enable = true;
              gofmt.enable = true;
            };
          }) shellHook;
        });
      });
      packages = forAllSystems ({ pkgs, ... }: { default = pkgs.gobar; });
      apps = forAllSystems ({ pkgs, ... }: { default = { type = "app"; program = "${pkgs.gobar}/bin/gobar"; }; });
    };
}
