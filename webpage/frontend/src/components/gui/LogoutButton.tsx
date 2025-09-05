import { useTranslation } from "react-i18next";

export default function LogoutButton() {
    const { t } = useTranslation();

    return (
        <>
            <a href="/p/logout" className="c-btn">
                {t("logout")}
            </a>
        </>
    );
}