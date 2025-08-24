import React, { PropsWithChildren, useEffect } from "react";
import { SessionType, useSession } from "../context/session";
import { useNavigate } from "react-router-dom";

export default function Protected(
    props: PropsWithChildren<object>
): React.JSX.Element {
    const session = useSession();
    const navigate = useNavigate();

    useEffect(() => {
        console.log(session);
        if (session.sessionType == SessionType.NO_SESSION) {
            // user is not authenticated
            // TODO fix this. We want to redirect outside of react router, in order to hit the backend login endpoint.
            // For now this is ok, it redirects to a react page with a single anchor that does take us to the backend
            navigate("/login");
        }
    }, [session, navigate]);

    return (
        <>{session.sessionType != SessionType.NO_SESSION && props.children}</>
    );
}
