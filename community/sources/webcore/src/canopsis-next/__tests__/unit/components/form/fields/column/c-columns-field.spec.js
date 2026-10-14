import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CColumnsField from '@/components/forms/fields/column/c-columns-field.vue';

const stubs = {
  'c-card-iterator-form': {
    props: {
      value: {
        type: Array,
        default: () => [],
      },
      addButtonJustify: {
        type: String,
        default: '',
      },
    },
    template: `
      <ul class="c-card-iterator-form" :data-add-button-justify="addButtonJustify">
        <li v-for="(item, index) in value" :key="item.key">
          <slot :item="item" :index="index" name="item" />
        </li>
        <button class="c-card-iterator-form__add" @click="$emit('add')" />
      </ul>
    `,
  },
  'column-field': true,
};

const selectAddCardButton = wrapper => wrapper.find('button.c-card-iterator-form__add');
const selectColumnFields = wrapper => wrapper.findAll('column-field-stub');
const selectColumnFieldByIndex = (wrapper, index) => selectColumnFields(wrapper).at(index);

describe('c-columns-field', () => {
  const columns = [
    { key: 'key-1', label: 'column-1' },
    { key: 'key-2', label: 'column-2' },
    { key: 'key-3', label: 'column-3' },
    { key: 'key-4', label: 'column-4' },
  ];

  const factory = generateShallowRenderer(CColumnsField, { stubs });

  test('Add button is aligned to the right', () => {
    const wrapper = factory();

    expect(wrapper.find('.c-card-iterator-form').attributes('data-add-button-justify')).toBe('end');
  });

  test('Column added after trigger add event', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: columns,
      },
    });

    await selectAddCardButton(wrapper).trigger('click');

    expect(wrapper).toEmitInput([
      ...columns,
      {
        column: '',
        key: expect.any(String),
        label: '',
        isHtml: false,
      },
    ]);
  });

  test('Column changed after trigger column field', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: columns,
      },
    });

    const newColumn = {
      ...columns[2],
      label: Faker.datatype.string(),
    };

    await selectColumnFieldByIndex(wrapper, 2).triggerCustomEvent('update:modelValue', newColumn);

    expect(wrapper).toEmitInput([columns[0], columns[1], newColumn, columns[3]]);
  });

  test('Column removed after trigger remove event', async () => {
    const columnToRemoveIndex = 2;
    const wrapper = factory({
      propsData: {
        modelValue: columns,
      },
    });

    await selectColumnFieldByIndex(wrapper, columnToRemoveIndex).triggerCustomEvent('remove');

    expect(wrapper).toEmitInput([columns[0], columns[1], columns[3]]);
  });
});
