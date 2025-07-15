{self, ...}: {
  perSystem = {
    pkgs,
    self',
    lib,
    ...
  }: {
    packages.default = pkgs.buildGoModule rec {
      pname = "wowsims-mop";
      version = "0.22.0";
      src = "${self}";

      vendorHash = "sha256-8SZj0GzWyoABvbYfQFII2POFRJZo4Y3NQc4Bu3WF7z0=";

      nativeBuildInputs = with pkgs; [
        nodejs_20
        protobuf
        protoc-gen-go
        gnumake
      ];

      overrideModAttrs = (_: {
        preBuild = ''
          # Generate protobuf files for vendoring
          protoc -I=./proto --go_out=./sim/core ./proto/*.proto
          
          # Create binary_dist structure needed for vendoring
          mkdir -p binary_dist/mop
          touch binary_dist/mop/embedded
          cp sim/web/dist.go.tmpl binary_dist/dist.go
        '';
      });

      preBuild = ''
        # Install Node dependencies
        export HOME=$TMPDIR
        npm ci

        # Generate protobuf files
        make proto
      '';

      buildPhase = ''
        runHook preBuild

        # Build for current platform only
        make clean
        make proto
        make wowsimmop
        go build -o wowsimcli --tags=with_db -ldflags="-X 'main.Version=${version}' -s -w" ./cmd/wowsimcli/cli_main.go

        runHook postBuild
      '';

      installPhase = ''
        runHook preInstall

        mkdir -p $out/bin

        # Install binaries
        cp wowsimmop $out/bin/
        cp wowsimcli $out/bin/

        # Make them executable
        chmod +x $out/bin/*

        runHook postInstall
      '';

      meta = with lib; {
        description = "WoW Mists of Pandaria Classic DPS Simulator";
        homepage = "https://github.com/wowsims/mop";
        license = licenses.mit;
        maintainers = [];
        platforms = platforms.linux ++ platforms.darwin;
      };
    };

    apps.default = {
      type = "app";
      program = "${lib.getExe self'.packages.default}";
    };
  };
}

