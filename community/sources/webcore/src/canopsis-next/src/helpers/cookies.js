import Cookies from 'js-cookie';

import { COOKIE_SESSION_KEY } from '@/config';

/**
 * Remove application cookie
 */
export const removeCookie = () => Cookies.remove(COOKIE_SESSION_KEY);

/**
 * Get application cookie
 *
 * @return {string}
 */
export const getCookie = () => Cookies.get(COOKIE_SESSION_KEY);

/**
 * Checking for the existence of a cookie
 *
 * @return {boolean}
 */
export const hasCookie = () => !!getCookie();
