import React from "react";
import * as Router from "react-router-dom";
import Landing from "./Landing";
import TestNotLanding from "./TestNotLanding";
import {
    emptySession,
    SessionContext,
} from "./context/session";
import Protected from "./route/Protected";

export default function Session(): React.JSX.Element {
    // Consider if we can use session or local storage for the session
    // const [session, setSession] = useState<ISession>();

    // const loadData = useCallback(() => {
    //     // TODO actually get a session
    //     setSession(emptySession)
    // }, []);

    // useEffect(() => {
    //     loadData();
    // }, [loadData]);


    return (
        <SessionContext.Provider value={emptySession}>
            <Router.Routes>
                <Router.Route path="/" Component={Landing} />
                <Router.Route path="/notlanding" Component={TestNotLanding} />
                <Router.Route Component={Protected} >
                    <Router.Route path="/secret/notlanding" element={<TestNotLanding/>}/>
                </Router.Route>
            </Router.Routes>
        </SessionContext.Provider>
    );
}
