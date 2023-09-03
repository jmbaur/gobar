{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-CTtVNYaZHk6NK7f6YrQ1/yEgYC6QjlFOg5CEOhhZXDg=";
  ldflags = [ "-s" "-w" ];
}
