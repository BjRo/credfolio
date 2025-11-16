import React from "react";

export default async function HomePage() {
	const backendBaseUrl = process.env.BACKEND_API_URL;
	let backendStatus = "<unreachable>";
	try {
		const res = await fetch(`${backendBaseUrl}/healthz`, {
			cache: "no-store",
		});
		backendStatus = await res.text();
	} catch {
		backendStatus = "<unreachable>";
	}

	return (
		<main className="prose mx-auto p-6">
			<h1>Credfolio</h1>
			<p>Welcome. Frontend is up and running.</p>
			<p>Backend status is: {backendStatus}</p>
		</main>
	);
}
