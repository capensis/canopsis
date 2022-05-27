import { mount, createVueInstance } from '@unit/utils/vue';
import { createModalWrapperStub } from '@unit/stubs/modal';

import AlarmsList from '@/components/modals/alarm/alarms-list.vue';

const localVue = createVueInstance();

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarms-list-widget': true,
};

const snapshotFactory = (options = {}) => mount(AlarmsList, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('alarms-list', () => {
  test('Renders `alarms-list`', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            widget: {},
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
