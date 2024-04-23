import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CAdvancedSearchField from '@/components/common/search/c-advanced-search-field.vue';

const mockData = {
  search: Faker.lorem.words(),
  tooltipText: Faker.lorem.words(),
  column: {
    text: 'Column label',
    value: 'column.value',
  },
  field: Faker.lorem.word(),
};

const stubs = {
  'advanced-search-field': true,
};
const snapshotStubs = {
  'advanced-search-field': true,
};

describe.skip('c-advanced-search-field', () => {
  const factory = generateShallowRenderer(CAdvancedSearchField, { stubs });
  const snapshotFactory = generateRenderer(CAdvancedSearchField, { stubs: snapshotStubs });

  it('Submit search button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { query: { search } } });

    wrapper.findRoot().$emit('submit', search);

    expect(wrapper).toEmit('update:query', {
      search,
      page: 1,
    });
  });

  it('Submit search with search by column without columns', () => {
    const { column, search } = mockData;
    const advancedSearch = `- ${column.text} = '${search}'`;

    const wrapper = factory({
      propsData: {
        query: { search: advancedSearch },
      },
    });
    wrapper.findRoot().$emit('submit', advancedSearch);

    expect(wrapper).toEmit('update:query', {
      search: `${column.text} = '${search}'`,
      page: 1,
    });
  });

  it('Submit search with search by column', () => {
    const { column, search } = mockData;
    const advancedSearch = `- ${column.text} = '${search}'`;

    const wrapper = factory({
      propsData: {
        query: { search: advancedSearch },
        columns: [column],
      },
    });

    wrapper.findRoot().$emit('submit', advancedSearch);

    expect(wrapper).toEmit('update:query', {
      search: `${column.value} = '${search}'`,
      page: 1,
    });
  });

  it('Clear search by click on button', () => {
    const { search } = mockData;
    const wrapper = factory({
      propsData: {
        query: { search },
      },
    });

    wrapper.findRoot().$emit('clear');

    expect(wrapper).toEmit('update:query', {
      page: 1,
    });
  });

  it('Submit search with custom field name', () => {
    const { search, field } = mockData;
    const wrapper = factory({
      propsData: {
        query: { [field]: '' },
        field,
      },
    });

    wrapper.findRoot().$emit('submit', search);

    expect(wrapper).toEmit('update:query', {
      page: 1,
      [field]: search,
    });
  });

  it('Clear search with custom field name', () => {
    const { search, field } = mockData;
    const wrapper = factory({
      propsData: {
        query: { [field]: search },
        field,
      },
    });

    wrapper.findRoot().$emit('clear');

    expect(wrapper).toEmit('update:query', {
      page: 1,
    });
  });

  it('Toggle pin event works ', () => {
    const { search, field } = mockData;
    const wrapper = factory({
      propsData: {
        query: { [field]: search },
        field,
      },
    });

    wrapper.findRoot().$emit('toggle-pin', search);

    expect(wrapper).toEmit('toggle-pin', search);
  });

  it('Remove event works ', () => {
    const { search, field } = mockData;
    const wrapper = factory({
      propsData: {
        query: { [field]: search },
        field,
      },
    });

    wrapper.findRoot().$emit('remove', search);

    expect(wrapper).toEmit('remove', search);
  });

  it('Renders `c-advanced-search-field` correctly', () => {
    const tooltip = `
      <span>Tooltip content</span>
    `;
    const wrapper = snapshotFactory({
      propsData: {
        query: { search: 'c-advanced-search-field' },
        tooltip,
      },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
