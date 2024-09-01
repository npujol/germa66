{ pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  )
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix
, ...
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
in

pkgs.mkShell {
  hardeningDisable = [ "all" ];
  packages = [
    goEnv
    gomod2nix
    pkgs.golangci-lint
    pkgs.go_1_23
    pkgs.gotools
    pkgs.go-junit-report
    pkgs.go-task
    pkgs.delve
    pkgs.go-swag
  ];
}
