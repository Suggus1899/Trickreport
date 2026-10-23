'use client';
import { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter, useSearchParams } from 'next/navigation';
import QRCode from 'qrcode';
import { LockKeyhole, ArrowUpRight, ShieldCheck } from 'lucide-react';
import { api, json } from './api';
import type { Product, User } from './contracts';
import { Button, Field, Feedback, Heading, Notice, useMutation, useText } from './ui';

export function AuthPage({ route, product }: { route: string; product: Product }) {
 const t=useText(),router=useRouter(),params=useSearchParams(),mutation=useMutation(); const [challenge,setChallenge]=useState('');
 const register=route==='/register',reset=route==='/password-reset',confirm=route==='/password-reset/confirm';
 const title=challenge?t('Verifica tu identidad','Verify your identity'):confirm?t('Elige una contraseña nueva','Choose a new password'):reset?t('Recupera tu acceso','Recover your access'):register?t('Crea tu cuenta','Create your account'):t('Tu equipo, conectado.','Your team, connected.');
 async function submit(form: FormData) { await mutation.run(async()=>{
   const email=String(form.get('email')||''),password=String(form.get('password')||''),tenant=String(form.get('tenant')||'default');
   if(reset){await api('/auth/password-reset',json('POST',{email}),tenant);return;}
   if(confirm){await api('/auth/password-reset/confirm',json('POST',{token:params.get('token'),new_password:password}),tenant);router.push('/login');return;}
   const result=await api<{requires_mfa?:boolean;mfa_token?:string}>(challenge?'/auth/mfa/login':register?'/auth/register':'/auth/login',json('POST',challenge?{mfa_token:challenge,code:form.get('code')}:{email,password,...(register?{name:form.get('name')}:{})}),tenant);
   if(result.requires_mfa&&result.mfa_token){setChallenge(result.mfa_token);return;} router.push('/dashboard');router.refresh();
 },reset?t('Si la cuenta existe, recibirás instrucciones por correo.','If the account exists, instructions will arrive by email.'):''); }
 return <div className="auth-layout"><section className="auth-story"><Image className="auth-mark" src="/trickreport-mark.svg" alt="" aria-hidden="true" width={48} height={48} unoptimized/><p className="eyebrow">TRICKREPORT / {product==='it'?'IT SERVICE':product==='msp'?'MSP':'CUSTOMER SERVICE'}</p><h1>{t('Cada solicitud.\nUn paso adelante.','Every request.\nOne step forward.')}</h1><p>{t('Un espacio claro para colaborar, resolver y cuidar a las personas que confían en tu equipo.','A clear space to collaborate, resolve and care for the people who rely on your team.')}</p><div className="auth-foot"><ShieldCheck size={20}/>{t('Acceso seguro · Espacio independiente','Secure access · Independent workspace')}</div></section><section className="auth-panel"><div className="auth-form"><LockKeyhole className="accent" size={26}/><h2>{title}</h2><p className="muted">{challenge?t('Introduce los seis dígitos de tu autenticador.','Enter the six digits from your authenticator.'):t('Accede al espacio de tu organización.','Sign in to your organization’s workspace.')}</p><form action={submit}>
 {!challenge&&!confirm&&<Field label={t('Organización','Organization')} name="tenant" defaultValue="default" required autoComplete="organization"/>}
 {register&&<Field label={t('Nombre completo','Full name')} name="name" autoComplete="name" required minLength={2}/>}
 {!challenge&&!confirm&&<Field label={t('Correo electrónico','Email address')} name="email" type="email" autoComplete="email" required/>}
 {!challenge&&!reset&&<Field label={t('Contraseña','Password')} name="password" type="password" autoComplete={register||confirm?'new-password':'current-password'} minLength={register||confirm?12:1} maxLength={128} required/>}
 {challenge&&<Field label={t('Código de verificación','Verification code')} name="code" inputMode="numeric" pattern="[0-9]{6}" maxLength={6} autoComplete="one-time-code" required/>}
 <Feedback {...mutation}/><Button type="submit" isDisabled={mutation.busy} className="full">{mutation.busy?t('Procesando…','Processing…'):reset?t('Enviar instrucciones','Send instructions'):confirm?t('Guardar contraseña','Save password'):register?t('Crear cuenta','Create account'):t('Entrar al espacio','Sign in')}<ArrowUpRight size={16}/></Button>
 </form><div className="auth-links"><Link href="/password-reset">{t('Olvidé mi contraseña','Forgot password')}</Link><Link href={register?'/login':'/register'}>{register?t('Ya tengo una cuenta','I have an account'):t('Crear una cuenta','Create an account')}</Link></div></div></section></div>;
}

export function Profile({user}:{user:User}) { const t=useText(),m=useMutation(),router=useRouter();const [setup,setSetup]=useState<{secret:string;qr_url:string}>();const [qr,setQR]=useState('');
 return <><Heading title={t('Tu perfil','Your profile')} description={t('Identidad y seguridad de tu cuenta.','Your account identity and security.')}/><div className="two-columns"><section className="panel"><h2>{user.name}</h2><p>{user.email}</p><dl className="details"><dt>{t('Rol','Role')}</dt><dd>{user.role}</dd><dt>{t('Organización','Organization')}</dt><dd className="mono">{user.tenant_id}</dd></dl></section><section className="panel"><h2>{t('Verificación en dos pasos','Two-step verification')}</h2><p className="muted">{user.mfa_enabled?t('Tu cuenta está protegida con un autenticador.','Your account is protected with an authenticator.'):t('Añade un código temporal al iniciar sesión.','Add a temporary code when signing in.')}</p><Feedback {...m}/>{!user.mfa_enabled&&!setup&&<Button onPress={()=>m.run(async()=>{const result=await api<{secret:string;qr_url:string}>('/auth/mfa/setup',json('POST',{}));setSetup(result);setQR(await QRCode.toDataURL(result.qr_url));})} isDisabled={m.busy}>{t('Configurar autenticador','Set up authenticator')}</Button>}{setup&&<><Image className="qr" src={qr} unoptimized width={180} height={180} alt={t('Código QR para el autenticador','Authenticator QR code')}/><Notice>{t('Clave manual: ','Manual key: ')}<code>{setup.secret}</code></Notice></>}{(user.mfa_enabled||setup)&&<form action={async form=>m.run(async()=>{await api(`/auth/mfa/${user.mfa_enabled?'disable':'enable'}`,json('POST',{code:form.get('code'),...(setup?{secret:setup.secret}:{})}));setSetup(undefined);router.refresh();},t('Seguridad actualizada.','Security updated.'))}><Field label={t('Código de seis dígitos','Six-digit code')} name="code" pattern="[0-9]{6}" maxLength={6} inputMode="numeric" required/><Button type="submit" isDisabled={m.busy}>{user.mfa_enabled?t('Desactivar verificación','Disable verification'):t('Activar verificación','Enable verification')}</Button></form>}</section></div></>;
}
