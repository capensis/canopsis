import { differenceKeys } from '@unit/utils/object';

import messages from '@/i18n/messages';

const { fr, en } = messages;

describe('I18n messages', () => {
  test('EN translation keys DOESN\'T exists inside FR', () => {
    expect(differenceKeys(en, fr)).toEqual({});
  });

  test('FR translation keys DOESN\'T exists inside EN', () => {
    expect(differenceKeys(fr, en)).toEqual({});
  });
});
