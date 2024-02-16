import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ENTITIES_STATUSES } from '@/constants';

import CEntityStatusField from '@/components/forms/fields/entity/c-entity-status-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('c-entity-status-field', () => {
  const factory = generateShallowRenderer(CEntityStatusField, { stubs });
  const snapshotFactory = generateRenderer(CEntityStatusField);

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: ENTITIES_STATUSES.closed,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.vm.$emit('input', ENTITIES_STATUSES.cancelled);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(ENTITIES_STATUSES.cancelled);
  });

  it('Renders `c-entity-status-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ENTITIES_STATUSES.stealthy,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-entity-status-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ENTITIES_STATUSES.flapping,
        label: 'Custom label',
        name: 'customAlarmStatusName',
        disabled: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
