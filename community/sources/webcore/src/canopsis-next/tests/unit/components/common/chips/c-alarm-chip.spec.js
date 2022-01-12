import { mount, createVueInstance } from '@unit/utils/vue';

import { ENTITIES_STATUSES, ENTITY_INFOS_TYPE, EVENT_ENTITY_TYPES } from '@/constants';
import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CAlarmChip, {
  localVue,

  ...options,
});

describe('c-alarm-chip', () => {
  it('Renders `c-alarm-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-chip` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 2,
        type: ENTITY_INFOS_TYPE.status,
        badgeValue: 3,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(Object.entries(ENTITIES_STATUSES))('Renders `c-alarm-chip` with status: %s', (_, value) => {
    const wrapper = snapshotFactory({
      propsData: {
        value,
        type: ENTITY_INFOS_TYPE.status,
        badgeValue: value + 10,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(Object.entries(EVENT_ENTITY_TYPES))('Renders `c-alarm-chip` with state: %s', (_, value) => {
    const wrapper = snapshotFactory({
      propsData: {
        value,
        type: ENTITY_INFOS_TYPE.state,
        badgeValue: value + 20,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
