#!/usr/bin/env python3

"""
BaconC2 Lua Plugin Deployment Script
Copies Lua plugins to the build directory
Usage: python build_plugins.py [--examples]
  --examples: Also copy example plugins
"""

import os
import sys
import shutil
import argparse
from pathlib import Path

plugin_count = 0
success_count = 0


def copy_lua_plugins(source_dir: Path, skip_examples: bool, output_dir: Path) -> None:
    global plugin_count, success_count
    
    if not source_dir.is_dir():
        return
    
    for item in sorted(source_dir.iterdir()):
        if item.is_dir():
            if item.name == "examples" and skip_examples:
                continue
            
            plugin_file = item / "plugin.lua"
            if plugin_file.exists():
                output_name = f"{item.name}.lua"
                print(f"Copying: {item.name}/plugin.lua -> {output_name}")
                shutil.copy2(plugin_file, output_dir / output_name)
                plugin_count += 1
                success_count += 1
            else:
                copy_lua_plugins(item, False, output_dir)


def show_help():
    print("""
BaconC2 Lua Plugin Deployment

USAGE:
    python build_plugins.py [OPTIONS]

OPTIONS:
    --examples    Copy example plugins in addition to production plugins
    --help        Show this help message

DESCRIPTION:
    Copies all Lua plugins from the plugins/ directory to the build directory.
    Lua plugins are cross-platform and work on all architectures including ARM.

EXAMPLES:
    python build_plugins.py              # Copy production plugins only
    python build_plugins.py --examples   # Copy all plugins including examples
""")


def main():
    global plugin_count, success_count
    
    parser = argparse.ArgumentParser(
        description='BaconC2 Lua Plugin Deployment',
        add_help=False
    )
    parser.add_argument('--examples', action='store_true', 
                       help='Also copy example plugins')
    parser.add_argument('--help', action='store_true',
                       help='Show help menu')
    args = parser.parse_args()
    
    if args.help:
        show_help()
        sys.exit(0)
    
    script_dir = Path(__file__).parent.resolve()
    project_root = script_dir.parent
    plugin_dir = project_root / "plugins"
    output_dir = project_root / "build_plugins"
    
    print("=== BaconC2 Lua Plugin Deployment ===")
    print(f"Plugin Directory: {plugin_dir}")
    print(f"Output Directory: {output_dir}")
    print(f"Include Examples: {args.examples}")
    print()
    
    output_dir.mkdir(parents=True, exist_ok=True)
    
    if not plugin_dir.is_dir():
        print(f"ERROR: Plugin directory does not exist: {plugin_dir}")
        sys.exit(1)
    
    if args.examples:
        print("Mode: Copying all plugins including examples")
        copy_lua_plugins(plugin_dir, False, output_dir)
    else:
        print("Mode: Copying production plugins only (skipping examples)")
        copy_lua_plugins(plugin_dir, True, output_dir)
    
    print()
    print("=== Deployment Summary ===")
    print(f"Total plugins copied: {success_count}")
    print()
    
    if success_count > 0:
        print("✓ All Lua plugins deployed successfully!")
        sys.exit(0)
    else:
        print("⚠ No Lua plugins found")
        sys.exit(0)


if __name__ == "__main__":
    main()
