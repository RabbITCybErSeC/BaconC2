#!/bin/bash

# BaconC2 Plugin Build Script
# Compiles all plugins in the plugins/ directory
# Usage: ./build_plugins [--examples]
#   --examples: Also compile example plugins

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PLUGIN_DIR="$PROJECT_ROOT/plugins"
OUTPUT_DIR="$PROJECT_ROOT/plugin_builds"
BUILD_EXAMPLES=false
GO_VERSION=$(go version | awk '{print $3}')

for arg in "$@"; do
    case $arg in
        --examples)
            BUILD_EXAMPLES=true
            shift
            ;;
    esac
done

echo "=== BaconC2 Plugin Builder ==="
echo "Go Version: $GO_VERSION"
echo "Plugin Directory: $PLUGIN_DIR"
echo "Output Directory: $OUTPUT_DIR"
echo "Build Examples: $BUILD_EXAMPLES"
echo ""

if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo "ERROR: Go plugins are not supported on Windows"
    exit 1
fi

mkdir -p "$OUTPUT_DIR"

plugin_count=0
success_count=0
fail_count=0

compile_plugins_from_dir() {
    local source_dir=$1
    local skip_examples=$2
    
    for plugin_path in "$source_dir"/*; do
        if [ -d "$plugin_path" ]; then
            plugin_name=$(basename "$plugin_path")
            
            if [ "$plugin_name" = "examples" ] && [ "$skip_examples" = "true" ]; then
                continue
            fi
            
            if [ "$plugin_name" = "examples" ]; then
                compile_plugins_from_dir "$plugin_path" "false"
                continue
            fi
            
            if ls "$plugin_path"/*.go 1> /dev/null 2>&1; then
                ((plugin_count++))
                
                echo "[$plugin_count] Compiling plugin: $plugin_name"
                
                if go build -buildmode=plugin \
                    -o "$OUTPUT_DIR/$plugin_name.so" \
                    "$plugin_path"/*.go 2>&1; then
                    
                    size=$(ls -lh "$OUTPUT_DIR/$plugin_name.so" | awk '{print $5}')
                    echo "    ✓ Success ($size)"
                    ((success_count++))
                else
                    echo "    ✗ Failed"
                    ((fail_count++))
                fi
                echo ""
            fi
        fi
    done
}

if [ "$BUILD_EXAMPLES" = "true" ]; then
    compile_plugins_from_dir "$PLUGIN_DIR" "false"
else
    compile_plugins_from_dir "$PLUGIN_DIR" "true"
fi

echo "=== Build Summary ==="
echo "Total plugins: $plugin_count"
echo "Successful: $success_count"
echo "Failed: $fail_count"
echo ""

if [ $fail_count -eq 0 ]; then
    echo "✓ All plugins compiled successfully!"
    exit 0
else
    echo "✗ Some plugins failed to compile"
    exit 1
fi
