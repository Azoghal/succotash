
import { SocialAuth } from '@supabase/auth-ui-react'
import { useSession } from './context/session'
import { ThemeMinimal } from '@supabase/auth-ui-shared';

export default function LoginPage(){
    const session = useSession();

    if (!session.supabaseClient){
        return <div>Error getting login service</div>
    } 

    return <SocialAuth supabaseClient={session.supabaseClient} providers={['spotify']} appearance={{theme: ThemeMinimal}}/>
}