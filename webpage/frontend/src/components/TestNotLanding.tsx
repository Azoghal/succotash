import React from "react";
import TestToolTabBody from "./gui/TestToolTabBody";
import LogInOutHeader from "./gui/LogInOutHeader";

export default function TestNotLanding(): React.JSX.Element {
    return (
        <>
            <LogInOutHeader/>
            <main className="c-page">
                <TestToolTabBody />
            </main>
        </>
    );
}
