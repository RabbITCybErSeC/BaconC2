#!/bin/bash

# BaconC2 Plugin Build Script
# Compiles all plugins in the plugins/ directory

set -e

PLUGIN_DIR="./plugins"
OUTPUT_DIR="./compiled_plugins"
GO_VERSION=$(go version | awk '{print $3}')

echo "=== BaconC2 Plugin Builder ==="
echo "Go Version: $GO_VERSION"
echo "Plugin Directory: $PLUGIN_DIR"
echo "Output Directory: $OUTPUT_DIR"
echo ""


if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo "ERROR: Go plugins are not supported on Windows"
    exit 1
fi


mkdir -p "$OUTPUT_DIR"

plugin_count=0
success_count=0
fail_count=0

for plugin_path in "$PLUGIN_DIR"/*; do
    if [ -d "$plugin_path" ]; then
        plugin_name=$(basename "$plugin_path")
        
        if [ ! -f "$plugin_path"/*.go ]; then
            continue
        fi
        
        ((plugin_count++))
        
        echo "[$plugin_count] Compiling plugin: $plugin_name"
        
        # Build the plugin
        if go build -buildmode=plugin \
            -o "$OUTPUT_DIR/$plugin_name.so" \
            "$plugin_path"/*.go 2>&1; then
            
            # Get file size
            size=$(ls -lh "$OUTPUT_DIR/$plugin_name.so" | awk '{print $5}')
            echo "    ✓ Success ($size)"
            ((success_count++))
        else
            echo "    ✗ Failed"
            ((fail_count++))
        fi
        echo ""
    fi
done

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
