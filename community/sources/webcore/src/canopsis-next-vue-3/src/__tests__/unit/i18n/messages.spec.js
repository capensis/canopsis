import { differenceKeys } from '@unit/utils/object';
import { createI18n } from 'vue-i18n';

import messages from '@/i18n/messages';

const { fr, en } = messages;

describe('I18n messages', () => {
  test('EN translation keys DOESN\'T exists inside FR', () => {
    expect(differenceKeys(en, fr, true)).toEqual({});
  });

  test('FR translation keys DOESN\'T exists inside EN', () => {
    expect(differenceKeys(fr, en, true)).toEqual({});
  });

  test.each(['en', 'fr'])('Ticket URL title placeholders are rendered literally in %s', (locale) => {
    const i18n = createI18n({
      legacy: false,
      locale,
      messages,
    });

    expect(i18n.global.t('declareTicket.ticketURLTitleHelp')).toContain(
      '&lt;a href=&quot;{LINK}&quot;&gt;{TITLE}&lt;/a&gt;',
    );
  });
});
