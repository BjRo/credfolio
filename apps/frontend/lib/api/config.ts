import { OpenAPI } from "./generated";

export const initApi = () => {
	OpenAPI.BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
};
