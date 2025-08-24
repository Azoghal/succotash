import { expect, describe, it } from "vitest";
import { render, screen } from "@testing-library/react";
// import {vi} from 'vitest';
import App from "./App";

describe("App", () => {
    it("renders", () => {
        render(<App />);
        const succotashLogo = screen.getByAltText("Succotash Logo");
        expect(succotashLogo).toBeInTheDocument();
    });
});
