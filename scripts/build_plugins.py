#!/usr/bin/env python3

"""
BaconC2 Plugin Build Script
Compiles all plugins in the plugins/ directory
Usage: python build_plugins.py [--examples]
  --examples: Also compile example plugins
"""

import os
import sys
import subprocess
import argparse
from pathlib import Path
from typing import Tuple

plugin_count = 0
success_count = 0
fail_count = 0


def get_go_version() -> str:
    try:
        result = subprocess.run(['go', 'version'], capture_output=True, text=True, check=True)
        return result.stdout.strip().split()[2]
    except (subprocess.CalledProcessError, FileNotFoundError, IndexError):
        return "unknown"


def check_os_compatibility() -> bool:
    if sys.platform.startswith('win'):
        print("ERROR: Go plugins are not supported on Windows")
        return False
    return True


def get_file_size(filepath: Path) -> str:
    size_bytes = filepath.stat().st_size
    for unit in ['B', 'KB', 'MB', 'GB']:
        if size_bytes < 1024.0:
            return f"{size_bytes:.1f}{unit}"
        size_bytes /= 1024.0
    return f"{size_bytes:.1f}TB"


def extract_metadata(script_dir: Path, output_dir: Path, plugin_name: str) -> bool:
    extractor_path = script_dir / "extract_plugin_metadata.go"
    if not extractor_path.exists():
        return True
    
    plugin_so = output_dir / f"{plugin_name}.so"
    metadata_json = output_dir / f"{plugin_name}.json"
    
    print("    → Extracting metadata...")
    try:
        result = subprocess.run(
            ['go', 'run', str(extractor_path), str(plugin_so), str(metadata_json)],
            capture_output=True,
            text=True,
            check=True
        )
        print("    ✓ Metadata extracted")
        return True
    except subprocess.CalledProcessError as e:
        print("    ⚠ Metadata extraction failed (plugin still usable)")
        if e.stderr:
            print("    Error output:")
            for line in e.stderr.splitlines():
                print(f"      {line}")
        return False


def compile_plugin(plugin_path: Path, plugin_name: str, output_dir: Path, script_dir: Path) -> bool:
    global success_count, fail_count
    
    print(f"\n[{plugin_count}] Found plugin: {plugin_name}")
    print(f"    Source: {plugin_path}")
    
    go_files = list(plugin_path.glob("*.go"))
    for go_file in go_files:
        print(f"      - {go_file.name}")
    
    print("    Compiling...")
    
    output_file = output_dir / f"{plugin_name}.so"
    try:
        result = subprocess.run(
            ['go', 'build', '-buildmode=plugin', '-o', str(output_file)] + [str(f) for f in go_files],
            capture_output=True,
            text=True,
            check=True,
            cwd=str(plugin_path)
        )
        
        size = get_file_size(output_file)
        print(f"    ✓ Compiled ({size})")
        
        extract_metadata(script_dir, output_dir, plugin_name)
        
        success_count += 1
        return True
        
    except subprocess.CalledProcessError as e:
        print("    ✗ Compilation failed")
        print("    Error output:")
        error_output = e.stderr if e.stderr else e.stdout
        if error_output:
            for line in error_output.splitlines():
                print(f"      {line}")
        fail_count += 1
        return False


def compile_plugins_from_dir(source_dir: Path, skip_examples: bool, output_dir: Path, script_dir: Path) -> None:
    global plugin_count
    
    print(f"Scanning directory: {source_dir}")
    
    if not source_dir.is_dir():
        return
    
    for plugin_path in sorted(source_dir.iterdir()):
        if not plugin_path.is_dir():
            continue
        
        plugin_name = plugin_path.name
        print(f"  Found directory: {plugin_name}")
        
        if plugin_name == "examples" and skip_examples:
            continue
        
        if plugin_name == "examples":
            compile_plugins_from_dir(plugin_path, False, output_dir, script_dir)
            continue
        
        go_files = list(plugin_path.glob("*.go"))
        if go_files:
            plugin_count += 1
            compile_plugin(plugin_path, plugin_name, output_dir, script_dir)


def show_help():
    print("""
BaconC2 Plugin Builder

USAGE:
    python build_plugins.py [OPTIONS]

OPTIONS:
    --examples    Build example plugins in addition to production plugins
    --help        Show this help message

DESCRIPTION:
    Compiles all Go plugins in the plugins/ directory into .so files.
    Plugins are built in buildmode=plugin and metadata is extracted
    to JSON files for plugin discovery.

ENVIRONMENT:
    DEBUG=1       Enable verbose debug output

EXAMPLES:
    python build_plugins.py              # Build production plugins only
    python build_plugins.py --examples   # Build all plugins including examples
    DEBUG=1 python build_plugins.py      # Build with debug output
""")


def main():
    global plugin_count, success_count, fail_count
    
    parser = argparse.ArgumentParser(
        description='BaconC2 Plugin Builder',
        add_help=False
    )
    parser.add_argument('--examples', action='store_true', 
                       help='Also compile example plugins')
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
    
    debug = os.environ.get('DEBUG', '0') == '1'
    if debug:
        print("DEBUG MODE ENABLED")
    
    go_version = get_go_version()
    print("=== BaconC2 Plugin Builder ===")
    print(f"Go Version: {go_version}")
    print(f"Plugin Directory: {plugin_dir}")
    print(f"Output Directory: {output_dir}")
    print(f"Build Examples: {args.examples}")
    print()
    
    if not check_os_compatibility():
        sys.exit(1)
    
    output_dir.mkdir(parents=True, exist_ok=True)
    
    print("Starting plugin compilation...")
    print(f"Checking if plugin directory exists: {plugin_dir}")
    if not plugin_dir.is_dir():
        print(f"ERROR: Plugin directory does not exist: {plugin_dir}")
        sys.exit(1)
    
    if args.examples:
        print("Mode: Building all plugins including examples")
        compile_plugins_from_dir(plugin_dir, False, output_dir, script_dir)
    else:
        print("Mode: Building production plugins only (skipping examples)")
        compile_plugins_from_dir(plugin_dir, True, output_dir, script_dir)
    
    print()
    print("=== Build Summary ===")
    print(f"Total plugins: {plugin_count}")
    print(f"Successful: {success_count}")
    print(f"Failed: {fail_count}")
    print()
    
    if fail_count == 0:
        print("✓ All plugins compiled successfully!")
        sys.exit(0)
    else:
        print("✗ Some plugins failed to compile")
        sys.exit(1)


if __name__ == "__main__":
    main()
