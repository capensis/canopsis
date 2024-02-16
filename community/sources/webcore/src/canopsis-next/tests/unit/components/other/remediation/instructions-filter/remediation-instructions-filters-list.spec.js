import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import RemediationInstructionsFiltersList from '@/components/other/remediation/instructions-filter/remediation-instructions-filters-list.vue';

const stubs = {
  'remediation-instructions-filters-item': true,
};

const selectRemediationInstructionsFiltersItemsField = wrapper => wrapper.findAll('remediation-instructions-filters-item-stub');

describe('remediation-instructions-filters-list', () => {
  const filters = [{
    with: true,
    all: true,
    auto: true,
    manual: true,
    instructions: [],
    _id: 'id1',
  }, {
    with: false,
    all: true,
    auto: false,
    manual: false,
    instructions: [{}],
    _id: 'id2',
  }];

  const factory = generateShallowRenderer(RemediationInstructionsFiltersList, { stubs });
  const snapshotFactory = generateRenderer(RemediationInstructionsFiltersList, {
    stubs,
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

  it('Instruction filters changed after trigger instruction filters item field', () => {
    const wrapper = factory({
      propsData: {
        filters,
      },
    });

    const remediationInstructionsFiltersItemField = selectRemediationInstructionsFiltersItemsField(wrapper)
      .at(0);

    const updatedFilter = {
      with: true,
      all: false,
      auto: true,
      manual: true,
      instructions: [],
      _id: 'id1',
    };

    remediationInstructionsFiltersItemField.vm.$emit('input', updatedFilter);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      updatedFilter,
      ...filters.slice(1),
    ]);
  });

  it('Instruction filter removed after trigger remove event', () => {
    const wrapper = factory({
      propsData: {
        filters,
      },
    });

    const remediationInstructionsFiltersItemField = selectRemediationInstructionsFiltersItemsField(wrapper)
      .at(1);

    remediationInstructionsFiltersItemField.vm.$emit('remove');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData).toEqual(filters.slice(0, -1));
  });

  it('Renders `remediation-instructions-filters-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `remediation-instructions-filters-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        editable: true,
        closable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
