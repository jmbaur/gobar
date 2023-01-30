{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.7";
  src = ./.;
  vendorSha256 = "sha256-FAarg+JNKBRkmgh+U32VXP4U8owZcyQSEw/lv57PNKo=";
  ldflags = [ "-s" "-w" ];
}
