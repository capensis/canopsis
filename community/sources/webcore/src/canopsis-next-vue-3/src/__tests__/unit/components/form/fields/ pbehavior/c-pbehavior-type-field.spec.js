import { createSelectInputStub } from '@unit/stubs/input';
import { createMockedStoreModules, createPbehaviorTypesModule } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CPbehaviorTypeField from '@/components/forms/fields/pbehavior/c-pbehavior-type-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-pbehavior-type-field', () => {
  const { pbehaviorTypesModule, fieldPbehaviorTypes } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([pbehaviorTypesModule]);

  const factory = generateShallowRenderer(CPbehaviorTypeField, { stubs, store });

  test('Type changed after trigger select field', () => {
    const type = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };
    fieldPbehaviorTypes.mockReturnValueOnce([type]);

    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('input', type._id);

    expect(wrapper).toEmitInput(type._id);
  });
});
