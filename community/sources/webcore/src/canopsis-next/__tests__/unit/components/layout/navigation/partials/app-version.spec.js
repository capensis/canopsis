import { flushPromises, shallowMount } from '@vue/test-utils';
import { createPinia } from 'pinia';

import { useAppInfoStore } from '@/stores/modules/app-info';
import { useTemplateVarsStore } from '@/stores/modules/template-vars';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import AppVersion from '@/components/layout/navigation/partials/app-version.vue';

vi.unmock('@/stores/modules/template-vars');

describe('app-version', () => {
  test('renders configured variables in the version tooltip', async () => {
    const pinia = createPinia();

    useTemplateVarsStore(pinia).items = {
      Location: 'Paris',
      System: { ENV_TEST: 'production' },
    };
    useAppInfoStore(pinia).appInfo = {
      version: '26.10',
      edition: 'pro',
      serial_name: 'test-serial',
      version_updated: '2026-10-01',
      version_description:
        '{{env.Location}} / {{env.System.ENV_TEST}} / {{version}} / {{edition}} / {{serialName}} / {{versionUpdated}}',
    };

    const wrapper = shallowMount(AppVersion, {
      global: {
        plugins: [pinia],
        provide: { $system: { theme: {}, dark: false } },
        stubs: {
          'c-tooltip': {
            template: '<div><slot name="activator" :props="{}" /><slot /></div>',
          },
          'c-compiled-template': CCompiledTemplate,
          'c-runtime-template': {
            inheritAttrs: false,
            props: ['template'],
            template: '<div v-html="template" />',
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.text()).toContain(
      'Paris / production / 26.10 / pro / test-serial / 2026-10-01',
    );
  });
});
