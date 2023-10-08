{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-oNI7OaJYorkVYfaHKJk2gWdZJDkKIW7+KsPfCSwStmk=";
  ldflags = [ "-s" "-w" ];
}
