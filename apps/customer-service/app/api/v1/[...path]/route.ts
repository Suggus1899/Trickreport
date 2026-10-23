import { proxy } from '@trickreport/workspace/proxy';
export const dynamic = 'force-dynamic';
export const runtime = 'nodejs';
const handler = (request: Request, context: { params: Promise<{ path: string[] }> }) => proxy(request, context, 'customer-service');
export { handler as GET, handler as POST, handler as PUT, handler as PATCH, handler as DELETE, handler as HEAD };
