import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_STATUSES, ENTITY_TYPES } from '@/constants';

import ServiceEntityHeader from '@/components/other/service/partials/service-entity-header.vue';

const stubs = {
  'c-no-events-icon': true,
};

const selectAlert = wrapper => wrapper.find('v-alert-stub');

describe('service-entity-header', () => {
  const snapshotFactory = generateRenderer(ServiceEntityHeader, { stubs });
  const factory = generateShallowRenderer(ServiceEntityHeader, { stubs });

  test('Alert removed after trigger alert', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        lastActionUnavailable: true,
      },
    });

    const alert = selectAlert(wrapper);

    alert.triggerCustomEvent('input');

    expect(wrapper).toHaveBeenEmit('remove:unavailable');
  });

  test('Renders `service-entity-header` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-entity-header` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {
          _id: 'service-id',
          alarm_id: 'service-alarm-id',
          source_type: ENTITY_TYPES.service,
          ack: {},
          ticket: {},
          status: {
            val: ALARM_STATUSES.cancelled,
          },
        },
        selected: true,
        selectable: true,
        lastActionUnavailable: true,
        entityNameField: 'entity_name_field',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
