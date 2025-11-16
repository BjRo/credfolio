import { render, screen } from "@testing-library/react";
import HomePage from "./page";

describe("HomePage", () => {
	it("renders title", async () => {
		const ui = await HomePage();
		render(ui as unknown as React.ReactElement);
		expect(screen.getByText("Credfolio")).toBeInTheDocument();
	});
});
