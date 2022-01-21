import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';

import flushPromises from 'flush-promises';
import AlarmColumnCellPopupBody from '@/components/widgets/alarm/columns-formatting/alarm-column-cell-popup-body.vue';

const localVue = createVueInstance();

const stubs = {
  'v-runtime-template': true,
};

const factory = (options = {}) => shallowMount(AlarmColumnCellPopupBody, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmColumnCellPopupBody, {
  localVue,
  stubs,

  ...options,
});

const selectCloseButton = wrapper => wrapper.find('v-btn-stub');

describe('alarm-column-cell-popup-body', () => {
  it('Popup template closed after click on the button', async () => {
    const wrapper = factory({
      propsData: {
        template: '',
        alarm: {},
      },
    });

    const closeButton = selectCloseButton(wrapper);

    closeButton.vm.$emit('click');

    const closeEvents = wrapper.emitted('close');

    expect(closeEvents).toHaveLength(1);
  });

  it('Renders `alarm-column-cell-popup-body` with full alarm', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<p>Test</p><p>{{ alarm.name }}</p><p>{{ entity.name }}</p>',
        alarm: {
          name: 'alarm-name',
          entity: {
            name: 'entity-name',
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
