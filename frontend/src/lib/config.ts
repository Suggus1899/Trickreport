/**
 * On-premise / self-hosted configuration.
 *
 * When ON_PREMISE is true the frontend expects to be served from the same
 * origin as the API, so API_URL is empty and all requests use relative paths
 * (e.g. `/api/v1/...`). No external resources (fonts, CDNs, analytics) are
 * loaded, keeping the deployment fully air-gapped.
 */
export const ON_PREMISE =
  import.meta.env.ON_PREMISE === 'true' ||
  import.meta.env.PUBLIC_ON_PREMISE === 'true';

/**
 * Base URL for the backend API.
 *
 * - On-premise: empty string → requests go to the same origin (`/api/v1/...`).
 * - Development/standalone: defaults to `http://localhost:8080`.
 */
export const API_URL = import.meta.env.API_URL || (ON_PREMISE ? '' : 'http://localhost:8080');
