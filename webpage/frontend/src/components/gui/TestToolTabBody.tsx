import { useCallback, useMemo, useState } from "react";
import TabContainer from "./TabContainer";

const tabTitles = ["Tool", "Settings"];

export default function TestToolTabBody(): JSX.Element {
    const [activeTab, setActiveTab] = useState<string>(tabTitles[0]);

    const [currentPlaying, setCurrentPlaying] = useState<string>();

    const getCurrentPlaying = useCallback(() => {
        setCurrentPlaying("Well this isn't implemented yet...")
    }, [setCurrentPlaying]);

    const body = useMemo(() => {
        switch (activeTab) {
            default:
            case tabTitles[0]:
                return (
                    <div className="card">
                        <div className="card__title">
                            a card that lives in the tab body
                        </div>
                        <div className="card__body">the text of the card</div>
                        <div>
                            <button
                                className="c-btn"
                                onClick={getCurrentPlaying}
                            >
                                Get current playing
                            </button>
                            {currentPlaying ?? ""}
                        </div>
                    </div>
                );
            case tabTitles[1]:
                return <div>This one's not got a card in it</div>;
        }
    }, [activeTab, currentPlaying, getCurrentPlaying]);

    return (
        <TabContainer
            tabTitles={tabTitles}
            activeTab={activeTab}
            onTabchange={setActiveTab}
        >
            {body}
        </TabContainer>
    );
}
