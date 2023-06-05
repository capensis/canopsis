import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { ENTITIES_STATUSES, ENTITY_TYPES } from '@/constants';

import ServiceEntityHeader from '@/components/other/service/partials/service-entity-header.vue';

const localVue = createVueInstance();

const stubs = {
  'c-no-events-icon': true,
};

const selectAlert = wrapper => wrapper.find('v-alert-stub');

describe('service-entity-header', () => {
  const snapshotFactory = generateRenderer(ServiceEntityHeader, { localVue, stubs });
  const factory = generateShallowRenderer(ServiceEntityHeader, { localVue, stubs });

  test('Alert removed after trigger alert', () => {
    const wrapper = factory({
      propsData: {
        entity: {},
        lastActionUnavailable: true,
      },
    });

    const alert = selectAlert(wrapper);

    alert.vm.$emit('input');

    expect(wrapper).toEmit('remove:unavailable');
  });

  test('Renders `service-entity-header` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        entity: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
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
            val: ENTITIES_STATUSES.cancelled,
          },
        },
        selected: true,
        selectable: true,
        lastActionUnavailable: true,
        entityNameField: 'entity_name_field',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
