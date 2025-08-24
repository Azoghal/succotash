import * as React from "react";
import Tab from "./Tab";

interface ITabContainerProps {
    tabTitles: string[];
    activeTab: string;
    onTabchange(newTab: string): void;
}

export default function TabContainer(
    props: React.PropsWithChildren<ITabContainerProps>
): JSX.Element {
    return (
        <div className="tab-container">
            <div className="tab-container__tabs">
                {/* show all the tabs */}
                {props.tabTitles.map((s) => (
                    <Tab
                        active={props.activeTab == s}
                        title={s}
                        onClick={() => props.onTabchange(s)}
                    />
                ))}
            </div>
            <div className="tab-container__body">
                <div className="tab-body">{props.children}</div>
            </div>
        </div>
    );
}
