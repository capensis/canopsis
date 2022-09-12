import { mount, createVueInstance } from '@unit/utils/vue';

import ExtraDetailsChildren from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-children.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsChildren, {
  localVue,

  ...options,
});

describe('extra-details-children', () => {
  const total = 3;
  const rule = {
    name: 'rule-name',
  };

  it('Renders `extra-details-children` with full children and rule', () => {
    const wrapper = snapshotFactory({
      propsData: {
        total,
        rule,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
