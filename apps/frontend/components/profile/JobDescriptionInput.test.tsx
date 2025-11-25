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
			screen.getByPlaceholderText(/paste the complete job description/i),
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

	it("when text too short, shows minimum character warning", () => {
		// Arrange & Act
		render(<JobDescriptionInput value="Short text" onChange={() => {}} />);

		// Assert
		expect(screen.getByText(/minimum 50 required/i)).toBeInTheDocument();
	});

	it("when valid text length, shows checkmark", () => {
		// Arrange & Act
		const validText =
			"This is a job description that is long enough to meet the minimum character requirement for validation purposes.";
		render(<JobDescriptionInput value={validText} onChange={() => {}} />);

		// Assert
		expect(screen.getByText(/characters ✓/i)).toBeInTheDocument();
	});

	it("when error provided, displays error message", () => {
		// Arrange & Act
		render(
			<JobDescriptionInput
				value=""
				onChange={() => {}}
				error="Something went wrong"
			/>,
		);

		// Assert
		expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
	});

	it("when onTailor provided and valid, shows tailor button", () => {
		// Arrange
		const onTailor = vi.fn();
		const validText =
			"This is a job description that is long enough to meet the minimum character requirement for validation purposes.";
		render(
			<JobDescriptionInput
				value={validText}
				onChange={() => {}}
				onTailor={onTailor}
			/>,
		);

		// Act
		const button = screen.getByRole("button", { name: /tailor my profile/i });
		fireEvent.click(button);

		// Assert
		expect(onTailor).toHaveBeenCalled();
	});

	it("when isLoading, shows loading state and disables button", () => {
		// Arrange & Act
		const validText =
			"This is a job description that is long enough to meet the minimum character requirement for validation purposes.";
		render(
			<JobDescriptionInput
				value={validText}
				onChange={() => {}}
				onTailor={() => {}}
				isLoading={true}
			/>,
		);

		// Assert
		expect(screen.getByText(/analyzing job description/i)).toBeInTheDocument();
		expect(screen.getByRole("button")).toBeDisabled();
	});
});
