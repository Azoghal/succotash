import { useTranslation } from "react-i18next";
import { useSession } from "../context/session";

export default function SessionHealth() {
    const session = useSession();
    const {t} = useTranslation();
    const hasError = session.getSessionError;

    return (
        <div className="session-health">
            <div className={`status-icon ${hasError ? "error" : "ok"}`}></div>
            <div className="info-panel">
                <div>Status: {hasError ? t("sessionHealth.error") : t("sessionHealth.ok")}</div>
                <div>Event: {session.lastAuthEvent || 'None'}</div>
            </div>
        </div>
    );
}