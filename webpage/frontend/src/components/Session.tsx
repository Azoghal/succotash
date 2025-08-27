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


import { createClient } from '@supabase/supabase-js'
import { Session as SupabaseSession } from '@supabase/supabase-js'
import LoginPage from "./Login";

const supabase = createClient('https://ocdegtteilykjvohsxrl.supabase.co', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im9jZGVndHRlaWx5a2p2b2hzeHJsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTU4NDE2MjEsImV4cCI6MjA3MTQxNzYyMX0.quthDJ311-fbbIccEywNCAwLYRXZO8rwlthhi4hFuzw')

export default function Session(): React.JSX.Element {
    const [supabaseSession, setSupabaseSession] = useState<SupabaseSession>();

    useEffect(() => {
        supabase.auth.getSession().
            then(({ data: { session }, error }) => {
                console.log("got session. error?: ", error)
                if (session != null){
                    console.log("session a:",session)
                    setSupabaseSession(session)
                }
            })      
        
        const {
            data: { subscription },      
        } = supabase.auth.onAuthStateChange((event, session) => {
            console.log("auth state change, event: ", event)
            if (session != null){
                console.log("session b:",session)
                setSupabaseSession(session)
            }
        }) 

        return () => subscription.unsubscribe()    
    }, [])

    const session: ISession = useMemo(()=>{
        if (!supabaseSession){
            return {
                ...emptySession,
                supabaseClient: supabase
            }
        }

        return {
            sessionType: SessionType.USER,
            name: supabaseSession.user.email ?? "Unkown user email",
            supabaseClient: supabase,
        }
    }, [supabaseSession])

    return (
        <SessionContext.Provider value={session}>
            <Router.Routes>
                <Router.Route path="/p/login" element={<LoginPage/>} />
                <Router.Route path="/p/landing" element={<Landing/>} />
                <Router.Route path="/p/notlanding" element={<TestNotLanding/>} />
                <Router.Route element={<Protected/>} >
                    <Router.Route path="/p/secret/notlanding" element={<TestNotLanding/>}/>
                </Router.Route>
            </Router.Routes>
        </SessionContext.Provider>
    );
}
