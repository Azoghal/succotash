import React, { useEffect, useMemo, useState } from "react";
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

import { Session as SupabaseSession } from '@supabase/supabase-js'
import LoginPage from "./Login";

import { supabaseClient } from "../main";
import LogoutPage from "./Logout";

export default function Session(): React.JSX.Element {
    const [supabaseSession, setSupabaseSession] = useState<SupabaseSession>();

    useEffect(() => {
        supabaseClient.auth.getSession().
            then(({ data: { session }, error }) => {
                console.log("got session. error?: ", error)
                if (session != null){
                    console.log("session a:",session)
                    document.cookie = `supasession=${session.access_token}; path=/; SameSite=Lax; expires=session`;
                    setSupabaseSession(session)
                }
            })      
        
        const {
            data: { subscription },      
        } = supabaseClient.auth.onAuthStateChange((event, session) => {
            console.log("auth state change, event: ", event)
            if (session != null){
                console.log("session b:",session)
                document.cookie = `supasession=${session.access_token}; path=/; SameSite=Lax; expires=session`;
                setSupabaseSession(session)
            }
        }) 

        return () => subscription.unsubscribe()    
    }, [])

    const session: ISession = useMemo(()=>{
        if (!supabaseSession){
            console.log("has not supabase session")
            return {
                ...emptySession,
                supabaseClient: supabaseClient
            }
        }

        console.log("has supabase session")

        return {
            sessionType: SessionType.USER,
            name: supabaseSession.user.email ?? "Unkown user email",
            supabaseClient: supabaseClient,
        }
    }, [supabaseSession])

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
