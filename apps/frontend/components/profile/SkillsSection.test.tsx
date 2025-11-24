import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import SkillsSection from "./SkillsSection";

describe("SkillsSection", () => {
	it("renders nothing when skills are empty", () => {
		const { container } = render(<SkillsSection skills={[]} />);
		expect(container.firstChild).toBeNull();
	});

	it("renders nothing when skills are undefined", () => {
		const { container } = render(<SkillsSection skills={undefined} />);
		expect(container.firstChild).toBeNull();
	});

	it("renders skills section with single skill", () => {
		render(<SkillsSection skills={["JavaScript"]} />);
		expect(screen.getByText("Skills")).toBeInTheDocument();
		expect(screen.getByText("JavaScript")).toBeInTheDocument();
	});

	it("renders skills section with multiple skills", () => {
		render(<SkillsSection skills={["JavaScript", "TypeScript", "React"]} />);
		expect(screen.getByText("Skills")).toBeInTheDocument();
		expect(screen.getByText("JavaScript")).toBeInTheDocument();
		expect(screen.getByText("TypeScript")).toBeInTheDocument();
		expect(screen.getByText("React")).toBeInTheDocument();
	});
});
