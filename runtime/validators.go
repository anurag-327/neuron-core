package runtime

import (
	"fmt"
	"strings"
	"unicode"
)

func ValidateAndSanitizeCpp(code string) error {
	// 1 Size limit
	if len(code) > 256*1024 {
		return fmt.Errorf("c++ code too large (>256KB)")
	}

	// 2 Non-printable characters
	for _, r := range code {
		if !unicode.IsPrint(r) && r != '\n' && r != '\t' {
			return fmt.Errorf("contains invalid characters")
		}
	}

	// 3 Basic language heuristics
	if !strings.Contains(code, "main(") {
		return fmt.Errorf("missing main() function")
	}
	if !strings.Contains(code, "#include") && !strings.Contains(code, "int main(") {
		return fmt.Errorf("not valid C++ source")
	}

	// Removed brittle string-based blacklisting.
	// We rely on container isolation (Network=None, ReadOnlyRootFS) for security.
	return nil
}

func ValidateAndSanitizeJS(code string) error {
	// 1. Size limit
	if len(code) > 256*1024 {
		return fmt.Errorf("javascript code too large (>256KB)")
	}

	// 2. Non-printable characters
	for _, r := range code {
		if !unicode.IsPrint(r) && r != '\n' && r != '\t' {
			return fmt.Errorf("contains invalid characters")
		}
	}

	return nil
}

func ValidateAndSanitizePython(code string) error {
	// 1. Size limit
	if len(code) > 256*1024 {
		return fmt.Errorf("python code too large (>256KB)")
	}

	// 2. Non-printable characters
	for _, r := range code {
		if !unicode.IsPrint(r) && r != '\n' && r != '\t' {
			return fmt.Errorf("contains invalid characters")
		}
	}

	return nil
}

func ValidateAndSanitizeJava(code string) error {
	// 1. Size limit
	if len(code) > 256*1024 {
		return fmt.Errorf("java code too large (>256KB)")
	}

	// 2. Non-printable characters
	for _, r := range code {
		if !unicode.IsPrint(r) && r != '\n' && r != '\t' {
			return fmt.Errorf("contains invalid characters")
		}
	}

	// 3. Must contain class and main
	if !strings.Contains(code, "class ") {
		return fmt.Errorf("missing class declaration")
	}
	if !strings.Contains(code, "public static void main") {
		return fmt.Errorf("missing main method")
	}

	return nil
}
