import { PUBLIC_API_URL } from '$env/static/public';
import { jwt } from './auth.js';

export const loadSession = async () =>
	await (
		await fetch(`${PUBLIC_API_URL}/session`, {
			headers: {
				Authorization: `Bearer ${jwt}`
			}
		})
	).json();
