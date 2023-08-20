{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-//4sxpT8ZlGBYPr6av8lKg9r8JZf6zT769+oHpKM6Lo=";
  ldflags = [ "-s" "-w" ];
}
