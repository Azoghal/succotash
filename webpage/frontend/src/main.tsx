import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./sass/main.scss";
import "./i18n/i18n.ts";
import { BrowserRouter } from "react-router-dom";
import { createClient, SupabaseClient } from "@supabase/supabase-js";

export let supabaseClient: SupabaseClient;
try {
    supabaseClient = createClient(
        import.meta.env.VITE_SUPABASE_URL,
        import.meta.env.VITE_SUPABASE_ANON_KEY
    )
} catch(e) {
    console.log(e)
}

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <BrowserRouter>
            <App />
        </BrowserRouter>
    </React.StrictMode>
);
