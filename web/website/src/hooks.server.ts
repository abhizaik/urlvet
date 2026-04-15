import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

const ADMIN_PASSWORD = env.ADMIN_PASSWORD ?? '';

export const handle: Handle = async ({ event, resolve }) => {
  if (event.url.pathname.startsWith('/admin')) {
    // Block access entirely if no password is configured.
    if (!ADMIN_PASSWORD) {
      return new Response('Admin access is disabled (ADMIN_PASSWORD not set).', {
        status: 503,
        headers: { 'Content-Type': 'text/plain' },
      });
    }

    const authHeader = event.request.headers.get('authorization') ?? '';
    const authenticated = checkBasicAuth(authHeader, ADMIN_PASSWORD);

    if (!authenticated) {
      return new Response('Unauthorized', {
        status: 401,
        headers: {
          'WWW-Authenticate': 'Basic realm="SafeSurf Admin", charset="UTF-8"',
          'Content-Type': 'text/plain',
        },
      });
    }
  }

  return resolve(event);
};

/**
 * Validates a Basic Auth header against the expected password.
 * Username is ignored — only the password matters.
 */
function checkBasicAuth(authHeader: string, expectedPassword: string): boolean {
  if (!authHeader.startsWith('Basic ')) return false;
  try {
    const decoded = atob(authHeader.slice(6));
    // Format is "username:password" — take everything after the first colon.
    const colonIdx = decoded.indexOf(':');
    if (colonIdx === -1) return false;
    const password = decoded.slice(colonIdx + 1);
    return password === expectedPassword;
  } catch {
    return false;
  }
}
