import { createContext, useContext } from "react";
import { SupabaseClient } from '@supabase/supabase-js'

export enum SessionType{
    NO_SESSION,
    USER,
}

export type ISession = {
    sessionType: SessionType;
    name: string;
    supabaseClient?: SupabaseClient;
};

export const emptySession: ISession = { sessionType: SessionType.NO_SESSION, name: "" };

const defaultSessionState = emptySession;

export const SessionContext = createContext(defaultSessionState);

export function useSession(): ISession {
    return useContext(SessionContext);
}
