import React from "react";
import { SessionType, useSession } from "../context/session";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

export default function Protected(): React.JSX.Element {
    const session = useSession();
    const {t} = useTranslation();

    if (session.sessionType == SessionType.NO_SESSION){
        return <>
            <div>{t("protected.page")}</div>
            <li>
                <a href="/p/login">{t("login")}</a>
            </li>
            <li>
                <a href="/p/landing">{t("home")}</a>
            </li>
        </>
    }

    return (
        <Outlet />
    );
}
