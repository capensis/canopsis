import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

describe('c-compiled-template', () => {
  const factory = generateShallowRenderer(CCompiledTemplate, {
    stubs: {
      'c-runtime-template': {
        inheritAttrs: false,
        props: ['template', 'parent'],
        template: '<div />',
      },
    },
    parentComponent: {
      provide: {
        $system: {
          theme: {},
          dark: false,
        },
      },
    },
  });

  it('sanitizes and normalizes compiled HTML without a separate normalization pass', async () => {
    const wrapper = factory({
      propsData: {
        template: '<div><span>{{ value }}</div>',
        context: { value: 'content' },
        parentElement: 'section',
      },
    });

    await flushPromises();

    expect(wrapper.vm.compiledTemplate).toBe(
      '<section><div><span>content</span></div></section>',
    );
  });
});
