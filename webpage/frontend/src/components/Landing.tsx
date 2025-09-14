import React, { useCallback, useState } from "react";
import { t } from "i18next";
import ExampleCard from "./gui/ExampleCard";
import { useSession } from "./context/session";
import LogInOutHeader from "./gui/LogInOutHeader";

interface ILandingProps  {
    apiBaseUrl: string
}

export default function Landing(props: ILandingProps): React.JSX.Element {
    const session = useSession();
    const {apiBaseUrl} = props

    const [bobState, setBobState] = useState<string>("");
    const authedReq = useCallback((input:string) => {
        console.log("ignoring input: ", input)
        const extra: RequestInit = {credentials: "include"}
        fetch(`${apiBaseUrl}/authed/test/bob`, extra)
        .then(async (resp: Response)=>{
            if (!resp.ok) {
                throw new Error(`HTTP error! status: ${resp.status}`);
            }
            const response = await resp.json()
            setBobState(JSON.stringify(response));
        }).catch((e)=>{
            console.log("error making api request: ", e)
            setBobState("failed req");
        });
    }, [apiBaseUrl, setBobState])

    const [testState, setTestState] = useState<string>("");
    const unauthedReq = (input:string) => {
        console.log("ignoring input: ", input)
        fetch(`${apiBaseUrl}/unauthed/test`)
        .then(async (resp: Response)=>{
            if (!resp.ok) {
                throw new Error(`HTTP error! status: ${resp.status}`);
            }
            const response = await resp.json()
            setTestState(JSON.stringify(response));
        }).catch((e)=>{
            console.log("error making api request: ", e)
            setTestState("failed req");
        });
    }

    return (
        <>
            <LogInOutHeader>
                <button
                    className="c-btn"
                    onClick={() => {
                        console.log("the session", session);
                    }}
                >
                    check session
                </button>

                {/* {session.sessionType == SessionType.USER && (
                    <a href={authUrl} className="c-btn">
                        sign into spotify
                    </a>
                )} */}

                <a className="c-btn" href="/p/notlanding">
                    Example Page
                </a>
                <a className="c-btn" href="/p/secret/notlanding">
                    Restricted Page
                </a>
            </LogInOutHeader>
            <main className="c-page">
                {/* https://bulma.io/documentation/columns/options/ */}
                {/* The above has a cool example where we can have wider or narrower things in each row */}
                <div className="columns is-multiline is-centered card-holder">
                    <div className="column is-3" key={1}>
                        <ExampleCard
                            title={t("tool.card.authed")}
                            actionDescription="Do Bob"
                            state={bobState}
                            onActionClick={authedReq}
                        />
                    </div>
                    <div className="column is-3" key={2}>
                        <ExampleCard
                            title={t("tool.card.unauthed")}
                            actionDescription="Do unauth test"
                            state={testState}
                            onActionClick={unauthedReq}
                        />
                    </div>
                </div>
            </main>
        </>
    );
}
