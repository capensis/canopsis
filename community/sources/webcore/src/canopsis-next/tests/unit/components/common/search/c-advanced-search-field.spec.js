import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { advancedSearchRuleItemToFormItem } from '@/helpers/search/advanced-search';

import CAdvancedSearchField from '@/components/common/search/c-advanced-search-field.vue';

const stubs = {
  'c-action-btn': true,
  'advanced-search-rules': true,
  'advanced-search-history-btn': true,
};
const snapshotStubs = {
  'c-action-btn': true,
  'advanced-search-rules': true,
  'advanced-search-history-btn': true,
};

describe('c-advanced-search-field', () => {
  const factory = generateShallowRenderer(CAdvancedSearchField, { stubs });
  const snapshotFactory = generateRenderer(CAdvancedSearchField, { stubs: snapshotStubs });

  it('Emit submit when submit method called and validation passes', async () => {
    const wrapper = factory({
      propsData: {
        rules: [{ ...advancedSearchRuleItemToFormItem(), text: 'test' }],
        attributes: [],
      },
    });

    await wrapper.vm.submit();

    expect(wrapper).toEmit('submit', expect.objectContaining({ search: 'test' }));
  });

  it('Emit reset when reset method called', async () => {
    const wrapper = factory({
      propsData: {
        rules: [advancedSearchRuleItemToFormItem()],
        attributes: [],
      },
    });

    await wrapper.vm.reset();

    expect(wrapper).toHaveBeenEmit('reset');
  });

  it('Emit remove:search when history item removed', () => {
    const id = Faker.datatype.uuid();

    const wrapper = factory({
      propsData: {
        rules: [advancedSearchRuleItemToFormItem()],
        attributes: [],
        withHistory: true,
        searches: [{ _id: id }],
      },
    });

    wrapper.vm.removeSearch(id);

    expect(wrapper).toEmit('remove:search', id);
  });

  it('Emit toggle-pin:search when history item pin toggled', () => {
    const id = Faker.datatype.uuid();

    const wrapper = factory({
      propsData: {
        rules: [advancedSearchRuleItemToFormItem()],
        attributes: [],
        withHistory: true,
        searches: [{ _id: id }],
      },
    });

    wrapper.vm.togglePinForSearch(id);

    expect(wrapper).toEmit('toggle-pin:search', id);
  });

  it('Renders `c-advanced-search-field` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rules: [advancedSearchRuleItemToFormItem()],
        attributes: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
