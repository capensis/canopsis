import { isString, get } from 'lodash';

import fr from '@/i18n/messages/fr';
import en from '@/i18n/messages/en';

const checkMessageKeysIsEqual = (firstMessages, secondMessages, path) => {
  const messages = path.length ? get(firstMessages, path) : firstMessages;

  Object.entries(messages)
    .forEach(([key, message]) => {
      if (isString(message)) {
        expect(secondMessages).toHaveProperty([...path, key]);
      } else {
        checkMessageKeysIsEqual(firstMessages, secondMessages, [...path, key]);
      }
    });
};

describe('I18n messages', () => {
  const enKeys = Object.keys(en);
  const frKeys = Object.keys(fr);

  test('Translation modules count is equal', () => {
    expect(enKeys).toEqual(frKeys);
  });

  test.each(enKeys)('Translation module: "%s" count is equal', (module) => {
    checkMessageKeysIsEqual(en[module], fr[module], []);
  });
});
