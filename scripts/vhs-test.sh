#!/bin/bash
# ABOUTME: VHS test runner script for claudex UI regression testing
# This script runs all VHS tape files and validates the output.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
VHS_DIR="testdata/vhs"
OUTPUT_DIR="testdata/vhs/output"
EXPECTED_DIR="testdata/vhs/expected"
BINARY="./bin/claudex"

# Parse arguments
VERBOSE=false
UPDATE_EXPECTED=false
TAPE_FILTER=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -u|--update-expected)
            UPDATE_EXPECTED=true
            shift
            ;;
        -t|--tape)
            TAPE_FILTER="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose         Show detailed output"
            echo "  -u, --update-expected Update expected screenshots from current output"
            echo "  -t, --tape NAME       Run only tapes matching NAME (e.g., 'startup')"
            echo "  -h, --help            Show this help message"
            echo ""
            echo "VHS Test Runner for claudex"
            echo "Runs VHS tape files to test TUI functionality and captures screenshots."
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Check for required tools
check_requirements() {
    echo -e "${YELLOW}Checking requirements...${NC}"

    if ! command -v vhs &> /dev/null; then
        echo -e "${RED}Error: VHS is not installed.${NC}"
        echo "Install VHS from: https://github.com/charmbracelet/vhs"
        echo ""
        echo "Installation options:"
        echo "  brew install vhs        # macOS"
        echo "  go install github.com/charmbracelet/vhs@latest  # Go"
        echo "  nix-env -iA nixpkgs.vhs # Nix"
        exit 1
    fi

    if [[ ! -f "$BINARY" ]]; then
        echo -e "${YELLOW}Binary not found. Building...${NC}"
        make build
    fi

    echo -e "${GREEN}Requirements satisfied.${NC}"
}

# Create output directories
setup_dirs() {
    mkdir -p "$OUTPUT_DIR"
    mkdir -p "$EXPECTED_DIR"
}

# Run a single VHS tape
run_tape() {
    local tape="$1"
    local name=$(basename "$tape" .tape)

    echo -e "\n${YELLOW}Running: $name${NC}"

    if [[ $VERBOSE == true ]]; then
        vhs "$tape"
    else
        vhs "$tape" 2>&1 | grep -E "(Screenshot|Output|Error)" || true
    fi

    # Check if screenshots were created
    local screenshots=$(ls "$VHS_DIR"/*.png 2>/dev/null | wc -l)
    if [[ $screenshots -gt 0 ]]; then
        echo -e "${GREEN}  ✓ Created $screenshots screenshots${NC}"

        # Move screenshots to output directory with tape name prefix
        for png in "$VHS_DIR"/*.png; do
            if [[ -f "$png" ]]; then
                local basename=$(basename "$png")
                mv "$png" "$OUTPUT_DIR/${name}-${basename}"
            fi
        done
    fi

    # Check if GIF was created
    if [[ -f "$VHS_DIR/$name.gif" ]]; then
        echo -e "${GREEN}  ✓ Created $name.gif${NC}"
        mv "$VHS_DIR/$name.gif" "$OUTPUT_DIR/"
    fi
}

# Compare screenshots with expected (basic comparison)
compare_screenshots() {
    local name="$1"
    local passed=0
    local failed=0

    for output in "$OUTPUT_DIR"/${name}-*.png; do
        if [[ ! -f "$output" ]]; then
            continue
        fi

        local basename=$(basename "$output")
        local expected="$EXPECTED_DIR/$basename"

        if [[ ! -f "$expected" ]]; then
            echo -e "${YELLOW}  ? No expected file for $basename (new screenshot)${NC}"
            continue
        fi

        # Use file size comparison as basic check
        # For more accurate comparison, use imagemagick's compare
        local output_size=$(stat -f%z "$output" 2>/dev/null || stat -c%s "$output" 2>/dev/null)
        local expected_size=$(stat -f%z "$expected" 2>/dev/null || stat -c%s "$expected" 2>/dev/null)

        if [[ "$output_size" == "$expected_size" ]]; then
            echo -e "${GREEN}  ✓ $basename matches expected${NC}"
            ((passed++))
        else
            echo -e "${RED}  ✗ $basename differs from expected (size: $output_size vs $expected_size)${NC}"
            ((failed++))
        fi
    done

    return $failed
}

# Update expected screenshots from current output
update_expected() {
    echo -e "${YELLOW}Updating expected screenshots...${NC}"
    cp "$OUTPUT_DIR"/*.png "$EXPECTED_DIR/" 2>/dev/null || true
    echo -e "${GREEN}Expected screenshots updated.${NC}"
}

# Main test runner
main() {
    echo "========================================"
    echo "       claudex VHS Test Runner         "
    echo "========================================"

    check_requirements
    setup_dirs

    # Find all tape files
    local tapes=()
    for tape in "$VHS_DIR"/*.tape; do
        if [[ -f "$tape" ]]; then
            if [[ -n "$TAPE_FILTER" ]]; then
                if [[ $(basename "$tape") == *"$TAPE_FILTER"* ]]; then
                    tapes+=("$tape")
                fi
            else
                tapes+=("$tape")
            fi
        fi
    done

    if [[ ${#tapes[@]} -eq 0 ]]; then
        echo -e "${RED}No tape files found in $VHS_DIR${NC}"
        exit 1
    fi

    echo -e "\nFound ${#tapes[@]} tape file(s) to run:"
    for tape in "${tapes[@]}"; do
        echo "  - $(basename "$tape")"
    done

    # Run each tape
    local total_passed=0
    local total_failed=0

    for tape in "${tapes[@]}"; do
        run_tape "$tape"

        local name=$(basename "$tape" .tape)
        if compare_screenshots "$name"; then
            ((total_passed++))
        else
            ((total_failed++))
        fi
    done

    # Update expected if requested
    if [[ $UPDATE_EXPECTED == true ]]; then
        update_expected
    fi

    # Summary
    echo ""
    echo "========================================"
    echo "              Summary                   "
    echo "========================================"
    echo -e "Tapes run: ${#tapes[@]}"
    echo -e "Passed: ${GREEN}$total_passed${NC}"
    echo -e "Failed: ${RED}$total_failed${NC}"

    if [[ $total_failed -gt 0 ]]; then
        echo -e "\n${RED}Some tests failed!${NC}"
        echo "Run with -u to update expected screenshots if the changes are intentional."
        exit 1
    else
        echo -e "\n${GREEN}All tests passed!${NC}"
    fi
}

# Run main
main "$@"
