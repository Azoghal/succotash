import React from "react";
import { t } from "i18next";
import ExampleCard from "./gui/ExampleCard";
import LoginButton from "./gui/LoginButton";
import { useSession } from "./context/session";
import Header from "./gui/Header";

export default function Landing(): React.JSX.Element {
    const session = useSession();

    return (
        <>
            <Header>
                <div className="c-header-right">
                    <LoginButton />
                    <button
                        className="c-btn"
                        onClick={() => {
                            console.log("the session", session);
                        }}
                    >
                        check session
                    </button>

                    {/* TODO: turn this into a hard link to another page, where that page requires user to be authenticated.
                        Will need session user to indicate this state. I think buttons that need spotify log in to work should be
                        green/greyed out in some way to indicate to the user what they can do without logging in. */}
                    {/* {session.sessionType == SessionType.USER && (
                        <a href={authUrl} className="c-btn">
                            sign into spotify
                        </a>
                    )} */}

                    <a className="c-btn" href="/p/notlanding">
                        not landing
                    </a>
                    <a className="c-btn" href="/p/secret/notlanding">
                        secret
                    </a>
                </div>
            </Header>
            <main className="c-page">
                {/* https://bulma.io/documentation/columns/options/ */}
                {/* The above has a cool example where we can have wider or narrower things in each row */}
                <div className="columns is-multiline is-centered card-holder">
                    <div className="column is-3" key={1}>
                        <ExampleCard
                            title={t("tool.card.title", { count: 1 })}
                            bobOrBill="bob"
                        />
                    </div>
                    <div className="column is-3" key={1}>
                        <ExampleCard
                            title={t("tool.card.title", { count: 2 })}
                            bobOrBill="bill"
                        />
                    </div>
                </div>
            </main>
        </>
    );
}
