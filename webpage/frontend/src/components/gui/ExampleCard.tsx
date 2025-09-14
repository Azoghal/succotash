import { useCallback, useState } from "react";
import { useTranslation } from "react-i18next";

interface IExampleCardProps {
    title: string;
    actionDescription: string;
    state: string;
    onActionClick: (input: string)=>void 
}

export default function ExampleCard(props: IExampleCardProps): JSX.Element {
    const {t} = useTranslation();
    const {onActionClick, state} = props
    
    const [exampleInput, setExampleInput] = useState("");

    const doAction = useCallback(async () => {
        onActionClick(exampleInput)
    },[onActionClick, exampleInput]); 

    return (
        <div className="card">
            <div className="card__title">{props.title}</div>
            <div className="card__body">
                <div>
                    <input
                        type="text"
                        onChange={(e) => setExampleInput(e.target.value)}
                        value={exampleInput}
                    />
                    {exampleInput}
                </div>
                <div>{t("tool.card.example.state", {value: state != "" ? state : "..."})}</div>
            </div>
            <div className="card__quick-action">
                <button
                    className="c-btn c-btn__alternate"
                    onClick={doAction}
                >
                    {props.actionDescription}
                </button>
            </div>
        </div>
    );
}
