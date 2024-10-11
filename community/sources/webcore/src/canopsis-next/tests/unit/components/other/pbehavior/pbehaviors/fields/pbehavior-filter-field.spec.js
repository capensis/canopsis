import Faker from 'faker';
import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import { MODALS } from '@/constants';

import PbehaviorFilterField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-filter-field.vue';

const stubs = {};

const selectCreateFilterButton = wrapper => wrapper.find('v-btn-stub');

describe('pbehavior-filter-field', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(PbehaviorFilterField, {

    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(PbehaviorFilterField, {

    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Filter created after trigger create button', () => {
    const patterns = {
      entity_pattern: {
        groups: [],
      },
    };
    const wrapper = factory({
      propsData: {
        patterns,
      },
    });

    selectCreateFilterButton(wrapper).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPatterns,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          patterns,
          withEntity: true,
          action: expect.any(Function),
        },
      },
    );
    const [{ config }] = $modals.show.mock.calls[0];

    const newPatterns = {
      entity_pattern: {
        groups: [{
          rules: [{
            attribute: Faker.datatype.string(),
            value: Faker.datatype.string(),
          }],
        }],
      },
    };

    config.action(newPatterns);

    expect(wrapper).toEmit('input', newPatterns);
  });

  test('Filter created after trigger create button', () => {
    const patterns = {
      entity_pattern: {
        groups: [],
      },
    };
    const wrapper = factory({
      propsData: {
        patterns,
      },
    });

    selectCreateFilterButton(wrapper).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPatterns,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          patterns,
          withEntity: true,
          action: expect.any(Function),
        },
      },
    );
    const [{ config }] = $modals.show.mock.calls[0];

    const newPatterns = {
      entity_pattern: {
        groups: [{
          rules: [{
            attribute: Faker.datatype.string(),
            value: Faker.datatype.string(),
          }],
        }],
      },
    };

    config.action(newPatterns);

    expect(wrapper).toEmit('input', newPatterns);
  });

  test('Renders `pbehavior-filter-field` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          entity_pattern: {
            groups: [],
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  test('Renders `pbehavior-filter-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          entity_pattern: {
            groups: [{
              rules: [{
                attribute: 'test',
                value: 'test-value',
              }],
            }],
          },
        },
        patternsFieldName: 'patterns-field-name',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  test('Renders `pbehavior-filter-field` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          entity_pattern: {
            groups: [{
              rules: [{
                attribute: 'test',
                value: 'test-value',
              }],
            }],
          },
        },
        patternsFieldName: 'patterns-validate',
      },
    });

    await wrapper.setProps({
      patterns: {
        entity_pattern: {
          groups: [],
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
