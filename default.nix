{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorHash = "sha256-o0oCEa47l5U361UNCfojyBhC2CG1Y8yuSMYK2owvW/8=";
  ldflags = [ "-s" "-w" ];
}
