import i18n from 'i18next';

import { getLang } from './locale_service';
import en from './locales/en.json';
import fr from './locales/fr.json';

// eslint-disable-next-line import/no-named-as-default-member
i18n.init({
	lng: getLang(),
	fallbackLng: 'en',
	debug: process.env.NODE_ENV === 'development',
	interpolation: {
		escapeValue: false,
	},
	// add locales here to enable them in the UI
	resources: {
		en: {
			translation: en,
		},
		fr: {
			translation: fr,
		},
	},
});

export default i18n;
