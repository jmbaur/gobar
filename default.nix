{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-irxmaB4Vuu8gJm7oIAUinAqM3nSZnEhWOBJUmIDODUg=";
  ldflags = [ "-s" "-w" ];
}
