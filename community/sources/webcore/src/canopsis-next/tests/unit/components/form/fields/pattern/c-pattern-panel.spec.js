import { createVueInstance, mount } from '@unit/utils/vue';

import CPatternPanel from '@/components/forms/fields/pattern/c-pattern-panel.vue';

const localVue = createVueInstance();

const stubs = {
  'c-collapse-panel': {
    template: `
      <div class="c-collapse-panel">
        <slot name="header" />
        <slot />
      </div>
    `,
  },
};

const snapshotFactory = (options = {}) => mount(CPatternPanel, {
  localVue,
  stubs,

  ...options,
});

const selectCollapsePanel = wrapper => wrapper.find('.c-collapse-panel');

describe('c-pattern-panel', () => {
  it('Renders `c-pattern-panel` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pattern-panel` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title with errors',
      },
    });

    const collapsePanel = selectCollapsePanel(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => collapsePanel.vm,
      vm: collapsePanel.vm,
    });

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});
