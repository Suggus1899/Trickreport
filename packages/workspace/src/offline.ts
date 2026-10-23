import type { Product, User } from './contracts';
import { api, json } from './api';
export interface Identity { product: Product; userId: string; tenantId: string }
export interface Draft { id: string; identity: Identity; data: { title: string; description: string; priority: string; category: string }; createdAt: number; error?: string; leaseUntil?: number }
export function ownsDraft(draft: Draft, identity: Identity): boolean { return draft.identity.product === identity.product && draft.identity.userId === identity.userId && draft.identity.tenantId === identity.tenantId; }
async function store<T>(mode: IDBTransactionMode, action: (store: IDBObjectStore) => Promise<T>): Promise<T> {
 const db = await new Promise<IDBDatabase>((resolve,reject)=>{const request=indexedDB.open('trickreport-workspace-drafts',1);request.onupgradeneeded=()=>request.result.createObjectStore('drafts',{keyPath:'id'});request.onsuccess=()=>resolve(request.result);request.onerror=()=>reject(request.error);});
 const tx=db.transaction('drafts',mode);const completed=new Promise<void>((resolve,reject)=>{tx.oncomplete=()=>resolve();tx.onabort=tx.onerror=()=>reject(tx.error);});
 try { const result=await action(tx.objectStore('drafts'));await completed;return result; } catch(error) {await completed.catch(()=>{});throw error;} finally {db.close();}
}
function read<T>(request:IDBRequest<T>):Promise<T>{return new Promise((resolve,reject)=>{request.onsuccess=()=>resolve(request.result);request.onerror=()=>reject(request.error);});}
export async function saveDraft(identity:Identity,data:Draft['data']):Promise<string>{if(!identity.userId||!identity.tenantId||!identity.product)throw new Error('Verified identity required');const draft:Draft={id:crypto.randomUUID(),identity,data,createdAt:Date.now()};await store('readwrite',s=>read(s.add(draft)));return draft.id;}
export async function getDrafts(identity:Identity):Promise<Draft[]>{const drafts=await store('readonly',s=>read(s.getAll()) as Promise<Draft[]>);return drafts.filter(d=>ownsDraft(d,identity)).sort((a,b)=>a.createdAt-b.createdAt);}
export async function deleteDraft(id:string,identity:Identity):Promise<void>{await store('readwrite',async s=>{const draft=await read(s.get(id)) as Draft|undefined;if(draft&&ownsDraft(draft,identity))await read(s.delete(id));});}
export async function syncDraft(id:string,identity:Identity):Promise<void>{
 const user=await api<User>('/auth/me');if(user.id!==identity.userId||user.tenant_id!==identity.tenantId)throw new Error('The current session does not own this draft.');
 const draft=await store('readwrite',async s=>{const d=await read(s.get(id)) as Draft|undefined;if(!d||!ownsDraft(d,identity)||Number(d.leaseUntil)>Date.now())return;d.leaseUntil=Date.now()+60_000;await read(s.put(d));return d;});
 if(!draft)return;
 try{await api('/tickets',json('POST',draft.data,draft.id));await deleteDraft(id,identity);}catch(error){draft.leaseUntil=0;draft.error=error instanceof Error?error.message:'Sync failed';await store('readwrite',s=>read(s.put(draft)));throw error;}
}
