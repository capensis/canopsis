import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import DateTimePickerMenu from '@/components/forms/fields/date-time-picker/date-time-picker-menu.vue';

const stubs = {
  'date-time-picker': true,
  'v-menu': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: `
      <div class="v-menu">
        <div class="v-menu__activator" @click="$emit('update:modelValue', true)">
          <slot name="activator" :props="{}" />
        </div>
        <div v-if="modelValue" class="v-menu__content">
          <slot />
        </div>
      </div>
    `,
  },
};

const selectMenuActivator = wrapper => wrapper.find('.v-menu__activator');
const selectDateTimePicker = wrapper => wrapper.find('date-time-picker-stub');

describe('date-time-picker-menu', () => {
  const factory = generateShallowRenderer(DateTimePickerMenu, { stubs });

  test('Value changed after trigger input on the date time picker', async () => {
    const date = new Date(1387435600000);
    const wrapper = factory({
      propsData: {
        value: date,
      },
    });

    await selectMenuActivator(wrapper).trigger('click');

    const newDate = new Date(1387536600000);

    await selectDateTimePicker(wrapper).triggerCustomEvent('update:value', newDate);

    expect(wrapper).toEmitInput(newDate);
  });

  test('Menu opened after trigger button', async () => {
    const wrapper = factory();

    expect(selectDateTimePicker(wrapper).exists()).toBe(false);

    await selectMenuActivator(wrapper).trigger('click');
    await flushPromises();

    expect(selectDateTimePicker(wrapper).exists()).toBe(true);
  });

  test('Menu closed after trigger close event on the date time picker', async () => {
    const wrapper = factory();

    await selectMenuActivator(wrapper).trigger('click');
    await flushPromises();

    selectDateTimePicker(wrapper).triggerCustomEvent('close');

    await flushPromises();

    expect(wrapper.vm.opened).toBe(false);
    expect(selectDateTimePicker(wrapper).exists()).toBe(false);
  });
});
