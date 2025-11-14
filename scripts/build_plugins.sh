#!/bin/bash

# BaconC2 Lua Plugin Deployment Script
# Copies Lua plugins to the build directory
# Usage: ./build_plugins.sh [--examples]
#   --examples: Also copy example plugins

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PLUGIN_DIR="$PROJECT_ROOT/plugins"
OUTPUT_DIR="$PROJECT_ROOT/build_plugins"
BUILD_EXAMPLES=false

for arg in "$@"; do
    case $arg in
        --examples)
            BUILD_EXAMPLES=true
            ;;
    esac
done

echo "=== BaconC2 Lua Plugin Deployment ==="
echo "Plugin Directory: $PLUGIN_DIR"
echo "Output Directory: $OUTPUT_DIR"
echo "Include Examples: $BUILD_EXAMPLES"
echo ""

mkdir -p "$OUTPUT_DIR"

plugin_count=0
success_count=0

copy_lua_plugins() {
    local source_dir=$1
    local skip_examples=$2
    
    for item in "$source_dir"/*; do
        if [ -f "$item" ] && [[ "$item" == *.lua ]]; then
            plugin_name=$(basename "$item")
            echo "Copying: $plugin_name"
            cp "$item" "$OUTPUT_DIR/"
            plugin_count=$((plugin_count + 1))
            success_count=$((success_count + 1))
        elif [ -d "$item" ]; then
            dir_name=$(basename "$item")
            
            if [ "$dir_name" = "examples" ] && [ "$skip_examples" = "true" ]; then
                continue
            fi
            
            copy_lua_plugins "$item" "false"
        fi
    done
}

if [ ! -d "$PLUGIN_DIR" ]; then
    echo "ERROR: Plugin directory does not exist: $PLUGIN_DIR"
    exit 1
fi

if [ "$BUILD_EXAMPLES" = "true" ]; then
    echo "Mode: Copying all plugins including examples"
    copy_lua_plugins "$PLUGIN_DIR" "false"
else
    echo "Mode: Copying production plugins only (skipping examples)"
    copy_lua_plugins "$PLUGIN_DIR" "true"
fi

echo ""
echo "=== Deployment Summary ==="
echo "Total plugins copied: $success_count"
echo ""

if [ $success_count -gt 0 ]; then
    echo "✓ All Lua plugins deployed successfully!"
    exit 0
else
    echo "⚠ No Lua plugins found"
    exit 0
fi
