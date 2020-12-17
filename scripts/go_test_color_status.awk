# Functions to make string coloring easier, using ANSI escape sequences
# "\033[" is ANSI escape for "ESC[", and the "%dm" specifies
# a terminal color using the passed in parameter,
# "0m" resets the terminal color.
function color(c, s) {
  return sprintf("\033[%dm%s\033[0m", c, s)
}

function red(s) {
  return color(31, s)
}

function green(s) {
  return color(32, s)
}

function cyan(s) {
  return color(36, s)
}

function yellow(s) {
  return color(33, s)
}

# Empty pattern, matches every line.
{
	# Wrap RUN in cyan
	sub("=== RUN",  cyan("=== RUN"))

	# Wrap PASS or ok in green
	sub("--- PASS", green("--- PASS"))
	sub(/^PASS/, green("PASS"))
	sub(/^ok/, green("ok"))

	# Wrap FAIL in red
	sub("--- FAIL", red("--- FAIL"))
	sub(/^FAIL/, red("FAIL"))

	# Wrap ? (no test files) in yellow
	# sub(/^\?/, yellow("?"))

	print
}
