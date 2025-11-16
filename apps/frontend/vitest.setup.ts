import "@testing-library/jest-dom";
import React from "react";

// Ensure React is available when evaluating compiled JSX in tests
// Some server components imported directly may not go through Next.js' JSX transform
// and can expect a global React.
// @ts-ignore add React to global for test runtime
(globalThis as unknown as { React: typeof React }).React = React;
