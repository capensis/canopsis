import { mount, flushPromises } from '@vue/test-utils';
import { h } from 'vue';
import { VDataTableServer, VDialog, VLocaleProvider } from 'vuetify/components';

import { vuetify } from '@/plugins/vuetify';

import { i18n } from '@/i18n';

describe('Vuetify locale', () => {
  afterEach(() => {
    i18n.global.locale.value = 'en';
  });

  it('translates table pagination and empty-state messages inside dialogs after a locale change', async () => {
    const wrapper = mount({
      render: () => h(VLocaleProvider, null, {
        default: () => h(VDialog, { modelValue: true, transition: false }, {
          default: () => h(VDataTableServer, {
            headers: [{ title: 'Name', key: 'name' }],
            items: [],
            itemsLength: 0,
          }),
        }),
      }),
    }, {
      global: { plugins: [i18n, vuetify] },
    });

    await flushPromises();
    expect(document.body.textContent).toContain('No data available');
    expect(document.body.textContent).toContain('Items per page');

    i18n.global.locale.value = 'fr';
    await flushPromises();

    expect(document.body.textContent).toContain('Aucune donnée disponible');
    expect(document.body.textContent).toContain('Éléments par page');
    expect(document.body.textContent).toContain('0-0 de 0');
    expect(document.body.textContent).not.toContain('No data available');
    wrapper.unmount();
  });
});
