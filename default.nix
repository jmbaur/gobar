{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-pnjrZo6LffjvnUlV14R6a9qgfHTmOZiwBkteSGaXSPw=";
  ldflags = [ "-s" "-w" ];
}
