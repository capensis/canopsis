import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import { MAX_LIMIT } from '@/constants';

import CEntityCategoryField from '@/components/forms/fields/entity/c-entity-category-field.vue';

const stubs = {
  'v-select': {
    props: ['value', 'items'],
    template: `
    <div>
    <select class="v-select" :value="value" @change="$listeners.input($event.target.value)">
      <option v-for="item in items" :value="item.value" :key="item.value">
        {{ item.value }}
      </option>
    </select>
    <slot name="append-item"></slot>
    </div>
  `,
  },
  'v-text-field': {
    props: ['value'],
    template: `
    <input
      :value="value"
      class="v-text-field"
      @input="$listeners.input($event.target.value)"
      @keyup="keyupHandler"
      @blur="blurHandler"
    />
  `,
    methods: {
      blurHandler(event) {
        if (this.$listeners.blur) {
          this.$listeners.blur(event);
        }
      },

      blur() {},

      keyupHandler(event) {
        if (this.$listeners.keyup) {
          this.$listeners.keyup(event);
        }
      },
    },
  },
  'c-help-icon': true,
};
const snapshotStubs = {
  'c-help-icon': true,
};

const entityCategories = [
  {
    _id: 'c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf',
    name: 'Category',
    author: 'root',
    created: 1614861888,
    updated: 1614861888,
  },
  {
    _id: '441a2a17-9036-48a3-9ff7-f393487395a9',
    name: 'Category 2',
    author: 'root',
    created: 1614863990,
    updated: 1614863990,
  },
  {
    _id: '1cae4b8a-f598-480a-ad0c-b0a89a5c2e93',
    name: 'Category 3',
    author: 'root',
    created: 1614864049,
    updated: 1614864049,
  },
  {
    _id: 'c46bffd9-8f5a-4c6c-b045-416e23ab1d44',
    name: 'New category',
    author: 'root',
    created: 1614857014,
    updated: 1614857014,
  },
  {
    _id: 'd2403af7-712d-4353-911e-376f7a8053a7',
    name: 'Second category',
    author: 'root',
    created: 1613620731,
    updated: 1613620731,
  },
  {
    _id: '9bbb623c-7537-4c3b-afc0-0ace4f25a76b',
    name: 'Test category',
    author: 'root',
    created: 1613620721,
    updated: 1613620721,
  },
  {
    _id: 'fd35fcc4-36b0-445d-85be-999cc939047a',
    name: 'categoryasd',
    author: 'root',
    created: 1614864251,
    updated: 1614864251,
  },
  {
    _id: '70bfae47-cfdf-4a2c-9b43-3427f6aabea2',
    name: 'fdsfdf',
    author: 'root',
    created: 1622781697,
    updated: 1622781697,
  },
  {
    _id: 'b3f67a16-019a-4694-9d74-ed762affaa04',
    name: 'fdsfdfs',
    author: 'root',
    created: 1615435601,
    updated: 1615435601,
  },
  {
    _id: 'e1f3e64a-dc99-42ed-af72-d8678f2e62bf',
    name: 'test',
    author: 'root',
    created: 1615442556,
    updated: 1615442556,
  },
  {
    _id: '15094f5a-9472-4700-b0cd-52305f754754',
    name: 'еуые',
    author: 'root',
    created: 1615440560,
    updated: 1615440560,
  },
];

describe('c-entity-category-field', () => {
  const name = 'category';

  const factory = generateShallowRenderer(CEntityCategoryField, {

    stubs,
    store: createMockedStoreModules([{ name: 'entityCategory' }]),

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });
  const snapshotFactory = generateRenderer(CEntityCategoryField, {

    stubs: snapshotStubs,
    store: createMockedStoreModules([{
      name: 'entityCategory',
      getters: {
        pending: false,
        items: entityCategories,
      },
      actions: {
        fetchList: jest.fn(),
        create: jest.fn(),
      },
    }]),

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Check fetch list call', () => {
    const fetchListMock = jest.fn();

    factory({
      store: createMockedStoreModules([{
        name: 'entityCategory',
        getters: {
          items: [],
          pending: true,
        },
        actions: {
          fetchList: fetchListMock,
        },
      }]),
    });

    expect(fetchListMock).toBeCalledWith(expect.any(Object), { params: { limit: MAX_LIMIT } }, undefined);
  });

  it('Check v-validate uniq-name', async () => {
    const [{ name: categoryName }] = entityCategories;
    const createMock = jest.fn();
    const fetchListMock = jest.fn();

    const wrapper = factory({
      propsData: {
        name,
        addable: true,
      },
      store: createMockedStoreModules([{
        name: 'entityCategory',
        getters: {
          items: entityCategories,
          pending: false,
        },
        actions: {
          fetchList: fetchListMock,
          create: createMock,
        },
      }]),
    });

    const validator = wrapper.getValidator();
    const select = wrapper.find('.v-select');

    select.trigger('click');

    const textField = wrapper.find('.v-text-field');

    textField.setValue(categoryName);
    textField.trigger('keyup.enter');

    await flushPromises();

    expect(validator.errors.has(`${name}.create`)).toBeTruthy();
    expect(wrapper).not.toEmit('input');

    wrapper.destroy();
  });

  it('Check creating', async () => {
    const categoryName = 'test';
    const [category] = entityCategories;
    const fetchListMock = jest.fn();
    const createMock = jest.fn(() => category);
    const wrapper = factory({
      propsData: {
        addable: true,
      },

      store: createMockedStoreModules([{
        name: 'entityCategory',
        getters: {
          items: [],
          pending: false,
        },
        actions: {
          fetchList: fetchListMock,
          create: createMock,
        },
      }]),
    });

    const select = wrapper.find('.v-select');

    await select.trigger('click');

    const textField = wrapper.find('.v-text-field');

    textField.setValue(categoryName);
    textField.trigger('keyup.enter');

    await flushPromises();

    expect(createMock).toBeCalledWith(expect.any(Object), { data: { name: categoryName } }, undefined);
    expect(fetchListMock).toBeCalledTimes(2);
    expect(wrapper).toEmit('input', category);
    expect(textField.vm.value).toBe('');
  });

  it('Check clearing by blur', async () => {
    const categoryName = 'test';
    const fetchListMock = jest.fn();
    const createMock = jest.fn();
    const wrapper = factory({
      propsData: {
        addable: true,
      },

      store: createMockedStoreModules([{
        name: 'entityCategory',
        getters: {
          items: [],
          pending: false,
        },
        actions: {
          fetchList: fetchListMock,
          create: createMock,
        },
      }]),
    });

    const select = wrapper.find('.v-select');

    await select.trigger('click');

    const textField = wrapper.find('.v-text-field');

    textField.setValue(categoryName);
    textField.trigger('blur');

    await flushPromises();

    expect(fetchListMock).toBeCalledTimes(1);
    expect(wrapper).not.toEmit('input');
    expect(textField.vm.value).toBe('');
  });

  it('Renders `c-entity-category-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-entity-category-field` with default props and pending', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'entityCategory',
        getters: {
          pending: true,
          items: entityCategories,
        },
        actions: {
          fetchList: jest.fn(),
          create: jest.fn(),
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-entity-category-field` with default props and open select', () => {
    const wrapper = snapshotFactory();

    const select = wrapper.find('.v-input__slot');
    const menuContent = wrapper.findMenu();

    select.trigger('click');

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-entity-category-field` with custom props and open select', () => {
    const wrapper = snapshotFactory({
      propsData: {
        addable: true,
      },
    });

    const select = wrapper.find('.v-input__slot');
    const menuContent = wrapper.findMenu();

    select.trigger('click');

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-entity-category-field` with validator error', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        required: true,
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});
