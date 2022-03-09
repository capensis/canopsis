import { mount, createVueInstance } from '@unit/utils/vue';

import ExtraDetailsConsequences from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-consequences.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsConsequences, {
  localVue,

  ...options,
});

describe('extra-details-consequences', () => {
  const consequences = {
    total: 3,
  };
  const rule = {
    name: 'rule-name',
  };

  it('Renders `extra-details-consequences` with full consequences and rule', () => {
    const wrapper = snapshotFactory({
      propsData: {
        consequences,
        rule,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
