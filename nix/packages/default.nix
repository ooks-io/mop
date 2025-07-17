{self, ...}: {
  perSystem = {pkgs, ...}: {
    packages = {
      wowsimcli = pkgs.callPackage ./wowsimcli.nix {inherit self;};
    };
  };
}
