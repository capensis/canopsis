import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';

import ExtraDetailsCauses from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-causes.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsCauses, {
  localVue,

  ...options,
});

describe('extra-details-causes', () => {
  const causes = {
    total: 3,
    rules: [
      {
        id: 'cause-rule-1-id',
        name: 'cause-rule-1-name',
      },
      {
        id: 'cause-rule-2-id',
        name: 'cause-rule-2-name',
      },
      {
        id: 'cause-rule-3-id',
        name: 'cause-rule-3-name',
      },
    ],
  };

  it('Renders `extra-details-causes` with full causes', () => {
    const wrapper = snapshotFactory({
      propsData: {
        causes,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-causes` without rules', () => {
    const wrapper = snapshotFactory({
      propsData: {
        causes: omit(causes, ['rules']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
