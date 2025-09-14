import React, { useCallback, useEffect, useMemo, useState } from "react";
import * as Router from "react-router-dom";
import Landing from "./Landing";
import TestNotLanding from "./TestNotLanding";
import {
    emptySession,
    ISession,
    SessionContext,
    SessionType,
} from "./context/session";
import Protected from "./route/Protected";

import { AuthChangeEvent, Session as SupabaseSession } from '@supabase/supabase-js'
import LoginPage from "./Login";

import { supabaseClient } from "../main";
import LogoutPage from "./Logout";

const SESSION_COOKIE_NAME = "supasession";

export default function Session(): React.JSX.Element {
    const [supabaseSession, setSupabaseSession] = useState<SupabaseSession>();
    const [supabaseGetSessionError,setSupabaseGetSessionError ] = useState<boolean>(false);
    const [mostRecentAuthEventState, setMostRecentAuthEventState] = useState<AuthChangeEvent| "NONE">("NONE")

    const setSessionAndCookie = useCallback((session: SupabaseSession)=>{
        document.cookie = `${SESSION_COOKIE_NAME}=${session.access_token}; path=/; SameSite=Lax; expires=session`;
        setSupabaseSession(session)
        console.log("spotify token: ", session.provider_token);
    }, [setSupabaseSession])

    useEffect(() => {
        supabaseClient.auth.getSession().
            then(({ data: { session }, error }) => {
                if (error != null) {
                    console.log("failed to get supabase session: ", error.code, error.message)
                    setSupabaseGetSessionError(true);
                }
                if (session != null){
                    setSessionAndCookie(session);
                    setSupabaseGetSessionError(false);
                }
            })      
        
        const {
            data: { subscription },      
        } = supabaseClient.auth.onAuthStateChange((event, session) => {
            setMostRecentAuthEventState(event);
            if (session != null){
                setSessionAndCookie(session);
            }
        }) 

        return () => subscription.unsubscribe()    
    }, [setSessionAndCookie])

    const session: ISession = useMemo(()=>{
        const userSession: ISession = {
            ...emptySession,
            supabaseClient: supabaseClient,
            getSessionError: supabaseGetSessionError,
            lastAuthEvent: mostRecentAuthEventState,
        }

        if (!supabaseSession){
            console.log("no supabase session found")
            return userSession;
        } 

        userSession.sessionType = SessionType.USER;
        userSession.name = supabaseSession.user.email ?? "Unknown user email";
        
        return userSession
    }, [mostRecentAuthEventState, supabaseGetSessionError, supabaseSession])

    return (
        <SessionContext.Provider value={session}>
            <Router.Routes>
                <Router.Route path="/p/login" element={<LoginPage/>} />
                <Router.Route path="/p/logout" element={<LogoutPage/>} />
                <Router.Route path="/p/landing" element={<Landing apiBaseUrl="http://localhost:6789/api/v1"/>} />
                <Router.Route path="/p/notlanding" element={<TestNotLanding/>} />
                <Router.Route element={<Protected/>} >
                    <Router.Route path="/p/secret/notlanding" element={<TestNotLanding/>}/>
                </Router.Route>
            </Router.Routes>
        </SessionContext.Provider>
    );
}
