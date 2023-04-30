{ buildGoModule, ... }:
buildGoModule {
  pname = "gobar";
  version = "0.1.9";
  src = ./.;
  vendorSha256 = "sha256-Jl+K/H7FKFYcpcSr5zdFKZAN7oXnUd8FGG4kvW+xAKc=";
  ldflags = [ "-s" "-w" ];
}
