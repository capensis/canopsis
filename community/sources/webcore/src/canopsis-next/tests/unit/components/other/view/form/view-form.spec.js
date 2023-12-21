import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { randomPeriodicRefresh } from '@unit/utils/duration';

import ViewForm from '@/components/other/view/form/view-form.vue';

const snapshotStubs = {
  'view-duplicate-private-field': true,
  'c-name-field': true,
  'c-enabled-field': true,
  'periodic-refresh-field': true,
  'view-tags-field': true,
  'view-group-field': true,
};
const stubs = {
  ...snapshotStubs,
  'v-text-field': createInputStub('v-text-field'),
};

const selectTitleField = wrapper => wrapper.find('c-name-field-stub');
const selectDescriptionField = wrapper => wrapper.find('input.v-text-field');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectPeriodicRefreshField = wrapper => wrapper.find('periodic-refresh-field-stub');
const selectTagsField = wrapper => wrapper.find('view-tags-field-stub');
const selectGroupField = wrapper => wrapper.find('view-group-field-stub');

describe('view-form', () => {
  const factory = generateShallowRenderer(ViewForm, { stubs });
  const snapshotFactory = generateRenderer(ViewForm, { stubs: snapshotStubs });
  const form = {
    title: Faker.datatype.string(),
    description: Faker.datatype.string(),
    enabled: Faker.datatype.boolean(),
    periodic_refresh: randomPeriodicRefresh(),
    tags: [Faker.datatype.string()],
  };

  test('Title changed after trigger title field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTitle = Faker.datatype.string();

    selectTitleField(wrapper).vm.$emit('input', newTitle);

    expect(wrapper).toEmit('input', { ...form, title: newTitle });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newDescription = Faker.datatype.string();

    selectDescriptionField(wrapper).vm.$emit('input', newDescription);

    expect(wrapper).toEmit('input', { ...form, description: newDescription });
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newEnabled = !form.enabled;

    selectEnabledField(wrapper).vm.$emit('input', newEnabled);

    expect(wrapper).toEmit('input', { ...form, enabled: newEnabled });
  });

  it('Periodic refresh changed after trigger periodic refresh field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newPeriodicRefresh = randomPeriodicRefresh();

    selectPeriodicRefreshField(wrapper).vm.$emit('input', newPeriodicRefresh);

    expect(wrapper).toEmit('input', { ...form, periodic_refresh: newPeriodicRefresh });
  });

  it('Tag changed after trigger tag field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTags = [Faker.datatype.string()];

    selectTagsField(wrapper).vm.$emit('input', newTags);

    expect(wrapper).toEmit('input', { ...form, tags: newTags });
  });

  it('Group changed after trigger group field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newGroup = {
      _id: Faker.datatype.string(),
    };

    selectGroupField(wrapper).vm.$emit('input', newGroup);

    expect(wrapper).toEmit('input', { ...form, group: newGroup });
  });

  test('Renders `view-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `view-form` with custom props', () => {
    const group = { _id: 'group' };

    const wrapper = snapshotFactory({
      propsData: {
        form: {
          title: 'View title',
          description: 'View description',
          enabled: true,
          periodic_refresh: {},
          tags: ['View tag'],
          group,
        },
        groups: [group],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
