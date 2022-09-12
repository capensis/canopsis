import { mount, createVueInstance } from '@unit/utils/vue';

import ExtraDetailsParents from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-parents.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsParents, {
  localVue,
  attachTo: document.body,

  ...options,
});

describe('extra-details-parents', () => {
  const total = 3;
  const rules = [
    {
      id: 'parent-rule-1-id',
      name: 'parent-rule-1-name',
    },
    {
      id: 'parent-rule-2-id',
      name: 'parent-rule-2-name',
    },
    {
      id: 'parent-rule-3-id',
      name: 'parent-rule-3-name',
    },
  ];

  it('Renders `extra-details-parents` with full parents', () => {
    const wrapper = snapshotFactory({
      propsData: {
        total,
        rules,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-parents` without rules', () => {
    const wrapper = snapshotFactory({
      propsData: {
        total,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
