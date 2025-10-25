#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              Politburo Test Suite Runner                       ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Run tests with coverage and save output
echo -e "${BLUE}Running tests...${NC}"
TEST_OUTPUT=$(go test -v -cover ./... 2>&1)
TEST_EXIT_CODE=$?

# Count results
TOTAL_TESTS=$(echo "$TEST_OUTPUT" | grep -E "^(PASS|FAIL):" | wc -l)
PASSED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^PASS:")
FAILED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^FAIL:")
SKIPPED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "SKIP")

# Extract coverage
COVERAGE=$(echo "$TEST_OUTPUT" | grep "coverage:" | tail -1 | grep -oP '\d+\.\d+%' | head -1)
if [ -z "$COVERAGE" ]; then
    COVERAGE="N/A"
fi

# Print test output
echo "$TEST_OUTPUT"
echo ""

# Print summary
echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                       Test Summary                             ║${NC}"
echo -e "${CYAN}╠════════════════════════════════════════════════════════════════╣${NC}"

# Total tests
echo -e "${CYAN}║${NC} Total Tests:    ${BLUE}${TOTAL_TESTS}${NC}"

# Passed tests
if [ "$PASSED_TESTS" -gt 0 ]; then
    echo -e "${CYAN}║${NC} ${GREEN}✓${NC} Passed:        ${GREEN}${PASSED_TESTS}${NC}"
fi

# Failed tests
if [ "$FAILED_TESTS" -gt 0 ]; then
    echo -e "${CYAN}║${NC} ${RED}✗${NC} Failed:        ${RED}${FAILED_TESTS}${NC}"
fi

# Skipped tests
if [ "$SKIPPED_TESTS" -gt 0 ]; then
    echo -e "${CYAN}║${NC} ${YELLOW}○${NC} Skipped:       ${YELLOW}${SKIPPED_TESTS}${NC}"
fi

# Coverage
echo -e "${CYAN}║${NC} Coverage:       ${BLUE}${COVERAGE}${NC}"

echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Final status
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}✗ Some tests failed!${NC}"
    echo ""
    exit 1
fi
