import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';

const setValue = jest.fn();
const setSearch = jest.fn();

const stubs = {
  'v-combobox': {
    props: ['value'],
    template: `
      <input
        :value="value"
        class="v-combobox"
        @input="$listeners.change($event.target.value)"
      />
    `,
    methods: {
      setValue,
      setSearch,
    },
  },
};

const selectCombobox = wrapper => wrapper.find('.v-combobox, .v-autocomplete');

describe('time-picker-field', () => {
  const factory = generateShallowRenderer(TimePickerField, { stubs });
  const snapshotFactory = generateRenderer(TimePickerField, { attachTo: document.body });

  afterEach(() => {
    setValue.mockClear();
    setSearch.mockClear();
  });

  test('Combobox filter return true with exist text', async () => {
    const wrapper = factory();

    const combobox = selectCombobox(wrapper);

    expect(combobox.vm.filter({}, '12:', '12:13')).toBeTruthy();
  });

  test('Combobox filter return false with excluded text', async () => {
    const wrapper = factory();

    const combobox = selectCombobox(wrapper);

    expect(combobox.vm.filter({}, '13:', '12:13')).toBeFalsy();
  });

  test('Value changed after trigger change with empty value on the combobox', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
      },
    });

    selectCombobox(wrapper).setValue('');

    expect(setValue).toBeCalledWith('12:00');
    expect(setSearch).toBeCalledWith('');
    expect(wrapper).not.toHaveBeenEmit('input');
  });

  test('List scrolled correctly without value', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '',
      },
    });

    const combobox = selectCombobox(wrapper);

    await wrapper.activateAllMenus();

    const content = wrapper.findMenu();

    content.element.scrollTop = 200;

    expect(combobox.vm.menuProps.scrollCalculator(content.element)).toBe(200);
  });

  test('List scrolled correctly with value', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '12:00',
      },
    });

    const combobox = selectCombobox(wrapper);

    await wrapper.activateAllMenus();

    const content = wrapper.findMenu();

    content.element.scrollTop = 100;

    expect(combobox.vm.menuProps.scrollCalculator(content.element)).toBe(0);
  });

  test('Value changed after trigger change with valid value on the combobox', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
      },
    });

    const combobox = selectCombobox(wrapper);

    const newValue = '13:00';

    combobox.setValue(newValue);

    expect(wrapper).toEmitInput(newValue);

    expect(setValue).not.toBeCalled();
    expect(setSearch).not.toBeCalled();
  });

  test('Value changed after trigger change with round hours on the combobox', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
        roundHours: true,
      },
    });

    const combobox = selectCombobox(wrapper);

    const newValue = '13:13';

    combobox.setValue(newValue);

    expect(wrapper).toEmitInput('13:00');

    expect(setValue).toBeCalledWith('13:00');
    expect(setSearch).toBeCalledWith('');
  });

  test('Renders `time-picker-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `time-picker-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '12:15',
        label: 'label',
        stepsInHours: 2,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `time-picker-field` with rounded hours', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '12:00',
        roundHours: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
