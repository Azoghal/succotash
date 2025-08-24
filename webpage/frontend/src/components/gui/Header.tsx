
import SuccotashLogo from "../../assets/logo.png";
import { t } from "i18next";

interface IHeaderProps {

}


export default function Header(props: React.PropsWithChildren<IHeaderProps>): JSX.Element{
    return (
        <header className="c-header">
            <div className="c-header-left">
                <a
                    href="/"
                >
                    <img
                        src={SuccotashLogo}
                        className="logo"
                        alt="Succotash Logo"
                    />
                </a>
                <h1>{t("title.succotash")}</h1>
            </div>
            {props.children && ( 
                <div className="c-header-right">
                    {props.children}
                </div>
            )}
           
        </header>
    )
}