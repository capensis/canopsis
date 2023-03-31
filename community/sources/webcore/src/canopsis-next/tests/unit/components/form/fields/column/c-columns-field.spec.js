import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';
import { ENTITIES_TYPES } from '@/constants';

import CColumnsField from '@/components/forms/fields/column/c-columns-field.vue';

const snapshotStubs = {
  'c-movable-card-iterator-field': {
    props: ['value'],
    template: `
      <div class="c-movable-card-iterator-field">
        <div v-for="(item, index) in value" :key="index">
          <slot name="item" :item="item" :index="index" />
        </div>
        <button class="add-item-btn" type="button" @click="$emit('add')" />
      </div>
    `,
  },
  'column-field': true,
};
const stubs = {
  ...snapshotStubs,
  'v-btn': createButtonStub('v-btn'),
};

const selectAddCardButton = wrapper => wrapper.find('.add-item-btn');
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
  const snapshotFactory = generateRenderer(CColumnsField, { stubs: snapshotStubs });

  test('Column added after trigger add event', () => {
    const wrapper = factory({
      propsData: {
        columns,
      },
    });

    selectAddCardButton(wrapper).trigger('click');

    expect(wrapper).toEmit('input', [
      ...columns,
      {
        column: '',
        key: expect.any(String),
        label: '',
      },
    ]);
  });

  test('Column changed after trigger column field', () => {
    const wrapper = factory({
      propsData: {
        columns,
      },
    });

    const newColumn = {
      ...columns[2],
      label: Faker.datatype.string(),
    };

    selectColumnFieldByIndex(wrapper, 2).vm.$emit('input', newColumn);

    expect(wrapper).toEmit('input', [
      columns[0],
      columns[1],
      newColumn,
      columns[3],
    ]);
  });

  test('Renders `c-columns-field` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-columns-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        columns: [
          { id: 'id-1', key: 'key-1', label: 'label-1' },
          { id: 'id-2', key: 'key-1', label: 'label-2' },
        ],
        type: ENTITIES_TYPES.entity,
        withTemplate: true,
        withHtml: true,
        withColorIndicator: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
