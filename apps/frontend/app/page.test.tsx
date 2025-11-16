import { render, screen } from "@testing-library/react";
import HomePage from "./page";

describe("HomePage", () => {
	it("renders title", () => {
		render(<HomePage />);
		expect(screen.getByText("Credfolio")).toBeInTheDocument();
	});
});
