import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
};

describe('c-compiled-template', () => {
  const snapshotFactory = generateRenderer(CCompiledTemplate, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `c-compiled-template` after mount', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<span>mount</span>',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-compiled-template` with simple template', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<span>template</span>',
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-compiled-template` with context', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<span>{{value1}}</span>',
        context: { value1: 'context-value-1' },
        parentElement: 'span',
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-compiled-template` after update template', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<span>template</span>',
      },
    });

    wrapper.setProps({ template: '<div>updated template</div>' });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-compiled-template` after update context', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<span>{{value1}}</span>',
        context: { value1: 'context-value-1' },
      },
    });

    wrapper.setProps({ context: { value1: 'updated-context-value-1' } });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-compiled-template` after update parentElement', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<span>template</span>',
      },
    });

    wrapper.setProps({ parentElement: 'section' });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
