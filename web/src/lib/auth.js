import { browser } from '$app/environment';
import { decodeJwt } from 'jose';

const jwt = browser ? localStorage.getItem('token') : null;
const decoded = jwt ? decodeJwt(jwt) : { exp: 0 };
export const token = decoded.exp * 1000 > Date.now() ? decoded : null;
