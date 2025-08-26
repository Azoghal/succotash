import { useTranslation } from "react-i18next";

export default function LoginButton() {
    const { t } = useTranslation();

    return (
        <>
            <a href="/p/login" className="c-btn">
                {t("login")}
            </a>
        </>
    );
}
