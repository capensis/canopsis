import { generateShallowRenderer } from '@unit/utils/vue';

import CColorPickerField from '@/components/forms/fields/color/c-color-picker-field.vue';

const stubs = {
  'v-input': {
    props: ['hideDetails'],
    template: '<div :data-hide-details="hideDetails"><slot /></div>',
  },
  'v-menu': {
    template: '<div><slot name="activator" :props="{}" /><slot /></div>',
  },
  'v-card': {
    template: '<div><slot /></div>',
  },
  'v-color-picker': true,
};

const selectCustomActivator = wrapper => wrapper.find('.custom-color-picker-activator');
const selectColorButton = wrapper => wrapper.find('c-btn-with-error-stub');
const selectColorButtonWrapper = wrapper => wrapper.find('.color-picker__button__wrapper');
const selectInput = wrapper => wrapper.find('[data-hide-details]');

describe('c-color-picker-field', () => {
  const factory = generateShallowRenderer(CColorPickerField, { stubs });

  test('Activator slot replaces the default button', () => {
    const wrapper = factory({
      propsData: {
        label: 'Critical',
        modelValue: '#c62828',
      },
      slots: {
        activator: `
          <template #activator="{ props, color, label, style, disabled }">
            <button
              v-bind="props"
              :style="style"
              :disabled="disabled"
              class="custom-color-picker-activator"
            >
              {{ label }} {{ color }}
            </button>
          </template>
        `,
      },
    });

    expect(selectCustomActivator(wrapper).text()).toBe('Critical #c62828');
    expect(selectColorButton(wrapper).exists()).toBeFalsy();
  });

  test('Default mode displays only the colored button', () => {
    const wrapper = factory({
      propsData: {
        label: 'Critical',
        modelValue: '#c62828',
      },
    });

    expect(selectCustomActivator(wrapper).exists()).toBeFalsy();
    expect(selectColorButton(wrapper).exists()).toBeTruthy();
  });

  test('Activator layout fills the available width when requested', () => {
    const wrapper = factory({
      propsData: {
        fillWidth: true,
      },
    });
    const colorButtonWrapper = selectColorButtonWrapper(wrapper);

    expect(colorButtonWrapper.classes()).toContain('w-100');
    expect(colorButtonWrapper.element.parentElement.classList).toContain('w-100');
  });

  test('Activator layout keeps its content width by default', () => {
    const wrapper = factory();
    const colorButtonWrapper = selectColorButtonWrapper(wrapper);

    expect(colorButtonWrapper.classes()).not.toContain('w-100');
    expect(colorButtonWrapper.element.parentElement.classList).not.toContain('w-100');
  });

  test('Hides input details when requested', () => {
    const wrapper = factory({
      propsData: {
        hideDetails: true,
      },
    });

    expect(selectInput(wrapper).attributes('data-hide-details')).toBe('true');
  });
});
