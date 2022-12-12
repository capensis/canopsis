import fr from '@/i18n/messages/fr';
import en from '@/i18n/messages/en';

describe('I18n messages', () => {
  const enKeys = Object.keys(en);
  const frKeys = Object.keys(fr);

  test('Translation modules count is equal', () => {
    expect(enKeys).toEqual(frKeys);
  });

  test.each(enKeys)('Translation module: "%s" count is equal', (module) => {
    expect(en[module]).toStructureEqual(fr[module]);
  });
});
