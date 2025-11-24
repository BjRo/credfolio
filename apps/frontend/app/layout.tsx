import "./globals.css";
import type { ReactNode } from "react";
import ApiInitializer from "../components/ApiInitializer";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <ApiInitializer />
        {children}
      </body>
    </html>
  );
}
