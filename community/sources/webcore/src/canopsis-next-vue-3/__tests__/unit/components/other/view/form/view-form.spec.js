import { createInputStub } from '@unit/stubs/input';
import { randomPeriodicRefresh } from '@unit/utils/duration';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import ViewForm from '@/components/other/view/form/view-form.vue';

const stubs = {
  'view-duplicate-private-field': true,
  'c-name-field': true,
  'c-enabled-field': true,
  'periodic-refresh-field': true,
  'view-tags-field': true,
  'view-group-field': true,
  'c-text-field': createInputStub('c-text-field'),
};

const selectTitleField = wrapper => wrapper.find('c-name-field-stub');
const selectDescriptionField = wrapper => wrapper.find('input.c-text-field');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectPeriodicRefreshField = wrapper => wrapper.find('periodic-refresh-field-stub');
const selectTagsField = wrapper => wrapper.find('view-tags-field-stub');
const selectGroupField = wrapper => wrapper.find('view-group-field-stub');

describe('view-form', () => {
  const factory = generateShallowRenderer(ViewForm, { stubs });
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

    selectTitleField(wrapper).triggerCustomEvent('input', newTitle);

    expect(wrapper).toEmitInput({ ...form, title: newTitle });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newDescription = Faker.datatype.string();

    selectDescriptionField(wrapper).triggerCustomEvent('input', newDescription);

    expect(wrapper).toEmit('update:form', { ...form, description: newDescription });
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newEnabled = !form.enabled;

    selectEnabledField(wrapper).triggerCustomEvent('input', newEnabled);

    expect(wrapper).toEmitInput({ ...form, enabled: newEnabled });
  });

  it('Periodic refresh changed after trigger periodic refresh field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newPeriodicRefresh = randomPeriodicRefresh();

    selectPeriodicRefreshField(wrapper).triggerCustomEvent(
      'update:periodic-refresh',
      newPeriodicRefresh,
    );

    expect(wrapper).toEmit('update:form', { ...form, periodic_refresh: newPeriodicRefresh });
  });

  it('Tag changed after trigger tag field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTags = [Faker.datatype.string()];

    selectTagsField(wrapper).triggerCustomEvent('input', newTags);

    expect(wrapper).toEmitInput({ ...form, tags: newTags });
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

    selectGroupField(wrapper).triggerCustomEvent('input', newGroup);

    expect(wrapper).toEmitInput({ ...form, group: newGroup });
  });
});
