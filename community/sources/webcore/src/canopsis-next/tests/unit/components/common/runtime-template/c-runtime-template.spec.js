import { flushPromises, createVueInstance, generateRenderer } from '@unit/utils/vue';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';

const localVue = createVueInstance();

localVue.mixin({
  data() {
    return {
      mixinDataProperty: 'mixin-data-property',
    };
  },
  computed: {
    mixinComputedProperty() {
      return 'mixin-computed-property';
    },
  },
  methods: {
    mixinMethod() {
      return 'mixin-method-property';
    },
  },
});

const CustomComponent = {
  inject: ['provideProperty'],
  template: '<div class="custom-component">{{ provideProperty }}</div>',
};

const Component = {
  components: { CRuntimeTemplate, CustomComponent },
  provide() {
    return {
      provideProperty: 'provide-property-text',
    };
  },
  props: ['prop-one', 'prop-two', 'prop-three'],
  data() {
    return {
      dataProperty: 'data-property',
      mixinDataProperty: 'override-mixin-data-property',
    };
  },
  computed: {
    computedProperty() {
      return `${this.propTwo}:${this.propThree}`;
    },

    mixinComputedProperty() {
      return 'override-mixin-computed-property';
    },

    template() {
      return `
        <div>
          {{ propOne }}
          {{ computedProperty }}
          {{ mixinComputedProperty }}
          {{ dataProperty }}
          {{ mixinDataProperty }}
          {{ templateProp }}
          {{ getText() }}
          {{ mixinMethod() }}
          <custom-component />
        </div>
      `;
    },
  },
  methods: {
    getText() {
      return 'method-text';
    },

    mixinMethod() {
      return 'override-mixin-method-property';
    },
  },
  template: `
    <c-runtime-template :template="template" :template-props="{ templateProp: 'template-property-text' }" />
  `,
};

describe('c-runtime-template', () => {
  const snapshotFactory = generateRenderer(Component, {
    localVue,
    propsData: {
      propOne: 'prop-one',
      propTwo: 'prop-two',
      propThree: 'prop-three',
    },
  });

  test('Renders `c-runtime-template` props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-runtime-template` updated props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    await wrapper.setProps({
      propOne: 'prop-one-updated',
      propTwo: 'prop-two-updated',
      propThree: 'prop-three-updated',
    });

    expect(wrapper).toMatchSnapshot();
  });
});
