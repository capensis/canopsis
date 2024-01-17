import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { ENTITIES_STATUSES, ENTITY_INFOS_TYPE, EVENT_ENTITY_TYPES } from '@/constants';

import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';

const selectChipContainer = wrapper => wrapper.find('.chip-container');

describe('c-alarm-chip', () => {
  const snapshotFactory = generateRenderer(CAlarmChip);
  const factory = generateShallowRenderer(CAlarmChip);

  it('Click event emitted', () => {
    const wrapper = factory();

    selectChipContainer(wrapper).trigger('click');

    expect(wrapper).toEmit('click');
  });

  it('Renders `c-alarm-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-alarm-chip` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 2,
        type: ENTITY_INFOS_TYPE.status,
        badgeValue: 3,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it.each(Object.entries(ENTITIES_STATUSES))('Renders `c-alarm-chip` with status: %s', (_, value) => {
    const wrapper = snapshotFactory({
      propsData: {
        value,
        type: ENTITY_INFOS_TYPE.status,
        badgeValue: value + 10,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it.each(Object.entries(EVENT_ENTITY_TYPES))('Renders `c-alarm-chip` with state: %s', (_, value) => {
    const wrapper = snapshotFactory({
      propsData: {
        value,
        type: ENTITY_INFOS_TYPE.state,
        badgeValue: value + 20,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
