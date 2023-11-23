import { generateRenderer } from '@unit/utils/vue';

import CCollapsePanel from '@/components/common/block/c-collapse-panel.vue';

describe('c-collapse-panel', () => {
  const snapshotFactory = generateRenderer(CCollapsePanel);

  it('Renders `c-collapse-panel` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-collapse-panel` with custom title', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-collapse-panel` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        color: 'grey',
        icon: 'custom_icon',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
