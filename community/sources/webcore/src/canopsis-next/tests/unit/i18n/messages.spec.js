import { differenceKeys } from '@unit/utils/object';

import en from '@/i18n/messages/en';
import fr from '@/i18n/messages/fr';

describe('I18n messages', () => {
  test('EN translation keys exists inside FR', () => {
    expect(differenceKeys(en, fr)).toEqual({});
  });

  test('FR translation keys exists inside EN', () => {
    expect(differenceKeys(fr, en)).toEqual({});
  });
});
