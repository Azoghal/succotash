import React from "react";
import TestToolTabBody from "./gui/TestToolTabBody";
import Header from "./gui/Header";

export default function TestNotLanding(): React.JSX.Element {
    return (
        <>
            <Header/>
            <main className="c-page">
                <TestToolTabBody />
            </main>
        </>
    );
}
