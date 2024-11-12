import { PUBLIC_API_URL } from '$env/static/public';

export const loadSession = async () => await (await fetch(`${PUBLIC_API_URL}/session`)).json();
