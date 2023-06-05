import { mount, createVueInstance } from '@unit/utils/vue';

import CCollapsePanel from '@/components/common/block/c-collapse-panel.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CCollapsePanel, {
  localVue,

  ...options,
});

describe('c-collapse-panel', () => {
  it('Renders `c-collapse-panel` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-collapse-panel` with custom title', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-collapse-panel` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title with errors',
      },
    });

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => wrapper.vm,
      vm: wrapper.vm,
    });

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-collapse-panel` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        color: 'grey',
        icon: 'custom_icon',
        lazy: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
