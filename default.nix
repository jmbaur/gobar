{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.8";
  src = ./.;
  vendorSha256 = "sha256-kn3DFC/xd2gi1S14hUGygIMsttj8L+LwvGIF0POOa8Y=";
  ldflags = [ "-s" "-w" ];
}
