export const formatDuration = (ns) => {
	const s = ns / 1e9;
	if (s > 60) {
		return (s / 60).toFixed(2) + 'm';
	}

	return s.toFixed(2) + 's';
};
