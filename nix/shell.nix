{
  perSystem = {pkgs, ...}: {
    devShells.default = pkgs.mkShellNoCC {
      name = "wowsims-mop devshell";
      packages = with pkgs; [
        # Go toolchain (>= 1.23)
        go
        gopls
        gotools
        go-tools
        
        # Node.js (>= 20)
        nodejs_20
        
        # Protocol Buffers
        protobuf
        protoc-gen-go
        
        # Optional: .NET SDK for database generation
        dotnet-sdk_9
        
        # Development tools
        air # for file watching
        gnumake
      ];
      
      shellHook = ''
        echo "WoW Sims MoP Development Environment"
        echo "Go version: $(go version)"
        echo "Node version: $(node --version)"
        echo "Protoc version: $(protoc --version)"
        
        # Set up Go environment
        export GOPATH=$HOME/go
        export PATH=$PATH:$GOPATH/bin
        
        # Install Go protobuf plugins if not present
        if ! command -v protoc-gen-go &> /dev/null; then
          echo "Installing protoc-gen-go..."
          go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        fi
        
        echo "Run 'make setup' to install pre-commit hooks and air"
        echo "Run 'npm install' to install Node dependencies"
        echo "Run 'make' to build the project"
      '';
    };
  };
}
