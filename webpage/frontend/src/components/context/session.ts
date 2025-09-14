import { createContext, useContext } from "react";
import { AuthChangeEvent, SupabaseClient } from '@supabase/supabase-js'

export enum SessionType{
    NO_SESSION,
    USER,
}

export type ISession = {
    sessionType: SessionType;
    name: string;
    supabaseClient?: SupabaseClient;
    getSessionError: boolean; // whether most recent session get encountered error
    lastAuthEvent: AuthChangeEvent | "NONE"; // most recent auth event
};

export const emptySession: ISession = { sessionType: SessionType.NO_SESSION, name: "", getSessionError: false, lastAuthEvent: "NONE" };

const defaultSessionState = emptySession;

export const SessionContext = createContext(defaultSessionState);

export function useSession(): ISession {
    return useContext(SessionContext);
}
