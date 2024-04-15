import { flushPromises, generateRenderer } from '@unit/utils/vue';

import ExtraDetailsChildren from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-children.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
};

const selectExtraDetailsIcon = wrapper => wrapper.find('c-alarm-extra-details-chip-stub');

describe('extra-details-children', () => {
  const total = 3;
  const opened = 2;
  const closed = 1;
  const rule = {
    name: 'rule-name',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsChildren, {
    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-children` with full children and rule', async () => {
    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      propsData: {
        total,
        opened,
        closed,
        rule,
      },
    });

    await flushPromises();

    selectExtraDetailsIcon(wrapper).trigger('mouseenter');

    jest.runAllTimers();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
