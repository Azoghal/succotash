import Header from "./Header";
import { SessionType, useSession } from "../context/session";
import { useMemo } from "react";
import LoginButton from "./LoginButton";
import LogoutButton from "./LogoutButton";
import SessionHealth from "./SessionHealth";

interface ILogInOutHeaderProps {}

export default function LogInOutHeader(props: React.PropsWithChildren<ILogInOutHeaderProps>): JSX.Element{
    const session = useSession()

    const loggedIn = useMemo(()=>{
        return session.sessionType == SessionType.USER
    },[session.sessionType])

    return (
        <Header>
            {props.children && props.children}
            {loggedIn ? <LogoutButton/> : <LoginButton/>}
            <SessionHealth />
        </Header>
    )
}