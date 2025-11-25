import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { JobDescriptionInput } from "./JobDescriptionInput";

describe("JobDescriptionInput", () => {
	// T111: Unit test for JobDescriptionInput component when text entered updates value
	it("when text entered, updates value", () => {
		// Arrange
		const onChange = vi.fn();
		render(<JobDescriptionInput value="" onChange={onChange} />);

		// Act
		const textarea = screen.getByRole("textbox");
		fireEvent.change(textarea, {
			target: { value: "Backend engineer with Go experience" },
		});

		// Assert
		expect(onChange).toHaveBeenCalledWith(
			"Backend engineer with Go experience",
		);
	});

	it("displays current value", () => {
		// Arrange & Act
		render(
			<JobDescriptionInput value="Software Engineer" onChange={() => {}} />,
		);

		// Assert
		const textarea = screen.getByRole("textbox");
		expect(textarea).toHaveValue("Software Engineer");
	});

	it("displays placeholder text", () => {
		// Arrange & Act
		render(<JobDescriptionInput value="" onChange={() => {}} />);

		// Assert
		expect(
			screen.getByPlaceholderText(/paste the job description/i),
		).toBeInTheDocument();
	});

	it("displays label", () => {
		// Arrange & Act
		render(<JobDescriptionInput value="" onChange={() => {}} />);

		// Assert
		expect(screen.getByLabelText(/job description/i)).toBeInTheDocument();
	});

	it("when disabled, shows disabled state", () => {
		// Arrange & Act
		render(<JobDescriptionInput value="" onChange={() => {}} disabled />);

		// Assert
		const textarea = screen.getByRole("textbox");
		expect(textarea).toBeDisabled();
	});
});
