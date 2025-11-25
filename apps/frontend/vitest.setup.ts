import "@testing-library/jest-dom";
import React from "react";
import { vi } from "vitest";

// Ensure React is available when evaluating compiled JSX in tests
// Some server components imported directly may not go through Next.js' JSX transform
// and can expect a global React.
// @ts-ignore add React to global for test runtime
(globalThis as unknown as { React: typeof React }).React = React;

// Silence expected console.error outputs during tests
// These are from testing error handling scenarios and jsdom limitations
const originalConsoleError = console.error;
console.error = vi.fn((...args: unknown[]) => {
	const message = args[0]?.toString() || "";
	// Filter out expected error messages from tests
	const expectedErrors = [
		"Download error:",
		"Error fetching profile:",
		"Not implemented: navigation",
	];
	if (!expectedErrors.some((err) => message.includes(err))) {
		originalConsoleError.apply(console, args);
	}
});
