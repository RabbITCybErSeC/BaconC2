#!/bin/bash

# BaconC2 Plugin Build Script
# Compiles all plugins in the plugins/ directory
# Usage: ./build_plugins [--examples]
#   --examples: Also compile example plugins

if [ "${DEBUG:-0}" = "1" ]; then
    set -x
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PLUGIN_DIR="$PROJECT_ROOT/plugins"
OUTPUT_DIR="$PROJECT_ROOT/build_plugins"
BUILD_EXAMPLES=false
GO_VERSION=$(go version | awk '{print $3}')

for arg in "$@"; do
    case $arg in
        --examples)
            BUILD_EXAMPLES=true
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
    
    echo "Scanning directory: $source_dir"
    
    for plugin_path in "$source_dir"/*; do
        if [ -d "$plugin_path" ]; then
            plugin_name=$(basename "$plugin_path")
            echo "  Found directory: $plugin_name"
            
            if [ "$plugin_name" = "examples" ] && [ "$skip_examples" = "true" ]; then
                continue
            fi
            
            if [ "$plugin_name" = "examples" ]; then
                compile_plugins_from_dir "$plugin_path" "false"
                continue
            fi
            
            if ls "$plugin_path"/*.go 1> /dev/null 2>&1; then
                plugin_count=$((plugin_count + 1))
                
                echo ""
                echo "[$plugin_count] Found plugin: $plugin_name"
                echo "    Source: $plugin_path"
                for go_file in "$plugin_path"/*.go; do
                    echo "      - $(basename "$go_file")"
                done
                
                echo "    Compiling..."
                
                build_output=$(go build -buildmode=plugin \
                    -o "$OUTPUT_DIR/$plugin_name.so" \
                    "$plugin_path"/*.go 2>&1)
                build_status=$?
                
                if [ $build_status -eq 0 ]; then
                    
                    size=$(ls -lh "$OUTPUT_DIR/$plugin_name.so" | awk '{print $5}')
                    echo "    ✓ Compiled ($size)"
                    
                    # Extract metadata if extractor exists
                    if [ -f "$SCRIPT_DIR/extract_plugin_metadata.go" ]; then
                        echo "    → Extracting metadata..."
                        metadata_output=$(go run "$SCRIPT_DIR/extract_plugin_metadata.go" \
                            "$OUTPUT_DIR/$plugin_name.so" \
                            "$OUTPUT_DIR/$plugin_name.json" 2>&1)
                        metadata_status=$?
                        
                        if [ $metadata_status -eq 0 ]; then
                            echo "    ✓ Metadata extracted"
                        else
                            echo "    ⚠ Metadata extraction failed (plugin still usable)"
                            echo "    Error output:"
                            echo "$metadata_output" | sed 's/^/      /'
                        fi
                    fi
                    success_count=$((success_count + 1))
                else
                    echo "    ✗ Compilation failed"
                    echo "    Error output:"
                    echo "$build_output" | sed 's/^/      /'
                    fail_count=$((fail_count + 1))
                fi
                echo ""
            fi
        fi
    done
}

echo "Starting plugin compilation..."
echo "Checking if plugin directory exists: $PLUGIN_DIR"
if [ ! -d "$PLUGIN_DIR" ]; then
    echo "ERROR: Plugin directory does not exist: $PLUGIN_DIR"
    exit 1
fi

if [ "$BUILD_EXAMPLES" = "true" ]; then
    echo "Mode: Building all plugins including examples"
    compile_plugins_from_dir "$PLUGIN_DIR" "false"
else
    echo "Mode: Building production plugins only (skipping examples)"
    compile_plugins_from_dir "$PLUGIN_DIR" "true"
fi

echo ""
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
