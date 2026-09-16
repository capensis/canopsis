import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { defineComponent } from 'vue';

import CRequestHelper from '@/components/common/runtime-template/c-request-helper.vue';

describe('c-request-helper', () => {
  const factory = generateRenderer(CRequestHelper);

  test('Does not render anything when template is empty', () => {
    const wrapper = factory();

    expect(wrapper.html()).toBe('');
  });

  test('Renders provided runtime template', async () => {
    const wrapper = factory({
      propsData: {
        template: '<strong class="runtime-value">Runtime text</strong>',
      },
    });

    await flushPromises();

    expect(wrapper.find('.runtime-value').text()).toBe('Runtime text');
  });

  test('Passes template props to runtime template', async () => {
    const wrapper = factory({
      propsData: {
        template: '<span class="runtime-name">{{ name }}</span>',
        templateProps: {
          name: 'Canopsis',
        },
      },
    });

    await flushPromises();

    expect(wrapper.find('.runtime-name').text()).toBe('Canopsis');
  });

  test('Passes parent data to runtime template', async () => {
    const parent = defineComponent({
      components: { CRequestHelper },
      data: () => ({
        status: 'ready',
      }),
      template:
        '<c-request-helper template="<span class=&quot;runtime-status&quot;>{{ status }}</span>" />',
    });
    const wrapper = generateRenderer(parent)();

    await flushPromises();

    expect(wrapper.find('.runtime-status').text()).toBe('ready');
  });

  test('Passes parent props to runtime template', async () => {
    const parent = defineComponent({
      components: { CRequestHelper },
      props: {
        title: {
          type: String,
          default: '',
        },
      },
      template:
        '<c-request-helper template="<span class=&quot;runtime-title&quot;>{{ title }}</span>" />',
    });
    const wrapper = generateRenderer(parent)({
      propsData: {
        title: 'Alarm title',
      },
    });

    await flushPromises();

    expect(wrapper.find('.runtime-title').text()).toBe('Alarm title');
  });

  test('Passes parent methods to runtime template', async () => {
    const parent = defineComponent({
      components: { CRequestHelper },
      methods: {
        getLabel() {
          return 'computed label';
        },
      },
      template:
        '<c-request-helper template="<span class=&quot;runtime-label&quot;>{{ getLabel() }}</span>" />',
    });
    const wrapper = generateRenderer(parent)();

    await flushPromises();

    expect(wrapper.find('.runtime-label').text()).toBe('computed label');
  });

  test('Template props override parent values', async () => {
    const parent = defineComponent({
      components: { CRequestHelper },
      data: () => ({
        label: 'parent label',
      }),
      template: `
        <c-request-helper
          template="<span class=&quot;runtime-label&quot;>{{ label }}</span>"
          :template-props="{ label: 'template label' }"
        />
      `,
    });
    const wrapper = generateRenderer(parent)();

    await flushPromises();

    expect(wrapper.find('.runtime-label').text()).toBe('template label');
  });
});
