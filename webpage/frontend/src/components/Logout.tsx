import React from "react";
import { useTranslation } from "react-i18next";

export default function LogoutPage(): React.JSX.Element {
    const {t} = useTranslation();
    return (
        <>
            <main className="c-page">
                <div>
                    {t("logout.success")}
                </div>
                <div>
                    {t("logout.success.action")}
                </div>
                <div>
                    <a href="/p/landing">
                        {t("logout.home.link")}
                    </a>
                </div>
            </main>
        </>
    );
}
