"use client";

import { useEffect } from "react";
import { initApi } from "../lib/api/config";

export default function ApiInitializer() {
	useEffect(() => {
		initApi();
	}, []);

	return null;
}

