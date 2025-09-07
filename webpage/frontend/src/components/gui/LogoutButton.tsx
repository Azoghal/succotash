import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { supabaseClient } from "../../main";

export default function LogoutButton() {
    const { t } = useTranslation();
    const navigate = useNavigate(); // Initialize the navigate function

    const handleLogout = async () => {
        document.cookie = "supasession=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

        const { error } = await supabaseClient.auth.signOut();
        if (error) {
            console.error("Error logging out from Supabase:", error.message);
        }

        navigate("/p/logout");
    };

    return (
        <button onClick={handleLogout} className="c-btn">
            {t("logout")}
        </button>
    );
}