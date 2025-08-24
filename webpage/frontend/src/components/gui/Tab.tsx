import * as React from "react";

interface ITabProps {
    title: string;
    active: boolean;
    onClick(): void;
}

export default function Tab(props: ITabProps): React.JSX.Element {
    return (
        <div
            className={props.active ? "tab tab__selected" : "tab"}
            onClick={props.onClick}
        >
            {props.title}
        </div>
    );
}
