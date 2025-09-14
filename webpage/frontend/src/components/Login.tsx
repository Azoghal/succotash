
import { SocialAuth } from '@supabase/auth-ui-react'
import { useSession } from './context/session'
import { ThemeMinimal } from '@supabase/auth-ui-shared';
import { useTranslation } from 'react-i18next';

export default function LoginPage(){
    const session = useSession();
    const {t} = useTranslation();

    if (!session.supabaseClient){
        return <div>{t("error.session")}</div>
    } 

    return <SocialAuth supabaseClient={session.supabaseClient} providers={['spotify']} appearance={{theme: ThemeMinimal}}/>
}