import React, { PropsWithChildren } from "react";
import { SessionType, useSession } from "../context/session";

export default function Protected(
    props: PropsWithChildren<object>
): React.JSX.Element {
    const session = useSession();

    if (session.sessionType == SessionType.NO_SESSION){
        return <div>need to <a href="/p/login">log in</a> mate</div>
    }

    return (
        <>{props.children}</>
    );
}
