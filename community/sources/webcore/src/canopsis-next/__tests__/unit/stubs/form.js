export const createFormStub = className => ({
  emits: ['submit'],
  template: `
    <form class="${className}" @submit.prevent="$emit('submit', $event)">
      <slot />
    </form>
  `,
});

export const getFormBlockArrayFieldStub = () => ({
  props: {
    form: {
      type: Array,
      default: () => [],
    },
    value: {
      type: Array,
      default: () => [],
    },
    itemToForm: {
      type: Function,
      required: true,
    },
  },
  emits: ['input', 'update:modelValue'],
  computed: {
    items() {
      return this.form?.length ? this.form : (this.value ?? []);
    },
  },
  template: `
    <div>
      <div v-for="(item, index) in items" :key="index">
        <slot name="item" :index="index" :remove="() => remove(index)" />
      </div>
      <button type="button" @click="add">add</button>
    </div>
  `,
  methods: {
    add() {
      const value = [...this.items, this.itemToForm()];

      this.$emit('input', value);
      this.$emit('update:modelValue', value);
    },
    remove(index) {
      const value = this.items.filter((_, itemIndex) => itemIndex !== index);

      this.$emit('input', value);
      this.$emit('update:modelValue', value);
    },
  },
});

export const getFormGeneralPatternsTabsStub = () => ({
  template: `
    <div>
      <slot name="general" :set-ref="noop" />
      <slot name="patterns" :set-ref="noop" />
      <slot name="additional" :set-ref="noop" />
      <slot name="test-query" />
    </div>
  `,
  methods: {
    noop() {},
  },
});
