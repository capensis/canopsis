import { generateShallowRenderer } from '@unit/utils/vue';

import { CANOPSIS_EDITION, TRIGGERS_TYPES } from '@/constants';

import { useAppInfoStore } from '@/stores/modules/app-info';

import CTriggersField from '@/components/forms/fields/c-triggers-field.vue';

describe('c-triggers-field', () => {
  const factory = generateShallowRenderer(CTriggersField);

  it('Prepares grouped trigger options with output change triggers', async () => {
    const wrapper = factory();
    const appInfoStore = useAppInfoStore();

    appInfoStore.appInfo = { edition: CANOPSIS_EDITION.pro };
    await wrapper.vm.$nextTick();

    const triggerItems = wrapper.vm.preparedTriggers.filter(({ header }) => !header);

    expect(wrapper.vm.preparedTriggers.filter(({ header }) => header)).toEqual([
      { header: 'General' },
      { header: 'Pbehavior' },
      { header: 'Tickets' },
      { header: 'Instructions' },
    ]);
    expect(triggerItems).toContainEqual(
      expect.objectContaining({
        type: TRIGGERS_TYPES.changedoutput,
        text: 'Output changed',
      }),
    );
    expect(triggerItems).toContainEqual(
      expect.objectContaining({
        type: TRIGGERS_TYPES.changedlongoutput,
        text: 'Long output changed',
      }),
    );
    expect(triggerItems).toContainEqual(
      expect.objectContaining({
        type: TRIGGERS_TYPES.pause,
        text: 'Paused',
      }),
    );
    expect(triggerItems).toContainEqual(
      expect.objectContaining({
        type: TRIGGERS_TYPES.resume,
        text: 'Unpaused',
      }),
    );
  });

  it('Uses the alarm wording for selected chips', () => {
    const wrapper = factory();

    expect(
      wrapper.vm.getSelectedText({
        type: TRIGGERS_TYPES.changedoutput,
        text: 'Output changed',
      }),
    ).toBe('Alarm output changed');
    expect(
      wrapper.vm.getSelectedText({
        type: TRIGGERS_TYPES.resume,
        text: 'Unpause',
      }),
    ).toBe('Alarm unpaused');
  });

  it('Keeps custom trigger lists ungrouped', () => {
    const wrapper = factory({
      propsData: {
        types: [TRIGGERS_TYPES.eventscount],
        translationKeyPrefix: 'common.repeatTriggers',
      },
    });

    expect(wrapper.vm.preparedTriggers).toEqual([
      expect.objectContaining({
        type: TRIGGERS_TYPES.eventscount,
        text: 'Alarm received new check event',
      }),
    ]);
  });

  it('Uses unique values for group subheaders', () => {
    const wrapper = factory();

    expect(wrapper.vm.getTriggerItemValue({
      type: 'subheader',
      value: '__subheader__0',
    })).toBe('__subheader__0');
    expect(wrapper.vm.getTriggerItemValue({ type: TRIGGERS_TYPES.create })).toBe(
      TRIGGERS_TYPES.create,
    );
  });
});
