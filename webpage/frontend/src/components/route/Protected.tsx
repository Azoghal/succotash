import React from "react";
import { SessionType, useSession } from "../context/session";
import { Outlet } from "react-router-dom";

export default function Protected(): React.JSX.Element {
    const session = useSession();

    if (session.sessionType == SessionType.NO_SESSION){
        return <div>need to <a href="/p/login">log in</a> mate</div>
    }

    return (
        <Outlet />
    );
}
