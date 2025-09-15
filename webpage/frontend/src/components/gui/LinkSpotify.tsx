import { useTranslation } from "react-i18next"

export default function LinkSpotify(): JSX.Element {
    const {t} = useTranslation();

    return (
        <a href="/spotify/link" className="c-btn">
            {t("spotify.link")}
        </a>
    )
    
}