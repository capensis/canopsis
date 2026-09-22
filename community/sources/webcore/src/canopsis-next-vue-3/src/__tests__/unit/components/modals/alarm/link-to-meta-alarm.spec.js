import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { ALARM_STATES } from '@/constants';
import ClickOutside from '@/services/click-outside';

import LinkToMetaAlarm from '@/components/modals/alarm/link-to-meta-alarm.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'alarm-general-table': true,
  'link-meta-alarm-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('link-to-meta-alarm', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(LinkToMetaAlarm, {
    stubs,
    attachTo: document.body,
    parentComponent: { provide: { $clickOutside: new ClickOutside() } },
    mocks: { $modals },
  });

  const alarm = { _id: 'alarm-id', v: { state: { val: ALARM_STATES.major } } };
  const createWrapper = (config = {}) => factory({
    propsData: {
      modal: {
        id: 'test-modal-id',
        config,
      },
    },
  });

  test('initializes the default link form', () => {
    const wrapper = createWrapper();

    expect(wrapper.vm.values).toEqual({
      metaAlarm: null,
      comment: '',
      component: '',
      resource: '',
      auto_resolve: false,
      tags: { copy_from_children: false, filter_by_label: [] },
      infos: [],
    });
  });

  test('maps the form and alarm ids before submitting', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action, items: [alarm] });

    wrapper.vm.updateForm({
      metaAlarm: { _id: 'meta-alarm-id' },
      comment: 'Comment',
      component: 'component',
      resource: 'resource',
      auto_resolve: true,
      tags: { copy_from_children: false, filter_by_label: [] },
      infos: [],
    });
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({
      id: 'meta-alarm-id',
      alarms: ['alarm-id'],
      comment: 'Comment',
      component: 'component',
      resource: 'resource',
      auto_resolve: true,
      tags: { copy_from_children: false, filter_by_label: [] },
      infos: [],
    });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('submits an empty alarm list when no alarms are eligible for linking', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action, items: [{ _id: 'ok-alarm', v: { state: { val: ALARM_STATES.ok } } }] });

    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith(expect.objectContaining({ alarms: [] }));
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('closes without an action when cancel is clicked', async () => {
    const wrapper = createWrapper();

    await wrapper.findAll('button.v-btn').at(0).trigger('click');
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });
});
