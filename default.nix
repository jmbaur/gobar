{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.6";
  src = ./.;
  vendorSha256 = "sha256-q5owd4bnrvTzSAtM6WM+d80XU952M8U7x4KbPN5f4fE=";
  ldflags = [ "-s" "-w" ];
}
